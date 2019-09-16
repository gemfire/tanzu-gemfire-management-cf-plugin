package pcc

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
)

// BasicPlugin declares the dataset that commands work on
type BasicPlugin struct {
	commandData domain.CommandData
	comm        common.CommandProcessor
}

// NewBasicPlugin provides the constructor for a BasicPlugin struct
func NewBasicPlugin(comm common.CommandProcessor) (BasicPlugin, error) {
	return BasicPlugin{comm: comm}, nil
}

// Run is the main entry point for the CF plugin interface
// It is run once for each CF plugin command executed
func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	var err error
	c.commandData.Target, c.commandData.UserCommand = common.GetTargetAndClusterCommand(args)
	if c.commandData.UserCommand.Command == "" {
		fmt.Println(common.GenericErrorMessage, "missing command")
		os.Exit(1)
	}

	pluginConnection, err := NewPluginConnectionProvider(cliConnection)
	if err != nil {
		fmt.Printf(common.GenericErrorMessage, err.Error())
		os.Exit(1)
	}
	err = pluginConnection.GetConnectionData(&c.commandData)
	if err != nil {
		fmt.Printf(common.GenericErrorMessage, err.Error())
		os.Exit(1)
	}

	// From this point common code can handle the processing of the command
	err = c.comm.ProcessCommand(&c.commandData)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

// GetMetadata provides metadata about the CF plugin including a helptext for the user
func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "pcc",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "pcc",
				HelpText: "Commands to interact with Geode cluster.\n",
				UsageDetails: plugin.Usage{
					Usage: "cf  pcc  [target]  <command>  [options] \n\n" +
						"\ttarget:\n\t\ta pcc_instance. \n" +
						"\t\tomit if 'GEODE_TARGET' environment variable is set \n" +
						"\tcommand:\n\t\tuse 'cf pcc <target> commands' to see a list of supported commands \n" +
						common.GeneralOptions + "\n" +
						"\thelp: use -h or --help for general help, and provide <command> -help for command specific help",
				},
			},
		},
	}
}
