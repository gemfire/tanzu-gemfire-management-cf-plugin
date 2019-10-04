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
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/geode"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/pcc"
)

func main() {
	helper := common.Requester{}
	commonCode, err := common.NewCommandProcessor(&helper)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// figure out who is calling
	if strings.Contains(os.Args[0], ".cf/plugins") {
		basicPlugin, err := pcc.NewBasicPlugin(commonCode)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		plugin.Start(&basicPlugin)
	} else {
		geodeCommand, err := geode.NewGeodeCommand(commonCode)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = geodeCommand.Run(os.Args)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
