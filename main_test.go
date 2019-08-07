package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	cloudcachemanagementcfpluginfakes "github.com/gemfire/cloudcache-management-cf-plugin/cloudcache-management-cf-pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var firstEndpoint = "http://localhost:7070/management/experimental/api-docs"

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
			expectedUsername :="cluster_operator_ygTWCaBfqtFHuTWxdaOMQ"
			expectedPassword := "W97ghWi4p2YF5MsfRCu6Eg"
			expectedEndpoint := "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/management/experimental/cli"
			fakeCf.CmdReturns(keyInfo, nil)
			username, password, endpoint, err := GetUsernamePasswordEndpoint(fakeCf)
			Expect(username).To(Equal(expectedUsername))
			Expect(password).To(Equal(expectedPassword))
			Expect(endpoint).To(Equal(expectedEndpoint))
			Expect(err).To(BeNil())
		})
		It("Returns an error.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("", errors.New("CF Command Error"))
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incorrect JSON.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("{", nil)
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incomplete JSON.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			keyInfo := `Getting key mykey for service instance jjack as admin...

{
 
`
			fakeCf.CmdReturns(keyInfo, nil)
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf)
			Expect(err).To(Not(BeNil()))
		})
	})
	Context("Retrieving Service Key from PCC Instance", func() {
		It("Returns correct information", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `Getting keys for service instance jjack as admin...

name
mykey

`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			response, err := GetServiceKeyFromPCCInstance(fakeCf, "jjack")
			Expect(err).To(BeNil())
			expectedResponse := "mykey"
			Expect(response).To(Equal(expectedResponse))
		})
		It("Handling a no service instance found", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `FAILED
Service instance jjackk not found
`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			_, err := GetServiceKeyFromPCCInstance(fakeCf, "jjack")
			Expect(err).To(Not(BeNil()))
		})
		It("Handling no service key available", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `Getting keys for service instance oowen as admin...
No service key for service instance oowen
`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			_, err := GetServiceKeyFromPCCInstance(fakeCf, "oowen")
			Expect(err).To(Not(BeNil()))
		})
	})


	Context("Safekeeping tests", func(){
		It("Validate table filling", func(){
			columnSize := 20
			value := "some string"
			filler := "-"
			response := Fill(columnSize, value, filler)
			expectedResponse := " some string--------"
			Expect(response).To(Equal(expectedResponse))
		})
	})
	Context("Input Mapping tests", func(){
		It("list members", func(){
			APICallStruct.action = "list"
			APICallStruct.target = "members"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/members"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("get member", func(){
			APICallStruct.action = "get"
			APICallStruct.target = "member"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/members/{id}"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("list regions", func(){
			APICallStruct.action = "list"
			APICallStruct.target = "regions"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/regions"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("get region", func(){
			APICallStruct.action = "get"
			APICallStruct.target = "region"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/regions/{id}"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("list indexes", func(){
			APICallStruct.action = "list"
			APICallStruct.target = "indexes"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/regions/{regionName}/indexes"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("get index", func(){
			APICallStruct.action = "get"
			APICallStruct.target = "index"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/regions/{regionName}/indexes/{id}"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("start rebalance", func(){
			APICallStruct.action = "create" //see synonymConverter("start")
			APICallStruct.target = "rebalance"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/operations/rebalances"))
			Expect(endpoint.HttpMethod).To(Equal("post"))
		})
		It("list rebalances", func(){
			APICallStruct.action = "list"
			APICallStruct.target = "rebalances"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/operations/rebalances"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("check rebalance", func(){
			APICallStruct.action = "get" //see synonymConverter("check")
			APICallStruct.target = "rebalance"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/operations/rebalances/{id}"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("create region", func(){
			APICallStruct.action = "create"
			APICallStruct.target = "region"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/regions"))
			Expect(endpoint.HttpMethod).To(Equal("post"))
		})
		It("delete region", func(){
			APICallStruct.action = "delete"
			APICallStruct.target = "region"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/regions/{id}"))
			Expect(endpoint.HttpMethod).To(Equal("delete"))
		})
		It("list cli", func(){
			APICallStruct.action = "list"
			APICallStruct.target = "cli"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/cli"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("list ping", func(){
			APICallStruct.action = "list"
			APICallStruct.target = "ping"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/ping"))
			Expect(endpoint.HttpMethod).To(Equal("get"))
		})
		It("configure pdx", func(){
			APICallStruct.action = "configure"
			APICallStruct.target = "pdx"
			endpoint := processUserCallInTest()
			Expect(endpoint.Url).To(Equal("/experimental/configurations/pdx"))
			Expect(endpoint.HttpMethod).To(Equal("post"))
		})



	})
})

func processUserCallInTest() (endpoint IndividualEndpoint){
	err := executeFirstRequest(firstEndpoint)
	Expect(err).To(BeNil())
	endpoint, err = mapUserInputToAvailableEndpoint()
	Expect(err).To(BeNil())
	return
}
