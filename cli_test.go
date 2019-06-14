package main

import (
	"testing"
)

func TestValidatePCCInstance(t *testing.T) {
	ourPCCInstance := "pcc1"
	pccInstancesAvailable := []string{"pcc1", "pcc2", "pcc3"}
	err := ValidatePCCInstance(ourPCCInstance, pccInstancesAvailable)
	if err != nil{
		t.Error("Unexpected Error")
	}
}

func TestGetAnswerFromUrlResponse(t *testing.T) {
	clusterCommand := "list-members"
	urlResponse := `{"statusCode":"OK","result":[{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","groups":["group3"],"id":"server3","host":"10.118.19.35","status":"online","pid":63340,"cacheServers":[{"port":40406,"maxConnections":800,"maxThreads":0}],"maxHeap":4096,"initialHeap":256,"usedHeap":55,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/server3/server3.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/server3","clientConnections":0,"locator":false,"coordinator":false,"uri":"/members/server3"},{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","groups":["group-1"],"id":"server1","host":"10.118.19.35","status":"online","pid":39383,"cacheServers":[{"port":53105,"maxConnections":800,"maxThreads":0}],"maxHeap":4096,"initialHeap":256,"usedHeap":37,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/server1/server1.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/server1","clientConnections":0,"locator":false,"coordinator":false,"uri":"/members/server1"},{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","groups":["group2"],"id":"server2","host":"10.118.19.35","status":"online","pid":62028,"cacheServers":[{"port":40404,"maxConnections":800,"maxThreads":0}],"maxHeap":4096,"initialHeap":256,"usedHeap":57,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/server2/server2.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/server2","clientConnections":0,"locator":false,"coordinator":false,"uri":"/members/server2"},{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","id":"locator1","host":"10.118.19.35","status":"online","pid":39363,"port":10334,"maxHeap":4096,"initialHeap":256,"usedHeap":132,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/locator1/locator1.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/locator1","clientConnections":0,"locator":true,"coordinator":true,"uri":"/members/locator1"}]}`
	response := GetAnswerFromUrlResponse(clusterCommand, urlResponse)
	expectedResponse := `Status Code: OK

 id                 | host               | status             | pid                |
 ------------------------------------------------------------------------------------
 server3            | 10.118.19.35       | online             | 63340              |
 server1            | 10.118.19.35       | online             | 39383              |
 server2            | 10.118.19.35       | online             | 62028              |
 locator1           | 10.118.19.35       | online             | 39363              |

Number of Members: 4`
	if response != expectedResponse{
		t.Error("Unexpected Error")
	}
}

func TestEditResponseOnGroup(t *testing.T) {
	urlResponse := `{"statusCode":"OK","result":[{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","groups":["group3","group2"],"regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjack","type":"REPLICATE","entryCount":0,"uri":"/regions/jjack"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","groups":["group2"],"regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjoris","type":"REPLICATE","entryCount":0,"uri":"/regions/jjoris"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjens","type":"REPLICATE","entryCount":0,"uri":"/regions/jjens"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjinmei","type":"REPLICATE","entryCount":0,"uri":"/regions/jjinmei"}]}`
	groups := []string{"group2", "group3"}
	clusterCommand := "list-regions"
	response := EditResponseOnGroup(urlResponse, groups, clusterCommand)
	expectedResponse := `{"statusCode":"OK","statusMessage":"","memberStatus":null,"result":[{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group3","group2"],"name":"jjack","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjack"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group2"],"name":"jjoris","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjoris"}]}`
	if response != expectedResponse{
		t.Error("Unexpected Error")
	}
}

func TestGetJsonFromUrlResponse(t *testing.T) {
	urlResponse := `{"statusCode":"OK","statusMessage":"","memberStatus":null,"result":[{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group3","group2"],"name":"jjack","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjack"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group2"],"name":"jjoris","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjoris"}]}`
	expectedResponse := `{
  "statusCode": "OK",
  "statusMessage": "",
  "memberStatus": null,
  "result": [
    {
      "class": "org.apache.geode.management.configuration.RuntimeRegionConfig",
      "entryCount": 0,
      "groups": [
        "group3",
        "group2"
      ],
      "name": "jjack",
      "regionAttributes": {
        "concurrencyChecksEnabled": true,
        "dataPolicy": "REPLICATE",
        "scope": "DISTRIBUTED_ACK"
      },
      "type": "REPLICATE",
      "uri": "/regions/jjack"
    },
    {
      "class": "org.apache.geode.management.configuration.RuntimeRegionConfig",
      "entryCount": 0,
      "groups": [
        "group2"
      ],
      "name": "jjoris",
      "regionAttributes": {
        "concurrencyChecksEnabled": true,
        "dataPolicy": "REPLICATE",
        "scope": "DISTRIBUTED_ACK"
      },
      "type": "REPLICATE",
      "uri": "/regions/jjoris"
    }
  ]
}`
	response := GetJsonFromUrlResponse(urlResponse)
	if response != expectedResponse{
		t.Error("Unexpected Error")
	}
}
