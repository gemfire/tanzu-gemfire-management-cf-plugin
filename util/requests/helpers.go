package requests

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
)

//go:generate counterfeiter . CfService

// CfService is the interface for executing CF commands
type CfService interface {
	Cmd(name string, options ...string) (string, error)
}

// GetServiceKeyFromPCCInstance retrieves the service key by name from the CF foundation
func GetServiceKeyFromPCCInstance(cf CfService, target string) (serviceKey string, err error) {
	servKeyOutput, err := cf.Cmd("service-keys", target)
	if err != nil {
		return "", err
	}
	splitKeys := strings.Split(servKeyOutput, "\n")
	hasKey := false
	if strings.Contains(splitKeys[1], "No service key for service instance") {
		return "", errors.New(util.NoServiceKeyMessage)
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
		err = errors.New(util.NoServiceKeyMessage)
	}
	return
}

// GetUsernamePasswordEndpoinFromServiceKey extracts the username, password and endpoint from a CF service key
func GetUsernamePasswordEndpoinFromServiceKey(cf CfService, target string, serviceKey string) (username string, password string, endpoint string, err error) {
	username = ""
	password = ""
	endpoint = ""
	keyInfo, err := cf.Cmd("service-key", target, serviceKey)
	if err != nil {
		return "", "", "", err
	}
	splitKeyInfo := strings.Split(keyInfo, "\n")
	if len(splitKeyInfo) < 2 {
		return "", "", "", errors.New(util.InvalidServiceKeyResponse)
	}
	splitKeyInfo = splitKeyInfo[2:] //take out first two lines of cf service-key ... output
	joinKeyInfo := strings.Join(splitKeyInfo, "\n")
	serviceKeyStruct := domain.ServiceKey{}

	err = json.Unmarshal([]byte(joinKeyInfo), &serviceKeyStruct)
	if err != nil {
		return "", "", "", err
	}
	endpoint = serviceKeyStruct.Urls.Management
	if endpoint == "" {
		endpoint = strings.TrimSuffix(serviceKeyStruct.Urls.Gfsh, "/gemfire/v1")
	}
	for _, user := range serviceKeyStruct.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			username = user.Username
			password = user.Password
		}
	}
	return
}

func executeCommand(endpointURL string, httpAction string, commandData *domain.CommandData) (urlResponse string, err error) {
	if httpAction == "POST" {
		return executePostCommand(endpointURL, commandData.JSONFile)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(httpAction, endpointURL, nil)
	req.SetBasicAuth(commandData.Username, commandData.Password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	return getURLOutput(resp)
}

func executePostCommand(endpointURL string, jsonFile string) (urlResponse string, err error) {
	if jsonFile == "" {
		return "", errors.New(util.NoJsonFileProvidedMessage)
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
	req, err = http.NewRequest("POST", endpointURL, f)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return getURLOutput(resp)
}

func getURLOutput(resp *http.Response) (urlResponse string, err error) {
	respInASCII, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	urlResponse = fmt.Sprintf("%s", respInASCII)
	return urlResponse, nil
}

// GetTargetAndClusterCommand extracts the target and command from the args and environment variables
func GetTargetAndClusterCommand(args []string) (target string, userCommand domain.UserCommand, err error) {
	if os.Getenv("CFPCC") != "" {
		target = os.Getenv("CFPCC")
	}

	if len(args) < 2 {
		err = errors.New(util.IncorrectUserInputMessage)
		return
	}
	var commands []string
	if args[1] == target {
		commands = args[2:]
	} else if target == "" {
		target = args[1]
		commands = args[2:]
	} else {
		commands = args[1:]
	}

	// find the command name before the options
	for _, command := range commands {
		if strings.HasPrefix(command, "-") {
			break
		}
		userCommand.Command += command + " "
	}
	userCommand.Command = strings.Trim(userCommand.Command, " ")
	return
}

// GetEndPoints retrieves available endpoint from the Swagger endpoint on the PCC manageability service
func GetEndPoints(commandData *domain.CommandData) error {
	urlResponse, err := executeCommand(commandData.ConnnectionData.LocatorAddress+"/management/experimental/api-docs", "GET", commandData)
	if err == nil {
		err = json.Unmarshal([]byte(urlResponse), &commandData.FirstResponse)
		for url, v := range commandData.FirstResponse.Paths {
			for methodType := range v {
				var endpoint domain.IndividualEndpoint
				endpoint.URL = url
				endpoint.HTTPMethod = methodType
				endpoint.CommandCall = commandData.FirstResponse.Paths[url][methodType].Summary
				commandData.AvailableEndpoints = append(commandData.AvailableEndpoints, endpoint)
			}
		}
	}
	return err
}

// RequestToEndPoint makes the request to a manageability service endpoint and returns a response
func RequestToEndPoint(commandData *domain.CommandData) (string, error) {
	secondEndpoint := commandData.ConnnectionData.LocatorAddress + "/management" + commandData.Endpoint.URL
	urlResponse, err := executeCommand(secondEndpoint, strings.ToUpper(commandData.Endpoint.HTTPMethod), commandData)
	return urlResponse, err
}

// HasIDifNeeded checks if an ID needs to be passed and if absent produces an error
func HasIDifNeeded(commandData *domain.CommandData) error {
	if strings.Contains(commandData.Endpoint.URL, "{id}") {
		if commandData.ID == "" {
			return errors.New(util.NoIDGivenMessage)
		}
		commandData.Endpoint.URL = strings.Replace(commandData.Endpoint.URL, "{id}", commandData.ID, 1)
	}
	return nil
}

// HasRegionIfNeeded checks if a Region needs to passed and if absent produces an error
func HasRegionIfNeeded(commandData *domain.CommandData) error {
	if strings.Contains(commandData.Endpoint.URL, "{regionName}") {
		if commandData.Region == "" {
			return errors.New(util.NoRegionGivenMessage)
		}
		commandData.Endpoint.URL = strings.Replace(commandData.Endpoint.URL, "{regionName}", commandData.Region, 1)
	}
	return nil
}
