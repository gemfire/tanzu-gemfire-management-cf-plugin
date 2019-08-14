package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	cloudcachemanagementcfpluginfakes "github.com/gemfire/cloudcache-management-cf-plugin/cloudcache-management-cf-pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cf cli plugin", func() {
	Context("Retrieving Username, Password, and Endpoint", func() {
		It("Returns correct information", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			keyInfo := `Getting key mykey for service instance jjack as admin...

{
 "distributed_system_id": "0",
 "gfsh_login_string": "connect --url=https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/gemfire/v1 --user=cluster_operator_ygTWCaBfqtFHuTWxdaOMQ --password=W97ghWi4p2YF5MsfRCu6Eg --skip-ssl-validation",
 "locators": [
  "10.0.8.6[55221]"
 ],
 "urls": {
  "gfsh": "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/gemfire/v1",
  "pulse": "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/pulse"
 },
 "users": [
  {
   "password": "W97ghWi4p2YF5MsfRCu6Eg",
   "roles": [
    "cluster_operator"
   ],
   "username": "cluster_operator_ygTWCaBfqtFHuTWxdaOMQ"
  },
  {
   "password": "vcM942IBtpZrL3MxWyyi6Q",
   "roles": [
    "developer"
   ],
   "username": "developer_T2ONcuzffmoQ3Zv1HIraGQ"
  }
 ],
 "wan": {
  "sender_credentials": {
   "active": {
    "password": "nOVwPCVF25SeVfWgVTgCKA",
    "username": "gateway_sender_J3wGFvXCzO1hGH7ESx7EA"
   }
  }
 }
}
`
			expectedUsername := "cluster_operator_ygTWCaBfqtFHuTWxdaOMQ"
			expectedPassword := "W97ghWi4p2YF5MsfRCu6Eg"
			expectedEndpoint := "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com"
			fakeCf.CmdReturns(keyInfo, nil)
			username, password, endpoint, err := GetUsernamePasswordEndpoinFromServiceKey(fakeCf)
			Expect(username).To(Equal(expectedUsername))
			Expect(password).To(Equal(expectedPassword))
			Expect(endpoint).To(Equal(expectedEndpoint))
			Expect(err).To(BeNil())
		})
		It("Returns an error.", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("", errors.New("CF Command Error"))
			_, _, _, err := GetUsernamePasswordEndpoinFromServiceKey(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incorrect JSON.", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("{", nil)
			_, _, _, err := GetUsernamePasswordEndpoinFromServiceKey(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incomplete JSON.", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			keyInfo := `Getting key mykey for service instance jjack as admin...

{
 
`
			fakeCf.CmdReturns(keyInfo, nil)
			_, _, _, err := GetUsernamePasswordEndpoinFromServiceKey(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
	})
	Context("Retrieving Service Key from PCC Instance", func() {
		It("Returns correct information", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `Getting keys for service instance jjack as admin...

name
mykey

`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			response, err := GetServiceKeyFromPCCInstance(fakeCf)
			Expect(err).To(BeNil())
			expectedResponse := "mykey"
			Expect(response).To(Equal(expectedResponse))
		})
		It("Handling a no service instance found", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `FAILED
Service instance jjackk not found
`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			_, err := GetServiceKeyFromPCCInstance(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
		It("Handling no service key available", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `Getting keys for service instance oowen as admin...
No service key for service instance oowen
`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			_, err := GetServiceKeyFromPCCInstance(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
	})

	Context("Safekeeping tests", func() {
		It("Validate table filling", func() {
			columnSize := 20
			value := "some string"
			filler := "-"
			response := Fill(columnSize, value, filler)
			expectedResponse := " some string--------"
			Expect(response).To(Equal(expectedResponse))
		})
	})
})
