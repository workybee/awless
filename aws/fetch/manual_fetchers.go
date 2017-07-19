package awsfetch

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/wallix/awless/aws/conv"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/graph"
)

func addManualInfraFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {
	elbv2API := elbv2.New(sess)
	ecsAPI := ecs.New(sess)

	funcs["containerinstance"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var objects []*ecs.ContainerInstance
		var resources []*graph.Resource

		var clusterArns ([]*string)
		if cached := ctx.Value("getClustersNames").([]*string); cached == nil {
			res, err := getClustersNames(ctx, ecsAPI)
			if err != nil {
				return resources, objects, err
			}
			clusterArns = res
			ctx = context.WithValue(ctx, "getClustersNames", res)
		} else {
			clusterArns = cached
		}

		for _, cluster := range clusterArns {
			var badResErr error
			err := ecsAPI.ListContainerInstancesPages(&ecs.ListContainerInstancesInput{Cluster: cluster}, func(out *ecs.ListContainerInstancesOutput, lastPage bool) (shouldContinue bool) {
				var containerInstancesOut *ecs.DescribeContainerInstancesOutput
				if len(out.ContainerInstanceArns) == 0 {
					return out.NextToken != nil
				}

				if containerInstancesOut, badResErr = ecsAPI.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{Cluster: cluster, ContainerInstances: out.ContainerInstanceArns}); badResErr != nil {
					return false
				}

				for _, inst := range containerInstancesOut.ContainerInstances {
					objects = append(objects, inst)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(inst); badResErr != nil {
						return false
					}
					res.Properties[properties.Cluster] = awssdk.StringValue(cluster)
					resources = append(resources, res)
					parent := graph.InitResource(cloud.ContainerCluster, awssdk.StringValue(cluster))

					parent.Relations[rdf.ParentOf] = append(parent.Relations[rdf.ParentOf], res)
				}
				return out.NextToken != nil
			})
			if err != nil {
				return resources, objects, err
			}
			if badResErr != nil {
				return resources, objects, badResErr
			}
		}
		return resources, objects, nil
	}

	funcs["container"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var objects []*ecs.Container
		var resources []*graph.Resource

		var tasks ([]*ecs.Task)
		if cached := ctx.Value("getAllTasks").([]*ecs.Task); cached == nil {
			res, err := getAllTasks(ctx, ecsAPI)
			if err != nil {
				return resources, objects, err
			}
			tasks = res
			ctx = context.WithValue(ctx, "getAllTasks", res)
		} else {
			tasks = cached
		}

		var err error

		for _, task := range tasks {
			for _, container := range task.Containers {
				var res *graph.Resource
				objects = append(objects, container)
				if res, err = awsconv.NewResource(container); err != nil {
					return nil, nil, err
				}
				if task.ClusterArn != nil {
					res.Properties[properties.Cluster] = awssdk.StringValue(task.ClusterArn)
				}
				if task.ContainerInstanceArn != nil {
					res.Properties[properties.ContainerInstance] = awssdk.StringValue(task.ContainerInstanceArn)
				}
				if task.CreatedAt != nil {
					res.Properties[properties.Created] = awssdk.TimeValue(task.CreatedAt)
				}
				if task.StartedAt != nil {
					res.Properties[properties.Launched] = awssdk.TimeValue(task.StartedAt)
				}
				if task.StoppedAt != nil {
					res.Properties[properties.Stopped] = awssdk.TimeValue(task.StoppedAt)
				}
				if task.TaskDefinitionArn != nil {
					res.Properties[properties.ContainerTask] = awssdk.StringValue(task.TaskDefinitionArn)
				}
				if task.Group != nil {
					res.Properties[properties.DeploymentName] = awssdk.StringValue(task.Group)
				}

				clusterParent := graph.InitResource(cloud.ContainerCluster, awssdk.StringValue(task.ClusterArn))
				clusterParent.Relations[rdf.ParentOf] = append(clusterParent.Relations[rdf.ParentOf], res)

				taskParent := graph.InitResource(cloud.ContainerTask, awssdk.StringValue(task.TaskDefinitionArn))
				taskParent.Relations[rdf.ApplyOn] = append(taskParent.Relations[rdf.ApplyOn], res)

				instParent := graph.InitResource(cloud.ContainerInstance, awssdk.StringValue(task.ContainerInstanceArn))
				instParent.Relations[rdf.ApplyOn] = append(instParent.Relations[rdf.ApplyOn], res)

				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["containertask"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var objects []*ecs.TaskDefinition
		var resources []*graph.Resource

		type resStruct struct {
			res   *ecs.TaskDefinition
			tasks []*ecs.Task
			err   error
		}

		var wg sync.WaitGroup
		resc := make(chan resStruct)

		err := ecsAPI.ListTaskDefinitionsPages(&ecs.ListTaskDefinitionsInput{}, func(out *ecs.ListTaskDefinitionsOutput, lastPage bool) (shouldContinue bool) {
			for _, arn := range out.TaskDefinitionArns {
				wg.Add(1)
				go func(taskDefArn *string) {
					defer wg.Done()
					tasksOut, err := ecsAPI.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{TaskDefinition: taskDefArn})
					if err != nil {
						resc <- resStruct{err: err}
						return
					}
					resc <- resStruct{res: tasksOut.TaskDefinition}
				}(arn)
			}
			return out.NextToken != nil
		})
		if err != nil {
			return resources, objects, err
		}

		go func() {
			wg.Wait()
			close(resc)
		}()

		var tasks ([]*ecs.Task)
		if cached := ctx.Value("getAllTasks").([]*ecs.Task); cached == nil {
			res, err := getAllTasks(ctx, ecsAPI)
			if err != nil {
				return resources, objects, err
			}
			tasks = res
			ctx = context.WithValue(ctx, "getAllTasks", res)
		} else {
			tasks = cached
		}

		var errors []string

		for res := range resc {
			if res.err != nil {
				errors = appendIfNotInSlice(errors, res.err.Error())
				continue
			}
			objects = append(objects, res.res)
			var graphres *graph.Resource
			if graphres, err = awsconv.NewResource(res.res); err != nil {
				errors = appendIfNotInSlice(errors, err.Error())
				continue
			}
			var deployments []*graph.KeyValue
			var runningServicesCount, stoppedServicesCount, runningTasksCount, stoppedTasksCount uint
			for _, t := range tasks {
				if awssdk.StringValue(t.TaskDefinitionArn) == awssdk.StringValue(res.res.TaskDefinitionArn) {
					group := awssdk.StringValue(t.Group)
					state := strings.ToLower(awssdk.StringValue(t.LastStatus))
					clusterArn := awssdk.StringValue(t.ClusterArn)
					if strings.HasPrefix(group, "service:") {
						switch state {
						case "stopped":
							stoppedServicesCount++
							deployments = append(deployments, &graph.KeyValue{arnToName(clusterArn), group[len("service:"):] + " (stopped service)"})
						case "running":
							runningServicesCount++
							deployments = append(deployments, &graph.KeyValue{arnToName(clusterArn), group[len("service:"):] + " (running service)"})
						}
					}
					if strings.HasPrefix(group, "family:") {
						switch state {
						case "stopped":
							deployments = append(deployments, &graph.KeyValue{arnToName(clusterArn), group[len("family:"):] + " (stopped task)"})
							stoppedTasksCount++
						case "running":
							deployments = append(deployments, &graph.KeyValue{arnToName(clusterArn), group[len("family:"):] + " (running task)"})
							runningTasksCount++
						}
					}
				}
			}
			if len(deployments) > 0 {
				graphres.Properties[properties.Deployments] = deployments
			}
			switch {
			case runningServicesCount+stoppedServicesCount+runningTasksCount+stoppedTasksCount == 0:
				if state := strings.ToLower(awssdk.StringValue(res.res.Status)); state == "active" {
					graphres.Properties[properties.State] = "ready"
				} else {
					graphres.Properties[properties.State] = state
				}
			default:
				var stateSl []string
				if runningServicesCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s running", runningServicesCount, pluralizeIfNeeded("service", runningServicesCount)))
				}
				if stoppedServicesCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s stopped", stoppedServicesCount, pluralizeIfNeeded("service", runningServicesCount)))
				}
				if runningTasksCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s running", runningTasksCount, pluralizeIfNeeded("task", runningServicesCount)))
				}
				if stoppedTasksCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s stopped", stoppedTasksCount, pluralizeIfNeeded("task", runningServicesCount)))
				}
				if len(stateSl) > 0 {
					graphres.Properties[properties.State] = strings.Join(stateSl, " ")
				}
			}

			resources = append(resources, graphres)
		}

		if len(errors) > 0 {
			err = fmt.Errorf(strings.Join(errors, "; "))
		}

		return resources, objects, err
	}

	funcs["containercluster"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ecs.Cluster

		var clusterNames ([]*string)
		if cached := ctx.Value("getClustersNames").([]*string); cached == nil {
			res, err := getClustersNames(ctx, ecsAPI)
			if err != nil {
				return resources, objects, err
			}
			clusterNames = res
			ctx = context.WithValue(ctx, "getClustersNames", res)
		} else {
			clusterNames = cached
		}

		for _, clusterArns := range sliceOfSlice(clusterNames, 100) {
			clustersOut, err := ecsAPI.DescribeClusters(&ecs.DescribeClustersInput{Clusters: clusterArns})
			if err != nil {
				return resources, objects, err
			}

			for _, cluster := range clustersOut.Clusters {
				objects = append(objects, cluster)
				res, err := awsconv.NewResource(cluster)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}
		return resources, objects, nil
	}

	funcs["listener"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var objects []*elbv2.Listener
		var resources []*graph.Resource

		errc := make(chan error)
		resultc := make(chan *elbv2.Listener)
		var wg sync.WaitGroup

		err := elbv2API.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
			func(out *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
				for _, lb := range out.LoadBalancers {
					wg.Add(1)
					go func(lb *elbv2.LoadBalancer) {
						defer wg.Done()
						err := elbv2API.DescribeListenersPages(&elbv2.DescribeListenersInput{LoadBalancerArn: lb.LoadBalancerArn},
							func(out *elbv2.DescribeListenersOutput, lastPage bool) (shouldContinue bool) {
								for _, listen := range out.Listeners {
									resultc <- listen
								}
								return out.NextMarker != nil
							})
						if err != nil {
							errc <- err
						}
					}(lb)
				}
				return out.NextMarker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		go func() {
			wg.Wait()
			close(resultc)
		}()

		for {
			select {
			case err := <-errc:
				if err != nil {
					return resources, objects, err
				}
			case listener, ok := <-resultc:
				if !ok {
					return resources, objects, nil
				}
				objects = append(objects, listener)
				res, err := awsconv.NewResource(listener)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}
	}
}

func addManualAccessFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {
	iamAPI := iam.New(sess)

	funcs["user"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*iam.UserDetail

		var wg sync.WaitGroup
		errc := make(chan error)

		wg.Add(1)
		go func() {
			defer wg.Done()
			var badResErr error
			err := iamAPI.GetAccountAuthorizationDetailsPages(&iam.GetAccountAuthorizationDetailsInput{
				Filter: []*string{
					awssdk.String(iam.EntityTypeUser),
				},
			},
				func(out *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool) {
					for _, output := range out.UserDetailList {
						objects = append(objects, output)
						var res *graph.Resource
						res, badResErr = awsconv.NewResource(output)
						if badResErr != nil {
							return false
						}
						resources = append(resources, res)
					}
					return out.Marker != nil
				})
			if err != nil {
				errc <- err
				return
			}
			if badResErr != nil {
				errc <- err
				return
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := iamAPI.ListUsersPages(&iam.ListUsersInput{}, func(page *iam.ListUsersOutput, lastPage bool) bool {
				for _, user := range page.Users {
					res, badResErr := awsconv.NewResource(user)
					if badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return page.Marker != nil
			})
			if err != nil {
				errc <- err
			}
		}()

		go func() {
			wg.Wait()
			close(errc)
		}()

		for err := range errc {
			if err != nil {
				return resources, objects, err
			}
		}

		return resources, objects, nil
	}

	funcs["policy"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*iam.Policy

		errc := make(chan error)
		policiesc := make(chan *iam.Policy)

		processPagePolicies := func(page *iam.ListPoliciesOutput) bool {
			for _, p := range page.Policies {
				policiesc <- p
				res, rerr := awsconv.NewResource(p)
				if rerr != nil {
					return false
				}
				resources = append(resources, res)
			}
			return page.Marker != nil
		}

		var wg sync.WaitGroup

		// Return all policies that are only attached
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := iamAPI.ListPoliciesPages(&iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)},
				func(out *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool) {
					return processPagePolicies(out)
				})
			if err != nil {
				errc <- err
			}
		}()

		// Return only self managed policies (local scope)
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := iamAPI.ListPoliciesPages(&iam.ListPoliciesInput{Scope: awssdk.String("Local")},
				func(out *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool) {
					return processPagePolicies(out)
				})
			if err != nil {
				errc <- err
			}
		}()

		go func() {
			wg.Wait()
			close(errc)
			close(policiesc)
		}()

		for {
			select {
			case err := <-errc:
				if err != nil {
					return resources, objects, err
				}
			case p, ok := <-policiesc:
				if !ok {
					return resources, objects, nil
				}
				objects = append(objects, p)
			}
		}
	}
}
func addManualStorageFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {
	s3API := s3.New(sess)

	funcs["bucket"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*s3.Bucket
		bucketM := &sync.Mutex{}

		err := forEachBucketParallel(ctx, s3API, func(b *s3.Bucket) error {
			bucketM.Lock()
			objects = append(objects, b)
			bucketM.Unlock()
			res, err := awsconv.NewResource(b)
			if err != nil {
				return fmt.Errorf("build resource for bucket `%s`: %s", awssdk.StringValue(b.Name), err)
			}
			resources = append(resources, res)
			return nil
		})
		return resources, objects, err
	}

	funcs["s3object"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var objects []*s3.Object
		var resources []*graph.Resource

		err := forEachBucketParallel(ctx, s3API, func(b *s3.Bucket) error {
			return fetchObjectsForBucket(ctx, s3API, b, &resources)
		})

		return resources, objects, err
	}
}
func addManualMessagingFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {
	sqsAPI := sqs.New(sess)

	funcs["queue"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var objects []*string
		var resources []*graph.Resource

		out, err := sqsAPI.ListQueues(&sqs.ListQueuesInput{})
		if err != nil {
			return nil, objects, err
		}
		errc := make(chan error)
		var wg sync.WaitGroup

		for _, output := range out.QueueUrls {
			objects = append(objects, output)
			wg.Add(1)
			go func(url *string) {
				defer wg.Done()
				res := graph.InitResource(cloud.Queue, awssdk.StringValue(url))
				res.Properties[properties.ID] = awssdk.StringValue(url)
				attrs, err := sqsAPI.GetQueueAttributes(&sqs.GetQueueAttributesInput{AttributeNames: []*string{awssdk.String("All")}, QueueUrl: url})
				if e, ok := err.(awserr.RequestFailure); ok && (e.Code() == sqs.ErrCodeQueueDoesNotExist || e.Code() == sqs.ErrCodeQueueDeletedRecently) {
					return
				}
				if err != nil {
					errc <- err
					return
				}
				for k, v := range attrs.Attributes {
					switch k {
					case "ApproximateNumberOfMessages":
						count, err := strconv.Atoi(awssdk.StringValue(v))
						if err != nil {
							errc <- err
						}
						res.Properties[properties.ApproximateMessageCount] = count
					case "CreatedTimestamp":
						if vv := awssdk.StringValue(v); vv != "" {
							timestamp, err := strconv.ParseInt(vv, 10, 64)
							if err != nil {
								errc <- err
							}
							res.Properties[properties.Created] = time.Unix(int64(timestamp), 0)
						}
					case "LastModifiedTimestamp":
						if vv := awssdk.StringValue(v); vv != "" {
							timestamp, err := strconv.ParseInt(vv, 10, 64)
							if err != nil {
								errc <- err
							}
							res.Properties[properties.Modified] = time.Unix(int64(timestamp), 0)
						}
					case "QueueArn":
						res.Properties[properties.Arn] = awssdk.StringValue(v)
					case "DelaySeconds":
						delay, err := strconv.Atoi(awssdk.StringValue(v))
						if err != nil {
							errc <- err
						}
						res.Properties[properties.Delay] = delay
					}

				}
				resources = append(resources, res)
			}(output)

		}

		go func() {
			wg.Wait()
			close(errc)
		}()

		for err := range errc {
			if err != nil {
				return resources, objects, err
			}
		}

		return resources, objects, nil
	}
}
func addManualDnsFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {
	dnsAPI := route53.New(sess)

	funcs["record"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var objects []*route53.ResourceRecordSet
		var resources []*graph.Resource

		zonec := make(chan *route53.HostedZone)
		errc := make(chan error)

		go func() {
			err := dnsAPI.ListHostedZonesPages(&route53.ListHostedZonesInput{},
				func(out *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
					for _, output := range out.HostedZones {
						zonec <- output
					}
					return out.NextMarker != nil
				})
			if err != nil {
				errc <- err
			}
			close(zonec)
		}()

		resultc := make(chan *route53.ResourceRecordSet)

		go func() {
			var wg sync.WaitGroup

			for zone := range zonec {
				wg.Add(1)
				go func(z *route53.HostedZone) {
					defer wg.Done()
					err := dnsAPI.ListResourceRecordSetsPages(&route53.ListResourceRecordSetsInput{HostedZoneId: z.Id},
						func(out *route53.ListResourceRecordSetsOutput, lastPage bool) (shouldContinue bool) {
							for _, output := range out.ResourceRecordSets {
								resultc <- output
								res, err := awsconv.NewResource(output)
								if err != nil {
									errc <- err
								}
								resources = append(resources, res)
								parent, err := awsconv.InitResource(z)
								if err != nil {
									errc <- err
								}
								parent.Relations[rdf.ParentOf] = append(parent.Relations[rdf.ParentOf], res)
							}
							return out.NextRecordName != nil
						})
					if err != nil {
						errc <- err
					}
				}(zone)
			}

			go func() {
				wg.Wait()
				close(resultc)
			}()
		}()

		for {
			select {
			case err := <-errc:
				if err != nil {
					return resources, objects, err
				}
			case record, ok := <-resultc:
				if !ok {
					return resources, objects, nil
				}
				objects = append(objects, record)
			}
		}
	}
}
func addManualLambdaFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)         {}
func addManualMonitoringFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)     {}
func addManualCdnFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)            {}
func addManualCloudformationFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {}
