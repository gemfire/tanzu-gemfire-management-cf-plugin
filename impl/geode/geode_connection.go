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
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"os"
	"strings"
)

type geodeConnection struct {
}

// NewGeodeConnectionProvider provides a constructor for the Geode standalone implementation of ConnectionProvider
func NewGeodeConnectionProvider() (impl.ConnectionProvider, error) {
	return &geodeConnection{}, nil
}

func (gc *geodeConnection) GetConnectionData(commandData *domain.CommandData) error {
	commandData.ConnnectionData = domain.ConnectionData{}

	// LocatorAddress, Username and Password may be provided as environment variables
	// but can be overridden on the command line
	commandData.ConnnectionData.LocatorAddress = strings.TrimSuffix(commandData.Target, "/")
	commandData.ConnnectionData.Username = common.GetOption(commandData.UserCommand.Parameters, []string{"--user", "-u"})
	if commandData.ConnnectionData.Username == "" {
		commandData.ConnnectionData.Username = os.Getenv("GEODE_USERNAME")
	}
	commandData.ConnnectionData.Password = common.GetOption(commandData.UserCommand.Parameters, []string{"--password", "-p"})
	if commandData.ConnnectionData.Password == "" {
		commandData.ConnnectionData.Password = os.Getenv("GEODE_PASSWORD")
	}

	return nil
}
