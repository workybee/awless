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

func addManualInfraFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["containerinstance"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []*ecs.ContainerInstance
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.containerinstance.sync") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containerinstance]")
			return resources, objects, nil
		}

		var clusterArns ([]*string)
		if cached, ok := cache.Get("getClustersNames").([]*string); ok && cached != nil {
			clusterArns = cached
		} else {
			res, err := getClustersNames(ctx, conf.APIs.Ecs)
			if err != nil {
				return resources, objects, err
			}
			clusterArns = res
			cache.Store("getClustersNames", res)
		}

		for _, cluster := range clusterArns {
			var badResErr error
			err := conf.APIs.Ecs.ListContainerInstancesPages(&ecs.ListContainerInstancesInput{Cluster: cluster}, func(out *ecs.ListContainerInstancesOutput, lastPage bool) (shouldContinue bool) {
				var containerInstancesOut *ecs.DescribeContainerInstancesOutput
				if len(out.ContainerInstanceArns) == 0 {
					return out.NextToken != nil
				}

				if containerInstancesOut, badResErr = conf.APIs.Ecs.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{Cluster: cluster, ContainerInstances: out.ContainerInstanceArns}); badResErr != nil {
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

	funcs["container"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []*ecs.Container
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.container.sync") {
			conf.Log.Verbose("sync: *disabled* for resource infra[container]")
			return resources, objects, nil
		}

		var tasks ([]*ecs.Task)
		if cached, ok := cache.Get("getAllTasks").([]*ecs.Task); ok && cached != nil {
			tasks = cached
		} else {
			res, err := getAllTasks(ctx, cache, conf.APIs.Ecs)
			if err != nil {
				return resources, objects, err
			}
			tasks = res
			cache.Store("getAllTasks", res)
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

	funcs["containertask"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []*ecs.TaskDefinition
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.containertask.sync") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containertask]")
			return resources, objects, nil
		}

		type resStruct struct {
			res   *ecs.TaskDefinition
			tasks []*ecs.Task
			err   error
		}

		var wg sync.WaitGroup
		resc := make(chan resStruct)

		err := conf.APIs.Ecs.ListTaskDefinitionsPages(&ecs.ListTaskDefinitionsInput{}, func(out *ecs.ListTaskDefinitionsOutput, lastPage bool) (shouldContinue bool) {
			for _, arn := range out.TaskDefinitionArns {
				wg.Add(1)
				go func(taskDefArn *string) {
					defer wg.Done()
					tasksOut, err := conf.APIs.Ecs.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{TaskDefinition: taskDefArn})
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
		if cached, ok := cache.Get("getAllTasks").([]*ecs.Task); ok && cached != nil {
			tasks = cached
		} else {
			res, err := getAllTasks(ctx, cache, conf.APIs.Ecs)
			if err != nil {
				return resources, objects, err
			}
			tasks = res
			cache.Store("getAllTasks", res)
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

	funcs["containercluster"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ecs.Cluster

		if !conf.getBoolDefaultTrue("aws.infra.containercluster.sync") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containercluster]")
			return resources, objects, nil
		}

		var clusterNames ([]*string)
		if cached, ok := cache.Get("getClustersNames").([]*string); ok && cached != nil {
			clusterNames = cached
		} else {
			res, err := getClustersNames(ctx, conf.APIs.Ecs)
			if err != nil {
				return resources, objects, err
			}
			clusterNames = res
			cache.Store("getClustersNames", res)
		}

		for _, clusterArns := range sliceOfSlice(clusterNames, 100) {
			clustersOut, err := conf.APIs.Ecs.DescribeClusters(&ecs.DescribeClustersInput{Clusters: clusterArns})
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

	funcs["listener"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []*elbv2.Listener
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.listener.sync") {
			conf.Log.Verbose("sync: *disabled* for resource infra[listener]")
			return resources, objects, nil
		}

		errc := make(chan error)
		resultc := make(chan *elbv2.Listener)
		var wg sync.WaitGroup

		err := conf.APIs.Elbv2.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
			func(out *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
				for _, lb := range out.LoadBalancers {
					wg.Add(1)
					go func(lb *elbv2.LoadBalancer) {
						defer wg.Done()
						err := conf.APIs.Elbv2.DescribeListenersPages(&elbv2.DescribeListenersInput{LoadBalancerArn: lb.LoadBalancerArn},
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

func addManualAccessFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["user"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*iam.UserDetail

		if !conf.getBoolDefaultTrue("aws.access.user.sync") {
			conf.Log.Verbose("sync: *disabled* for resource access[user]")
			return resources, objects, nil
		}

		var wg sync.WaitGroup
		errc := make(chan error)

		wg.Add(1)
		go func() {
			defer wg.Done()
			var badResErr error
			err := conf.APIs.Iam.GetAccountAuthorizationDetailsPages(&iam.GetAccountAuthorizationDetailsInput{
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

			err := conf.APIs.Iam.ListUsersPages(&iam.ListUsersInput{}, func(page *iam.ListUsersOutput, lastPage bool) bool {
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

	funcs["policy"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*iam.Policy

		if !conf.getBoolDefaultTrue("aws.access.policy.sync") {
			conf.Log.Verbose("sync: *disabled* for resource access[policy]")
			return resources, objects, nil
		}

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
			err := conf.APIs.Iam.ListPoliciesPages(&iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)},
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
			err := conf.APIs.Iam.ListPoliciesPages(&iam.ListPoliciesInput{Scope: awssdk.String("Local")},
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
func addManualStorageFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["bucket"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*s3.Bucket

		if !conf.getBoolDefaultTrue("aws.storage.bucket.sync") {
			conf.Log.Verbose("sync: *disabled* for resource storage[bucket]")
			return resources, objects, nil
		}

		bucketM := &sync.Mutex{}

		err := forEachBucketParallel(ctx, cache, conf.APIs.S3, func(b *s3.Bucket) error {
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

	funcs["s3object"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []*s3.Object
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.storage.s3object.sync") {
			conf.Log.Verbose("sync: *disabled* for resource storage[s3object]")
			return resources, objects, nil
		}

		err := forEachBucketParallel(ctx, cache, conf.APIs.S3, func(b *s3.Bucket) error {
			return fetchObjectsForBucket(ctx, conf.APIs.S3, b, &resources)
		})

		return resources, objects, err
	}
}
func addManualMessagingFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["queue"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []*string
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.messaging.queue.sync") {
			conf.Log.Verbose("sync: *disabled* for resource messaging[queue]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Sqs.ListQueues(&sqs.ListQueuesInput{})
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
				attrs, err := conf.APIs.Sqs.GetQueueAttributes(&sqs.GetQueueAttributesInput{AttributeNames: []*string{awssdk.String("All")}, QueueUrl: url})
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
func addManualDnsFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["record"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []*route53.ResourceRecordSet
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.dns.record.sync") {
			conf.Log.Verbose("sync: *disabled* for resource dns[record]")
			return resources, objects, nil
		}

		zonec := make(chan *route53.HostedZone)
		errc := make(chan error)

		go func() {
			err := conf.APIs.Route53.ListHostedZonesPages(&route53.ListHostedZonesInput{},
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
					err := conf.APIs.Route53.ListResourceRecordSetsPages(&route53.ListResourceRecordSetsInput{HostedZoneId: z.Id},
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
func addManualLambdaFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualMonitoringFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualCdnFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualCloudformationFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
