package main

import (
	"bytes"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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


const incorrectUserInputMessage string = `Your request was denied.
You are missing a username, password, or the correct endpoint.`
const invalidPCCInstanceMessage string = `The PCC instance you provided is not a deployed PCC instance.
To deploy this instance, run: cf create-service p-cloudcache your_instance_name`
const noServiceKeyMessage string = `Please create a service key for %s.
To create a key enter: cf create-service-key %s your_key_name
`

func collectCloudCacheServices() (cloudCachesAvailable []string) {
	cmd := exec.Command("cf", "services")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	servicesTable := &out
	tableStr := servicesTable.String()
	splitTable := strings.Split(tableStr, "\n")
	for _, value := range splitTable {
		line := strings.Fields(value)
		if len(line) > 0 && line[1] == "p-cloudcache" {
			cloudCachesAvailable = append(cloudCachesAvailable, line[0])
		}
	}
	return
}

func getServiceKeyFromPCCInstance(pccService string) (serviceKey string, err error) {
	cmd := exec.Command("cf", "service-keys", pccService)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	servKeyOutput := &out
	keysStr := servKeyOutput.String()
	splitKeys := strings.Split(keysStr, "\n")
	hasKey := false
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
		return serviceKey, errors.New(noServiceKeyMessage)
	}
	return
}

func getUsernamePasswordEndpoint(pccService string, key string) (username string, password string, endpoint string) {
	username = ""
	password = ""
	endpoint = ""
	cmd := exec.Command("cf", "service-key", pccService, key)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	servKeyOutput := &out
	keyInfo := servKeyOutput.String()
	splitKeyInfo := strings.Split(keyInfo, "\n")
	splitKeyInfo = splitKeyInfo[2:] //take out first two lines of cf service-key ... output
	joinKeyInfo := strings.Join(splitKeyInfo, "\n")

	serviceKey := ServiceKey{}

	err = json.Unmarshal([]byte(joinKeyInfo), &serviceKey)
	if err != nil {
		log.Fatal(err)
	}
	endpoint = serviceKey.Urls.Gfsh
	for _ , user := range serviceKey.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			username = user.Username
			password = user.Password
		}
	}
	return
}

func ValidatePCCInstance(ourPCCInstance string, pccInstancesAvailable []string) (error){
	for _, pccInst := range pccInstancesAvailable {
		if ourPCCInstance == pccInst {
			return nil
		}
	}
	return errors.New(invalidPCCInstanceMessage)
}

func getEndpoint(clusterCommand string) (endpoint string){
	urlEnding := ""
	if clusterCommand == "list-regions"{
		urlEnding = "regions"
	} else if clusterCommand == "list-members"{
		urlEnding = "members"
	}
	endpoint = "http://localhost:7070/geode-management/v2/" + urlEnding //TODO: must change !!!
	return
}

func getUrlOutput(endpointUrl string) (urlResponse string){
	resp, err := http.Get(endpointUrl)
	if err != nil{
		log.Fatal(err)
	}
	respInAscii, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil{
		log.Fatal(1)
	}
	urlResponse = fmt.Sprintf("%s", respInAscii)
	return
}
func fill(columnSize int, value string, filler string) (response string){
	if len(value) > columnSize - 1{
		response = " " + value[:columnSize-1]
		return
	}
	numFillerChars := columnSize - len(value) - 1
	response = " " + value + strings.Repeat(filler, numFillerChars)
	return
}


func getTableHeadersFromClusterCommand(clusterCommand string) (tableHeaders []string){
	if clusterCommand == "list-regions"{
		tableHeaders = append(tableHeaders, "name", "type", "groups", "entryCount", "regionAttributes")
	} else if clusterCommand =="list-members"{
		tableHeaders = append(tableHeaders, "id", "host", "status", "pid")
	} else{
		return
	}
	return
}

func getAnswerFromUrlResponse(clusterCommand string, urlResponse string) (response string){
	urlOutput := ClusterManagementResult{}
	err := json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		log.Fatal(err)
	}
	response = "Status Code: " + urlOutput.StatusCode + "\n"
	if urlOutput.StatusMessage != ""{
		response += "Status Message: " + urlOutput.StatusMessage + "\n"
	}
	response += "\n"

	tableHeaders := getTableHeadersFromClusterCommand(clusterCommand)
	for _, header := range tableHeaders {
		response += fill(20, header, " ") + "|"
	}
	response += "\n" + fill (20 * len(tableHeaders) + 5, "", "-") + "\n"

	memberCount := 0
	for _, result := range urlOutput.Result{
		memberCount++
		for _, key := range tableHeaders {
			if result[key] == nil {
				response += fill(20, "", " ") + "|"
			} else {
				resultVal := result[key]
				if fmt.Sprintf("%T", result[key]) == "float64"{
					resultVal = fmt.Sprintf("%.0f", result[key])
				}
				response += fill(20, fmt.Sprintf("%s",resultVal), " ") + "|"
			}
		}
		response += "\n"
	}
	if clusterCommand == "list-regions"{
		response += "\nNumber of Regions: " + strconv.Itoa(memberCount)
	} else if clusterCommand == "list-members"{
		response += "\nNumber of Members: " + strconv.Itoa(memberCount)
	}
	return
}


func getJsonFromUrlResponse(urlResponse string) (jsonOutput string){
	urlOutput := ClusterManagementResult{}
	err := json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		log.Fatal(err)
	}
	jsonExtracted, err := json.MarshalIndent(urlOutput, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	jsonOutput = string(jsonExtracted)
	return
}

func isRegionInGroups(regionInGroup bool, resultVal interface{}, groups []string) (isInGroups bool){
	if regionInGroup{
		isInGroups = true
		return
	}
	for _, group := range groups{
		resultValToStr := fmt.Sprintf("%s", resultVal)
		for  _, regionName := range strings.Split(resultValToStr[1:len(resultValToStr)-1], " "){
			if regionName == group{
				isInGroups = true
				return
			}
		}
	}
	isInGroups = false
	return
}

func editResponseOnGroup(urlResponse string, groups []string, clusterCommand string) (editedUrlResponse string){
	urlOutput := ClusterManagementResult{}
	err := json.Unmarshal([]byte(urlResponse), &urlOutput)
	if err != nil {
		log.Fatal(err)
	}
	var newUrlOutputResult string
	var newResult []map[string]interface{}
	for _, result := range urlOutput.Result{
		regionInGroups:= false
		for _, key :=range getTableHeadersFromClusterCommand(clusterCommand){
			if key == "groups"{
				regionInGroups = isRegionInGroups(regionInGroups, result[key], groups)
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
		log.Fatal(err)
	}
	editedUrlResponse = string(byteResponse)
	return
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL"{
		return
	}
	var username, password, endpoint, pccInUse, clusterCommand string
	var groups []string
	if len(args) >= 3 {
		pccInUse = args[2]
		clusterCommand = args[1]
	} else{
		fmt.Println(incorrectUserInputMessage)
		return
	}
	endpointLink := getEndpoint(clusterCommand)
	urlResponse := getUrlOutput(endpointLink)
	hasJ := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "--groups="){
			groups = strings.Split(arg[9:], ",")
			urlResponse = editResponseOnGroup(urlResponse, groups, clusterCommand)
		}
		if arg == "--j"{
			hasJ = true
		}
	}
	if hasJ{
		fmt.Println(getJsonFromUrlResponse(urlResponse))
		return
	}
	fmt.Println("PCC in use: " + pccInUse)
	if os.Getenv("CFLOGIN") != "" && os.Getenv("CFPASSWORD") != "" && os.Getenv("CFENDPOINT") != "" {
		username = os.Getenv("CFLOGIN")
		password = os.Getenv("CFPASSWORD")
		endpoint = os.Getenv("CFENDPOINT")
	} else {
		pccServicesAvailable := collectCloudCacheServices()
		if err := ValidatePCCInstance(pccInUse, pccServicesAvailable); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		serviceKey, err := getServiceKeyFromPCCInstance(pccInUse)
		if err != nil{
			fmt.Printf(err.Error(), pccInUse, pccInUse)
			os.Exit(1)
		}
		fmt.Println("Service key: " + serviceKey)
		username, password, endpoint = getUsernamePasswordEndpoint(pccInUse, serviceKey)
	}
	successMessage := fmt.Sprintf("Cluster Command: %s \nEndpoint: %s \nUsername: %s \nPassword: %s \n",
		clusterCommand, endpoint, username, password)
	if username != "" && password != "" && clusterCommand != "" && endpoint != "" {
		answer := getAnswerFromUrlResponse(clusterCommand, urlResponse)
		fmt.Println()
		fmt.Println(answer)
		fmt.Println()
		fmt.Println(successMessage)
	} else {
		fmt.Println(incorrectUserInputMessage)
	}
}


func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "CLI_InDev",
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
				Name:     "cli",
				HelpText: "cli's help text",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "   cf cli [action] [pcc_instance] [*options] (* = optional)\n" +
						"	Actions: \n" +
						"		list-regions, list-members\n" +
						"	Options: \n" +
						"		--h : this help screen\n" +
						"		--j : json output of API endpoint\n",
				},
			},
		},
	}
}


func main() {
	plugin.Start(new(BasicPlugin))
}
