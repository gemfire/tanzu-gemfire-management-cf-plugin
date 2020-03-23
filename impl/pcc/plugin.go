/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for additional information regarding
 * copyright ownership. The ASF licenses this file to You under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License. You may obtain a
 * copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package pcc

import (
	"errors"
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common/format"
)

// BasicPlugin declares the dataset that commands work on
type BasicPlugin struct {
	commandData domain.CommandData
	comm        impl.CommandProcessor
}

// NewBasicPlugin provides the constructor for a BasicPlugin struct
func NewBasicPlugin(comm impl.CommandProcessor) (plugin.Plugin, error) {
	if comm == nil {
		return nil, errors.New("command processor is not valid")
	}
	return &BasicPlugin{comm: comm}, nil
}

// Run is the main entry point for the CF plugin interface
// It is run once for each CF plugin command executed
func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	var err error
	c.commandData.Target, c.commandData.UserCommand = common.GetTargetAndClusterCommand(args)
	if c.commandData.UserCommand.Command == "" {
		fmt.Println("missing command")
		os.Exit(1)
	}

	pluginConnection, err := New(cliConnection)
	if err != nil {
		fmt.Printf(format.GenericErrorMessage, err.Error())
		os.Exit(1)
	}
	err = pluginConnection.GetConnectionData(&c.commandData)
	if err != nil {
		fmt.Printf(format.GenericErrorMessage, err.Error())
		os.Exit(1)
	}

	// From this point common code can handle the processing of the command
	err = c.comm.ProcessCommand(&c.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return
}

// GetMetadata provides metadata about the CF plugin including a helptext for the user
func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "pcc",
		Version: domain.VersionType,
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "pcc",
				HelpText: "Commands to interact with Geode cluster.\n",
				UsageDetails: plugin.Usage{
					Usage: "cf  pcc  [target]  <command>  [options] \n\n" +
						"\ttarget:\n\t\ta pcc_instance. \n" +
						"\t\tomit if 'GEODE_TARGET' environment variable is set \n" +
						"\tcommand:\n\t\tuse 'cf pcc <target> commands' to see a list of supported commands \n" +
						"\toptions:\n\t\tuse 'cf pcc <target> command -help' to see options for individual command." +
						format.GeneralOptions + "\n" +
						"\thelp\nt\t\t: use -h or --help for general help, and provide <command> -help for command specific help",
				},
			},
		},
	}
}
