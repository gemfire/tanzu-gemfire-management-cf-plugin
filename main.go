package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"os"
	"strings"
)


var username, password, endpoint, pccInUse, clusterCommand, serviceKey, region, regionJSONfile, group, ca_cert, httpAction string
var hasGroup, isJSONOutput = false, false

var parameters map[string]string
var APICallStruct RestAPICall


func main() {
	plugin.Start(new(BasicPlugin))
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient := &cfservice.Cf{}
	if args[0] == "CLI-MESSAGE-UNINSTALL"{
		return
	}
	var err error
	err = getPCCInUseAndClusterCommand(args)
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
	APICallStruct.parameters = make(map[string]string)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-g="){
			APICallStruct.parameters["g"] = "true"
			hasGroup = true
			group = arg[3:]
			if err != nil{
				fmt.Println(err.Error())
				os.Exit(1)
			}
		} else if arg == "-j"{
			APICallStruct.parameters["j"] ="true"
			isJSONOutput = true
		} else if strings.HasPrefix(arg, "-r="){
			APICallStruct.parameters["r"] = "true"
			region = arg[3:]
		} else if strings.HasPrefix(arg, "-u="){
			username = arg[3:]
		} else if strings.HasPrefix(arg, "-p="){
			password = arg[3:]
		} else if strings.HasPrefix(arg, "-cacert="){
			ca_cert=arg[8:]
			APICallStruct.parameters["cacert"] = ca_cert
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

	firstEndpoint := "http://localhost:7070/management/v2/cli" +"?command="+APICallStruct.command
	urlResponse, err := getUrlOutput(firstEndpoint, username, password, "GET")

	response := ResponseFromAPI{}
	err = json.Unmarshal([]byte(urlResponse), &response)
	if err != nil {
		fmt.Println(IncorrectUserInputMessage)
		os.Exit(1)
	}
	secondEndpoint := "http://localhost:7070/management/v2/" + response.Url
	urlResponse, err = getUrlOutput(secondEndpoint, username, password, response.HttpMethod)
	if err != nil {
		fmt.Println("her222222")
		fmt.Println(IncorrectUserInputMessage)
		os.Exit(1)
	}

	jsonToBePrinted, err := GetJsonFromUrlResponse(urlResponse)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(jsonToBePrinted)
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
						"	Actions: list, create, get, delete\n\n" +
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
