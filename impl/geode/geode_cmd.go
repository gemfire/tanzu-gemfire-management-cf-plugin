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

package geode

import (
	"fmt"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
)

type geodeCommand struct {
	commandData domain.CommandData
	comm        common.CommandProcessor
}

// NewGeodeCommand provides a constructor for the Geode standalone implementation for the client
func NewGeodeCommand(comm common.CommandProcessor) (geodeCommand, error) {
	return geodeCommand{comm: comm}, nil
}

// Run is the main entry point for the standalone Geode command line interface
// It is run once for each command executed
func (gc *geodeCommand) Run(args []string) (err error) {

	gc.commandData.Target, gc.commandData.UserCommand = common.GetTargetAndClusterCommand(args)

	if common.HasOption(gc.commandData.UserCommand.Parameters, []string{"-v", "--version"}) {
		fmt.Printf("Version: %d.%d.%d\n", domain.VersionType.Major, domain.VersionType.Minor, domain.VersionType.Build)
		return
	}

	// if no user command and args contains -h or --help
	if gc.commandData.UserCommand.Command == "" {
		printHelp()
		return
	}

	geodeConnection, err := NewGeodeConnectionProvider()
	if err != nil {
		printHelp()
		return
	}

	err = geodeConnection.GetConnectionData(&gc.commandData)
	if err != nil {
		printHelp()
		return
	}

	// From this point common code can handle the processing of the command
	err = gc.comm.ProcessCommand(&gc.commandData)

	return
}

func printHelp() {
	fmt.Println("Commands to interact with geode cluster.")
	fmt.Println("")
	fmt.Println("Usage: pcc <target> <command> [options]")
	fmt.Println("")
	fmt.Println("\ttarget: \n\t\turl to a geode locator in the form of : http(s)://host:port")
	fmt.Println("\t\tomit if 'GEODE_TARGET' environment variable is set")
	fmt.Println("\tcommand:\n\t\tuse 'pcc <target> commands' to see a list of supported commands")
	fmt.Println("\toptions:\n\t\tuse 'pcc <target> <command> -h' to see options for individual command.")
	fmt.Println(common.GeneralOptions)
	fmt.Println("\thelp:\n\t\tuse -h or --help for general help, and provide <command> for command specific help.")
}
