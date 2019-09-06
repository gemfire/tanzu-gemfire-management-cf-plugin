package common

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/format"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/requests"
)

// ProcessCommand handles the common steps for executing a command against the Geode cluster
func ProcessCommand(commandData *domain.CommandData) {
	var err error

	err = requests.GetEndPoints(commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	userCommand := commandData.UserCommand.Command
	if userCommand == "commands" {
		for _, command := range commandData.AvailableEndpoints {
			fmt.Println(Describe(command))
		}
		os.Exit(0)
	}

	restEndPoint, avalable := commandData.AvailableEndpoints[userCommand]
	if !avalable {
		fmt.Println("Invalid command: " + userCommand)
		os.Exit(1)
	}

	err = checkRequiredParam(restEndPoint, commandData.UserCommand)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	url := commandData.ConnnectionData.LocatorAddress + "/management" + restEndPoint.URL
	urlResponse, err := requests.ExecuteCommand(url, strings.ToUpper(restEndPoint.HTTPMethod), commandData)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jqFilter := commandData.UserCommand.Parameters["-t"]
	jsonToBePrinted, err := format.GetJSONFromURLResponse(urlResponse, string(jqFilter))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(jsonToBePrinted)
}

func checkRequiredParam(restEndPoint domain.RestEndPoint, command domain.UserCommand) error {
	for _, s := range restEndPoint.Parameters {
		if s.Required {
			value := command.Parameters["-"+s.Name]
			if value == "" {
				return errors.New("Required Parameter is missing: " + s.Name)
			}
		}
	}
	return nil
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// Describe an end point with command name and required/optional parameters
func Describe(endPoint domain.RestEndPoint) string {
	var buffer bytes.Buffer
	buffer.WriteString(endPoint.CommandName + " ")
	// show the required options first
	for _, param := range endPoint.Parameters {
		if param.Required {
			buffer.WriteString(getOption(param))
		}
	}

	for _, param := range endPoint.Parameters {
		if !param.Required {
			buffer.WriteString("[" + strings.Trim(getOption(param), " ") + "] ")
		}
	}
	return buffer.String()
}

func getOption(param domain.RestAPIParam) string {
	if param.In == "body" {
		return "-body  "
	}
	return "-" + param.Name + " "
}
