package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	cloudcachemanagementcfpluginfakes "github.com/gemfire/cloudcache-management-cf-plugin/cloudcache-management-cf-pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cf gf plugin", func() {
	Context("Retrieving Username, Password, and Endpoint", func() {
		It("Returns correct information", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			pccService := "jjack"
			key := "mykey"
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
			expectedEndpoint := "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/management/v2"
			fakeCf.CmdReturns(keyInfo, nil)
			username, password, endpoint, err := GetUsernamePasswordEndpoint(fakeCf, pccService, key)
			Expect(username).To(Equal(expectedUsername))
			Expect(password).To(Equal(expectedPassword))
			Expect(endpoint).To(Equal(expectedEndpoint))
			Expect(err).To(BeNil())
		})
		It("Returns an error.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("", errors.New("CF Command Error"))
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf, "", "")
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incorrect JSON.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("{", nil)
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf, "", "")
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incomplete JSON.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			pccService := "jjack"
			key := "mykey"
			keyInfo := `Getting key mykey for service instance jjack as admin...

{
 
`
			fakeCf.CmdReturns(keyInfo, nil)
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf, pccService, key)
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

	Context("Handling Regions and Listing Indexes", func(){
		It("Handling unrecognized regions", func(){
			clusterCommand := "list indexes"
			urlResponse := `{"statusCode":"ENTITY_NOT_FOUND","statusMessage":"RegionConfig with id = sdljf not found.","result":[]}`
			_, err := GetTableFromUrlResponse(clusterCommand, urlResponse)
			Expect(err).To(Not(BeNil()))
		})
		It("No region provided when listing indexes", func(){
			endpoint := " https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/geode-management/v2"
			clusterCommand := "list indexes"
			_, err := getCompleteEndpoint(endpoint, clusterCommand)
			Expect(err).To(Not(BeNil()))
		})
	})
	Context("Safekeeping tests", func(){
		It("Handling unauthenticated requests", func(){
			clusterCommand := "list indexes"
			urlResponse := `{"statusCode":"UNAUTHENTICATED","statusMessage":"Authentication error. Please check your credentials.","result":[]}`
			_, err := GetTableFromUrlResponse(clusterCommand, urlResponse)
			Expect(err).To(Not(BeNil()))
		})
		It("Validate answer in table form", func(){
			clusterCommand := "list members"
			urlResponse := `{"statusCode":"OK","result":[{"config":{"class":"org.apache.geode.management.configuration.MemberConfig"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.MemberInformation","name":"cacheserver-bc54e683-3e01-4767-9efd-5ac1394212f4","id":"bc54e683-3e01-4767-9efd-5ac1394212f4(cacheserver-bc54e683-3e01-4767-9efd-5ac1394212f4:1)<v1>:56153","workingDirPath":"/var/vcap/store/gemfire-server","groups":"cacheserver-bc54e683-3e01-4767-9efd-5ac1394212f4","logFilePath":"/var/vcap/sys/log/gemfire-server/gemfire/server.log","statArchiveFilePath":"/var/vcap/store/gemfire-server/statistics.gfs","locators":"bc54e683-3e01-4767-9efd-5ac1394212f4.locator-server.applevalley-services-subnet.service-instance-6c5d4877-45b6-4213-b0b5-c479a039e37f.bosh[55221]","status":"online","heapUsage":78,"maxHeapSize":3059,"initHeapSize":3090,"cacheXmlFilePath":"/var/vcap/store/gemfire-server/cache.xml","host":"bc54e683-3e01-4767-9efd-5ac1394212f4.locator-server.applevalley-services-subnet.service-instance-6c5d4877-45b6-4213-b0b5-c479a039e37f.bosh","processId":1,"locatorPort":0,"httpServicePort":7070,"httpServiceBindAddress":"bc54e683-3e01-4767-9efd-5ac1394212f4.locator-server.applevalley-services-subnet.service-instance-6c5d4877-45b6-4213-b0b5-c479a039e37f.bosh","clientCount":0,"cpuUsage":0.0,"hostedRegions":["TEST2","r2","region1","TEST3","testing_example3","TEST4","testing_example1","testing_example2","example_partition_region","r0","r1","test_","TEST1"],"webSSL":true,"server":true,"coordinator":false,"cacheServerInfo":[{"port":40404,"maxConnections":800,"maxThreads":0,"running":true}],"secured":false},{"class":"org.apache.geode.management.runtime.MemberInformation","name":"locator-bc54e683-3e01-4767-9efd-5ac1394212f4","id":"bc54e683-3e01-4767-9efd-5ac1394212f4(locator-bc54e683-3e01-4767-9efd-5ac1394212f4:1:locator)<ec><v0>:56152","workingDirPath":"/var/vcap/store/gemfire-locator","logFilePath":"/var/vcap/sys/log/gemfire-locator/gemfire/locator.log","statArchiveFilePath":"/var/vcap/store/gemfire-locator/statistics.gfs","locators":"bc54e683-3e01-4767-9efd-5ac1394212f4.locator-server.applevalley-services-subnet.service-instance-6c5d4877-45b6-4213-b0b5-c479a039e37f.bosh[55221],10.0.8.6[55221]","status":"online","heapUsage":126,"maxHeapSize":494,"initHeapSize":512,"cacheXmlFilePath":"/var/vcap/store/gemfire-locator","host":"bc54e683-3e01-4767-9efd-5ac1394212f4.locator-server.applevalley-services-subnet.service-instance-6c5d4877-45b6-4213-b0b5-c479a039e37f.bosh","processId":1,"locatorPort":55221,"httpServicePort":8080,"httpServiceBindAddress":"bc54e683-3e01-4767-9efd-5ac1394212f4.locator-server.applevalley-services-subnet.service-instance-6c5d4877-45b6-4213-b0b5-c479a039e37f.bosh","clientCount":0,"cpuUsage":0.0,"webSSL":true,"server":false,"coordinator":true,"secured":true}]}]}`
			response, _ := GetTableFromUrlResponse(clusterCommand, urlResponse)
			expectedResponse := `Status Code: OK

 id                 | host               | status             | pid                |
 ------------------------------------------------------------------------------------
 bc54e683-3e01-4767…| bc54e683-3e01-4767…| online             |                    |

Number of Results: 1
To see the full output, append -j to your command.`
			Expect(response).To(Equal(expectedResponse))
		})
		It("Validate JSON from URL response", func(){
			urlResponse := `{"statusCode":"OK","result":[{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","groups":["cacheserver-bc54e683-3e01-4767-9efd-5ac1394212f4"],"regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"region1","type":"REPLICATE","uri":"/management/v2/regions/region1"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"partitionAttributes":{"redundantCopies":"1"},"evictionAttributes":{"lruHeapPercentage":{"action":"LOCAL_DESTROY"}},"dataPolicy":"PARTITION","concurrencyChecksEnabled":true},"name":"example_partition_region","type":"PARTITION_REDUNDANT_HEAP_LRU","uri":"/management/v2/regions/example_partition_region"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"indexes":[{"class":"org.apache.geode.cache.configuration.RegionConfig$Index","name":"index1","expression":"id","fromClause":"/r0","keyIndex":true,"type":"key"}],"name":"r0","type":"REPLICATE","uri":"/management/v2/regions/r0"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"r1","type":"REPLICATE","uri":"/management/v2/regions/r1"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"r2","type":"REPLICATE","uri":"/management/v2/regions/r2"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"testing_example1","type":"REPLICATE","uri":"/management/v2/regions/testing_example1"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"testing_example2","type":"REPLICATE","uri":"/management/v2/regions/testing_example2"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"test_","type":"REPLICATE","uri":"/management/v2/regions/test_"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"TEST1","type":"REPLICATE","uri":"/management/v2/regions/TEST1"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"TEST2","type":"REPLICATE","uri":"/management/v2/regions/TEST2"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"TEST3","type":"REPLICATE","uri":"/management/v2/regions/TEST3"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"TEST4","type":"REPLICATE","uri":"/management/v2/regions/TEST4"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]},{"config":{"class":"org.apache.geode.cache.configuration.RegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"testing_example3","type":"REPLICATE","uri":"/management/v2/regions/testing_example3"},"runtimeInfo":[{"class":"org.apache.geode.management.runtime.RuntimeRegionInfo","entryCount":0}]}]}`
			expectedResponse := `{
  "statusCode": "OK",
  "statusMessage": "",
  "memberStatus": null,
  "result": [
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "groups": [
          "cacheserver-bc54e683-3e01-4767-9efd-5ac1394212f4"
        ],
        "name": "region1",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/region1"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "example_partition_region",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "PARTITION",
          "evictionAttributes": {
            "lruHeapPercentage": {
              "action": "LOCAL_DESTROY"
            }
          },
          "partitionAttributes": {
            "redundantCopies": "1"
          }
        },
        "type": "PARTITION_REDUNDANT_HEAP_LRU",
        "uri": "/management/v2/regions/example_partition_region"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "indexes": [
          {
            "class": "org.apache.geode.cache.configuration.RegionConfig$Index",
            "expression": "id",
            "fromClause": "/r0",
            "keyIndex": true,
            "name": "index1",
            "type": "key"
          }
        ],
        "name": "r0",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/r0"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "r1",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/r1"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "r2",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/r2"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "testing_example1",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/testing_example1"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "testing_example2",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/testing_example2"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "test_",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/test_"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "TEST1",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/TEST1"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "TEST2",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/TEST2"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "TEST3",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/TEST3"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "TEST4",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/TEST4"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    },
    {
      "config": {
        "class": "org.apache.geode.cache.configuration.RegionConfig",
        "name": "testing_example3",
        "regionAttributes": {
          "concurrencyChecksEnabled": true,
          "dataPolicy": "REPLICATE",
          "scope": "DISTRIBUTED_ACK"
        },
        "type": "REPLICATE",
        "uri": "/management/v2/regions/testing_example3"
      },
      "runtimeInfo": [
        {
          "class": "org.apache.geode.management.runtime.RuntimeRegionInfo",
          "entryCount": 0
        }
      ]
    }
  ]
}`
			response, _:= GetJsonFromUrlResponse(urlResponse)
			Expect(response).To(Equal(expectedResponse))
		})
		It("Validate table filling", func(){
			columnSize := 20
			value := "some string"
			filler := "-"
			response := Fill(columnSize, value, filler)
			expectedResponse := " some string--------"
			Expect(response).To(Equal(expectedResponse))
		})
	})
})
