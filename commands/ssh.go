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

package commands

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/ssh"
)

var keyPathFlag, proxyInstanceThroughFlag string
var sshPortFlag int
var printSSHConfigFlag bool
var printSSHCLIFlag bool
var privateIPFlag bool
var disableStrictHostKeyCheckingFlag bool

func init() {
	RootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringVarP(&keyPathFlag, "identity", "i", "", "Set path or name toward the identity (key file) to use to connect through SSH")
	sshCmd.Flags().IntVar(&sshPortFlag, "port", 22, "Set SSH port")
	sshCmd.Flags().StringVar(&proxyInstanceThroughFlag, "through", "", "Name of instance to proxy through to connect to a destination host")
	sshCmd.Flags().BoolVar(&printSSHConfigFlag, "print-config", false, "Print SSH configuration for ~/.ssh/config file.")
	sshCmd.Flags().BoolVar(&printSSHCLIFlag, "print-cli", false, "Print the CLI one-liner to connect with SSH. (/usr/bin/ssh user@ip -i ...)")
	sshCmd.Flags().BoolVar(&privateIPFlag, "private", false, "Use private ip to connect to host")
	sshCmd.Flags().BoolVar(&disableStrictHostKeyCheckingFlag, "disable-strict-host-keychecking", false, "Disable the remote host key check from ~/.ssh/known_hosts or ~/.awless/known_hosts file")
}

var sshCmd = &cobra.Command{
	Use:   "ssh [USER@]INSTANCE",
	Short: "Launch a SSH (Secure Shell) session to an instance given an id or alias",
	Example: `  awless ssh i-8d43b21b                       # using the instance id
  awless ssh ec2-user@redis-prod              # using the instance name and specify a user
  awless ssh redis-prod -i keyname # using a key stored in ~/.ssh/keyname.pem
  awless ssh redis-prod -i ./path/toward/key # with a keyfile`,
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("instance required")
		}

		var err error
		var connectionCtx *instanceConnectionContext

		if proxyInstanceThroughFlag != "" {
			connectionCtx, err = initInstanceConnectionContext(proxyInstanceThroughFlag, keyPathFlag)
		} else {
			connectionCtx, err = initInstanceConnectionContext(args[0], keyPathFlag)
		}
		exitOn(err)

		firsHopClient, err := ssh.InitClient(connectionCtx.keypath, config.KeysDir, filepath.Join(os.Getenv("HOME"), ".ssh"))
		exitOn(err)

		if err != nil && strings.Contains(err.Error(), "cannot find SSH key") && keyPathFlag == "" {
			logger.Info("you may want to specify a key filepath with `-i /path/to/key.pem`")
		}
		exitOn(err)

		firsHopClient.SetLogger(logger.DefaultLogger)
		firsHopClient.SetStrictHostKeyChecking(!disableStrictHostKeyCheckingFlag)
		firsHopClient.InteractiveTerminalFunc = console.InteractiveTerminal
		firsHopClient.Port = sshPortFlag

		if privateIPFlag {
			if priv := connectionCtx.privip; priv != "" {
				firsHopClient.IP = connectionCtx.privip
			} else {
				exitOn(fmt.Errorf(
					"no private IP resolved for instance %s (state '%s')",
					connectionCtx.instance.Id(), connectionCtx.state,
				))
			}
		} else {
			if pub := connectionCtx.ip; pub != "" {
				firsHopClient.IP = connectionCtx.ip
			} else {
				logger.Infof("`--private` flag can be used to connect through instance's private IP '%s'", connectionCtx.privip)
				exitOn(fmt.Errorf("no public IP resolved for instance %s (state '%s')", connectionCtx.instance.Id(), connectionCtx.state))
			}
		}

		if connectionCtx.user != "" {
			err = firsHopClient.DialWithUsers(connectionCtx.user)
		} else {
			err = firsHopClient.DialWithUsers(awsconfig.DefaultAMIUsers...)
		}

		if isConnectionRefusedErr(err) {
			logger.Warning("cannot connect to this instance, maybe the system is still booting?")
			exitOn(err)
			return nil
		}

		if err != nil {
			if e := connectionCtx.checkInstanceAccessible(); e != nil {
				logger.Error(e.Error())
			}
			exitOn(err)
		}

		targetClient := firsHopClient

		if proxyInstanceThroughFlag != "" {
			destInstanceCtx, err := initInstanceConnectionContext(args[0], keyPathFlag)
			exitOn(err)
			if destInstanceCtx.user != "" {
				targetClient, err = firsHopClient.NewClientWithProxy(destInstanceCtx.privip, destInstanceCtx.user)
			} else {
				targetClient, err = firsHopClient.NewClientWithProxy(destInstanceCtx.privip, awsconfig.DefaultAMIUsers...)
			}
		}

		if printSSHConfigFlag {
			fmt.Println(targetClient.SSHConfigString(connectionCtx.instanceName))
			return nil
		}

		if printSSHCLIFlag {
			fmt.Println(targetClient.ConnectString())
			return nil
		}

		exitOn(targetClient.Connect())
		return nil
	},
}

func isConnectionRefusedErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "connection refused")
}

type instanceConnectionContext struct {
	ip, privip          string
	myip                net.IP
	user, keypath       string
	state, instanceName string
	instance            *graph.Resource
	resourcesGraph      *graph.Graph
}

func initInstanceConnectionContext(userhost, keypath string) (*instanceConnectionContext, error) {
	ctx := &instanceConnectionContext{}

	if strings.Contains(userhost, "@") {
		ctx.user = strings.Split(userhost, "@")[0]
		ctx.instanceName = strings.Split(userhost, "@")[1]
	} else {
		ctx.instanceName = userhost
	}

	ctx.fetchConnectionInfo()

	instanceResolvers := []graph.Resolver{&graph.ByProperty{Key: "Name", Value: ctx.instanceName}, &graph.ByType{Typ: cloud.Instance}}
	resources, err := ctx.resourcesGraph.ResolveResources(&graph.And{Resolvers: instanceResolvers})
	exitOn(err)
	switch len(resources) {
	case 0:
		// No instance with that name, use the id
		ctx.instance, err = findResource(ctx.resourcesGraph, ctx.instanceName, cloud.Instance)
		exitOn(err)
	case 1:
		ctx.instance = resources[0]
	default:
		idStatus := graph.Resources(resources).Map(func(r *graph.Resource) string {
			return fmt.Sprintf("%s (%s)", r.Id(), r.Properties[properties.State])
		})
		logger.Infof("Found %d resources with name '%s': %s", len(resources), ctx.instanceName, strings.Join(idStatus, ", "))

		var running []*graph.Resource
		running, err = ctx.resourcesGraph.ResolveResources(&graph.And{Resolvers: append(instanceResolvers, &graph.ByProperty{Key: properties.State, Value: "running"})})
		exitOn(err)

		switch len(running) {
		case 0:
			logger.Warning("None of them is running, cannot connect through SSH")
			return ctx, errors.New("non running instances")
		case 1:
			logger.Infof("Found only one instance running: %s. Will connect to this instance.", running[0].Id())
			ctx.instance = running[0]
		default:
			logger.Warning("Connect through the running ones using their id:")
			for _, res := range running {
				var up string
				if uptime, ok := res.Properties[properties.Launched].(time.Time); ok {
					up = fmt.Sprintf("\t\t(uptime: %s)", console.HumanizeTime(uptime))
				}
				logger.Warningf("\t`awless ssh %s`%s", res.Id(), up)
			}
			return ctx, errors.New("use instances ids")
		}
	}

	ctx.privip, _ = ctx.instance.Properties[properties.PrivateIP].(string)
	ctx.ip, _ = ctx.instance.Properties[properties.PublicIP].(string)
	ctx.state, _ = ctx.instance.Properties[properties.State].(string)

	if keypath != "" {
		ctx.keypath = keypath
	} else {
		keypair, ok := ctx.instance.Properties[properties.KeyPair].(string)
		if ok {
			ctx.keypath = fmt.Sprint(keypair)
		}
	}

	return ctx, nil
}

func (ctx *instanceConnectionContext) fetchConnectionInfo() {
	var resourcesGraph, sgroupsGraph *graph.Graph
	var myip net.IP
	var wg sync.WaitGroup
	var errc = make(chan error)

	wg.Add(1)
	go func() {
		var err error
		defer wg.Done()
		resourcesGraph, err = aws.InfraService.FetchByType(cloud.Instance)
		if err != nil {
			errc <- err
		}
	}()

	wg.Add(1)
	go func() {
		var err error
		defer wg.Done()
		sgroupsGraph, err = aws.InfraService.FetchByType(cloud.SecurityGroup)
		if err != nil {
			errc <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		myip = getMyIP()
	}()
	go func() {
		wg.Wait()
		close(errc)
	}()
	for err := range errc {
		if err != nil {
			exitOn(err)
		}
	}
	resourcesGraph.AddGraph(sgroupsGraph)

	ctx.resourcesGraph = resourcesGraph
	ctx.myip = myip
	return
}

func (ctx *instanceConnectionContext) checkInstanceAccessible() (err error) {
	if st := ctx.state; st != "running" {
		logger.Warningf("this instance is '%s' (cannot ssh to a non running state)", st)
		if st == "stopped" {
			logger.Warningf("you can start it with `awless -f start instance id=%s`", ctx.instance.Id())
		}
		return errors.New("instance not accessible")
	}

	sgroups, ok := ctx.instance.Properties[properties.SecurityGroups].([]string)
	if ok {
		var sshPortOpen, myIPAllowed bool
		for _, id := range sgroups {
			var sgroup *graph.Resource
			sgroup, err = findResource(ctx.resourcesGraph, id, cloud.SecurityGroup)
			if err != nil {
				logger.Errorf("cannot get securitygroup '%s' for instance '%s': %s", id, ctx.instance.Id(), err)
				break
			}

			rules, ok := sgroup.Properties[properties.InboundRules].([]*graph.FirewallRule)
			if ok {
				for _, r := range rules {
					if r.PortRange.Contains(22) {
						sshPortOpen = true
					}
					if ctx.myip != nil && r.Contains(ctx.myip.String()) {
						myIPAllowed = true
					}
				}
			}
		}

		if !sshPortOpen {
			logger.Warning("port 22 is not open on this instance")
			return errors.New("instance not accessible")
		}

		if !myIPAllowed && ctx.myip != nil {
			logger.Warningf("your ip %s is not authorized for this instance. You might want to update the securitygroup with:", ctx.myip)
			var group = "mygroup"
			if len(sgroups) == 1 {
				group = sgroups[0]
			}
			logger.Warningf("`awless update securitygroup id=%s inbound=authorize protocol=tcp cidr=%s/32 portrange=22`", group, ctx.myip)
			return errors.New("instance not accessible")
		}
	}

	return nil
}

func findResource(g *graph.Graph, id, typ string) (*graph.Resource, error) {
	if found, err := g.FindResource(id); found == nil || err != nil {
		return nil, fmt.Errorf("instance '%s' not found", id)
	}

	return g.GetResource(typ, id)
}
