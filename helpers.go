package main

import (
	"bytes"
	"code.cloudfoundry.org/cli/cf/errors"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetServiceKeyFromPCCInstance(cf cfservice.CfService, pccService string) (serviceKey string, err error) {
	servKeyOutput, err := cf.Cmd("service-keys", pccService)
	if err != nil{
		return "", err
	}
	splitKeys := strings.Split(servKeyOutput, "\n")
	hasKey := false
	if strings.Contains(splitKeys[1], "No service key for service instance"){
		return "", errors.New(NoServiceKeyMessage)
	}
	for _, value := range splitKeys {
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
		return serviceKey, errors.New(NoServiceKeyMessage)
	}
	return
}

func GetUsernamePasswordEndpoint(cf cfservice.CfService, pccService string, key string) (username string, password string, endpoint string, err error) {
	username = ""
	password = ""
	endpoint = ""
	keyInfo, err := cf.Cmd("service-key", pccService, key)
	if err != nil {
		return "", "", "", err
	}
	splitKeyInfo := strings.Split(keyInfo, "\n")
	if len(splitKeyInfo) < 2{
		return "", "", "", errors.New(InvalidServiceKeyResponse)
	}
	splitKeyInfo = splitKeyInfo[2:] //take out first two lines of cf service-key ... output
	joinKeyInfo := strings.Join(splitKeyInfo, "\n")
	serviceKey := ServiceKey{}

	err = json.Unmarshal([]byte(joinKeyInfo), &serviceKey)
	if err != nil {
		return "", "", "", err
	}
	endpoint = serviceKey.Urls.Management
	endpoint = strings.TrimSuffix(serviceKey.Urls.Gfsh, "gemfire/v1") + "management/v2/cli"
	for _ , user := range serviceKey.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			username = user.Username
			password = user.Password
		}
	}
	return
}

func getCompleteEndpoint(endpoint string, clusterCommand string) (string, error){
	urlEnding := ""
	switch clusterCommand{
	case "list regions":
		urlEnding = "/regions"
	case "list members":
		urlEnding = "/members"
	case "list gateway-receivers":
		urlEnding = "/gateways/receivers"
	case "list indexes":
		if region == ""{
			return "", errors.New(NoRegionGivenMessage)
		}
		if strings.HasPrefix(region, "/"){
			region = region[1:]
		}
		urlEnding = "/regions/" + region + "/indexes"
	case "post region":
		urlEnding = "/regions"
	default:
		return endpoint, nil
	}
	endpoint = endpoint + urlEnding + "?group=" + group
	return endpoint, nil
}

func getTableHeadersFromClusterCommand(clusterCommand string) (tableHeaders []string){
	switch clusterCommand {
	case "list regions":
		tableHeaders = []string{"name", "type", "groups", "entryCount", "regionAttributes"}
	case "list members":
		tableHeaders = []string{"id", "host", "status", "pid"}
	case "list gateway-receivers":
		tableHeaders = []string{"hostnameForSenders", "uri", "group", "class"}
	case "list indexes":
		tableHeaders = []string{"name", "type", "fromClause", "expression"}
	default:
		return
	}
	return
}

func getHttpRequestMethod(call RestAPICall) (requestMethod string, err error){
	switch call.action {
	case "list":
		return "GET", nil
	case "get":
		return "GET", nil
	case "delete":
		return "DELETE", nil
	case "create":
		return "POST", nil
	case "post":
		return "POST", nil
	default:
		return call.action, errors.New(UnsupportedClusterCommandMessage)
	}
}

func getUrlOutput(endpointUrl string, username string, password string, httpAction string) (urlResponse string, err error){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	requestBody, err := json.Marshal(APICallStruct)

	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(httpAction, endpointUrl, bytes.NewBuffer(requestBody))
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil{
		return "", err
	}

	respInAscii, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil{
		return "", err
	}

	urlResponse = fmt.Sprintf("%s", respInAscii)
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


func isUsingPCCfromEnvironmentVariables(args []string) bool{
	if os.Getenv("CFPCC") != "" && len(args) >= 3 && args[1] != os.Getenv("CFPCC"){
		return true
	}
	return false
}

func getPCCInUseAndClusterCommand(args []string) (error){
	if isUsingPCCfromEnvironmentVariables(args){
		pccInUse = os.Getenv("CFPCC")
		APICallStruct.action = args[1]
		APICallStruct.target = args[2]
		APICallStruct.command = APICallStruct.action + "_" + APICallStruct.target
	} else if len(args) >= 4 {
		pccInUse = args[1]
		APICallStruct.action = args[2]
		APICallStruct.target = args[3]
		APICallStruct.command = APICallStruct.action + "_" + APICallStruct.target
	} else{
		return errors.New(IncorrectUserInputMessage)
	}
	return nil
}
