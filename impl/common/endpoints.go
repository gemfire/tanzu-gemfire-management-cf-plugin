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
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
)

// GetEndPoints retrieves available endpoint from the Swagger endpoint on the Geode/PCC locator
func GetEndPoints(commandData *domain.CommandData, requester impl.RequestHelper) error {
	var urlResponse, apiDocURL string
	var statusCode int
	var err error
	var responseMap map[string]interface{}
	fallbackCodes := "401 403 404 407"
	apiDocURLs := []string{
		commandData.ConnnectionData.LocatorAddress + "/management/",
		commandData.ConnnectionData.LocatorAddress + "/management/v1/api-docs",
		commandData.ConnnectionData.LocatorAddress + "/management/experimental/api-docs",
	}

	for pos, URL := range apiDocURLs {
		apiDocURL = URL
		urlResponse, statusCode, err = requester.Exchange(URL, "GET", nil, nil)
		if err != nil {
			return errors.New("unable to reach " + URL + ". Error: " + err.Error())
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
						apiDocURL = getString(latestURL)
						urlResponse, statusCode, err = requester.Exchange(apiDocURL, "GET", nil, nil)
						if err != nil {
							return errors.New("unable to reach " + apiDocURL + ": " + err.Error())
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
	for url, v := range apiPaths.Paths {
		for methodType := range v {
			var endpoint domain.RestEndPoint
			endpoint.URL = url
			endpoint.HTTPMethod = methodType
			endpoint.CommandName = apiPaths.Paths[url][methodType].CommandName
			endpoint.JQFilter = apiPaths.Paths[url][methodType].JQFilter
			endpoint.Parameters = []domain.RestAPIParam{}
			endpoint.Parameters = apiPaths.Paths[url][methodType].Parameters
			for index, parameter := range endpoint.Parameters {
				if parameter.In == "body" {
					schemaName := strings.ReplaceAll(parameter.Schema["$ref"], "#/definitions/", "")
					if schemaName != "" {
						endpoint.Parameters[index].BodyDefinition = buildStructure(apiPaths.Definitions[schemaName].Properties, apiPaths.Definitions)
					}
				}
			}

			commandData.AvailableEndpoints[endpoint.CommandName] = endpoint
		}
	}
	return nil
}

func buildStructure(propertiesMap map[string]domain.PropertyDetail, definitions map[string]domain.DefinitionDetail) (structure map[string]interface{}) {
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
			structure[key] = generateSampleArray(property.Items, definitions)
		case "":
			if len(property.Ref) > 0 {
				refName := strings.ReplaceAll(property.Ref, "#/definitions/", "")
				if refName != "" {
					if refName == "DeclarableType" {
						structure[key] = "DeclarableType"
					} else {
						subStructure := buildStructure(definitions[refName].Properties, definitions)
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

func generateSampleArray(itemMap map[string]string, definitions map[string]domain.DefinitionDetail) interface{} {
	switch itemMap["type"] {
	case "string":
		return []string{"stringOne", "stringTwo"}
	case "integer":
		return []int{41, 42}
	case "boolean":
		return []bool{true, false}
	case "":
		if len(itemMap["$ref"]) > 0 {
			refName := strings.ReplaceAll(itemMap["$ref"], "#/definitions/", "")
			if refName != "" {
				if refName == "DeclarableType" {
					return []string{"DeclarableType", "DeclarableType"}
				}
				subStructure := buildStructure(definitions[refName].Properties, definitions)
				return []interface{}{subStructure}
			}
		}
	default:
		return []interface{}{}
	}
	return []interface{}{}
}
