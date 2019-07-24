package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
)


func getTableHeadersFromClusterCommand(clusterCommand string) (tableHeaders []string){
	switch clusterCommand {
	case "list_regions":
		tableHeaders = []string{"name", "type", "groups", "entryCount", "regionAttributes"}
	case "list_members":
		tableHeaders = []string{"memberName", "host", "status", "pid"}
	case "list_gateway-receivers":
		tableHeaders = []string{"hostnameForSenders", "uri", "group", "class"}
	case "list_indexes":
		tableHeaders = []string{"name", "type", "fromClause", "expression"}
	default:
		return
	}
	return
}


func Fill(columnSize int, value string, filler string) (response string){
	if len(value) > columnSize - 1{
		response = " " + value[:columnSize-len([]rune(Ellipsis)) -1] + Ellipsis
		return
	}
	numFillerChars := columnSize - len(value) - 1
	response = " " + value + strings.Repeat(filler, numFillerChars)
	return
}

func GetTableFromUrlResponse(clusterCommand string, urlResponse string) (response string, err error){
	urlOutput := ClusterManagementResults{}
	err = json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		return "", err
	}
	if urlOutput.StatusCode == "UNAUTHENTICATED"{
		return "", errors.New(NotAuthenticatedMessage)
	} else if urlOutput.StatusCode == "ENTITY_NOT_FOUND"{
		return "", errors.New(NonExistentRegionMessage)
	}
	response = "Status Code: " + urlOutput.StatusCode + "\n"
	if urlOutput.StatusMessage != ""{
		response += "Status Message: " + urlOutput.StatusMessage + "\n"
	}
	response += "\n"

	tableHeaders := getTableHeadersFromClusterCommand(clusterCommand)
	for _, header := range tableHeaders {
		response += Fill(20, header, " ") + "|"
	}
	response += "\n" + Fill (20 * len(tableHeaders) + 5, "", "-") + "\n"

	memberCount := 0
	runtimeInfoIsNil := false
	for _, result := range urlOutput.Results{
		memberCount++
		if err != nil {
			return "", err
		}
		if result.RuntimeInfo == nil{
			runtimeInfoIsNil = true
		}
		for _, key := range tableHeaders {
			if result.Config[key] == nil && (runtimeInfoIsNil || result.RuntimeInfo[0][key] == nil){
				response += Fill(20, "", " ") + "|"
			} else {
				resultVal := result.Config[key]
				if resultVal == nil && !runtimeInfoIsNil{
					resultVal = result.RuntimeInfo[0][key]
				}
				if fmt.Sprintf("%T", resultVal) == "float64"{
					resultVal = fmt.Sprintf("%.0f", resultVal)
				}
				response += Fill(20, fmt.Sprintf("%s",resultVal), " ") + "|"
			}
		}
		response += "\n"
	}

	response += "\nNumber of Results: " + strconv.Itoa(memberCount)
	if strings.Contains(response, Ellipsis){
		response += "\nTo see the full output, append -j to your command."
	}
	return
}



func GetJsonFromUrlResponse(urlResponse string) (jsonOutput string, err error){
	urlOutput := ClusterManagementResults{}
	err = json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		return "", err
	}
	jsonExtracted, err := json.MarshalIndent(urlOutput, "", "  ")
	if err != nil {
		return "", err
	}
	jsonOutput = string(jsonExtracted)
	return
}
