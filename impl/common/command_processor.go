package common

import (
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
	err = GetEndPoints(commandData, c.requester)
	if err != nil {
		return
	}

	userCommand := commandData.UserCommand.Command
	restEndPoint, available := commandData.AvailableEndpoints[userCommand]

	if userCommand == "commands" || !available {
		commandNames := sortCommandNames(commandData)
		for _, commandName := range commandNames {
			fmt.Println(DescribeEndpoint(commandData.AvailableEndpoints[commandName], false))
		}
		if userCommand != "commands" {
			err = errors.New("Invalid command: " + userCommand)
		}
		return
	}

	if HasOption(commandData.UserCommand.Parameters, []string{"-h", "--help", "-help"}) {
		fmt.Println(DescribeEndpoint(restEndPoint, true))
		return
	}

	err = CheckRequiredParam(restEndPoint, commandData.UserCommand)
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

func CheckRequiredParam(restEndPoint domain.RestEndPoint, command domain.UserCommand) error {
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

func executeCommand(commandData *domain.CommandData, requester impl.RequestHelper) (urlResponse string, err error) {
	var bodyReader io.Reader

	restEndPoint, _ := commandData.AvailableEndpoints[commandData.UserCommand.Command]
	httpAction := strings.ToUpper(restEndPoint.HTTPMethod)
	endpointURL, body := prepareRequest(restEndPoint, commandData)

	if body != "" {
		bodyReader, err = getBodyReader(body)
		if err != nil {
			return "", err
		}
	}
	return requester.Exchange(endpointURL, httpAction, bodyReader, commandData.ConnnectionData.Username,
		commandData.ConnnectionData.Password)
}

func getBodyReader(jsonFile string) (bodyReader io.Reader, err error) {
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

func prepareRequest(restEndPoint domain.RestEndPoint, commandData *domain.CommandData) (requestURL string, body string) {
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
			case "body":
				body = value
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
