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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common/format"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
)

type pluginConnection struct {
	cliConnection plugin.CliConnection
}

// New provides a constructor for the PCC implementation of ConnectionProvider
func New(connection plugin.CliConnection) (impl.ConnectionProvider, error) {
	if connection == nil {
		return nil, errors.New("cliConnection is not valid")
	}
	return &pluginConnection{cliConnection: connection}, nil
}

// GetConnectionData provides the connection data from a PCC cluster using the CF CLI
func (pc *pluginConnection) GetConnectionData(commandData *domain.CommandData) error {
	commandData.ConnnectionData = domain.ConnectionData{}
	serviceKey, err := pc.getServiceKey(commandData.Target)
	if err != nil {
		return err
	}

	return pc.getServiceKeyDetails(commandData, serviceKey)

}

func (pc *pluginConnection) getServiceKey(target string) (serviceKey string, err error) {
	results, err := pc.cliConnection.CliCommandWithoutTerminalOutput("service-keys", target)
	if err != nil {
		return "", err
	}
	hasKey := false
	if strings.Contains(results[1], "No service key for service instance") {
		return "", fmt.Errorf(format.NoServiceKeyMessage, target, target)
	}
	for _, value := range results {
		line := strings.Fields(value)
		if len(line) > 0 {
			if hasKey {
				serviceKey = line[0]
				return
			} else if line[0] == "name" {
				hasKey = true
			}
		}
	}
	if serviceKey == "" {
		err = fmt.Errorf(format.NoServiceKeyMessage, target, target)
	}
	return
}

func (pc *pluginConnection) getServiceKeyDetails(commandData *domain.CommandData, serviceKey string) (err error) {
	keyInfo, err := pc.cliConnection.CliCommandWithoutTerminalOutput("service-key", commandData.Target, serviceKey)
	if err != nil {
		return err
	}

	if len(keyInfo) < 2 {
		return errors.New(format.InvalidServiceKeyResponse)
	}
	keyInfo = keyInfo[2:] //take out first two lines of cf service-key ... output
	joinKeyInfo := strings.Join(keyInfo, "\n")
	serviceKeyStruct := domain.ServiceKey{}

	err = json.Unmarshal([]byte(joinKeyInfo), &serviceKeyStruct)
	if err != nil {
		return err
	}

	commandData.ConnnectionData.LocatorAddress = strings.TrimSuffix(serviceKeyStruct.Urls.Gfsh, "/gemfire/v1")

	// use the credentials in the command line
	if commandData.ConnnectionData.Username == "" || commandData.ConnnectionData.Password == "" {
		commandData.ConnnectionData.Username = common.GetOption(commandData.UserCommand.Parameters, []string{"--user", "-u"})
		commandData.ConnnectionData.Password = common.GetOption(commandData.UserCommand.Parameters, []string{"--password", "-p"})
	}

	// find credentials in the service key
	for _, user := range serviceKeyStruct.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			commandData.ConnnectionData.Username = user.Username
			commandData.ConnnectionData.Password = user.Password
		}
	}

	// find the access token if any
	token, _ := pc.cliConnection.AccessToken()
	commandData.ConnnectionData.Token = token

	return
}
