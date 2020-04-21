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

package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common/builder"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common/filter"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common/format"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/gemfire"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/geode"
	"os"
	"strings"
)

func main() {
	processRequest := common.Exchange
	formatter, err := format.New(filter.GOJQFilter)
	checkError(err)
	commonCode, err := common.NewCommandProcessor(processRequest, formatter, builder.BuildRequest)
	checkError(err)

	// figure out who is calling. If invoked as a standalone cli
	if (strings.HasSuffix(os.Args[0], "main_go") ||
		strings.HasSuffix(os.Args[0], "gemfire")) &&
		!strings.Contains(os.Args[0], "cf/plugins") {
		geodeCommand, err := geode.New(commonCode)
		checkError(err)
		err = geodeCommand.Run(os.Args)
		checkError(err)
		// otherwise assume invoked as a plugin
	} else {
		basicPlugin, err := gemfire.NewBasicPlugin(commonCode)
		checkError(err)
		plugin.Start(basicPlugin)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
