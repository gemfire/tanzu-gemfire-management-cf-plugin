package geode

import (
	"fmt"
	"os"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/requests"
)

type geodeCommand struct {
	commandData domain.CommandData
}

// NewGeodeCommand provides a constructor for the Geode standalone implementation for the client
func NewGeodeCommand() (geodeCommand, error) {
	return geodeCommand{}, nil
}

// Run is the main entry point for the standalone Geode command line interface
// It is run once for each command executed
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

	// From this point common code can handle the processing of the command
	common.ProcessCommand(&gc.commandData, args)

	return
}
