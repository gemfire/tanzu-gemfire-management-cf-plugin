package common

import (
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/input"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/output"
)

type Common struct {
	helper impl.RequestHelper
}

func NewCommon(helper impl.RequestHelper) (Common, error) {
	return Common{helper: helper}, nil
}

// ProcessCommand handles the common steps for executing a command against the Geode cluster
func (c *Common) ProcessCommand(commandData *domain.CommandData) {
	var err error

	err = c.helper.GetEndPoints(commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	userCommand := commandData.UserCommand.Command
	if userCommand == "commands" {
		for _, command := range commandData.AvailableEndpoints {
			fmt.Println(output.Describe(command))
		}
		os.Exit(0)
	}

	restEndPoint, avalable := commandData.AvailableEndpoints[userCommand]
	if !avalable {
		fmt.Println("Invalid command: " + userCommand)
		os.Exit(1)
	}

	if input.HasOption(commandData, "-h") || input.HasOption(commandData, "--help") {
		for _, command := range commandData.AvailableEndpoints {
			if command.CommandName == userCommand {
				fmt.Println(output.Describe(command))
			}
		}
		os.Exit(0)
	}

	err = checkRequiredParam(restEndPoint, commandData.UserCommand)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	url := commandData.ConnnectionData.LocatorAddress + "/management" + restEndPoint.URL
	urlResponse, err := c.helper.ExecuteCommand(url, strings.ToUpper(restEndPoint.HTTPMethod), commandData)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jqFilter := commandData.UserCommand.Parameters["-t"]
	jsonToBePrinted, err := output.GetJSONFromURLResponse(urlResponse, jqFilter)
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
