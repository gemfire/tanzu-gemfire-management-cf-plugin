package main

import (
	"errors"
	"strings"
)

func mapUserInputToAvailableEndpoint() (IndividualEndpoint, error){
	convertTargetToPlural()
	switch APICallStruct.action{
	case "list":
		for _, ep := range availableEndpoints{
			if ep.HttpMethod == "get" && !strings.Contains(ep.Url, "{id}") &&
				hasCorrectUrlEnding(ep.Url){
				return ep, nil
			}
		}
	case "get":
		for _, ep := range availableEndpoints{
			if ep.HttpMethod == "get" && strings.Contains(ep.Url, "{id}") &&
				hasCorrectUrlEnding(ep.Url){
				return ep, nil
			}
		}
	case "delete":
		for _, ep := range availableEndpoints{
			if ep.HttpMethod == "delete" &&
				hasCorrectUrlEnding(ep.Url){
				return ep, nil
			}
		}
	case "create":
		for _, ep := range availableEndpoints{
			if ep.HttpMethod == "post" &&
				hasCorrectUrlEnding(ep.Url){
				return ep, nil
			}
		}
	case "configure":
		for _, ep := range availableEndpoints{
			if ep.HttpMethod == "post" &&
				strings.Contains(ep.Url, "configurations/"+APICallStruct.target){
				return ep, nil
			}
		}
	}
	return IndividualEndpoint{}, errors.New(NoEndpointFoundMessage)
}

func convertTargetToPlural(){
	if APICallStruct.target[len(APICallStruct.target)-1] == 's' {
		return
	} else if APICallStruct.target == "cli" || APICallStruct.target == "pdx" || APICallStruct.target == "ping"{
		return
	} else if APICallStruct.target[len(APICallStruct.target)-1] == 'x'{
		APICallStruct.target = APICallStruct.target + "es"
	} else{
		APICallStruct.target = APICallStruct.target + "s"
	}
}

func hasCorrectUrlEnding(fullUrl string)(bool){
	ending := ""
	if strings.Contains(fullUrl, "/{id}"){
		ending =  fullUrl[len(fullUrl) - len(APICallStruct.target) - len("/{id}"):]
	} else{
		ending = fullUrl[len(fullUrl) - len(APICallStruct.target):]
	}
	return strings.Contains(ending, APICallStruct.target)
}

func synonymConverter(initialWord string) (string){
	if initialWord == "post" || initialWord == "start"{
		return "create"
	} else if initialWord == "check" {
		return "get"
	} else{
		return initialWord
	}
}

