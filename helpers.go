package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetServiceKeyFromPCCInstance(cf cfservice.CfService) (serviceKey string, err error) {
	servKeyOutput, err := cf.Cmd("service-keys", target)
	if err != nil {
		return "", err
	}
	splitKeys := strings.Split(servKeyOutput, "\n")
	hasKey := false
	if strings.Contains(splitKeys[1], "No service key for service instance") {
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

func GetUsernamePasswordEndpoinFromServiceKey(cf cfservice.CfService) (username string, password string, endpoint string, err error) {
	username = ""
	password = ""
	endpoint = ""
	keyInfo, err := cf.Cmd("service-key", target, serviceKey)
	if err != nil {
		return "", "", "", err
	}
	splitKeyInfo := strings.Split(keyInfo, "\n")
	if len(splitKeyInfo) < 2 {
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
	if endpoint == "" {
		endpoint = strings.TrimSuffix(serviceKey.Urls.Gfsh, "gemfire/v1")
	}
	for _, user := range serviceKey.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			username = user.Username
			password = user.Password
		}
	}
	return
}

func executeCommand(endpointUrl string, httpAction string) (urlResponse string, err error) {
	if httpAction == "POST" {
		return executePostCommand(endpointUrl)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(httpAction, endpointUrl, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	return getUrlOutput(resp)
}

func executePostCommand(endpointUrl string) (urlResponse string, err error) {
	if jsonFile == "" {
		return "", errors.New(NoJsonFileProvidedMessage)
	}
	var f io.Reader
	var req *http.Request
	if jsonFile[0] == '@' && len(jsonFile) > 1 {
		f, err = os.Open(jsonFile[1:])
		if err != nil {
			return "", err
		}
	} else {
		f = strings.NewReader(jsonFile)
	}
	req, err = http.NewRequest("POST", endpointUrl, f)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return getUrlOutput(resp)
}

func getUrlOutput(resp *http.Response) (urlResponse string, err error) {
	respInAscii, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	urlResponse = fmt.Sprintf("%s", respInAscii)
	return urlResponse, nil
}

func isUsingPCCfromEnvironmentVariables(args []string) bool {
	if os.Getenv("CFPCC") != "" && len(args) >= 3 && args[1] != os.Getenv("CFPCC") {
		return true
	}
	return false
}

func getTargetAndClusterCommand(args []string) error {
	if isUsingPCCfromEnvironmentVariables(args) {
		target = os.Getenv("CFPCC")
		userCommand.command = args[1]
		if len(args) > 2 && !strings.HasPrefix(args[2], "-") {
			userCommand.command += " " + args[2]
		}
	} else if len(args) >= 3 {
		target = args[1]
		userCommand.command = args[2]
		if len(args) > 3 && !strings.HasPrefix(args[3], "-") {
			userCommand.command += " " + args[3]
		}
	} else {
		return errors.New(IncorrectUserInputMessage)
	}
	return nil
}

func getEndPoints() error {
	urlResponse, err := executeCommand(locatorAddress+"/management/experimental/api-docs", "GET")
	err = json.Unmarshal([]byte(urlResponse), &firstResponse)
	for url, v := range firstResponse.Paths {
		for methodType := range v {
			var endpoint IndividualEndpoint
			endpoint.Url = url
			endpoint.HttpMethod = methodType
			endpoint.CommandCall = firstResponse.Paths[url][methodType].Summary
			availableEndpoints = append(availableEndpoints, endpoint)
		}
	}
	return err
}

func requestToEndPoint() (string, error) {
	secondEndpoint := locatorAddress + "/management" + endPoint.Url
	urlResponse, err := executeCommand(secondEndpoint, strings.ToUpper(endPoint.HttpMethod))
	return urlResponse, err
}

func hasIDifNeeded() error {
	if strings.Contains(endPoint.Url, "{id}") {
		if id == "" {
			return errors.New(NoIDGivenMessage)
		}
		endPoint.Url = strings.Replace(endPoint.Url, "{id}", id, 1)
	}
	return nil
}

func hasRegionIfNeeded() error {
	if strings.Contains(endPoint.Url, "{regionName}") {
		if region == "" {
			return errors.New(NoRegionGivenMessage)
		}
		endPoint.Url = strings.Replace(endPoint.Url, "{regionName}", region, 1)
	}
	return nil
}
