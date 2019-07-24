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

func GetUsernamePasswordEndpoint(cf cfservice.CfService) (username string, password string, endpoint string, err error) {
	username = ""
	password = ""
	endpoint = ""
	keyInfo, err := cf.Cmd("service-key", pccInUse, serviceKey)
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


func getUrlOutput(endpointUrl string, httpAction string) (urlResponse string, err error){
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

func executeFirstRequest(endpoint string) (error){
	urlResponse, err := getUrlOutput(endpoint, "GET")
	fmt.Println(endpoint)
	err = json.Unmarshal([]byte(urlResponse), &firstResponse)
	return err
}

func executeSecondRequest() (string, error){
	secondEndpoint := "http://localhost:7070/management/v2/" + firstResponse.Url
	fmt.Println(secondEndpoint)
	urlResponse, err := getUrlOutput(secondEndpoint, firstResponse.HttpMethod)
	return urlResponse, err
}

func hasIDifNeeded() (error){
	if strings.Contains(firstResponse.Url, "{id}"){
		if id == ""{
			return errors.New(NoIDGivenMessage)
		}
		firstResponse.Url = strings.Replace(firstResponse.Url, "{id}", id, 1)
	}
	return nil
}

func hasRegionIfNeeded() (error){
	if strings.Contains(firstResponse.Url, "{regionName}"){
		if region == ""{
			return errors.New(NoRegionGivenMessage)
		}
		firstResponse.Url = strings.Replace(firstResponse.Url, "{regionName}", region, 1)
	}
	return nil
}
