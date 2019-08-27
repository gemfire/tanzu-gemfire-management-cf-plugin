package geode

import (
	"fmt"
	"os"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/format"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/requests"
)

type geodeCommand struct {
	commandData domain.CommandData
}

// NewGeodeCommand provides a constructor for the Geode standalone implementation for the client
func NewGeodeCommand() (geodeCommand, error) {
	return geodeCommand{}, nil
}

func (gc *geodeCommand) Run(args []string) {
	var err error
	gc.commandData.Target, gc.commandData.UserCommand, err = requests.GetTargetAndClusterCommand(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	geodeConnection, err := NewGeodeConnectionProvider()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	gc.commandData.ConnnectionData, err = geodeConnection.GetConnectionData(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = util.ParseArguments(args, &gc.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = requests.GetEndPoints(&gc.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if gc.commandData.UserCommand.Command == "commands" {
		for _, command := range gc.commandData.AvailableEndpoints {
			fmt.Println(command.CommandCall)
		}
		os.Exit(0)
	}

	err = requests.MapUserInputToAvailableEndpoint(&gc.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = requests.HasIDifNeeded(&gc.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = requests.HasRegionIfNeeded(&gc.commandData)
	if err != nil {
		fmt.Printf(err.Error(), gc.commandData.Target)
		os.Exit(1)
	}

	urlResponse, err := requests.RequestToEndPoint(&gc.commandData)
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

	return
}
