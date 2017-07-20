// Auto generated implementation for the AWS cloud service

/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsservices

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"context"
	"errors"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/aws/fetch"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

const accessDenied = "Access Denied"

var ServiceNames = []string{
	"infra",
	"access",
	"storage",
	"messaging",
	"dns",
	"lambda",
	"monitoring",
	"cdn",
	"cloudformation",
}

var ResourceTypes = []string{
	"instance",
	"subnet",
	"vpc",
	"keypair",
	"securitygroup",
	"volume",
	"internetgateway",
	"natgateway",
	"routetable",
	"availabilityzone",
	"image",
	"importimagetask",
	"elasticip",
	"snapshot",
	"loadbalancer",
	"targetgroup",
	"listener",
	"database",
	"dbsubnetgroup",
	"launchconfiguration",
	"scalinggroup",
	"scalingpolicy",
	"repository",
	"containercluster",
	"containertask",
	"container",
	"containerinstance",
	"user",
	"group",
	"role",
	"policy",
	"accesskey",
	"instanceprofile",
	"bucket",
	"s3object",
	"subscription",
	"topic",
	"queue",
	"zone",
	"record",
	"function",
	"metric",
	"alarm",
	"distribution",
	"stack",
}

var ServicePerAPI = map[string]string{
	"ec2":         "infra",
	"elbv2":       "infra",
	"rds":         "infra",
	"autoscaling": "infra",
	"ecr":         "infra",
	"ecs":         "infra",
	"applicationautoscaling": "infra",
	"iam":            "access",
	"sts":            "access",
	"s3":             "storage",
	"sns":            "messaging",
	"sqs":            "messaging",
	"route53":        "dns",
	"lambda":         "lambda",
	"cloudwatch":     "monitoring",
	"cloudfront":     "cdn",
	"cloudformation": "cloudformation",
}

var ServicePerResourceType = map[string]string{
	"instance":            "infra",
	"subnet":              "infra",
	"vpc":                 "infra",
	"keypair":             "infra",
	"securitygroup":       "infra",
	"volume":              "infra",
	"internetgateway":     "infra",
	"natgateway":          "infra",
	"routetable":          "infra",
	"availabilityzone":    "infra",
	"image":               "infra",
	"importimagetask":     "infra",
	"elasticip":           "infra",
	"snapshot":            "infra",
	"loadbalancer":        "infra",
	"targetgroup":         "infra",
	"listener":            "infra",
	"database":            "infra",
	"dbsubnetgroup":       "infra",
	"launchconfiguration": "infra",
	"scalinggroup":        "infra",
	"scalingpolicy":       "infra",
	"repository":          "infra",
	"containercluster":    "infra",
	"containertask":       "infra",
	"container":           "infra",
	"containerinstance":   "infra",
	"user":                "access",
	"group":               "access",
	"role":                "access",
	"policy":              "access",
	"accesskey":           "access",
	"instanceprofile":     "access",
	"bucket":              "storage",
	"s3object":            "storage",
	"subscription":        "messaging",
	"topic":               "messaging",
	"queue":               "messaging",
	"zone":                "dns",
	"record":              "dns",
	"function":            "lambda",
	"metric":              "monitoring",
	"alarm":               "monitoring",
	"distribution":        "cdn",
	"stack":               "cloudformation",
}

var APIPerResourceType = map[string]string{
	"instance":            "ec2",
	"subnet":              "ec2",
	"vpc":                 "ec2",
	"keypair":             "ec2",
	"securitygroup":       "ec2",
	"volume":              "ec2",
	"internetgateway":     "ec2",
	"natgateway":          "ec2",
	"routetable":          "ec2",
	"availabilityzone":    "ec2",
	"image":               "ec2",
	"importimagetask":     "ec2",
	"elasticip":           "ec2",
	"snapshot":            "ec2",
	"loadbalancer":        "elbv2",
	"targetgroup":         "elbv2",
	"listener":            "elbv2",
	"database":            "rds",
	"dbsubnetgroup":       "rds",
	"launchconfiguration": "autoscaling",
	"scalinggroup":        "autoscaling",
	"scalingpolicy":       "autoscaling",
	"repository":          "ecr",
	"containercluster":    "ecs",
	"containertask":       "ecs",
	"container":           "ecs",
	"containerinstance":   "ecs",
	"user":                "iam",
	"group":               "iam",
	"role":                "iam",
	"policy":              "iam",
	"accesskey":           "iam",
	"instanceprofile":     "iam",
	"bucket":              "s3",
	"s3object":            "s3",
	"subscription":        "sns",
	"topic":               "sns",
	"queue":               "sqs",
	"zone":                "route53",
	"record":              "route53",
	"function":            "lambda",
	"metric":              "cloudwatch",
	"alarm":               "cloudwatch",
	"distribution":        "cloudfront",
	"stack":               "cloudformation",
}

var GlobalServices = []string{
	"access",
	"dns",
	"cdn",
}

type Infra struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	ec2iface.EC2API
	elbv2iface.ELBV2API
	rdsiface.RDSAPI
	autoscalingiface.AutoScalingAPI
	ecriface.ECRAPI
	ecsiface.ECSAPI
	applicationautoscalingiface.ApplicationAutoScalingAPI
}

func NewInfra(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	ec2API := ec2.New(sess)
	elbv2API := elbv2.New(sess)
	rdsAPI := rds.New(sess)
	autoscalingAPI := autoscaling.New(sess)
	ecrAPI := ecr.New(sess)
	ecsAPI := ecs.New(sess)
	applicationautoscalingAPI := applicationautoscaling.New(sess)

	fetchConfig := awsfetch.NewConfig(
		ec2API,
		elbv2API,
		rdsAPI,
		autoscalingAPI,
		ecrAPI,
		ecsAPI,
		applicationautoscalingAPI,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Infra{
		EC2API:         ec2API,
		ELBV2API:       elbv2API,
		RDSAPI:         rdsAPI,
		AutoScalingAPI: autoscalingAPI,
		ECRAPI:         ecrAPI,
		ECSAPI:         ecsAPI,
		ApplicationAutoScalingAPI: applicationautoscalingAPI,
		fetcher:                   fetch.NewFetcher(awsfetch.BuildInfraFetchFuncs(fetchConfig)),
		config:                    awsconf,
		region:                    region,
		log:                       log,
	}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Region() string {
	return s.region
}

func (s *Infra) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewEc2Driver(s.EC2API),
		awsdriver.NewElbv2Driver(s.ELBV2API),
		awsdriver.NewRdsDriver(s.RDSAPI),
		awsdriver.NewAutoscalingDriver(s.AutoScalingAPI),
		awsdriver.NewEcrDriver(s.ECRAPI),
		awsdriver.NewEcsDriver(s.ECSAPI),
		awsdriver.NewApplicationautoscalingDriver(s.ApplicationAutoScalingAPI),
	}
}

func (s *Infra) ResourceTypes() []string {
	return []string{
		"instance",
		"subnet",
		"vpc",
		"keypair",
		"securitygroup",
		"volume",
		"internetgateway",
		"natgateway",
		"routetable",
		"availabilityzone",
		"image",
		"importimagetask",
		"elasticip",
		"snapshot",
		"loadbalancer",
		"targetgroup",
		"listener",
		"database",
		"dbsubnetgroup",
		"launchconfiguration",
		"scalinggroup",
		"scalingpolicy",
		"repository",
		"containercluster",
		"containertask",
		"container",
		"containerinstance",
	}
}

func (s *Infra) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.infra.instance.sync", true) {
		list, ok := s.fetcher.Get("instance_objects").([]*ec2.Instance)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Instance' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["instance"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.Instance) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.subnet.sync", true) {
		list, ok := s.fetcher.Get("subnet_objects").([]*ec2.Subnet)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Subnet' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["subnet"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.Subnet) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.vpc.sync", true) {
		list, ok := s.fetcher.Get("vpc_objects").([]*ec2.Vpc)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Vpc' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["vpc"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.Vpc) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.keypair.sync", true) {
		list, ok := s.fetcher.Get("keypair_objects").([]*ec2.KeyPairInfo)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.KeyPairInfo' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["keypair"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.KeyPairInfo) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.securitygroup.sync", true) {
		list, ok := s.fetcher.Get("securitygroup_objects").([]*ec2.SecurityGroup)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.SecurityGroup' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["securitygroup"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.SecurityGroup) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.volume.sync", true) {
		list, ok := s.fetcher.Get("volume_objects").([]*ec2.Volume)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Volume' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["volume"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.Volume) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.internetgateway.sync", true) {
		list, ok := s.fetcher.Get("internetgateway_objects").([]*ec2.InternetGateway)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.InternetGateway' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["internetgateway"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.InternetGateway) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.natgateway.sync", true) {
		list, ok := s.fetcher.Get("natgateway_objects").([]*ec2.NatGateway)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.NatGateway' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["natgateway"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.NatGateway) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.routetable.sync", true) {
		list, ok := s.fetcher.Get("routetable_objects").([]*ec2.RouteTable)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.RouteTable' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["routetable"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.RouteTable) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.availabilityzone.sync", true) {
		list, ok := s.fetcher.Get("availabilityzone_objects").([]*ec2.AvailabilityZone)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.AvailabilityZone' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["availabilityzone"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.AvailabilityZone) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.image.sync", true) {
		list, ok := s.fetcher.Get("image_objects").([]*ec2.Image)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Image' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["image"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.Image) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.importimagetask.sync", true) {
		list, ok := s.fetcher.Get("importimagetask_objects").([]*ec2.ImportImageTask)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.ImportImageTask' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["importimagetask"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.ImportImageTask) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.elasticip.sync", true) {
		list, ok := s.fetcher.Get("elasticip_objects").([]*ec2.Address)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Address' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["elasticip"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.Address) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.snapshot.sync", true) {
		list, ok := s.fetcher.Get("snapshot_objects").([]*ec2.Snapshot)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Snapshot' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["snapshot"] {
				wg.Add(1)
				go func(f addParentFn, res *ec2.Snapshot) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.loadbalancer.sync", true) {
		list, ok := s.fetcher.Get("loadbalancer_objects").([]*elbv2.LoadBalancer)
		if !ok {
			return gph, errors.New("cannot cast to '[]*elbv2.LoadBalancer' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["loadbalancer"] {
				wg.Add(1)
				go func(f addParentFn, res *elbv2.LoadBalancer) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.targetgroup.sync", true) {
		list, ok := s.fetcher.Get("targetgroup_objects").([]*elbv2.TargetGroup)
		if !ok {
			return gph, errors.New("cannot cast to '[]*elbv2.TargetGroup' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["targetgroup"] {
				wg.Add(1)
				go func(f addParentFn, res *elbv2.TargetGroup) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.listener.sync", true) {
		list, ok := s.fetcher.Get("listener_objects").([]*elbv2.Listener)
		if !ok {
			return gph, errors.New("cannot cast to '[]*elbv2.Listener' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["listener"] {
				wg.Add(1)
				go func(f addParentFn, res *elbv2.Listener) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.database.sync", true) {
		list, ok := s.fetcher.Get("database_objects").([]*rds.DBInstance)
		if !ok {
			return gph, errors.New("cannot cast to '[]*rds.DBInstance' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["database"] {
				wg.Add(1)
				go func(f addParentFn, res *rds.DBInstance) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.dbsubnetgroup.sync", true) {
		list, ok := s.fetcher.Get("dbsubnetgroup_objects").([]*rds.DBSubnetGroup)
		if !ok {
			return gph, errors.New("cannot cast to '[]*rds.DBSubnetGroup' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["dbsubnetgroup"] {
				wg.Add(1)
				go func(f addParentFn, res *rds.DBSubnetGroup) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.launchconfiguration.sync", true) {
		list, ok := s.fetcher.Get("launchconfiguration_objects").([]*autoscaling.LaunchConfiguration)
		if !ok {
			return gph, errors.New("cannot cast to '[]*autoscaling.LaunchConfiguration' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["launchconfiguration"] {
				wg.Add(1)
				go func(f addParentFn, res *autoscaling.LaunchConfiguration) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.scalinggroup.sync", true) {
		list, ok := s.fetcher.Get("scalinggroup_objects").([]*autoscaling.Group)
		if !ok {
			return gph, errors.New("cannot cast to '[]*autoscaling.Group' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["scalinggroup"] {
				wg.Add(1)
				go func(f addParentFn, res *autoscaling.Group) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.scalingpolicy.sync", true) {
		list, ok := s.fetcher.Get("scalingpolicy_objects").([]*autoscaling.ScalingPolicy)
		if !ok {
			return gph, errors.New("cannot cast to '[]*autoscaling.ScalingPolicy' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["scalingpolicy"] {
				wg.Add(1)
				go func(f addParentFn, res *autoscaling.ScalingPolicy) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.repository.sync", true) {
		list, ok := s.fetcher.Get("repository_objects").([]*ecr.Repository)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ecr.Repository' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["repository"] {
				wg.Add(1)
				go func(f addParentFn, res *ecr.Repository) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.containercluster.sync", true) {
		list, ok := s.fetcher.Get("containercluster_objects").([]*ecs.Cluster)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ecs.Cluster' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["containercluster"] {
				wg.Add(1)
				go func(f addParentFn, res *ecs.Cluster) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.containertask.sync", true) {
		list, ok := s.fetcher.Get("containertask_objects").([]*ecs.TaskDefinition)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ecs.TaskDefinition' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["containertask"] {
				wg.Add(1)
				go func(f addParentFn, res *ecs.TaskDefinition) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.container.sync", true) {
		list, ok := s.fetcher.Get("container_objects").([]*ecs.Container)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ecs.Container' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["container"] {
				wg.Add(1)
				go func(f addParentFn, res *ecs.Container) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.infra.containerinstance.sync", true) {
		list, ok := s.fetcher.Get("containerinstance_objects").([]*ecs.ContainerInstance)
		if !ok {
			return gph, errors.New("cannot cast to '[]*ecs.ContainerInstance' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["containerinstance"] {
				wg.Add(1)
				go func(f addParentFn, res *ecs.ContainerInstance) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Infra) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Infra) IsSyncDisabled() bool {
	return !s.config.getBool("aws.infra.sync", true)
}

type Access struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	iamiface.IAMAPI
	stsiface.STSAPI
}

func NewAccess(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := "global"
	iamAPI := iam.New(sess)
	stsAPI := sts.New(sess)

	fetchConfig := awsfetch.NewConfig(
		iamAPI,
		stsAPI,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Access{
		IAMAPI:  iamAPI,
		STSAPI:  stsAPI,
		fetcher: fetch.NewFetcher(awsfetch.BuildAccessFetchFuncs(fetchConfig)),
		config:  awsconf,
		region:  region,
		log:     log,
	}
}

func (s *Access) Name() string {
	return "access"
}

func (s *Access) Region() string {
	return s.region
}

func (s *Access) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewIamDriver(s.IAMAPI),
		awsdriver.NewStsDriver(s.STSAPI),
	}
}

func (s *Access) ResourceTypes() []string {
	return []string{
		"user",
		"group",
		"role",
		"policy",
		"accesskey",
		"instanceprofile",
	}
}

func (s *Access) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.access.user.sync", true) {
		list, ok := s.fetcher.Get("user_objects").([]*iam.UserDetail)
		if !ok {
			return gph, errors.New("cannot cast to '[]*iam.UserDetail' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["user"] {
				wg.Add(1)
				go func(f addParentFn, res *iam.UserDetail) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.access.group.sync", true) {
		list, ok := s.fetcher.Get("group_objects").([]*iam.GroupDetail)
		if !ok {
			return gph, errors.New("cannot cast to '[]*iam.GroupDetail' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["group"] {
				wg.Add(1)
				go func(f addParentFn, res *iam.GroupDetail) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.access.role.sync", true) {
		list, ok := s.fetcher.Get("role_objects").([]*iam.RoleDetail)
		if !ok {
			return gph, errors.New("cannot cast to '[]*iam.RoleDetail' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["role"] {
				wg.Add(1)
				go func(f addParentFn, res *iam.RoleDetail) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.access.policy.sync", true) {
		list, ok := s.fetcher.Get("policy_objects").([]*iam.Policy)
		if !ok {
			return gph, errors.New("cannot cast to '[]*iam.Policy' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["policy"] {
				wg.Add(1)
				go func(f addParentFn, res *iam.Policy) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.access.accesskey.sync", true) {
		list, ok := s.fetcher.Get("accesskey_objects").([]*iam.AccessKeyMetadata)
		if !ok {
			return gph, errors.New("cannot cast to '[]*iam.AccessKeyMetadata' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["accesskey"] {
				wg.Add(1)
				go func(f addParentFn, res *iam.AccessKeyMetadata) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.access.instanceprofile.sync", true) {
		list, ok := s.fetcher.Get("instanceprofile_objects").([]*iam.InstanceProfile)
		if !ok {
			return gph, errors.New("cannot cast to '[]*iam.InstanceProfile' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["instanceprofile"] {
				wg.Add(1)
				go func(f addParentFn, res *iam.InstanceProfile) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Access) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Access) IsSyncDisabled() bool {
	return !s.config.getBool("aws.access.sync", true)
}

type Storage struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	s3iface.S3API
}

func NewStorage(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	s3API := s3.New(sess)

	fetchConfig := awsfetch.NewConfig(
		s3API,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Storage{
		S3API:   s3API,
		fetcher: fetch.NewFetcher(awsfetch.BuildStorageFetchFuncs(fetchConfig)),
		config:  awsconf,
		region:  region,
		log:     log,
	}
}

func (s *Storage) Name() string {
	return "storage"
}

func (s *Storage) Region() string {
	return s.region
}

func (s *Storage) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewS3Driver(s.S3API),
	}
}

func (s *Storage) ResourceTypes() []string {
	return []string{
		"bucket",
		"s3object",
	}
}

func (s *Storage) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.storage.bucket.sync", true) {
		list, ok := s.fetcher.Get("bucket_objects").([]*s3.Bucket)
		if !ok {
			return gph, errors.New("cannot cast to '[]*s3.Bucket' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["bucket"] {
				wg.Add(1)
				go func(f addParentFn, res *s3.Bucket) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.storage.s3object.sync", true) {
		list, ok := s.fetcher.Get("s3object_objects").([]*s3.Object)
		if !ok {
			return gph, errors.New("cannot cast to '[]*s3.Object' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["s3object"] {
				wg.Add(1)
				go func(f addParentFn, res *s3.Object) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Storage) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Storage) IsSyncDisabled() bool {
	return !s.config.getBool("aws.storage.sync", true)
}

type Messaging struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	snsiface.SNSAPI
	sqsiface.SQSAPI
}

func NewMessaging(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	snsAPI := sns.New(sess)
	sqsAPI := sqs.New(sess)

	fetchConfig := awsfetch.NewConfig(
		snsAPI,
		sqsAPI,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Messaging{
		SNSAPI:  snsAPI,
		SQSAPI:  sqsAPI,
		fetcher: fetch.NewFetcher(awsfetch.BuildMessagingFetchFuncs(fetchConfig)),
		config:  awsconf,
		region:  region,
		log:     log,
	}
}

func (s *Messaging) Name() string {
	return "messaging"
}

func (s *Messaging) Region() string {
	return s.region
}

func (s *Messaging) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewSnsDriver(s.SNSAPI),
		awsdriver.NewSqsDriver(s.SQSAPI),
	}
}

func (s *Messaging) ResourceTypes() []string {
	return []string{
		"subscription",
		"topic",
		"queue",
	}
}

func (s *Messaging) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.messaging.subscription.sync", true) {
		list, ok := s.fetcher.Get("subscription_objects").([]*sns.Subscription)
		if !ok {
			return gph, errors.New("cannot cast to '[]*sns.Subscription' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["subscription"] {
				wg.Add(1)
				go func(f addParentFn, res *sns.Subscription) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.messaging.topic.sync", true) {
		list, ok := s.fetcher.Get("topic_objects").([]*sns.Topic)
		if !ok {
			return gph, errors.New("cannot cast to '[]*sns.Topic' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["topic"] {
				wg.Add(1)
				go func(f addParentFn, res *sns.Topic) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.messaging.queue.sync", true) {
		list, ok := s.fetcher.Get("queue_objects").([]*string)
		if !ok {
			return gph, errors.New("cannot cast to '[]*string' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["queue"] {
				wg.Add(1)
				go func(f addParentFn, res *string) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Messaging) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Messaging) IsSyncDisabled() bool {
	return !s.config.getBool("aws.messaging.sync", true)
}

type Dns struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	route53iface.Route53API
}

func NewDns(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := "global"
	route53API := route53.New(sess)

	fetchConfig := awsfetch.NewConfig(
		route53API,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Dns{
		Route53API: route53API,
		fetcher:    fetch.NewFetcher(awsfetch.BuildDnsFetchFuncs(fetchConfig)),
		config:     awsconf,
		region:     region,
		log:        log,
	}
}

func (s *Dns) Name() string {
	return "dns"
}

func (s *Dns) Region() string {
	return s.region
}

func (s *Dns) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewRoute53Driver(s.Route53API),
	}
}

func (s *Dns) ResourceTypes() []string {
	return []string{
		"zone",
		"record",
	}
}

func (s *Dns) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.dns.zone.sync", true) {
		list, ok := s.fetcher.Get("zone_objects").([]*route53.HostedZone)
		if !ok {
			return gph, errors.New("cannot cast to '[]*route53.HostedZone' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["zone"] {
				wg.Add(1)
				go func(f addParentFn, res *route53.HostedZone) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.dns.record.sync", true) {
		list, ok := s.fetcher.Get("record_objects").([]*route53.ResourceRecordSet)
		if !ok {
			return gph, errors.New("cannot cast to '[]*route53.ResourceRecordSet' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["record"] {
				wg.Add(1)
				go func(f addParentFn, res *route53.ResourceRecordSet) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Dns) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Dns) IsSyncDisabled() bool {
	return !s.config.getBool("aws.dns.sync", true)
}

type Lambda struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	lambdaiface.LambdaAPI
}

func NewLambda(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	lambdaAPI := lambda.New(sess)

	fetchConfig := awsfetch.NewConfig(
		lambdaAPI,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Lambda{
		LambdaAPI: lambdaAPI,
		fetcher:   fetch.NewFetcher(awsfetch.BuildLambdaFetchFuncs(fetchConfig)),
		config:    awsconf,
		region:    region,
		log:       log,
	}
}

func (s *Lambda) Name() string {
	return "lambda"
}

func (s *Lambda) Region() string {
	return s.region
}

func (s *Lambda) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewLambdaDriver(s.LambdaAPI),
	}
}

func (s *Lambda) ResourceTypes() []string {
	return []string{
		"function",
	}
}

func (s *Lambda) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.lambda.function.sync", true) {
		list, ok := s.fetcher.Get("function_objects").([]*lambda.FunctionConfiguration)
		if !ok {
			return gph, errors.New("cannot cast to '[]*lambda.FunctionConfiguration' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["function"] {
				wg.Add(1)
				go func(f addParentFn, res *lambda.FunctionConfiguration) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Lambda) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Lambda) IsSyncDisabled() bool {
	return !s.config.getBool("aws.lambda.sync", true)
}

type Monitoring struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	cloudwatchiface.CloudWatchAPI
}

func NewMonitoring(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	cloudwatchAPI := cloudwatch.New(sess)

	fetchConfig := awsfetch.NewConfig(
		cloudwatchAPI,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Monitoring{
		CloudWatchAPI: cloudwatchAPI,
		fetcher:       fetch.NewFetcher(awsfetch.BuildMonitoringFetchFuncs(fetchConfig)),
		config:        awsconf,
		region:        region,
		log:           log,
	}
}

func (s *Monitoring) Name() string {
	return "monitoring"
}

func (s *Monitoring) Region() string {
	return s.region
}

func (s *Monitoring) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewCloudwatchDriver(s.CloudWatchAPI),
	}
}

func (s *Monitoring) ResourceTypes() []string {
	return []string{
		"metric",
		"alarm",
	}
}

func (s *Monitoring) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.monitoring.metric.sync", true) {
		list, ok := s.fetcher.Get("metric_objects").([]*cloudwatch.Metric)
		if !ok {
			return gph, errors.New("cannot cast to '[]*cloudwatch.Metric' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["metric"] {
				wg.Add(1)
				go func(f addParentFn, res *cloudwatch.Metric) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
	if s.config.getBool("aws.monitoring.alarm.sync", true) {
		list, ok := s.fetcher.Get("alarm_objects").([]*cloudwatch.MetricAlarm)
		if !ok {
			return gph, errors.New("cannot cast to '[]*cloudwatch.MetricAlarm' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["alarm"] {
				wg.Add(1)
				go func(f addParentFn, res *cloudwatch.MetricAlarm) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Monitoring) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Monitoring) IsSyncDisabled() bool {
	return !s.config.getBool("aws.monitoring.sync", true)
}

type Cdn struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	cloudfrontiface.CloudFrontAPI
}

func NewCdn(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := "global"
	cloudfrontAPI := cloudfront.New(sess)

	fetchConfig := awsfetch.NewConfig(
		cloudfrontAPI,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Cdn{
		CloudFrontAPI: cloudfrontAPI,
		fetcher:       fetch.NewFetcher(awsfetch.BuildCdnFetchFuncs(fetchConfig)),
		config:        awsconf,
		region:        region,
		log:           log,
	}
}

func (s *Cdn) Name() string {
	return "cdn"
}

func (s *Cdn) Region() string {
	return s.region
}

func (s *Cdn) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewCloudfrontDriver(s.CloudFrontAPI),
	}
}

func (s *Cdn) ResourceTypes() []string {
	return []string{
		"distribution",
	}
}

func (s *Cdn) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.cdn.distribution.sync", true) {
		list, ok := s.fetcher.Get("distribution_objects").([]*cloudfront.DistributionSummary)
		if !ok {
			return gph, errors.New("cannot cast to '[]*cloudfront.DistributionSummary' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["distribution"] {
				wg.Add(1)
				go func(f addParentFn, res *cloudfront.DistributionSummary) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Cdn) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Cdn) IsSyncDisabled() bool {
	return !s.config.getBool("aws.cdn.sync", true)
}

type Cloudformation struct {
	fetcher fetch.Fetcher
	region  string
	config  config
	log     *logger.Logger
	cloudformationiface.CloudFormationAPI
}

func NewCloudformation(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	cloudformationAPI := cloudformation.New(sess)

	fetchConfig := awsfetch.NewConfig(
		cloudformationAPI,
	)
	fetchConfig.Extra = awsconf
	fetchConfig.Log = log

	return &Cloudformation{
		CloudFormationAPI: cloudformationAPI,
		fetcher:           fetch.NewFetcher(awsfetch.BuildCloudformationFetchFuncs(fetchConfig)),
		config:            awsconf,
		region:            region,
		log:               log,
	}
}

func (s *Cloudformation) Name() string {
	return "cloudformation"
}

func (s *Cloudformation) Region() string {
	return s.region
}

func (s *Cloudformation) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewCloudformationDriver(s.CloudFormationAPI),
	}
}

func (s *Cloudformation) ResourceTypes() []string {
	return []string{
		"stack",
	}
}

func (s *Cloudformation) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
	gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	if s.config.getBool("aws.cloudformation.stack.sync", true) {
		list, ok := s.fetcher.Get("stack_objects").([]*cloudformation.Stack)
		if !ok {
			return gph, errors.New("cannot cast to '[]*cloudformation.Stack' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["stack"] {
				wg.Add(1)
				go func(f addParentFn, res *cloudformation.Stack) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Cloudformation) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
	return s.fetcher.FetchByType(ctx, t)
}

func (s *Cloudformation) IsSyncDisabled() bool {
	return !s.config.getBool("aws.cloudformation.sync", true)
}
