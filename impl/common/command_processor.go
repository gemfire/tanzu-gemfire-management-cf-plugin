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

package common

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Formatter

// Formatter interface provides response and other output formatting
type Formatter interface {
	DescribeEndpoint(domain.RestEndPoint, bool) string
	FormatResponse(string, string, bool) (string, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RequestBuilder

// RequestBuilder is function type generating a request
type RequestBuilder func(endpoint domain.RestEndPoint, commandData *domain.CommandData) (request *http.Request, err error)

// CommandProcessor struct holds the implementation for the RequestHelper interface
type commandProcessor struct {
	processRequest impl.RequestHelper
	formatter      Formatter
	buildRequest   RequestBuilder
}

// NewCommandProcessor provides the constructor for the CommandProcessor
func NewCommandProcessor(requester impl.RequestHelper, formatter Formatter, requestBuilder RequestBuilder) (impl.CommandProcessor, error) {
	var errorString []string
	if requester == nil {
		errorString = append(errorString, "requester")
	}
	if formatter == nil {
		errorString = append(errorString, "formatter")
	}
	if requestBuilder == nil {
		errorString = append(errorString, "requestBuilder")
	}
	if len(errorString) > 0 {
		return nil, errors.New(strings.Join(errorString, " and ") + " must not be nil")
	}
	return &commandProcessor{processRequest: requester, formatter: formatter, buildRequest: requestBuilder}, nil
}

// ProcessCommand handles the common steps for executing a command against the Geode cluster
func (c *commandProcessor) ProcessCommand(commandData *domain.CommandData) (err error) {
	err = GetEndPoints(commandData, c.processRequest)
	if err != nil {
		return
	}

	userCommand := commandData.UserCommand.Command
	restEndPoint, available := commandData.AvailableEndpoints[userCommand]

	if userCommand == "commands" || !available {
		commandNames := sortCommandNames(commandData)
		for _, commandName := range commandNames {
			fmt.Println(c.formatter.DescribeEndpoint(commandData.AvailableEndpoints[commandName], false))
		}
		if userCommand != "commands" {
			err = errors.New("Invalid command: " + userCommand)
		}
		return
	}

	if HasOption(commandData.UserCommand.Parameters, []string{"-h", "--help", "-help"}) {
		fmt.Println(c.formatter.DescribeEndpoint(restEndPoint, true))
		return
	}

	err = CheckRequiredParam(restEndPoint, commandData.UserCommand)
	if err != nil {
		return
	}

	urlResponse, err := c.executeCommand(commandData)
	if err != nil {
		return
	}

	var jqFilter string
	var userFilter bool
	if HasOption(commandData.UserCommand.Parameters, []string{"-t", "--table"}) {
		jqFilter = GetOption(commandData.UserCommand.Parameters, []string{"--table", "-t"})
		userFilter = true
		// if no jqFilter is specified by the user, use the default defined by the rest end point
		if jqFilter == "" {
			jqFilter = restEndPoint.JQFilter
			if jqFilter == "" {
				jqFilter = "."
			}
			userFilter = false
		}
	}

	jsonToBePrinted, err := c.formatter.FormatResponse(urlResponse, jqFilter, userFilter)
	if err != nil {
		return
	}
	fmt.Println(jsonToBePrinted)

	return
}

// CheckRequiredParam checks if required parameters have been provided
func CheckRequiredParam(restEndPoint domain.RestEndPoint, command domain.UserCommand) error {
	for _, s := range restEndPoint.Parameters {
		if s.Required {
			value := command.Parameters["--"+s.Name]
			if value == "" {
				return errors.New("Required Parameter is missing: " + s.Name)
			}
		}
	}
	return nil
}

func (c *commandProcessor) executeCommand(commandData *domain.CommandData) (urlResponse string, err error) {
	restEndPoint, _ := commandData.AvailableEndpoints[commandData.UserCommand.Command]
	request, err := c.buildRequest(restEndPoint, commandData)
	if err != nil {
		return "", err
	}
	urlResponse, _, err = c.processRequest(request)
	return
}

func sortCommandNames(commandData *domain.CommandData) (commandNames []string) {
	commandNames = make([]string, 0, len(commandData.AvailableEndpoints))
	for _, command := range commandData.AvailableEndpoints {
		commandNames = append(commandNames, command.CommandName)
	}
	sort.Strings(commandNames)
	return
}
