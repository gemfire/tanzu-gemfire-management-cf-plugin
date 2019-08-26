package pcc

import (
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/format"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/requests"
)

// BasicPlugin declares the dataset that commands work on
type BasicPlugin struct {
	commandData domain.CommandData
}

// Run is the main entry point for the CF plugin interface
// It is run once for each CF plugin command executed
func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient := &cfservice.Cf{}
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	var err error
	c.commandData.Target, c.commandData.UserCommand, err = requests.GetTargetAndClusterCommand(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// first get credentials from environment
	c.commandData.Username = os.Getenv("CFLOGIN")                                                                                          // not needed for cf cli
	c.commandData.Password = os.Getenv("CFPASSWORD")                                                                                       // not needed for cf cli
	c.commandData.ExplicitTarget = strings.Contains(c.commandData.Target, "http://") || strings.Contains(c.commandData.Target, "https://") // not needed cf cli
	if c.commandData.ExplicitTarget {
		c.commandData.LocatorAddress = c.commandData.Target
	} else {
		// at this point, we have a valid clusterCommand
		c.commandData.ServiceKey, err = requests.GetServiceKeyFromPCCInstance(cfClient, c.commandData.Target)
		if err != nil {
			fmt.Printf(err.Error(), c.commandData.Target)
			os.Exit(1)
		}

		serviceKeyUser, serviceKeyPswd, url, err := requests.GetUsernamePasswordEndpoinFromServiceKey(cfClient, c.commandData.Target, c.commandData.ServiceKey)
		if err != nil {
			fmt.Println(util.GenericErrorMessage, err.Error())
			os.Exit(1)
		}

		c.commandData.LocatorAddress = url

		// then get the credentials from the serviceKey
		if serviceKeyUser != "" && serviceKeyPswd != "" {
			c.commandData.Username = serviceKeyUser
			c.commandData.Password = serviceKeyPswd
		}

		// Code below should be all we need once we rejig things downstream
		pluginConnection := pluginConnection{cliConnection: cliConnection}
		c.commandData.ConnnectionData, err = pluginConnection.GetConnectionData(c.commandData.Target)
		if err != nil {
			fmt.Println(util.GenericErrorMessage, err.Error())
			os.Exit(1)
		}
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	c.commandData.UserCommand.Parameters = make(map[string]string)

	// lastly get the credentials from the command line
	err = util.ParseArguments(args, &c.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// if pcc instance is specified, username/password are required
	if !c.commandData.ExplicitTarget {
		if c.commandData.Username == "" && c.commandData.Password == "" {
			fmt.Printf(util.NeedToProvideUsernamePassWordMessage, c.commandData.Target, c.commandData.UserCommand.Command)
			os.Exit(1)
		} else if c.commandData.Username != "" && c.commandData.Password == "" {
			err = errors.New(util.ProvidedUsernameAndNotPassword)
			fmt.Printf(err.Error(), c.commandData.Target, c.commandData.UserCommand.Command, c.commandData.Username)
			os.Exit(1)
		} else if c.commandData.Username == "" && c.commandData.Password != "" {
			err = errors.New(util.ProvidedPasswordAndNotUsername)
			fmt.Printf(err.Error(), c.commandData.Target, c.commandData.UserCommand.Command, c.commandData.Password)
			os.Exit(1)
		}
	}

	err = requests.GetEndPoints(&c.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if c.commandData.UserCommand.Command == "commands" {
		for _, command := range c.commandData.AvailableEndpoints {
			fmt.Println(command.CommandCall)
		}
		os.Exit(0)
	}

	err = requests.MapUserInputToAvailableEndpoint(&c.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = requests.HasIDifNeeded(&c.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = requests.HasRegionIfNeeded(&c.commandData)
	if err != nil {
		fmt.Printf(err.Error(), c.commandData.Target)
		os.Exit(1)
	}

	urlResponse, err := requests.RequestToEndPoint(&c.commandData)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jsonToBePrinted, err := format.GetJSONFromURLResponse(urlResponse)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Println(jsonToBePrinted)

	return
}

// GetMetadata provides metadata about the CF plugin including a helptext for the user
func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "pcc",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "pcc",
				HelpText: "Commands to interact with geode cluster.\n",
				UsageDetails: plugin.Usage{
					Usage: "	cf  pcc  <*target>  <command>  [*options]  (* = optional)\n" +
						"\nSupported commands:	use 'cf pcc <*target> commands' to see a list of supported commands \n" +
						"\nNote: target is either a pcc_instance or an explicit locator url in the form of: http(s)://host:port" +
						"\nIt can be saved at [$CFPCC], then omit <*target> from command ",
					Options: map[string]string{
						"h":  "this help screen\n",
						"u":  "followed by equals username (-u=<your_username>) [$CFLOGIN]\n",
						"p":  "followed by equals password (-p=<your_password>) [$CFPASSWORD]\n",
						"r":  "followed by equals region (-r=<your_region>)\n",
						"id": "followed by an identifier required for any get command\n",
						"d": "followed by @<json_file_path> OR single quoted JSON input \n" +
							"	     JSON required for creating/post commands\n",
					},
				},
			},
		},
	}
}
