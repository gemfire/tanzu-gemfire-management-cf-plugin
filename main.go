package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"os"
	"strings"
	"time"
)


var username, password, endpoint, pccInUse, clusterCommand, serviceKey, region, group, regionJSONfile, ca_cert, httpAction string
var hasGroup, isJSONOutput = false, false



func main() {
	plugin.Start(new(BasicPlugin))
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient := &cfservice.Cf{}
	start := time.Now()
	if args[0] == "CLI-MESSAGE-UNINSTALL"{
		return
	}
	var err error
	err = getPCCInUseAndClusterCommand(args)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = isSupportedClusterCommand(clusterCommand)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// at this point, we have a valid clusterCommand
	serviceKey, err =  GetServiceKeyFromPCCInstance(cfClient, pccInUse)
	if err != nil{
		fmt.Printf(err.Error(), pccInUse, pccInUse)
		os.Exit(1)
	}
	if os.Getenv("CFLOGIN") != "" && os.Getenv("CFPASSWORD") != ""{
		username = os.Getenv("CFLOGIN")
		password = os.Getenv("CFPASSWORD")
		_, _, endpoint, err = GetUsernamePasswordEndpoint(cfClient, pccInUse, serviceKey)
	} else {
		username, password, endpoint, err = GetUsernamePasswordEndpoint(cfClient, pccInUse, serviceKey)
		if err != nil{
			fmt.Println(GenericErrorMessage, err.Error())
			os.Exit(1)
		}
	}

	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-g="){
			hasGroup = true
			group = arg[3:]
			if err != nil{
				fmt.Println(err.Error())
				os.Exit(1)
			}
		} else if arg == "-j"{
			isJSONOutput = true
		} else if strings.HasPrefix(arg, "-r="){
			region = arg[3:]
		} else if strings.HasPrefix(arg, "-u="){
			username = arg[3:]
		} else if strings.HasPrefix(arg, "-p="){
			password = arg[3:]
		} else if strings.HasPrefix(arg, "-j="){
			regionJSONfile = arg[3:]
		} else if strings.HasPrefix(arg, "-cacert="){
			ca_cert=arg[8:]
		}
	}

	if username == "" && password == "" {
		fmt.Println(NeedToProvideUsernamePassWordMessage, pccInUse, clusterCommand)
		os.Exit(1)
	} else if username != "" && password == "" {
		err = errors.New(ProvidedUsernameAndNotPassword)
		fmt.Printf(err.Error(), pccInUse, clusterCommand, username)
		os.Exit(1)
	} else if username=="" && password!="" {
		err = errors.New(ProvidedPasswordAndNotUsername)
		fmt.Printf(err.Error(), pccInUse, clusterCommand, password)
		os.Exit(1)
	}

	// at this point, we should have non-empty username and password
	endpoint, err = getCompleteEndpoint(endpoint, clusterCommand)
	if err != nil{
		fmt.Printf(err.Error(), pccInUse, pccInUse)
		os.Exit(1)
	}

	//preform post commands
	if strings.HasPrefix(clusterCommand, "post"){
		httpAction = "POST"
	} else {
		httpAction = "GET"

	}
	urlResponse, err := getUrlOutput(endpoint, username, password, httpAction)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if !isJSONOutput{
		answer, err := GetTableFromUrlResponse(clusterCommand, urlResponse)

		if err != nil{
			if err.Error() == NotAuthenticatedMessage{
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			fmt.Printf(err.Error(), pccInUse)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println(answer)
		fmt.Println()
		t := time.Now()
		fmt.Println(t.Sub(start))
	} else {
		jsonToBePrinted, err := GetJsonFromUrlResponse(urlResponse)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(jsonToBePrinted)
	}
	return
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "PCC_InDev",
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
					Usage: "	cf  pcc  <*pcc_instance>  <action>  <data_type>  [*options]  (* = optional)\n\n" +
						"	Actions: list, post\n\n" +
						"	Data Types: regions, members, gateway-receivers, indexes\n\n" +
						"	Note: pcc_instance can be saved at [$CFPCC], then omit pcc_instance from command ",
					Options: map[string]string{
						"h" : "this help screen\n",
						"u" : "followed by equals username (-u=<your_username>) [$CFLOGIN]\n",
						"p" : "followed by equals password (-p=<your_password>) [$CFPASSWORD]\n",
						"r" : "followed by equals region (-r=<your_region>)\n",
						"j" : "json input for region post\n",
						"cacert" : "ca-certification needed to post region\n",
					},
				},
			},
		},
	}
}
