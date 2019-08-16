package pcc

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"os"
	"strings"
)

var username, password, locatorAddress, target, serviceKey, region, jsonFile, group, id string
var hasGroup, isJSONOutput, explicitTarget = false, false, false

var userCommand UserCommand
var firstResponse SwaggerInfo
var availableEndpoints []IndividualEndpoint
var endPoint IndividualEndpoint

func main() {
	plugin.Start(new(BasicPlugin))
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient := &cfservice.Cf{}
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	var err error
	err = getTargetAndClusterCommand(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// first get credentials from environment
	username = os.Getenv("CFLOGIN")
	password = os.Getenv("CFPASSWORD")
	explicitTarget = strings.Contains(target, "http://") || strings.Contains(target, "https://")
	if explicitTarget {
		locatorAddress = target
	} else {
		// at this point, we have a valid clusterCommand
		serviceKey, err = GetServiceKeyFromPCCInstance(cfClient)
		if err != nil {
			fmt.Printf(err.Error(), target, target)
			os.Exit(1)
		}

		serviceKeyUser, serviceKeyPswd, url, err := GetUsernamePasswordEndpoinFromServiceKey(cfClient)
		if err != nil {
			fmt.Println(GenericErrorMessage, err.Error())
			os.Exit(1)
		}

		locatorAddress = url

		// then get the credentials from the serviceKey
		if serviceKeyUser != "" && serviceKeyPswd != "" {
			username = serviceKeyUser
			password = serviceKeyPswd
		}
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	userCommand.parameters = make(map[string]string)

	// lastly get the credentials from the command line
	err = parseArguments(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// if pcc instance is specified, username/password are required
	if !explicitTarget {
		if username == "" && password == "" {
			fmt.Printf(NeedToProvideUsernamePassWordMessage, target, userCommand.command)
			os.Exit(1)
		} else if username != "" && password == "" {
			err = errors.New(ProvidedUsernameAndNotPassword)
			fmt.Printf(err.Error(), target, userCommand.command, username)
			os.Exit(1)
		} else if username == "" && password != "" {
			err = errors.New(ProvidedPasswordAndNotUsername)
			fmt.Printf(err.Error(), target, userCommand.command, password)
			os.Exit(1)
		}
	}

	err = getEndPoints()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if userCommand.command == "commands" {
		for _, command := range availableEndpoints {
			fmt.Println(command.CommandCall)
		}
		os.Exit(0)
	}

	endPoint, err = mapUserInputToAvailableEndpoint()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = hasIDifNeeded()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = hasRegionIfNeeded()
	if err != nil {
		fmt.Printf(err.Error(), target)
		os.Exit(1)
	}

	urlResponse, err := requestToEndPoint()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jsonToBePrinted, err := GetJsonFromUrlResponse(urlResponse)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Println(jsonToBePrinted)

	return
}

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
