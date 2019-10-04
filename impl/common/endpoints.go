package common

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
)

// GetEndPoints retrieves available endpoint from the Swagger endpoint on the Geode/PCC locator
func GetEndPoints(commandData *domain.CommandData, requester impl.RequestHelper) error {
	apiDocURL := commandData.ConnnectionData.LocatorAddress + "/management/experimental/api-docs"
	urlResponse, err := requester.Exchange(apiDocURL, "GET", nil, nil)

	if err != nil {
		return errors.New("unable to reach " + apiDocURL + ": " + err.Error())
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
					subStructure := buildStructure(definitions[refName].Properties, definitions)
					structure[key] = subStructure
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
				subStructure := buildStructure(definitions[refName].Properties, definitions)
				return []interface{}{subStructure}
			}
		}
	default:
		return []interface{}{}
	}
	return []interface{}{}
}
