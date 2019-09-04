package common

import (
	"fmt"
	"os"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/format"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/requests"
)

// ProcessCommand handles the common steps for executing a command against the Geode cluster
func ProcessCommand(commandData *domain.CommandData, args []string) {
	var err error

	commandData.UserCommand.Parameters = make(map[string]string)

	err = util.ParseArguments(args, commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = requests.GetEndPoints(commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if commandData.UserCommand.Command == "commands" {
		for _, command := range commandData.AvailableEndpoints {
			fmt.Println(command.CommandCall)
		}
		os.Exit(0)
	}

	err = requests.MapUserInputToAvailableEndpoint(commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = requests.HasIDifNeeded(commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = requests.HasRegionIfNeeded(commandData)
	if err != nil {
		fmt.Printf(err.Error(), commandData.Target)
		os.Exit(1)
	}

	urlResponse, err := requests.RequestToEndPoint(commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jsonToBePrinted, err := format.GetJSONFromURLResponse(urlResponse)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(jsonToBePrinted)

}
