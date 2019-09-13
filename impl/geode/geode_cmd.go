package geode

import (
	"errors"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
)

type geodeCommand struct {
	commandData domain.CommandData
	comm        common.CommandProcessor
}

// NewGeodeCommand provides a constructor for the Geode standalone implementation for the client
func NewGeodeCommand(comm common.CommandProcessor) (geodeCommand, error) {
	return geodeCommand{comm: comm}, nil
}

// Run is the main entry point for the standalone Geode command line interface
// It is run once for each command executed
func (gc *geodeCommand) Run(args []string) (err error) {

	gc.commandData.Target, gc.commandData.UserCommand = common.GetTargetAndClusterCommand(args)

	// if no user command and args contains -h or --help
	if gc.commandData.UserCommand.Command == "" {
		if common.HasOption(gc.commandData.UserCommand.Parameters, []string{"--help", "-h"}) {
			printHelp()
			return
		} else {
			err = errors.New("Invalid command")
			return
		}
	}

	geodeConnection, err := NewGeodeConnectionProvider()
	if err != nil {
		return
	}

	err = geodeConnection.GetConnectionData(&gc.commandData)
	if err != nil {
		return
	}

	// From this point common code can handle the processing of the command
	err = gc.comm.ProcessCommand(&gc.commandData)

	return
}

func printHelp() {
	fmt.Println("Commands to interact with geode cluster.")
	fmt.Println("")
	fmt.Println("Usage: pcc <target> <command> [options]")
	fmt.Println("")
	fmt.Println("\ttarget: \n\t\turl to a geode locator in the form of : http(s)://host:port")
	fmt.Println("\t\tomit if 'GEODE_TARGET' environment variable is set")
	fmt.Println("\tcommand:\n\t\tuse 'pcc <target> commands' to see a list of supported commands")
	fmt.Println("\toptions:\n\t\tsee help for individual commands for options.")
	fmt.Println(common.GeneralOptions)
	fmt.Println("\thelp:\n\t\tuse -h or --help for general help, and provide <command> for command specific help.")
}
