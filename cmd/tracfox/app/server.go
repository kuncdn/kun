/*
Copyright 2019 The Labchan Authors.
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

package app

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"tracfox.io/tracfox/cmd/tracfox/app/options"
	"tracfox.io/tracfox/internal/util"
	"tracfox.io/tracfox/pkg/tracfox/config"
	"tracfox.io/tracfox/pkg/tracfox/server"
)

const (
	// componenttracway component name
	componentTracway = "tracfox"
)

// NewTracwayCommand 新建 tracwayCommand
func NewTracwayCommand(stopCh <-chan struct{}) *cobra.Command {
	cleanFlagSet := pflag.NewFlagSet(componentTracway, pflag.ContinueOnError)
	tracwayFlags := options.NewTracwayFlags()

	tracwayConfiguration := &config.TracwayConfiguration{} // 携带默认值的配置

	cmd := &cobra.Command{
		Use:                componentTracway,
		Short:              "tracfox service, is the api gateway micro service component of labchan",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cleanFlagSet.Parse(args); err != nil {
				cmd.Usage()
				glog.Exitln(err.Error())
			}

			// check if there are non-flag arguments in the command line
			cmds := cleanFlagSet.Args()
			if len(cmds) > 0 {
				cmd.Usage()
				glog.Exitf("unknown command: %s", cmds[0])
			}

			help, err := cleanFlagSet.GetBool("help")
			if err != nil {
				glog.Exitln(`"help" flag is non-bool, programmer error, please correct`)
			}
			if help {
				cmd.Help()
				return
			}

			if errs := options.ValidateTracwayFlags(tracwayFlags); len(errs) != 0 {
				glog.Exitln(util.NewAggregateError(errs))
			}

			if configFile := tracwayFlags.TracwayConfig; len(configFile) != 0 {
				loadConfigFile(configFile)
				tracwayConfiguration, err = loadConfigFile(configFile)
				if err != nil {
					glog.Exitln(err.Error())
				}
				if err := tracwayFlagPrecedence(tracwayConfiguration, args); err != nil {
					glog.Exitln(err.Error())
				}
				if errs := options.ValidateTracwayConfiguration(tracwayConfiguration); len(errs) != 0 {
					glog.Exitln("config file is incorrect error msg:", util.NewAggregateError(errs))
				}
			}
			if tracwayFlags.DryRun {
				glog.Warningln("The configuration file is correct. You have enabled the dry run parameter so exit")
				os.Exit(0)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			srv, err := server.NewManager(ctx, *tracwayConfiguration)
			if err != nil {
				glog.Exitln(err)
			}
			go func() {
				if err := srv.Run(); err != nil {
					panic(err)
				}
			}()
			<-stopCh
			ctx, cancel = context.WithTimeout(ctx, time.Duration(tracwayConfiguration.Default.GraceTimeOut)*time.Second)
			defer cancel()
			srv.GracefulStop(ctx)
			return
		},
	}

	cleanFlagSet.BoolP("help", "h", false, fmt.Sprintf("show more information about %s", cmd.Name()))
	tracwayFlags.AddFlags(cleanFlagSet)
	cleanFlagSet.AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})
	options.AddTracwayConfigurationFlags(cleanFlagSet, tracwayConfiguration)
	cmd.Flags().AddFlagSet(cleanFlagSet)
	return cmd
}

func tracwayFlagPrecedence(tracwayConfiguration *config.TracwayConfiguration, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.AddGoFlagSet(flag.CommandLine)
	tracwayFlags := options.NewTracwayFlags()

	options.AddTracwayConfigurationFlags(fs, tracwayConfiguration)
	tracwayFlags.AddFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	return nil
}

func loadConfigFile(filename string) (*config.TracwayConfiguration, error) {
	tracwayConfiguration := config.TracwayConfiguration{}
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stream, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(stream, &tracwayConfiguration)
	if err != nil {
		return nil, err
	}

	return &tracwayConfiguration, nil
}
