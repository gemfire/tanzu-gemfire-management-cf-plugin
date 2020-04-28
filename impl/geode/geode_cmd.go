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
	"errors"
	"fmt"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common/format"
)

// Command is the basic struct that the command works on
type command struct {
	commandData domain.CommandData
	comm        impl.CommandProcessor
}

// New provides a constructor for the Geode standalone implementation for the client
func New(comm impl.CommandProcessor) (command, error) {
	if comm == nil {
		return command{}, errors.New("command processor is not valid")
	}
	return command{comm: comm}, nil
}

// Run is the main entry point for the standalone Geode command line interface
// It is run once for each command executed
func (gc *command) Run(args []string) (err error) {

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

	geodeConnection := &GeodeConnection{}

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
	fmt.Println("Commands to interact with a Geode cluster.")
	fmt.Println("")
	fmt.Println("Usage: gemfire <target> <command> [options]")
	fmt.Println("")
	fmt.Println("\ttarget: \n\t\tURL to a Geode locator in the form of: http(s)://host:port")
	fmt.Println("\t\tOptional if 'GEODE_TARGET' environment variable is set")
	fmt.Println("\tcommand:\n\t\t'gemfire <target> commands' lists available commands")
	fmt.Println("\toptions:\n\t\t'gemfire <target> <command> -h' lists options for an individual command")
	fmt.Println(format.GeneralOptions)
	fmt.Println("\thelp:\n\t\t--help, -h for general help, and provide <target> and <command> for command-specific help")
}
