package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
)

// CommandProcessor struct holds the implementation for the RequestHelper interface
type CommandProcessor struct {
	requester impl.RequestHelper
}

// NewCommandProcessor provides the constructor for the CommandProcessor
func NewCommandProcessor(requester impl.RequestHelper) (CommandProcessor, error) {
	return CommandProcessor{requester: requester}, nil
}

// ProcessCommand handles the common steps for executing a command against the Geode cluster
func (c *CommandProcessor) ProcessCommand(commandData *domain.CommandData) (err error) {
	err = getEndPoints(commandData, c.requester)
	if err != nil {
		return
	}

	userCommand := commandData.UserCommand.Command
	if userCommand == "commands" {
		commandNames := sortCommandNames(commandData)
		for _, commandName := range commandNames {
			fmt.Println(DescribeEndpoint(commandData.AvailableEndpoints[commandName]))
		}
		return
	}

	restEndPoint, available := commandData.AvailableEndpoints[userCommand]
	if !available {
		err = errors.New("Invalid command: " + userCommand)
		return
	}

	if HasOption(commandData.UserCommand.Parameters, []string{"-h", "--help", "-help"}) {
		for _, command := range commandData.AvailableEndpoints {
			if command.CommandName == userCommand {
				fmt.Println(DescribeEndpoint(command))
				fmt.Println(GeneralOptions)
			}
		}
		return
	}

	err = checkRequiredParam(restEndPoint, commandData.UserCommand)
	if err != nil {
		return
	}

	urlResponse, err := executeCommand(commandData, c.requester)
	if err != nil {
		return
	}

	var jqFilter string
	if HasOption(commandData.UserCommand.Parameters, []string{"-t", "--table"}) {
		jqFilter = GetOption(commandData.UserCommand.Parameters, []string{"--table", "-t"})
		// if no jqFilter is specified by the user, use the default defined by the rest end point
		if jqFilter == "" {
			jqFilter = restEndPoint.JQFilter
		}
		// if no default jqFilter is configured, then use the entire json
		if jqFilter == "" {
			jqFilter = "."
		}
	}

	jsonToBePrinted, err := FormatResponse(urlResponse, jqFilter)
	if err != nil {
		return
	}
	fmt.Println(jsonToBePrinted)

	return
}

func checkRequiredParam(restEndPoint domain.RestEndPoint, command domain.UserCommand) error {
	for _, s := range restEndPoint.Parameters {
		if s.Required {
			value := command.Parameters["--"+s.Name]
			if value == "" {
				return errors.New("Required Parameter is missing: " + s.Name)
			}
		}
	}
	return nil
}

// GetEndPoints retrieves available endpoint from the Swagger endpoint on the PCC manageability service
func getEndPoints(commandData *domain.CommandData, requester impl.RequestHelper) error {
	apiDocURL := commandData.ConnnectionData.LocatorAddress + "/management/experimental/api-docs"
	urlResponse, err := requester.Exchange(apiDocURL, "GET", nil, commandData.ConnnectionData.Username,
		commandData.ConnnectionData.Password)

	if err != nil {
		return errors.New("unable to reach " + apiDocURL + ": " + err.Error())
	}

	var apiPaths domain.RestAPI
	err = json.Unmarshal([]byte(urlResponse), &apiPaths)

	if err != nil {
		return errors.New("invalid response " + urlResponse)
	}

	commandData.AvailableEndpoints = make(map[string]domain.RestEndPoint)
	for url, v := range apiPaths.Paths {
		for methodType := range v {
			var endpoint domain.RestEndPoint
			endpoint.URL = url
			endpoint.HTTPMethod = methodType
			endpoint.CommandName = apiPaths.Paths[url][methodType].CommandName
			endpoint.JQFilter = apiPaths.Paths[url][methodType].JQFilter
			endpoint.Parameters = []domain.RestAPIParam{}
			endpoint.Parameters = apiPaths.Paths[url][methodType].Parameters
			commandData.AvailableEndpoints[endpoint.CommandName] = endpoint
		}
	}
	return nil
}

func executeCommand(commandData *domain.CommandData, requester impl.RequestHelper) (urlResponse string, err error) {
	var bodyReader io.Reader

	restEndPoint, _ := commandData.AvailableEndpoints[commandData.UserCommand.Command]
	httpAction := strings.ToUpper(restEndPoint.HTTPMethod)
	endpointURL := makeURL(restEndPoint, commandData)

	if httpAction == "POST" {
		bodyReader, err = getBodyReader(commandData.UserCommand.Parameters["--body"])
		if err != nil {
			return "", err
		}
	}
	return requester.Exchange(endpointURL, httpAction, bodyReader, commandData.ConnnectionData.Username,
		commandData.ConnnectionData.Password)
}

func getBodyReader(jsonFile string) (bodyReader io.Reader, err error) {
	if jsonFile == "" {
		err = errors.New(NoJSONFileProvidedMessage)
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

func makeURL(restEndPoint domain.RestEndPoint, commandData *domain.CommandData) (requestURL string) {
	requestURL = commandData.ConnnectionData.LocatorAddress + "/management" + restEndPoint.URL
	var query string
	for _, param := range restEndPoint.Parameters {
		value, ok := commandData.UserCommand.Parameters["--"+param.Name]
		if ok {
			switch param.In {
			case "path":
				requestURL = strings.ReplaceAll(requestURL, "{"+param.Name+"}", url.PathEscape(value))
			case "query":
				if len(query) == 0 {
					query = "?" + param.Name + "=" + url.PathEscape(value)
				} else {
					query = query + "&" + param.Name + "=" + url.PathEscape(value)
				}
			}
		}
	}

	requestURL = requestURL + query
	return
}

func sortCommandNames(commandData *domain.CommandData) (commandNames []string) {
	commandNames = make([]string, 0, len(commandData.AvailableEndpoints))
	for _, command := range commandData.AvailableEndpoints {
		commandNames = append(commandNames, command.CommandName)
	}
	sort.Strings(commandNames)
	return
}
