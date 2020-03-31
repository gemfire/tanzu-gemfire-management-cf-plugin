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
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl"
	"io/ioutil"

	"code.cloudfoundry.org/cli/cf/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	. "github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/implfakes"
)

var _ = Describe("Endpoints", func() {

	Describe("GetEndPoints", func() {

		var (
			requester             *implfakes.FakeRequestHelper
			processRequest        impl.RequestHelper
			commandData           domain.CommandData
			fakeResponse          string
			fakeResponseGemfire99 string
		)

		BeforeEach(func() {
			requester = new(implfakes.FakeRequestHelper)
			processRequest = requester.Spy
			commandData = domain.CommandData{}
			JSONBytes, err := ioutil.ReadFile("../../testdata/api-docs.json")
			Expect(err).To(BeNil())
			fakeResponse = string(JSONBytes)
			JSONBytes, err = ioutil.ReadFile("../../testdata/api-docs-gemfire-99.json")
			Expect(err).To(BeNil())
			fakeResponseGemfire99 = string(JSONBytes)
		})

		It("Builds AvailableEndpoints when swagger data is received from Gemfire 9.9", func() {
			requester.ReturnsOnCall(0, "", 404, nil)
			requester.ReturnsOnCall(1, fakeResponseGemfire99, 200, nil)
			err := GetEndPoints(&commandData, processRequest)
			Expect(err).To(BeNil())
			Expect(len(commandData.AvailableEndpoints)).To(Equal(15))
			Expect(commandData.ConnnectionData.UseToken).To(BeFalse())

			checkRebalance, available := commandData.AvailableEndpoints["checkRebalanceStatus"]
			Expect(available).To(BeTrue())
			Expect(len(checkRebalance.Parameters)).To(Equal(1))
			Expect(checkRebalance.Parameters[0].Name).To(Equal("id"))
			Expect(checkRebalance.Parameters[0].Required).To(BeTrue())
			Expect(checkRebalance.Parameters[0].BodyDefinition).To(BeEmpty())

			listIndexes, available := commandData.AvailableEndpoints["listIndex"]
			Expect(available).To(BeTrue())
			Expect(len(listIndexes.Parameters)).To(Equal(2))
			Expect(listIndexes.Parameters[0].Name).To(Equal("id"))
			Expect(listIndexes.Parameters[0].Required).To(BeFalse())
			Expect(listIndexes.Parameters[0].BodyDefinition).To(BeEmpty())

			startRebalance, available := commandData.AvailableEndpoints["startRebalance"]
			Expect(available).To(BeTrue())
			Expect(len(startRebalance.Parameters)).To(Equal(1))
			Expect(startRebalance.Parameters[0].Name).To(Equal("operation"))
			Expect(startRebalance.Parameters[0].Required).To(BeTrue())
			Expect(startRebalance.Parameters[0].BodyDefinition).NotTo(BeEmpty())
			Expect(len(startRebalance.Parameters[0].BodyDefinition)).To(Equal(3))

			excludeRegions, available := startRebalance.Parameters[0].BodyDefinition["excludeRegions"]
			Expect(len(excludeRegions.([]string))).To(Equal(2))
			includeRegions, available := startRebalance.Parameters[0].BodyDefinition["includeRegions"]
			Expect(len(includeRegions.([]string))).To(Equal(2))
			simulate, available := startRebalance.Parameters[0].BodyDefinition["simulate"]
			Expect(simulate.(bool)).To(BeTrue())

			createRegion, available := commandData.AvailableEndpoints["create regions"]
			Expect(available).To(BeTrue())
			Expect(len(createRegion.Parameters)).To(Equal(1))
			Expect(createRegion.Parameters[0].Name).To(Equal("regionConfig"))
			Expect(createRegion.Parameters[0].Required).To(BeTrue())
			Expect(createRegion.Parameters[0].BodyDefinition).NotTo(BeEmpty())
			Expect(len(createRegion.Parameters[0].BodyDefinition)).To(Equal(7))
		})

		It("Builds AvailableEndpoints when swagger data is received", func() {
			requester.ReturnsOnCall(0, "", 404, nil)
			requester.ReturnsOnCall(1, fakeResponse, 200, nil)
			err := GetEndPoints(&commandData, processRequest)
			Expect(err).To(BeNil())
			Expect(len(commandData.AvailableEndpoints)).To(Equal(17))
			Expect(commandData.ConnnectionData.UseToken).To(BeTrue())

			checkRebalance, available := commandData.AvailableEndpoints["check rebalance"]
			Expect(available).To(BeTrue())
			Expect(len(checkRebalance.Parameters)).To(Equal(1))
			Expect(checkRebalance.Parameters[0].Name).To(Equal("id"))
			Expect(checkRebalance.Parameters[0].Required).To(BeTrue())
			Expect(checkRebalance.Parameters[0].BodyDefinition).To(BeEmpty())

			listIndexes, available := commandData.AvailableEndpoints["list indexes"]
			Expect(available).To(BeTrue())
			Expect(len(listIndexes.Parameters)).To(Equal(1))
			Expect(listIndexes.Parameters[0].Name).To(Equal("id"))
			Expect(listIndexes.Parameters[0].Required).To(BeFalse())
			Expect(listIndexes.Parameters[0].BodyDefinition).To(BeEmpty())

			startRebalance, available := commandData.AvailableEndpoints["start rebalance"]
			Expect(available).To(BeTrue())
			Expect(len(startRebalance.Parameters)).To(Equal(1))
			Expect(startRebalance.Parameters[0].Name).To(Equal("operation"))
			Expect(startRebalance.Parameters[0].Required).To(BeTrue())
			Expect(startRebalance.Parameters[0].BodyDefinition).NotTo(BeEmpty())
			Expect(len(startRebalance.Parameters[0].BodyDefinition)).To(Equal(3))

			excludeRegions, available := startRebalance.Parameters[0].BodyDefinition["excludeRegions"]
			Expect(len(excludeRegions.([]string))).To(Equal(2))
			includeRegions, available := startRebalance.Parameters[0].BodyDefinition["includeRegions"]
			Expect(len(includeRegions.([]string))).To(Equal(2))
			simulate, available := startRebalance.Parameters[0].BodyDefinition["simulate"]
			Expect(simulate.(bool)).To(BeTrue())

			createRegion, available := commandData.AvailableEndpoints["create region"]
			Expect(available).To(BeTrue())
			Expect(len(createRegion.Parameters)).To(Equal(1))
			Expect(createRegion.Parameters[0].Name).To(Equal("regionConfig"))
			Expect(createRegion.Parameters[0].Required).To(BeTrue())
			Expect(createRegion.Parameters[0].BodyDefinition).NotTo(BeEmpty())
			Expect(len(createRegion.Parameters[0].BodyDefinition)).To(Equal(9))

			expirations, available := createRegion.Parameters[0].BodyDefinition["expirations"]
			Expect(len(expirations.([]interface{}))).To(Equal(1))
			expiration := expirations.([]interface{})[0]
			Expect(expiration).NotTo(BeNil())
			action := expiration.(map[string]interface{})["action"]
			Expect(action).NotTo(BeNil())
			Expect(action).To(Equal("ENUM, one of: DESTROY, INVALIDATE, LEGACY"))
			expirationType := expiration.(map[string]interface{})["type"]
			Expect(expirationType).NotTo(BeNil())
			Expect(expirationType).To(Equal("ENUM, one of: ENTRY_TIME_TO_LIVE, ENTRY_IDLE_TIME, LEGACY"))
			timeInSeconds := expiration.(map[string]interface{})["timeInSeconds"]
			Expect(timeInSeconds).NotTo(BeNil())
			Expect(timeInSeconds).To(Equal(42))
		})

		It("Returns an error when Exchange call returns an error", func() {
			requester.Returns("", 0, errors.New("Failed call"))
			err := GetEndPoints(&commandData, processRequest)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Unable to reach /management/. Error: Failed call"))
		})

		It("Returns an error when swagger output cannot be parsed", func() {
			requester.ReturnsOnCall(0, "", 404, nil)
			requester.ReturnsOnCall(1, "", 200, nil)
			err := GetEndPoints(&commandData, processRequest)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("invalid response : unexpected end of JSON input"))
		})
	})

})
