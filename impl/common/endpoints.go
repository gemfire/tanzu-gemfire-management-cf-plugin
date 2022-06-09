/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for additional information regarding
 * copyright ownership. The ASF licenses this file to You under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License. You may obtain a
 * copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common/format"
)

// For backwards compatibility purposes before 1.15 (openapi specification)
var swaggerBodyParams = map[string]string{
	"configure pdx":            "pdxType",
	"create disk-store":        "diskStoreConfig",
	"create gateway-receiver":  "gatewayReceiverConfig",
	"create index":             "indexConfig",
	"create region":            "regionConfig",
	"create region index":      "indexConfig",
	"start rebalance":          "operation",
	"start restore-redundancy": "operation",
}

const (
	contentTypeJson      = "application/json"
	contentTypeMultiForm = "multipart/form-data"
)

// GetEndPoints retrieves available endpoint from the Swagger endpoint on the Geode/PCC locator
func GetEndPoints(commandData *domain.CommandData, processRequest impl.RequestHelper) error {
	var urlResponse, apiDocURL string
	var statusCode int
	var err error
	var responseMap map[string]interface{}
	fallbackCodes := "401 403 404 407"
	apiDocURLs := []string{
		commandData.ConnnectionData.LocatorAddress + "/management/",
		commandData.ConnnectionData.LocatorAddress + "/management/v1/api-docs",
		commandData.ConnnectionData.LocatorAddress + "/management/experimental/api-docs",
		commandData.ConnnectionData.LocatorAddress + "/management/v3/api-docs",
	}

	for pos, URL := range apiDocURLs {
		apiDocURL = URL
		request, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return err
		}
		urlResponse, statusCode, err = processRequest(request)
		if err != nil {
			return errors.New("Unable to reach " + URL + ". Error: " + err.Error())
		}
		if !strings.Contains(fallbackCodes, strconv.Itoa(statusCode)) {
			if statusCode == 200 {
				if pos == 0 {
					err = json.Unmarshal([]byte(urlResponse), &responseMap)
					if err != nil {
						return errors.New("Unable to parse response: " + urlResponse + ". Error: " + err.Error())
					}
					latestURL, Ok := responseMap["latest"]
					if Ok {
						apiDocURL = format.GetString(latestURL)
						request, err := http.NewRequest("GET", apiDocURL, nil)
						if err != nil {
							return err
						}
						urlResponse, statusCode, err = processRequest(request)
						if err != nil {
							return errors.New("Unable to reach " + apiDocURL + ": " + err.Error())
						}
					} else {
						return errors.New("Unable to determine latest API endpoint: " + urlResponse + ".")
					}
				}
				break
			}
			return errors.New("Unable to reach " + URL + ". Status Code: " + strconv.Itoa(statusCode))
		}
	}

	if statusCode != 200 {
		return errors.New("Unable to reach " + apiDocURL + ". Status Code: " + strconv.Itoa(statusCode))
	}

	var apiPaths domain.RestAPI
	err = json.Unmarshal([]byte(urlResponse), &apiPaths)

	if err != nil {
		return errors.New("invalid response " + urlResponse + ": " + err.Error())
	}
	commandData.ConnnectionData.UseToken = apiPaths.Info.TokenEnabled == "true"
	commandData.AvailableEndpoints = make(map[string]domain.RestEndPoint)
	// Openapi specification is used if definitions are not set
	isOpenApi := false
	if apiPaths.Definitions == nil {
		isOpenApi = true
		apiPaths.Definitions = apiPaths.Components["schemas"]
	}
	for url, v := range apiPaths.Paths {
		for methodType := range v {
			var endpoint domain.RestEndPoint
			endpoint.URL = url
			endpoint.HTTPMethod = methodType
			endpoint.Consumes = apiPaths.Paths[url][methodType].Consumes
			endpoint.CommandName = apiPaths.Paths[url][methodType].CommandName
			endpoint.JQFilter = apiPaths.Paths[url][methodType].JQFilter
			endpoint.Parameters = apiPaths.Paths[url][methodType].Parameters
			for index, parameter := range endpoint.Parameters {
				if parameter.In == "body" {
					definitionPath := "#/definitions/"
					schemaName := strings.ReplaceAll(parameter.Schema["$ref"], definitionPath, "")
					if schemaName != "" {
						endpoint.Parameters[index].BodyDefinition = buildStructure(apiPaths.Definitions[schemaName].Properties, apiPaths.Definitions, definitionPath)
					}
				}
			}
			if isOpenApi {
				transformEndpointFromOpenApi(&endpoint, methodType, apiPaths, url)
			}
			commandData.AvailableEndpoints[endpoint.CommandName] = endpoint
		}
	}
	return nil
}

func transformEndpointFromOpenApi(endpoint *domain.RestEndPoint, methodType string, apiPaths domain.RestAPI, url string) {
	if strings.ToLower(methodType) == "post" || strings.ToLower(methodType) == "put" {
		requestBody := apiPaths.Paths[url][methodType].RequestBody
		applicationJson, ok := requestBody.Content[contentTypeJson].(map[string]interface{})
		if ok {
			endpoint.Consumes = []string{contentTypeJson}
			schemaName, ok := applicationJson["schema"].(map[string]interface{})["$ref"].(string)
			parameterName, ok := swaggerBodyParams[endpoint.CommandName]
			definitionPath := "#/components/schemas/"
			schemaName = strings.ReplaceAll(schemaName, definitionPath, "")
			if !ok {
				parameterName = "body"
			}
			param := domain.RestAPIParam{
				Name:        parameterName,
				Required:    requestBody.Required,
				Description: parameterName,
				In:          "body",
			}
			if schemaName != "" {
				param.BodyDefinition = buildStructure(apiPaths.Definitions[schemaName].Properties, apiPaths.Definitions, definitionPath)
			}
			endpoint.Parameters = append(endpoint.Parameters, param)
		}
		multipartForm, ok := requestBody.Content[contentTypeMultiForm].(map[string]interface{})
		if ok {
			endpoint.Consumes = []string{contentTypeMultiForm}
			required := multipartForm["schema"].(map[string]interface{})["required"].([]interface{})
			for _, req := range required {
				strReq := fmt.Sprintf("%v", req)
				endpoint.Parameters = append(endpoint.Parameters, domain.RestAPIParam{
					Name:        strReq,
					Required:    true,
					Description: "filePath",
					Type:        "file",
					In:          "formData",
				})
			}
		}
	}
}

func buildStructure(propertiesMap map[string]domain.PropertyDetail, definitions map[string]domain.DefinitionDetail, definitionPath string) (structure map[string]interface{}) {
	structure = make(map[string]interface{})
	for key, property := range propertiesMap {
		switch property.Type {
		case "string":
			if len(property.Enum) > 0 {
				structure[key] = "ENUM, one of: " + strings.Join(property.Enum, ", ")
			} else {
				structure[key] = "string-value"
			}
		case "integer":
			structure[key] = 42
		case "boolean":
			structure[key] = true
		case "object":
			structure[key] = map[string]string{"name": "value"}
		case "array":
			structure[key] = generateSampleArray(property.Items, definitions, definitionPath)
		case "":
			if len(property.Ref) > 0 {
				refName := strings.ReplaceAll(property.Ref, definitionPath, "")
				if refName != "" {
					if refName == "DeclarableType" {
						structure[key] = "DeclarableType"
					} else {
						subStructure := buildStructure(definitions[refName].Properties, definitions, definitionPath)
						structure[key] = subStructure
					}
				}
			}
		default:
			structure[key] = "unknown"
		}
	}
	return
}

func generateSampleArray(itemMap map[string]string, definitions map[string]domain.DefinitionDetail, definitionPath string) interface{} {
	switch itemMap["type"] {
	case "string":
		return []string{"stringOne", "stringTwo"}
	case "integer":
		return []int{41, 42}
	case "boolean":
		return []bool{true, false}
	case "":
		if len(itemMap["$ref"]) > 0 {
			refName := strings.ReplaceAll(itemMap["$ref"], definitionPath, "")
			if refName != "" {
				if refName == "DeclarableType" {
					return []string{"DeclarableType", "DeclarableType"}
				}
				subStructure := buildStructure(definitions[refName].Properties, definitions, definitionPath)
				return []interface{}{subStructure}
			}
		}
	default:
		return []interface{}{}
	}
	return []interface{}{}
}
