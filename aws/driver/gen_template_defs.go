/* Copyright 2017 WALLIX

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

// DO NOT EDIT
// This file was automatically generated with go generate
package awsdriver

import (
	"github.com/wallix/awless/template"
)

var APIPerTemplateDefName = map[string]string{
	"createvpc":                 "ec2",
	"deletevpc":                 "ec2",
	"createsubnet":              "ec2",
	"updatesubnet":              "ec2",
	"deletesubnet":              "ec2",
	"createinstance":            "ec2",
	"updateinstance":            "ec2",
	"deleteinstance":            "ec2",
	"startinstance":             "ec2",
	"stopinstance":              "ec2",
	"checkinstance":             "ec2",
	"createsecuritygroup":       "ec2",
	"updatesecuritygroup":       "ec2",
	"deletesecuritygroup":       "ec2",
	"checksecuritygroup":        "ec2",
	"attachsecuritygroup":       "ec2",
	"detachsecuritygroup":       "ec2",
	"copyimage":                 "ec2",
	"importimage":               "ec2",
	"deleteimage":               "ec2",
	"createvolume":              "ec2",
	"deletevolume":              "ec2",
	"attachvolume":              "ec2",
	"detachvolume":              "ec2",
	"createsnapshot":            "ec2",
	"deletesnapshot":            "ec2",
	"copysnapshot":              "ec2",
	"createinternetgateway":     "ec2",
	"deleteinternetgateway":     "ec2",
	"attachinternetgateway":     "ec2",
	"detachinternetgateway":     "ec2",
	"createroutetable":          "ec2",
	"deleteroutetable":          "ec2",
	"attachroutetable":          "ec2",
	"detachroutetable":          "ec2",
	"createroute":               "ec2",
	"deleteroute":               "ec2",
	"createtag":                 "ec2",
	"deletetag":                 "ec2",
	"createkeypair":             "ec2",
	"deletekeypair":             "ec2",
	"createelasticip":           "ec2",
	"deleteelasticip":           "ec2",
	"attachelasticip":           "ec2",
	"detachelasticip":           "ec2",
	"createloadbalancer":        "elbv2",
	"deleteloadbalancer":        "elbv2",
	"checkloadbalancer":         "elbv2",
	"createlistener":            "elbv2",
	"deletelistener":            "elbv2",
	"createtargetgroup":         "elbv2",
	"deletetargetgroup":         "elbv2",
	"attachinstance":            "elbv2",
	"detachinstance":            "elbv2",
	"createlaunchconfiguration": "autoscaling",
	"deletelaunchconfiguration": "autoscaling",
	"createscalinggroup":        "autoscaling",
	"updatescalinggroup":        "autoscaling",
	"deletescalinggroup":        "autoscaling",
	"checkscalinggroup":         "autoscaling",
	"createscalingpolicy":       "autoscaling",
	"deletescalingpolicy":       "autoscaling",
	"createdatabase":            "rds",
	"deletedatabase":            "rds",
	"createdbsubnetgroup":       "rds",
	"deletedbsubnetgroup":       "rds",
	"createuser":                "iam",
	"deleteuser":                "iam",
	"attachuser":                "iam",
	"detachuser":                "iam",
	"createaccesskey":           "iam",
	"deleteaccesskey":           "iam",
	"createloginprofile":        "iam",
	"updateloginprofile":        "iam",
	"deleteloginprofile":        "iam",
	"creategroup":               "iam",
	"deletegroup":               "iam",
	"createrole":                "iam",
	"deleterole":                "iam",
	"attachrole":                "iam",
	"detachrole":                "iam",
	"createinstanceprofile":     "iam",
	"deleteinstanceprofile":     "iam",
	"createpolicy":              "iam",
	"deletepolicy":              "iam",
	"attachpolicy":              "iam",
	"detachpolicy":              "iam",
	"createbucket":              "s3",
	"deletebucket":              "s3",
	"creates3object":            "s3",
	"deletes3object":            "s3",
	"createtopic":               "sns",
	"deletetopic":               "sns",
	"createsubscription":        "sns",
	"deletesubscription":        "sns",
	"createqueue":               "sqs",
	"deletequeue":               "sqs",
	"createzone":                "route53",
	"deletezone":                "route53",
	"createrecord":              "route53",
	"deleterecord":              "route53",
	"createfunction":            "lambda",
	"deletefunction":            "lambda",
	"createalarm":               "cloudwatch",
	"deletealarm":               "cloudwatch",
	"startalarm":                "cloudwatch",
	"stopalarm":                 "cloudwatch",
	"attachalarm":               "cloudwatch",
	"detachalarm":               "cloudwatch",
}

var AWSTemplatesDefinitions = map[string]template.Definition{
	"createvpc": {
		Action:         "create",
		Entity:         "vpc",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "cidr"}},
		ExtraParams:    []*template.Param{{Name: "name"}},
	},
	"deletevpc": {
		Action:         "delete",
		Entity:         "vpc",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createsubnet": {
		Action:         "create",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "cidr"}, {Name: "vpc"}},
		ExtraParams:    []*template.Param{{Name: "availabilityzone"}, {Name: "name"}},
	},
	"updatesubnet": {
		Action:         "update",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{{Name: "public"}},
	},
	"deletesubnet": {
		Action:         "delete",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createinstance": {
		Action:         "create",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "count"}, {Name: "image"}, {Name: "name"}, {Name: "subnet"}, {Name: "type"}},
		ExtraParams:    []*template.Param{{Name: "ip"}, {Name: "keypair"}, {Name: "lock"}, {Name: "role"}, {Name: "securitygroup"}, {Name: "userdata"}},
	},
	"updateinstance": {
		Action:         "update",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{{Name: "lock"}, {Name: "type"}},
	},
	"deleteinstance": {
		Action:         "delete",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"startinstance": {
		Action:         "start",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"stopinstance": {
		Action:         "stop",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"checkinstance": {
		Action:         "check",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "state"}, {Name: "timeout"}},
		ExtraParams:    []*template.Param{},
	},
	"createsecuritygroup": {
		Action:         "create",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "description"}, {Name: "name"}, {Name: "vpc"}},
		ExtraParams:    []*template.Param{},
	},
	"updatesecuritygroup": {
		Action:         "update",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "cidr"}, {Name: "id"}, {Name: "protocol"}},
		ExtraParams:    []*template.Param{{Name: "inbound"}, {Name: "outbound"}, {Name: "portrange"}},
	},
	"deletesecuritygroup": {
		Action:         "delete",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"checksecuritygroup": {
		Action:         "check",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "state"}, {Name: "timeout"}},
		ExtraParams:    []*template.Param{},
	},
	"attachsecuritygroup": {
		Action:         "attach",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{{Name: "instance"}},
	},
	"detachsecuritygroup": {
		Action:         "detach",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{{Name: "instance"}},
	},
	"copyimage": {
		Action:         "copy",
		Entity:         "image",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "name"}, {Name: "source-id"}, {Name: "source-region"}},
		ExtraParams:    []*template.Param{{Name: "description"}, {Name: "encrypted"}},
	},
	"importimage": {
		Action:         "import",
		Entity:         "image",
		Api:            "ec2",
		RequiredParams: []*template.Param{},
		ExtraParams:    []*template.Param{{Name: "architecture"}, {Name: "bucket"}, {Name: "description"}, {Name: "license"}, {Name: "platform"}, {Name: "role"}, {Name: "s3object"}, {Name: "snapshot"}, {Name: "url"}},
	},
	"deleteimage": {
		Action:         "delete",
		Entity:         "image",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "delete-snapshots"}, {Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createvolume": {
		Action:         "create",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "availabilityzone"}, {Name: "size"}},
		ExtraParams:    []*template.Param{},
	},
	"deletevolume": {
		Action:         "delete",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"attachvolume": {
		Action:         "attach",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "device"}, {Name: "id"}, {Name: "instance"}},
		ExtraParams:    []*template.Param{},
	},
	"detachvolume": {
		Action:         "detach",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "device"}, {Name: "id"}, {Name: "instance"}},
		ExtraParams:    []*template.Param{{Name: "force"}},
	},
	"createsnapshot": {
		Action:         "create",
		Entity:         "snapshot",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "volume"}},
		ExtraParams:    []*template.Param{{Name: "description"}},
	},
	"deletesnapshot": {
		Action:         "delete",
		Entity:         "snapshot",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"copysnapshot": {
		Action:         "copy",
		Entity:         "snapshot",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "source-id"}, {Name: "source-region"}},
		ExtraParams:    []*template.Param{{Name: "description"}, {Name: "encrypted"}},
	},
	"createinternetgateway": {
		Action:         "create",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []*template.Param{},
		ExtraParams:    []*template.Param{},
	},
	"deleteinternetgateway": {
		Action:         "delete",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"attachinternetgateway": {
		Action:         "attach",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "vpc"}},
		ExtraParams:    []*template.Param{},
	},
	"detachinternetgateway": {
		Action:         "detach",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "vpc"}},
		ExtraParams:    []*template.Param{},
	},
	"createroutetable": {
		Action:         "create",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "vpc"}},
		ExtraParams:    []*template.Param{},
	},
	"deleteroutetable": {
		Action:         "delete",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"attachroutetable": {
		Action:         "attach",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "subnet"}},
		ExtraParams:    []*template.Param{},
	},
	"detachroutetable": {
		Action:         "detach",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "association"}},
		ExtraParams:    []*template.Param{},
	},
	"createroute": {
		Action:         "create",
		Entity:         "route",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "cidr"}, {Name: "gateway"}, {Name: "table"}},
		ExtraParams:    []*template.Param{},
	},
	"deleteroute": {
		Action:         "delete",
		Entity:         "route",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "cidr"}, {Name: "table"}},
		ExtraParams:    []*template.Param{},
	},
	"createtag": {
		Action:         "create",
		Entity:         "tag",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "key"}, {Name: "resource"}, {Name: "value"}},
		ExtraParams:    []*template.Param{},
	},
	"deletetag": {
		Action:         "delete",
		Entity:         "tag",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "key"}, {Name: "resource"}, {Name: "value"}},
		ExtraParams:    []*template.Param{},
	},
	"createkeypair": {
		Action:         "create",
		Entity:         "keypair",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{{Name: "encrypted"}},
	},
	"deletekeypair": {
		Action:         "delete",
		Entity:         "keypair",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createelasticip": {
		Action:         "create",
		Entity:         "elasticip",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "domain"}},
		ExtraParams:    []*template.Param{},
	},
	"deleteelasticip": {
		Action:         "delete",
		Entity:         "elasticip",
		Api:            "ec2",
		RequiredParams: []*template.Param{},
		ExtraParams:    []*template.Param{{Name: "id"}, {Name: "ip"}},
	},
	"attachelasticip": {
		Action:         "attach",
		Entity:         "elasticip",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{{Name: "allow-reassociation"}, {Name: "instance"}, {Name: "networkinterface"}, {Name: "privateip"}},
	},
	"detachelasticip": {
		Action:         "detach",
		Entity:         "elasticip",
		Api:            "ec2",
		RequiredParams: []*template.Param{{Name: "association"}},
		ExtraParams:    []*template.Param{},
	},
	"createloadbalancer": {
		Action:         "create",
		Entity:         "loadbalancer",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "name"}, {Name: "subnets"}},
		ExtraParams:    []*template.Param{{Name: "iptype"}, {Name: "scheme"}, {Name: "securitygroups"}},
	},
	"deleteloadbalancer": {
		Action:         "delete",
		Entity:         "loadbalancer",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"checkloadbalancer": {
		Action:         "check",
		Entity:         "loadbalancer",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "state"}, {Name: "timeout"}},
		ExtraParams:    []*template.Param{},
	},
	"createlistener": {
		Action:         "create",
		Entity:         "listener",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "actiontype"}, {Name: "loadbalancer"}, {Name: "port"}, {Name: "protocol"}, {Name: "target"}},
		ExtraParams:    []*template.Param{{Name: "certificate"}, {Name: "sslpolicy"}},
	},
	"deletelistener": {
		Action:         "delete",
		Entity:         "listener",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createtargetgroup": {
		Action:         "create",
		Entity:         "targetgroup",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "name"}, {Name: "port"}, {Name: "protocol"}, {Name: "vpc"}},
		ExtraParams:    []*template.Param{{Name: "healthcheckinterval"}, {Name: "healthcheckpath"}, {Name: "healthcheckport"}, {Name: "healthcheckprotocol"}, {Name: "healthchecktimeout"}, {Name: "healthythreshold"}, {Name: "matcher"}, {Name: "unhealthythreshold"}},
	},
	"deletetargetgroup": {
		Action:         "delete",
		Entity:         "targetgroup",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"attachinstance": {
		Action:         "attach",
		Entity:         "instance",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "targetgroup"}},
		ExtraParams:    []*template.Param{{Name: "port"}},
	},
	"detachinstance": {
		Action:         "detach",
		Entity:         "instance",
		Api:            "elbv2",
		RequiredParams: []*template.Param{{Name: "id"}, {Name: "targetgroup"}},
		ExtraParams:    []*template.Param{},
	},
	"createlaunchconfiguration": {
		Action:         "create",
		Entity:         "launchconfiguration",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "image"}, {Name: "name"}, {Name: "type"}},
		ExtraParams:    []*template.Param{{Name: "keypair"}, {Name: "public"}, {Name: "role"}, {Name: "securitygroups"}, {Name: "spotprice"}, {Name: "userdata"}},
	},
	"deletelaunchconfiguration": {
		Action:         "delete",
		Entity:         "launchconfiguration",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"createscalinggroup": {
		Action:         "create",
		Entity:         "scalinggroup",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "launchconfiguration"}, {Name: "max-size"}, {Name: "min-size"}, {Name: "name"}, {Name: "subnets"}},
		ExtraParams:    []*template.Param{{Name: "cooldown"}, {Name: "desired-capacity"}, {Name: "healthcheck-grace-period"}, {Name: "healthcheck-type"}, {Name: "new-instances-protected"}, {Name: "targetgroups"}},
	},
	"updatescalinggroup": {
		Action:         "update",
		Entity:         "scalinggroup",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{{Name: "cooldown"}, {Name: "desired-capacity"}, {Name: "healthcheck-grace-period"}, {Name: "healthcheck-type"}, {Name: "launchconfiguration"}, {Name: "max-size"}, {Name: "min-size"}, {Name: "new-instances-protected"}, {Name: "subnets"}},
	},
	"deletescalinggroup": {
		Action:         "delete",
		Entity:         "scalinggroup",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{{Name: "force"}},
	},
	"checkscalinggroup": {
		Action:         "check",
		Entity:         "scalinggroup",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "count"}, {Name: "name"}, {Name: "timeout"}},
		ExtraParams:    []*template.Param{},
	},
	"createscalingpolicy": {
		Action:         "create",
		Entity:         "scalingpolicy",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "adjustment-scaling"}, {Name: "adjustment-type"}, {Name: "name"}, {Name: "scalinggroup"}},
		ExtraParams:    []*template.Param{{Name: "adjustment-magnitude"}, {Name: "cooldown"}, {Name: "metric-aggregation"}},
	},
	"deletescalingpolicy": {
		Action:         "delete",
		Entity:         "scalingpolicy",
		Api:            "autoscaling",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createdatabase": {
		Action:         "create",
		Entity:         "database",
		Api:            "rds",
		RequiredParams: []*template.Param{{Name: "engine"}, {Name: "id"}, {Name: "password"}, {Name: "size"}, {Name: "type"}, {Name: "username"}},
		ExtraParams:    []*template.Param{{Name: "autoupgrade"}, {Name: "availabilityzone"}, {Name: "backupretention"}, {Name: "backupwindow"}, {Name: "cluster"}, {Name: "dbname"}, {Name: "dbsecuritygroups"}, {Name: "domain"}, {Name: "encrypted"}, {Name: "iamrole"}, {Name: "iops"}, {Name: "license"}, {Name: "maintenancewindow"}, {Name: "multiaz"}, {Name: "optiongroup"}, {Name: "parametergroup"}, {Name: "port"}, {Name: "public"}, {Name: "storagetype"}, {Name: "subnetgroup"}, {Name: "timezone"}, {Name: "version"}, {Name: "vpcsecuritygroups"}},
	},
	"deletedatabase": {
		Action:         "delete",
		Entity:         "database",
		Api:            "rds",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{{Name: "skipsnapshot"}, {Name: "snapshotid"}},
	},
	"createdbsubnetgroup": {
		Action:         "create",
		Entity:         "dbsubnetgroup",
		Api:            "rds",
		RequiredParams: []*template.Param{{Name: "description"}, {Name: "name"}, {Name: "subnets"}},
		ExtraParams:    []*template.Param{},
	},
	"deletedbsubnetgroup": {
		Action:         "delete",
		Entity:         "dbsubnetgroup",
		Api:            "rds",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createuser": {
		Action:         "create",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"deleteuser": {
		Action:         "delete",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"attachuser": {
		Action:         "attach",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "group"}, {Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"detachuser": {
		Action:         "detach",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "group"}, {Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"createaccesskey": {
		Action:         "create",
		Entity:         "accesskey",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "user"}},
		ExtraParams:    []*template.Param{},
	},
	"deleteaccesskey": {
		Action:         "delete",
		Entity:         "accesskey",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createloginprofile": {
		Action:         "create",
		Entity:         "loginprofile",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "password"}, {Name: "username"}},
		ExtraParams:    []*template.Param{{Name: "password-reset"}},
	},
	"updateloginprofile": {
		Action:         "update",
		Entity:         "loginprofile",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "password"}, {Name: "username"}},
		ExtraParams:    []*template.Param{{Name: "password-reset"}},
	},
	"deleteloginprofile": {
		Action:         "delete",
		Entity:         "loginprofile",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "username"}},
		ExtraParams:    []*template.Param{},
	},
	"creategroup": {
		Action:         "create",
		Entity:         "group",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"deletegroup": {
		Action:         "delete",
		Entity:         "group",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"createrole": {
		Action:         "create",
		Entity:         "role",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{{Name: "principal-account"}, {Name: "principal-service"}, {Name: "principal-user"}, {Name: "sleep-after"}},
	},
	"deleterole": {
		Action:         "delete",
		Entity:         "role",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"attachrole": {
		Action:         "attach",
		Entity:         "role",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "instanceprofile"}, {Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"detachrole": {
		Action:         "detach",
		Entity:         "role",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "instanceprofile"}, {Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"createinstanceprofile": {
		Action:         "create",
		Entity:         "instanceprofile",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"deleteinstanceprofile": {
		Action:         "delete",
		Entity:         "instanceprofile",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"createpolicy": {
		Action:         "create",
		Entity:         "policy",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{{Name: "action"}, {Name: "description"}, {Name: "effect"}, {Name: "resource"}},
	},
	"deletepolicy": {
		Action:         "delete",
		Entity:         "policy",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "arn"}},
		ExtraParams:    []*template.Param{},
	},
	"attachpolicy": {
		Action:         "attach",
		Entity:         "policy",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "arn"}},
		ExtraParams:    []*template.Param{{Name: "group"}, {Name: "role"}, {Name: "user"}},
	},
	"detachpolicy": {
		Action:         "detach",
		Entity:         "policy",
		Api:            "iam",
		RequiredParams: []*template.Param{{Name: "arn"}},
		ExtraParams:    []*template.Param{{Name: "group"}, {Name: "role"}, {Name: "user"}},
	},
	"createbucket": {
		Action:         "create",
		Entity:         "bucket",
		Api:            "s3",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"deletebucket": {
		Action:         "delete",
		Entity:         "bucket",
		Api:            "s3",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"creates3object": {
		Action:         "create",
		Entity:         "s3object",
		Api:            "s3",
		RequiredParams: []*template.Param{{Name: "bucket"}, {Name: "file"}},
		ExtraParams:    []*template.Param{{Name: "name"}},
	},
	"deletes3object": {
		Action:         "delete",
		Entity:         "s3object",
		Api:            "s3",
		RequiredParams: []*template.Param{{Name: "bucket"}, {Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"createtopic": {
		Action:         "create",
		Entity:         "topic",
		Api:            "sns",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"deletetopic": {
		Action:         "delete",
		Entity:         "topic",
		Api:            "sns",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createsubscription": {
		Action:         "create",
		Entity:         "subscription",
		Api:            "sns",
		RequiredParams: []*template.Param{{Name: "endpoint"}, {Name: "protocol"}, {Name: "topic"}},
		ExtraParams:    []*template.Param{},
	},
	"deletesubscription": {
		Action:         "delete",
		Entity:         "subscription",
		Api:            "sns",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createqueue": {
		Action:         "create",
		Entity:         "queue",
		Api:            "sqs",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{{Name: "delay"}, {Name: "maxMsgSize"}, {Name: "msgWait"}, {Name: "policy"}, {Name: "redrivePolicy"}, {Name: "retentionPeriod"}, {Name: "visibilityTimeout"}},
	},
	"deletequeue": {
		Action:         "delete",
		Entity:         "queue",
		Api:            "sqs",
		RequiredParams: []*template.Param{{Name: "url"}},
		ExtraParams:    []*template.Param{},
	},
	"createzone": {
		Action:         "create",
		Entity:         "zone",
		Api:            "route53",
		RequiredParams: []*template.Param{{Name: "callerreference"}, {Name: "name"}},
		ExtraParams:    []*template.Param{{Name: "comment"}, {Name: "delegationsetid"}, {Name: "isprivate"}, {Name: "vpcid"}, {Name: "vpcregion"}},
	},
	"deletezone": {
		Action:         "delete",
		Entity:         "zone",
		Api:            "route53",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{},
	},
	"createrecord": {
		Action:         "create",
		Entity:         "record",
		Api:            "route53",
		RequiredParams: []*template.Param{{Name: "name"}, {Name: "ttl"}, {Name: "type"}, {Name: "value"}, {Name: "zone"}},
		ExtraParams:    []*template.Param{{Name: "comment"}},
	},
	"deleterecord": {
		Action:         "delete",
		Entity:         "record",
		Api:            "route53",
		RequiredParams: []*template.Param{{Name: "name"}, {Name: "ttl"}, {Name: "type"}, {Name: "value"}, {Name: "zone"}},
		ExtraParams:    []*template.Param{},
	},
	"createfunction": {
		Action:         "create",
		Entity:         "function",
		Api:            "lambda",
		RequiredParams: []*template.Param{{Name: "handler"}, {Name: "name"}, {Name: "role"}, {Name: "runtime"}},
		ExtraParams:    []*template.Param{{Name: "bucket"}, {Name: "description"}, {Name: "memory"}, {Name: "object"}, {Name: "objectversion"}, {Name: "publish"}, {Name: "timeout"}, {Name: "zipfile"}},
	},
	"deletefunction": {
		Action:         "delete",
		Entity:         "function",
		Api:            "lambda",
		RequiredParams: []*template.Param{{Name: "id"}},
		ExtraParams:    []*template.Param{{Name: "version"}},
	},
	"createalarm": {
		Action:         "create",
		Entity:         "alarm",
		Api:            "cloudwatch",
		RequiredParams: []*template.Param{{Name: "evaluation-periods"}, {Name: "metric"}, {Name: "name"}, {Name: "namespace"}, {Name: "operator"}, {Name: "period"}, {Name: "statistic-function"}, {Name: "threshold"}},
		ExtraParams:    []*template.Param{{Name: "alarm-actions"}, {Name: "description"}, {Name: "dimensions"}, {Name: "enabled"}, {Name: "insufficientdata-actions"}, {Name: "ok-actions"}, {Name: "unit"}},
	},
	"deletealarm": {
		Action:         "delete",
		Entity:         "alarm",
		Api:            "cloudwatch",
		RequiredParams: []*template.Param{{Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"startalarm": {
		Action:         "start",
		Entity:         "alarm",
		Api:            "cloudwatch",
		RequiredParams: []*template.Param{{Name: "names"}},
		ExtraParams:    []*template.Param{},
	},
	"stopalarm": {
		Action:         "stop",
		Entity:         "alarm",
		Api:            "cloudwatch",
		RequiredParams: []*template.Param{{Name: "names"}},
		ExtraParams:    []*template.Param{},
	},
	"attachalarm": {
		Action:         "attach",
		Entity:         "alarm",
		Api:            "cloudwatch",
		RequiredParams: []*template.Param{{Name: "action-arn"}, {Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
	"detachalarm": {
		Action:         "detach",
		Entity:         "alarm",
		Api:            "cloudwatch",
		RequiredParams: []*template.Param{{Name: "action-arn"}, {Name: "name"}},
		ExtraParams:    []*template.Param{},
	},
}

func DriverSupportedActions() map[string][]string {
	supported := make(map[string][]string)
	supported["create"] = append(supported["create"], "vpc")
	supported["delete"] = append(supported["delete"], "vpc")
	supported["create"] = append(supported["create"], "subnet")
	supported["update"] = append(supported["update"], "subnet")
	supported["delete"] = append(supported["delete"], "subnet")
	supported["create"] = append(supported["create"], "instance")
	supported["update"] = append(supported["update"], "instance")
	supported["delete"] = append(supported["delete"], "instance")
	supported["start"] = append(supported["start"], "instance")
	supported["stop"] = append(supported["stop"], "instance")
	supported["check"] = append(supported["check"], "instance")
	supported["create"] = append(supported["create"], "securitygroup")
	supported["update"] = append(supported["update"], "securitygroup")
	supported["delete"] = append(supported["delete"], "securitygroup")
	supported["check"] = append(supported["check"], "securitygroup")
	supported["attach"] = append(supported["attach"], "securitygroup")
	supported["detach"] = append(supported["detach"], "securitygroup")
	supported["copy"] = append(supported["copy"], "image")
	supported["import"] = append(supported["import"], "image")
	supported["delete"] = append(supported["delete"], "image")
	supported["create"] = append(supported["create"], "volume")
	supported["delete"] = append(supported["delete"], "volume")
	supported["attach"] = append(supported["attach"], "volume")
	supported["detach"] = append(supported["detach"], "volume")
	supported["create"] = append(supported["create"], "snapshot")
	supported["delete"] = append(supported["delete"], "snapshot")
	supported["copy"] = append(supported["copy"], "snapshot")
	supported["create"] = append(supported["create"], "internetgateway")
	supported["delete"] = append(supported["delete"], "internetgateway")
	supported["attach"] = append(supported["attach"], "internetgateway")
	supported["detach"] = append(supported["detach"], "internetgateway")
	supported["create"] = append(supported["create"], "routetable")
	supported["delete"] = append(supported["delete"], "routetable")
	supported["attach"] = append(supported["attach"], "routetable")
	supported["detach"] = append(supported["detach"], "routetable")
	supported["create"] = append(supported["create"], "route")
	supported["delete"] = append(supported["delete"], "route")
	supported["create"] = append(supported["create"], "tag")
	supported["delete"] = append(supported["delete"], "tag")
	supported["create"] = append(supported["create"], "keypair")
	supported["delete"] = append(supported["delete"], "keypair")
	supported["create"] = append(supported["create"], "elasticip")
	supported["delete"] = append(supported["delete"], "elasticip")
	supported["attach"] = append(supported["attach"], "elasticip")
	supported["detach"] = append(supported["detach"], "elasticip")
	supported["create"] = append(supported["create"], "loadbalancer")
	supported["delete"] = append(supported["delete"], "loadbalancer")
	supported["check"] = append(supported["check"], "loadbalancer")
	supported["create"] = append(supported["create"], "listener")
	supported["delete"] = append(supported["delete"], "listener")
	supported["create"] = append(supported["create"], "targetgroup")
	supported["delete"] = append(supported["delete"], "targetgroup")
	supported["attach"] = append(supported["attach"], "instance")
	supported["detach"] = append(supported["detach"], "instance")
	supported["create"] = append(supported["create"], "launchconfiguration")
	supported["delete"] = append(supported["delete"], "launchconfiguration")
	supported["create"] = append(supported["create"], "scalinggroup")
	supported["update"] = append(supported["update"], "scalinggroup")
	supported["delete"] = append(supported["delete"], "scalinggroup")
	supported["check"] = append(supported["check"], "scalinggroup")
	supported["create"] = append(supported["create"], "scalingpolicy")
	supported["delete"] = append(supported["delete"], "scalingpolicy")
	supported["create"] = append(supported["create"], "database")
	supported["delete"] = append(supported["delete"], "database")
	supported["create"] = append(supported["create"], "dbsubnetgroup")
	supported["delete"] = append(supported["delete"], "dbsubnetgroup")
	supported["create"] = append(supported["create"], "user")
	supported["delete"] = append(supported["delete"], "user")
	supported["attach"] = append(supported["attach"], "user")
	supported["detach"] = append(supported["detach"], "user")
	supported["create"] = append(supported["create"], "accesskey")
	supported["delete"] = append(supported["delete"], "accesskey")
	supported["create"] = append(supported["create"], "loginprofile")
	supported["update"] = append(supported["update"], "loginprofile")
	supported["delete"] = append(supported["delete"], "loginprofile")
	supported["create"] = append(supported["create"], "group")
	supported["delete"] = append(supported["delete"], "group")
	supported["create"] = append(supported["create"], "role")
	supported["delete"] = append(supported["delete"], "role")
	supported["attach"] = append(supported["attach"], "role")
	supported["detach"] = append(supported["detach"], "role")
	supported["create"] = append(supported["create"], "instanceprofile")
	supported["delete"] = append(supported["delete"], "instanceprofile")
	supported["create"] = append(supported["create"], "policy")
	supported["delete"] = append(supported["delete"], "policy")
	supported["attach"] = append(supported["attach"], "policy")
	supported["detach"] = append(supported["detach"], "policy")
	supported["create"] = append(supported["create"], "bucket")
	supported["delete"] = append(supported["delete"], "bucket")
	supported["create"] = append(supported["create"], "s3object")
	supported["delete"] = append(supported["delete"], "s3object")
	supported["create"] = append(supported["create"], "topic")
	supported["delete"] = append(supported["delete"], "topic")
	supported["create"] = append(supported["create"], "subscription")
	supported["delete"] = append(supported["delete"], "subscription")
	supported["create"] = append(supported["create"], "queue")
	supported["delete"] = append(supported["delete"], "queue")
	supported["create"] = append(supported["create"], "zone")
	supported["delete"] = append(supported["delete"], "zone")
	supported["create"] = append(supported["create"], "record")
	supported["delete"] = append(supported["delete"], "record")
	supported["create"] = append(supported["create"], "function")
	supported["delete"] = append(supported["delete"], "function")
	supported["create"] = append(supported["create"], "alarm")
	supported["delete"] = append(supported["delete"], "alarm")
	supported["start"] = append(supported["start"], "alarm")
	supported["stop"] = append(supported["stop"], "alarm")
	supported["attach"] = append(supported["attach"], "alarm")
	supported["detach"] = append(supported["detach"], "alarm")
	return supported
}
