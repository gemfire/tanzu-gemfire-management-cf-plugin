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
	"io/ioutil"

	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl"
	. "github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common/commonfakes"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/implfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommandProcessor", func() {

	var (
		requester        *implfakes.FakeRequestHelper
		formatter        *commonfakes.FakeFormatter
		requestBuilder   *commonfakes.FakeRequestBuilder
		commandProcessor impl.CommandProcessor
		err              error
		commandData      domain.CommandData
		fakeResponse     string
	)

	BeforeEach(func() {
		requester = new(implfakes.FakeRequestHelper)
		formatter = new(commonfakes.FakeFormatter)
		requestBuilder = new(commonfakes.FakeRequestBuilder)
		commandProcessor, err = NewCommandProcessor(requester.Spy, formatter, requestBuilder.Spy)
		Expect(err).NotTo(HaveOccurred())
		commandData = domain.CommandData{}
	})

	Context("NewCommandProcessor", func() {
		Context("When dependencies are missing", func() {
			It("Returns an error indicating missing dependencies", func() {
				cp, err := NewCommandProcessor(nil, nil, nil)
				Expect(cp).To(BeNil())
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("requester"))
				Expect(err.Error()).To(ContainSubstring("formatter"))
				Expect(err.Error()).To(ContainSubstring("requestBuilder"))
				Expect(err.Error()).To(ContainSubstring("must not be nil"))
			})
		})
	})

	Context("ProcessCommand", func() {
		Context("Exercise the request helper", func() {
			Context("Root and v1 not found, error on experimental", func() {
				It("Returns an error indicating API docs are not found for experimental", func() {
					requester.ReturnsOnCall(0, "", 404, nil)
					requester.ReturnsOnCall(1, "", 404, nil)
					requester.ReturnsOnCall(2, "", 500, errors.New("unable to get endpoints"))
					err = commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Unable to reach /management/experimental/api-docs. Error: unable to get endpoints"))
					Expect(len(commandData.AvailableEndpoints)).To(BeZero())
					Expect(requester.CallCount()).To(Equal(3))
				})
			})

			Context("Internal error on Root call", func() {
				It("Returns an error indicating /management/ is unreacheable", func() {
					requester.ReturnsOnCall(0, "", 500, nil)
					err = commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Unable to reach /management/. Status Code: 500"))
					Expect(len(commandData.AvailableEndpoints)).To(BeZero())
					Expect(requester.CallCount()).To(Equal(1))
				})
			})

			Context("First call not found, second call ok", func() {
				It("Returns an empty test JSON API description", func() {
					requester.ReturnsOnCall(0, "", 404, nil)
					requester.ReturnsOnCall(1, "{}", 200, nil)
					err = commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid command: "))
					Expect(len(commandData.AvailableEndpoints)).To(BeZero())
					Expect(requester.CallCount()).To(Equal(2))
				})
			})

			Context("First call unauthorized, second call ok", func() {
				It("Returns an empty test JSON API description", func() {
					requester.ReturnsOnCall(0, "", 401, nil)
					requester.ReturnsOnCall(1, "{}", 200, nil)
					err = commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid command: "))
					Expect(len(commandData.AvailableEndpoints)).To(BeZero())
					Expect(requester.CallCount()).To(Equal(2))
				})
			})

			Context("First call forbidden, second call ok", func() {
				It("Returns an empty test JSON API description", func() {
					requester.ReturnsOnCall(0, "", 403, nil)
					requester.ReturnsOnCall(1, "{}", 200, nil)
					err = commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid command: "))
					Expect(len(commandData.AvailableEndpoints)).To(BeZero())
					Expect(requester.CallCount()).To(Equal(2))
				})
			})

			Context("First call returns proxy error, second call ok", func() {
				It("Returns an empty test JSON API description", func() {
					requester.ReturnsOnCall(0, "", 407, nil)
					requester.ReturnsOnCall(1, "{}", 200, nil)
					err = commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid command: "))
					Expect(len(commandData.AvailableEndpoints)).To(BeZero())
					Expect(requester.CallCount()).To(Equal(2))
				})
			})
		})

		It("Returns an error if the command is not in the list of available commands", func() {
			// fake getEndPoint returns empty end points
			requester.ReturnsOnCall(0, "", 404, nil)
			requester.ReturnsOnCall(1, "{}", 200, nil)
			commandData.UserCommand.Command = "badcommand"
			commandData.AvailableEndpoints = make(map[string]domain.RestEndPoint)
			err = commandProcessor.ProcessCommand(&commandData)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid command: badcommand"))
			Expect(requester.CallCount()).To(Equal(2))
		})

		Context("Help output", func() {

			BeforeEach(func() {
				JSONBytes, err := ioutil.ReadFile("../../testdata/api-docs.json")
				Expect(err).NotTo(HaveOccurred())
				fakeResponse = string(JSONBytes)
				requester.ReturnsOnCall(0, "", 404, nil)
				requester.ReturnsOnCall(1, fakeResponse, 200, nil)
			})

			Context("'commands' command is given", func() {

				BeforeEach(func() {
					commandData.UserCommand.Command = "commands"
				})

				It("Describes endpoints in short form", func() {
					err := commandProcessor.ProcessCommand(&commandData)
					Expect(err).NotTo(HaveOccurred())
					Expect(formatter.DescribeEndpointCallCount()).To(Equal(17))
					_, provideDetails := formatter.DescribeEndpointArgsForCall(0)
					Expect(provideDetails).To(BeFalse())
				})
			})

			Context("'--help' or '-h' is given with specific command", func() {

				BeforeEach(func() {
					commandData.UserCommand.Command = "list indexes"
					commandData.UserCommand.Parameters = map[string]string{"--help": ""}
				})

				It("Describes endpoints in short form", func() {
					err := commandProcessor.ProcessCommand(&commandData)
					Expect(err).NotTo(HaveOccurred())
					Expect(formatter.DescribeEndpointCallCount()).To(Equal(1))
					_, provideDetails := formatter.DescribeEndpointArgsForCall(0)
					Expect(provideDetails).To(BeTrue())
				})
			})
		})

		Context("Check required params", func() {

			BeforeEach(func() {
				JSONBytes, err := ioutil.ReadFile("../../testdata/api-docs.json")
				Expect(err).NotTo(HaveOccurred())
				fakeResponse = string(JSONBytes)
				requester.ReturnsOnCall(0, "", 404, nil)
				requester.ReturnsOnCall(1, fakeResponse, 200, nil)

				commandData.UserCommand.Command = "delete region"
			})

			Context("When parameters not present", func() {

				BeforeEach(func() {
					commandData.UserCommand.Parameters = map[string]string{}
				})

				It("Returns an error indicating missing parameter", func() {
					err := commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Required Parameter is missing: "))
				})
			})
		})

		Context("Execute command", func() {

			BeforeEach(func() {
				JSONBytes, err := ioutil.ReadFile("../../testdata/api-docs.json")
				Expect(err).NotTo(HaveOccurred())
				fakeResponse = string(JSONBytes)
				requester.ReturnsOnCall(0, "", 404, nil)
				requester.ReturnsOnCall(1, fakeResponse, 200, nil)

				commandData.UserCommand.Command = "delete region"
				commandData.UserCommand.Parameters = map[string]string{"--id": "regionId"}
			})

			Context("When buildRequest fails", func() {

				BeforeEach(func() {
					requestBuilder.Returns("", nil, errors.New("Build request failed"))
				})

				It("Returns an error indicating buildRequest failed", func() {
					err := commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Build request failed"))
				})
			})

			Context("When processReqeust fails", func() {

				BeforeEach(func() {
					requester.Returns("", 404, errors.New("Process request failed"))
				})

				It("Returns an error indicating buildRequest failed", func() {
					err := commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Process request failed"))
				})
			})
		})

		Context("Table format", func() {

			BeforeEach(func() {
				JSONBytes, err := ioutil.ReadFile("../../testdata/api-docs.json")
				Expect(err).NotTo(HaveOccurred())
				fakeResponse = string(JSONBytes)
				requester.ReturnsOnCall(0, "", 404, nil)
				requester.ReturnsOnCall(1, fakeResponse, 200, nil)

				commandData.UserCommand.Command = "list members"
			})

			Context("When user provides JQ string", func() {

				BeforeEach(func() {
					commandData.UserCommand.Parameters = map[string]string{"--table": "."}
				})

				It("Calls the formatter with the user provided JQ string", func() {
					err := commandProcessor.ProcessCommand(&commandData)
					Expect(err).NotTo(HaveOccurred())
					Expect(formatter.FormatResponseCallCount()).To(Equal(1))
					_, query, userProvided := formatter.FormatResponseArgsForCall(0)
					Expect(userProvided).To(BeTrue())
					Expect(query).To(Equal("."))
				})
			})

			Context("When user does not provide JQ string", func() {

				BeforeEach(func() {
					commandData.UserCommand.Parameters = map[string]string{"--table": ""}
				})

				Context("When default JQ string is provided", func() {

					It("Calls the formatter with the default JQ string", func() {
						err := commandProcessor.ProcessCommand(&commandData)
						Expect(err).NotTo(HaveOccurred())
						_, query, userProvided := formatter.FormatResponseArgsForCall(0)
						Expect(userProvided).To(BeFalse())
						Expect(query).To(Equal(".result[] | .runtimeInfo[] | {name:.memberName,status:.status}"))
					})
				})

				Context("When default JQ is not provided", func() {

					BeforeEach(func() {
						commandData.UserCommand.Command = "list gateway-receivers"
					})

					It("Calls the formatter with hard-coded '.' JQ string", func() {
						err := commandProcessor.ProcessCommand(&commandData)
						Expect(err).NotTo(HaveOccurred())
						_, query, userProvided := formatter.FormatResponseArgsForCall(0)
						Expect(userProvided).To(BeFalse())
						Expect(query).To(Equal("."))
					})
				})
			})
		})

		Context("Format Response", func() {

			BeforeEach(func() {
				JSONBytes, err := ioutil.ReadFile("../../testdata/api-docs.json")
				Expect(err).NotTo(HaveOccurred())
				fakeResponse = string(JSONBytes)
				requester.ReturnsOnCall(0, "", 404, nil)
				requester.ReturnsOnCall(1, fakeResponse, 200, nil)

				commandData.UserCommand.Command = "list members"
			})

			Context("When format fails", func() {

				BeforeEach(func() {
					formatter.FormatResponseReturns("", errors.New("Format failed"))
				})

				It("Returns an error indicating format failed", func() {
					err := commandProcessor.ProcessCommand(&commandData)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Format failed"))
				})
			})
		})
	})

	Context("Required params", func() {
		var (
			endPoint domain.RestEndPoint
			param    domain.RestAPIParam
			command  domain.UserCommand
			err      error
		)

		BeforeEach(func() {
			endPoint.CommandName = "test"
			endPoint.Parameters = make([]domain.RestAPIParam, 1)

			param.In = "query"
			param.Name = "id"
			param.Description = "id"
			param.Required = true
			endPoint.Parameters[0] = param

			command.Parameters = make(map[string]string)
		})

		Context("A parameter is required but not found", func() {
			It("Returns an error indicating the missing parameter", func() {
				err = CheckRequiredParam(endPoint, command)
				Expect(err.Error()).To(Equal("Required Parameter is missing: id"))
			})
		})

		Context("A parameter is required and present", func() {
			It("Returns no error", func() {
				command.Parameters["--id"] = "value"
				err = CheckRequiredParam(endPoint, command)
				Expect(err).To(BeNil())
			})
		})

		Context("There are no required parameters", func() {
			It("Returns no error", func() {
				param.Required = false
				endPoint.Parameters[0] = param
				err = CheckRequiredParam(endPoint, command)
				Expect(err).To(BeNil())
			})
		})
	})

})
