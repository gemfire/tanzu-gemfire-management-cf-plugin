package main

import (
	"bytes"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	//"flag"
)

type BasicPlugin struct{}

type ServiceKeyUsers struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type ServiceKeyUrls struct {
	Management string `json:"management"`
	Gfsh string `json:"gfsh"`
}

type ServiceKey struct {
	Urls  ServiceKeyUrls    `json:"urls"`
	Users []ServiceKeyUsers `json:"users"`
}

type ClusterManagementResults struct {
	StatusCode string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
	MemberStatus []MemberStatus `json:"memberStatus"`
	Results []IndividualClusterManagementResult `json:"result"`
}

type IndividualClusterManagementResult struct {
	Config map[string]interface{} `json:"config"`
	RuntimeInfo []map[string]interface{} `json:"runtimeInfo"`
}

type MemberStatus struct {
	ServerName string
	Success bool
	Message string
}

type PostJson struct {
	name string `json:"name"`
	_type string `json:"type"`
}

const MissingInformationMessage string = `Your request was denied.
You are missing a username, password, or the correct endpoint.

For help see: cf pcc --help
`
const IncorrectUserInputMessage string = `Your request was denied.
The format of your request is incorrect.

For help see: cf pcc --help
`

const NoServiceKeyMessage string = `Please create a service key for %s.
To create a key enter: 

	cf create-service-key %s <your_key_name>
	
For help see: cf create-service-key --help
`
const GenericErrorMessage string = `Cannot retrieve credentials. Error: %s`
const InvalidServiceKeyResponse string = `The cf service-key response is invalid.`
const ProvidedUsernameAndNotPassword string = `You did not specify your password.
Please enter username and password:

	cf pcc %s %s -u=%s -p=<your_password>

For help see: cf pcc --help
`
const ProvidedPasswordAndNotUsername string = `You did not specify your username.
Please enter username and password:

	cf pcc %s %s -u=<your_username> -p=%s

For help see: cf pcc --help
`
const NoRegionGivenMessage string = `You need to provide a region to list your indexes from.
The proper format is:

	cf pcc %s list regions -r=<your_region>

To see your available regions:
	
	cf pcc %s list regions

For help see: cf pcc --help
`
const NotAuthenticatedMessage string = `The username and password is incorrect.

For help see: cf pcc --help
`
const NonExistentRegionMessage string = `The region you selected does not exist.
To see your active regions, enter:
	
	cf pcc %s list regions

For help see: cf pcc --help
`
const NeedToProvideUsernamePassWordMessage string = `You need to provide your username and password.
The proper format is: cf pcc %s %s -u=<your_username> -p=<your_password>

For help see: cf pcc --help
`
const UnsupportedClusterCommandMessage string = `You entered an unsupported cluster command.

For help see: cf pcc --help`

const Ellipsis string = "…"

var username, password, endpoint, pccInUse, clusterCommand, serviceKey, region, group, regionJSONfile, ca_cert, httpAction string
var hasGroup, isJSONOutput = false, false


func GetServiceKeyFromPCCInstance(cf cfservice.CfService, pccService string) (serviceKey string, err error) {
	servKeyOutput, err := cf.Cmd("service-keys", pccService)
	if err != nil{
		return "", err
	}
	splitKeys := strings.Split(servKeyOutput, "\n")
	hasKey := false
	if strings.Contains(splitKeys[1], "No service key for service instance"){
		return "", errors.New(NoServiceKeyMessage)
	}
	for _, value := range splitKeys {
		line := strings.Fields(value)
		if len(line) > 0 {
			if hasKey {
				serviceKey = line[0]
				return
			} else if line[0] == "name" {
				hasKey = true
			}
		}
	}
	if serviceKey == "" {
		return serviceKey, errors.New(NoServiceKeyMessage)
	}
	return
}

func GetUsernamePasswordEndpoint(cf cfservice.CfService, pccService string, key string) (username string, password string, endpoint string, err error) {
	username = ""
	password = ""
	endpoint = ""
	keyInfo, err := cf.Cmd("service-key", pccService, key)
	if err != nil {
		return "", "", "", err
	}
	splitKeyInfo := strings.Split(keyInfo, "\n")
	if len(splitKeyInfo) < 2{
		return "", "", "", errors.New(InvalidServiceKeyResponse)
	}
	splitKeyInfo = splitKeyInfo[2:] //take out first two lines of cf service-key ... output
	joinKeyInfo := strings.Join(splitKeyInfo, "\n")

	serviceKey := ServiceKey{}

	err = json.Unmarshal([]byte(joinKeyInfo), &serviceKey)
	if err != nil {
		return "", "", "", err
	}
	endpoint = serviceKey.Urls.Management
	endpoint = strings.TrimSuffix(serviceKey.Urls.Gfsh, "gemfire/v1") + "management/v2"
	for _ , user := range serviceKey.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			username = user.Username
			password = user.Password
		}
	}
	return
}

func getCompleteEndpoint(endpoint string, clusterCommand string) (string, error){
	urlEnding := ""
	switch clusterCommand{
	case "list regions":
		urlEnding = "/regions"
	case "list members":
		urlEnding = "/members"
	case "list gateway-receivers":
		urlEnding = "/gateways/receivers"
	case "list indexes":
		if region == ""{
			return "", errors.New(NoRegionGivenMessage)
		}
		if strings.HasPrefix(region, "/"){
			region = region[1:]
		}
		urlEnding = "/regions/" + region + "/indexes"
	case "post region":
		urlEnding = "/regions"
	default:
		return endpoint, nil
	}
	endpoint = endpoint + urlEnding + "?group=" + group
	return endpoint, nil
}

func getTableHeadersFromClusterCommand(clusterCommand string) (tableHeaders []string){
	switch clusterCommand {
	case "list regions":
		tableHeaders = []string{"name", "type", "groups", "entryCount", "regionAttributes"}
	case "list members":
		tableHeaders = []string{"id", "host", "status", "pid"}
	case "list gateway-receivers":
		tableHeaders = []string{"hostnameForSenders", "uri", "group", "class"}
	case "list indexes":
		tableHeaders = []string{"name", "type", "fromClause", "expression"}
	default:
		return
	}
	return
}

func getUrlOutput(endpointUrl string, username string, password string, httpAction string) (urlResponse string, err error){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//urlOutput := ClusterManagementResult{}
	//err = json.Unmarshal([]byte(urlResponse), &urlOutput)
	postJsonResults := PostJson{}
	if httpAction == "POST"{
		err = json.Unmarshal([]byte(regionJSONfile), &postJsonResults)
	}
	requestBody, err := json.Marshal(postJsonResults)
	if err != nil {
		return "",err
	}
	req, err := http.NewRequest(httpAction, endpointUrl, bytes.NewBuffer(requestBody))
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil{
		return "", err
	}
	respInAscii, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil{
		return "", err
	}
	urlResponse = fmt.Sprintf("%s", respInAscii)
	return
}

func Fill(columnSize int, value string, filler string) (response string){
	if len(value) > columnSize - 1{
		response = " " + value[:columnSize-len([]rune(Ellipsis)) -1] + Ellipsis
		return
	}
	numFillerChars := columnSize - len(value) - 1
	response = " " + value + strings.Repeat(filler, numFillerChars)
	return
}


func GetTableFromUrlResponse(clusterCommand string, urlResponse string) (response string, err error){
	urlOutput := ClusterManagementResults{}
	err = json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		return "", err
	}
	if urlOutput.StatusCode == "UNAUTHENTICATED"{
		return "", errors.New(NotAuthenticatedMessage)
	} else if urlOutput.StatusCode == "ENTITY_NOT_FOUND"{
		return "", errors.New(NonExistentRegionMessage)
	}
	response = "Status Code: " + urlOutput.StatusCode + "\n"
	if urlOutput.StatusMessage != ""{
		response += "Status Message: " + urlOutput.StatusMessage + "\n"
	}
	response += "\n"

	tableHeaders := getTableHeadersFromClusterCommand(clusterCommand)
	for _, header := range tableHeaders {
		response += Fill(20, header, " ") + "|"
	}
	response += "\n" + Fill (20 * len(tableHeaders) + 5, "", "-") + "\n"

	memberCount := 0
	for _, result := range urlOutput.Results{
		memberCount++
		if err != nil {
			return "", err
		}
		for _, key := range tableHeaders {
			if result.RuntimeInfo[0][key] == nil && result.Config[key] == nil {
				response += Fill(20, "", " ") + "|"
			} else {
				resultVal := result.Config[key]
				if resultVal == nil{
					resultVal = result.RuntimeInfo[0][key]
				}
				if fmt.Sprintf("%T", resultVal) == "float64"{
					resultVal = fmt.Sprintf("%.0f", resultVal)
				}
				response += Fill(20, fmt.Sprintf("%s",resultVal), " ") + "|"
			}
		}
		response += "\n"
	}

	response += "\nNumber of Results: " + strconv.Itoa(memberCount)
	if strings.Contains(response, Ellipsis){
		response += "\nTo see the full output, append -j to your command."
	}
	return
}


func GetJsonFromUrlResponse(urlResponse string) (jsonOutput string, err error){
	urlOutput := ClusterManagementResults{}
	err = json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		return "", err
	}
	jsonExtracted, err := json.MarshalIndent(urlOutput, "", "  ")
	if err != nil {
		return "", err
	}
	jsonOutput = string(jsonExtracted)
	return
}


func isSupportedClusterCommand(clusterCommandFromUser string) (error){
	clusterCommandsWeSupport := []string{"list members", "list regions", "list gateway-receivers", "list indexes", "post region"}
	for _,command := range clusterCommandsWeSupport{
		if clusterCommandFromUser == command{
			return nil
		}
	}
	return errors.New(UnsupportedClusterCommandMessage)
}


func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient := &cfservice.Cf{}
	start := time.Now()
	if args[0] == "CLI-MESSAGE-UNINSTALL"{
		return
	}
	var err error
	//var username, password, endpoint, pccInUse, clusterCommand, serviceKey, region string
	if len(args) >= 4 {
		pccInUse = args[1]
		clusterCommand = args[2] + " " + args[3]
	} else{
		fmt.Println(IncorrectUserInputMessage)
		return
	}

	err = isSupportedClusterCommand(clusterCommand)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// at this point, we have a valid clusterCommand

	if os.Getenv("CFLOGIN") != "" && os.Getenv("CFPASSWORD") != ""{
		username = os.Getenv("CFLOGIN")
		password = os.Getenv("CFPASSWORD")
	} else {
		var err error
		serviceKey, err =  GetServiceKeyFromPCCInstance(cfClient, pccInUse)
		if err != nil{
			fmt.Printf(err.Error(), pccInUse, pccInUse)
			os.Exit(1)
		}
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
		fmt.Println(endpoint)
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
					Usage: "	cf  pcc  <pcc_instance>  <action>  <data_type>  [*options]  (* = optional)\n\n" +
						"	Actions: list, post" +
						"	Data Types: regions, members, gateway-receivers, indexes\n",
					Options: map[string]string{
						"h" : "this help screen\n",
						"g" : "followed by equals group(s), split by comma, only data within those groups\n"+
						"		(example: cf pcc <instance> list-regions -g=group1)\n",
						"u" : "followed by equals username (-u=<your_username>) [$CFLOGIN]\n",
						"p" : "followed by equals password (-p=<your_password>) [$CFPASSWORD]\n",
						"r" : "followed by equals region (-r=<your_region>)\n",
						"j" : "json input for region post",
						"cacert" : "ca-certification needed to post region",
					},
				},
			},
		},
	}
}


func main() {
	plugin.Start(new(BasicPlugin))
}
