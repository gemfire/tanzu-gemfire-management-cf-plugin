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

func executeCommand(endpointURL string, httpAction string, commandData *domain.CommandData) (urlResponse string, err error) {
	var bodyReader io.Reader

	if httpAction == "POST" {
		bodyReader, err = getBodyReader(commandData.UserCommand.Parameters["-j"])
		if err != nil {
			return "", err
		}
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest(httpAction, endpointURL, bodyReader)
	req.SetBasicAuth(commandData.ConnnectionData.Username, commandData.ConnnectionData.Password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	return getURLOutput(resp)
}

func getBodyReader(jsonFile string) (bodyReader io.Reader, err error) {
	if jsonFile == "" {
		err = errors.New(util.NoJsonFileProvidedMessage)
		return
	}
	if jsonFile[0] == '@' && len(jsonFile) > 1 {
		bodyReader, err = os.Open(jsonFile[1:])
		if err != nil {
			return
		}
	} else {
		bodyReader = strings.NewReader(jsonFile)
	}
	return
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
func GetTargetAndClusterCommand(args []string) (target string, userCommand domain.UserCommand) {
	target = os.Getenv("CFPCC")
	if len(args) < 2 {
		return
	}

	commandStart := 2
	if target == "" {
		target = args[1]
	} else if target != args[1] {
		commandStart = 1
	}

	userCommand.Parameters = make(map[string]string)
	// find the command name before the options
	var option = ""
	for i := commandStart; i < len(args); i++ {
		token := args[i]
		if strings.HasPrefix(token, "-") {
			if option != "" {
				userCommand.Parameters[option] = "true"
			}
			option = token
		} else if option == "" {
			userCommand.Command += token + " "
		} else {
			userCommand.Parameters[option] = token
			option = ""
		}
	}
	userCommand.Command = strings.Trim(userCommand.Command, " ")
	if option != "" {
		userCommand.Parameters[option] = "true"
	}
	return
}

// GetEndPoints retrieves available endpoint from the Swagger endpoint on the PCC manageability service
func GetEndPoints(commandData *domain.CommandData) error {
	apiDocUrl := commandData.ConnnectionData.LocatorAddress + "/management/experimental/api-docs"
	urlResponse, err := executeCommand(apiDocUrl, "GET", commandData)

	if err != nil {
		return errors.New("unable to reach " + apiDocUrl + ": " + err.Error())
	}

	err = json.Unmarshal([]byte(urlResponse), &commandData.FirstResponse)

	if err != nil {
		return errors.New("invalid response " + urlResponse)
	}

	for url, v := range commandData.FirstResponse.Paths {
		for methodType := range v {
			var endpoint domain.IndividualEndpoint
			endpoint.URL = url
			endpoint.HTTPMethod = methodType
			endpoint.CommandCall = commandData.FirstResponse.Paths[url][methodType].Summary
			commandData.AvailableEndpoints = append(commandData.AvailableEndpoints, endpoint)
		}
	}

	return nil
}

// RequestToEndPoint makes the request to a manageability service endpoint and returns a response
func RequestToEndPoint(commandData *domain.CommandData) (string, error) {
	secondEndpoint := commandData.ConnnectionData.LocatorAddress + "/management" + commandData.Endpoint.URL
	urlResponse, err := executeCommand(secondEndpoint, strings.ToUpper(commandData.Endpoint.HTTPMethod), commandData)
	return urlResponse, err
}
