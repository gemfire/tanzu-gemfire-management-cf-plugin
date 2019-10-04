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

package common_test

import (
	"errors"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	. "github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/implfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommandProcessor", func() {

	Context("ProcessCommand", func() {
		var (
			helper           *implfakes.FakeRequestHelper
			commandProcessor CommandProcessor
			err              error
			commandData      domain.CommandData
		)

		BeforeEach(func() {
			helper = new(implfakes.FakeRequestHelper)
			commandProcessor, err = NewCommandProcessor(helper)
			Expect(err).NotTo(HaveOccurred())
			commandData = domain.CommandData{}
		})

		It("Returns an error if RequestHelper cannot get endpoints", func() {
			helper.ExchangeReturns("", errors.New("Unable to get endpoints"))
			err = commandProcessor.ProcessCommand(&commandData)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unable to reach /management/experimental/api-docs: Unable to get endpoints"))
			Expect(len(commandData.AvailableEndpoints)).To(BeZero())
			Expect(helper.ExchangeCallCount()).To(Equal(1))
		})

		It("Returns an error if the command is not in the list of available commands", func() {
			// fake getEndPoint returns empty end points
			helper.ExchangeReturns("{}", nil)
			commandData.UserCommand.Command = "badcommand"
			commandData.AvailableEndpoints = make(map[string]domain.RestEndPoint)
			err = commandProcessor.ProcessCommand(&commandData)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid command: badcommand"))
			Expect(helper.ExchangeCallCount()).To(Equal(1))
		})
	})

	Context("Check required params", func() {
		It("Returns an error if required param is not found", func() {
			var endPoint domain.RestEndPoint
			endPoint.CommandName = "test"
			endPoint.Parameters = make([]domain.RestAPIParam, 1)

			var param domain.RestAPIParam
			param.In = "query"
			param.Name = "id"
			param.Description = "id"
			param.Required = true
			endPoint.Parameters[0] = param

			var command domain.UserCommand
			command.Parameters = make(map[string]string)
			err := CheckRequiredParam(endPoint, command)
			Expect(err.Error()).To(Equal("Required Parameter is missing: id"))

			param.In = "body"
			param.Name = "config"
			endPoint.Parameters[0] = param
			err = CheckRequiredParam(endPoint, command)
			Expect(err.Error()).To(Equal("Required Parameter is missing: config"))

			param.Required = false
			endPoint.Parameters[0] = param
			err = CheckRequiredParam(endPoint, command)
			Expect(err).To(BeNil())
		})
	})
})
