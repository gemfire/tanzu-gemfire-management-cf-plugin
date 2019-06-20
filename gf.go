package main

import (
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
)

type BasicPlugin struct{}

type ServiceKeyUsers struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type ServiceKeyUrls struct {
	Gfsh string `json:"gfsh"`
}

type ServiceKey struct {
	Urls  ServiceKeyUrls    `json:"urls"`
	Users []ServiceKeyUsers `json:"users"`
}

type ClusterManagementResult struct {
	StatusCode string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
	MemberStatus []MemberStatus `json:"memberStatus"`
	Result []map[string]interface{} `json:"result"`
}

type MemberStatus struct {
	ServerName string
	Success bool
	Message string
}

const MissingInformationMessage string = `Your request was denied.
You are missing a username, password, or the correct endpoint.

For help see: cf gf --help
`
const IncorrectUserInputMessage string = `Your request was denied.
The format of your request is incorrect.

For help see: cf gf --help
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

	cf gf %s %s -u=%s -p=<your_password>

For help see: cf gf --help
`
const ProvidedPasswordAndNotUsername string = `You did not specify your username.
Please enter username and password:

	cf gf %s %s -u=<your_username> -p=%s

For help see: cf gf --help
`
const NoRegionGivenMessage string = `You need to provide a region to list your indexes from.
The proper format is:

	cf gf list-indexes %s -r=<your_region>

To see your available regions:
	
	cf gf list-regions %s

For help see: cf gf --help
`
const NotAuthenticatedMessage string = `The username and password is incorrect.

For help see: cf gf --help
`
const NonExistentRegionMessage string = `The region you selected does not exist.
To see your active regions, enter:
	
	cf gf list-regions %s

For help see: cf gf --help
`
const NeedToProvideUsernamePassWordMessage string = `You need to provide your username and password.
The proper format is: cf gf %s %s -u=<your_username> -p=<your_password>

For help see: cf gf --help
`
const UnsupportedClusterCommandMessage string = `You entered an unsupported cluster command.

For help see: cf gf --help`

const Ellipsis string = "…"

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
	endpoint = serviceKey.Urls.Gfsh
	endpoint = strings.Replace(endpoint, "gemfire/v1", "geode-management/v2", 1)
	for _ , user := range serviceKey.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			username = user.Username
			password = user.Password
		}
	}
	return
}

func getCompleteEndpoint(endpoint string, clusterCommand string, region string) (string, error){
	urlEnding := ""
	switch clusterCommand{
	case "list-regions":
		urlEnding = "/regions"
	case "list-members":
		urlEnding = "/members"
	case "list-gateway-receivers":
		urlEnding = "/gateways/receivers"
	case "list-indexes":
		if region == ""{
			return "", errors.New(NoRegionGivenMessage)
		}
		if strings.HasPrefix(region, "/"){
			region = region[1:]
		}
		urlEnding = "/regions/" + region + "/indexes"
	default:
		return endpoint, nil
	}
	endpoint = endpoint + urlEnding
	return endpoint, nil
}

func getUrlOutput(endpointUrl string, username string, password string) (urlResponse string, err error){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", endpointUrl, nil)
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


func getTableHeadersFromClusterCommand(clusterCommand string) (tableHeaders []string){
	switch clusterCommand {
	case "list-regions":
		tableHeaders = []string{"name", "type", "groups", "entryCount", "regionAttributes"}
	case "list-members":
		tableHeaders = []string{"id", "host", "status", "pid"}
	case "list-gateway-receivers":
		tableHeaders = []string{"hostnameForSenders", "uri", "group", "class"}
	case "list-indexes":
		tableHeaders = []string{"name", "type", "fromClause", "expression"}
	default:
		return
	}
	return
}

func GetAnswerFromUrlResponse(clusterCommand string, urlResponse string) (response string, err error){
	urlOutput := ClusterManagementResult{}
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
	for _, result := range urlOutput.Result{
		memberCount++
		for _, key := range tableHeaders {
			if result[key] == nil {
				response += Fill(20, "", " ") + "|"
			} else {
				resultVal := result[key]
				if fmt.Sprintf("%T", result[key]) == "float64"{
					resultVal = fmt.Sprintf("%.0f", result[key])
				}
				response += Fill(20, fmt.Sprintf("%s",resultVal), " ") + "|"
			}
		}
		response += "\n"
	}

	if clusterCommand == "list-regions"{
		response += "\nNumber of Regions: " + strconv.Itoa(memberCount)
	} else if clusterCommand == "list-members"{
		response += "\nNumber of Members: " + strconv.Itoa(memberCount)
	} else if clusterCommand == "list-indexes" {
		response += "\nNumber of Indices in your Region: " + strconv.Itoa(memberCount)
	}
	if strings.Contains(response, Ellipsis){
		response += "\nTo see the full output, append -j to your command."
	}
	return
}


func GetJsonFromUrlResponse(urlResponse string) (jsonOutput string, err error){
	urlOutput := ClusterManagementResult{}
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

func isRegionInGroups(regionInGroup bool, groupsWeHave []string, groupsWeWant []string) (isInGroups bool){
	if regionInGroup{
		isInGroups = true
		return
	}
	for _, group := range groupsWeWant{
		for  _, regionName := range groupsWeHave {
			if regionName == group{
				isInGroups = true
				return
			}
		}
	}
	isInGroups = false
	return
}

func EditResponseOnGroup(urlResponse string, groups []string, clusterCommand string) (editedUrlResponse string, err error){
	urlOutput := ClusterManagementResult{}
	err = json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		return "", err
	}
	var newUrlOutputResult string
	var newResult []map[string]interface{}
	for _, result := range urlOutput.Result{
		regionInGroups:= false
		for _, key :=range getTableHeadersFromClusterCommand(clusterCommand){
			if key == "groups"{
				regionInGroups = isRegionInGroups(regionInGroups, toSlice(result[key]), groups)
			}
			if regionInGroups{
				break
			}
		}
		if regionInGroups{
			newUrlOutputResult += fmt.Sprintf("%s",result)
			newResult = append(newResult, result)
		}
	}
	urlOutput.Result = newResult
	byteResponse, err := json.Marshal(urlOutput)
	if err != nil {
		return "", err
	}
	editedUrlResponse = string(byteResponse)
	return
}

// Convert an implicit slice of strings, represented by interface{}, into an actual []string
func toSlice(input interface{}) []string {
	result := make([]string, 0)
	if input != nil {
		for _, entry := range input.([]interface{}) {
			result = append(result, entry.(string))
		}
	}
	return result
}

func isSupportedClusterCommand(clusterCommandFromUser string) (error){
	clusterCommandsWeSupport := []string{"list-members", "list-regions", "list-gateway-receivers", "list-indexes"}
	for _,command := range clusterCommandsWeSupport{
		if clusterCommandFromUser == command{
			return nil
		}
	}
	return errors.New(UnsupportedClusterCommandMessage)
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	start := time.Now()
	if args[0] == "CLI-MESSAGE-UNINSTALL"{
		return
	}
	var err error
	var username, password, endpoint, pccInUse, clusterCommand, serviceKey, region string
	var groups []string
	if len(args) >= 3 {
		pccInUse = args[2]
		clusterCommand = args[1]
	} else{
		fmt.Println(IncorrectUserInputMessage)
		return
	}
	err = isSupportedClusterCommand(clusterCommand)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if os.Getenv("CFLOGIN") != "" && os.Getenv("CFPASSWORD") != "" && os.Getenv("CFENDPOINT") != "" {
		username = os.Getenv("CFLOGIN")
		password = os.Getenv("CFPASSWORD")
		endpoint = os.Getenv("CFENDPOINT")
	} else {
		var err error
		cfClient := &cfservice.Cf{}
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
	hasJ, hasG, hasU, hasP := false, false, false, false
	for _, arg := range args {
		if strings.HasPrefix(arg, "-g="){
			hasG = true
			groups = strings.Split(arg[3:], ",")
			if err != nil{
				fmt.Println(err.Error())
				os.Exit(1)
			}
		} else if arg == "-j"{
			hasJ = true
		} else if strings.HasPrefix(arg, "-r="){
			region = arg[3:]
		} else if strings.HasPrefix(arg, "-u="){
			hasU = true
			username = arg[3:]
		} else if strings.HasPrefix(arg, "-p="){
			password = arg[3:]
			hasP = true
		}
	}
	if (hasU && !hasP) {
		err = errors.New(ProvidedUsernameAndNotPassword)
		fmt.Printf(err.Error(), clusterCommand, pccInUse, username)
		os.Exit(1)
	} else if(!hasU && hasP){
		err = errors.New(ProvidedPasswordAndNotUsername)
		fmt.Printf(err.Error(), clusterCommand, pccInUse, password)
		os.Exit(1)
	}
	endpoint, err = getCompleteEndpoint(endpoint, clusterCommand, region)
	if err != nil{
		fmt.Printf(err.Error(), pccInUse, pccInUse)
		os.Exit(1)
	}
	urlResponse, err := getUrlOutput(endpoint, username, password)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if hasG{
		urlResponse, err = EditResponseOnGroup(urlResponse, groups, clusterCommand)
	}

	if hasJ{
		jsonToBePrinted, err := GetJsonFromUrlResponse(urlResponse)
		if err != nil{
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(jsonToBePrinted)
		return
	}
	if username != "" && password != "" && clusterCommand != "" && endpoint != "" {
		answer, err := GetAnswerFromUrlResponse(clusterCommand, urlResponse)
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
	} else if username == "" && password == "" &&clusterCommand != "" && endpoint != ""{
		fmt.Println(NeedToProvideUsernamePassWordMessage, clusterCommand, pccInUse)
	} else {
		fmt.Println(MissingInformationMessage)
	}
	t := time.Now()
	fmt.Println(t.Sub(start))
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "GF_InDev",
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
				Name:     "gf",
				HelpText: "gf's help text",
				UsageDetails: plugin.Usage{
					Usage: "	cf  gf  <action>  <pcc_instance>  [*options]  (* = optional)\n" +
						"	Actions: \n" +
						"		list-regions, list-members, list-gateway-receivers, list-indexes\n" +
						"	Options: \n" +
						"		-h : this help screen\n" +
						"		-j : json output of API endpoint\n" +
						"		-g : followed by equals group(s), split by comma, only data within those groups\n" +
						"			(example: cf gf list-regions -g=group1,group2)\n" +
						"		-u : followed by equals username (-u=<your_username>)\n" +
						"		-p : followed by equals password (-p=<your_password>)",
				},
			},
		},
	}
}


func main() {
	plugin.Start(new(BasicPlugin))
}
