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
	comm        common.Common
}

// NewGeodeCommand provides a constructor for the Geode standalone implementation for the client
func NewGeodeCommand(comm common.Common) (geodeCommand, error) {
	return geodeCommand{comm: comm}, nil
}

// Run is the main entry point for the standalone Geode command line interface
// It is run once for each command executed
func (gc *geodeCommand) Run(args []string) {
	var err error
	gc.commandData.Target, gc.commandData.UserCommand = requests.GetTargetAndClusterCommand(args)

	// if no user command and args contains -h or --help
	if gc.commandData.UserCommand.Command == "" {
		if hasOption(gc, "-h") || hasOption(gc, "--help") {
			printHelp()
			os.Exit(0)
		} else {
			fmt.Println("Invalid command")
			os.Exit(1)
		}
	}

	geodeConnection, err := NewGeodeConnectionProvider()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = geodeConnection.GetConnectionData(&gc.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// From this point common code can handle the processing of the command
	gc.comm.ProcessCommand(&gc.commandData)

	return
}

func hasOption(gc *geodeCommand, option string) bool {
	return gc.commandData.UserCommand.Parameters[option] != "" || gc.commandData.Target == option
}

func printHelp() {
	fmt.Println("Commands to interact with geode cluster.")
	fmt.Println("")
	fmt.Println("Usage: pcc <target> <command> [options]")
	fmt.Println("")
	fmt.Println("\ttarget: url to a geode locator in the form of : http(s)://host:port")
	fmt.Println("\tcommand: use 'pcc <target> commands' to see a list of supported commands")
	fmt.Println("\toptions: see help for individual commands for options.")
	fmt.Println("\thelp: use -h or --help for general help, and provide <command> for command specific help.")
}
