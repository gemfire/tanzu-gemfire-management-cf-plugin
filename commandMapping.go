package main

import (
	"errors"
)

func mapUserInputToAvailableEndpoint() (IndividualEndpoint, error) {
	for _, ep := range availableEndpoints {
		if ep.CommandCall == APICallStruct.command {
			return ep, nil
		}
	}
	return IndividualEndpoint{}, errors.New(NoEndpointFoundMessage)
}

func printAvailableCommands()(string){
	executeFirstRequest()
	toPrint := ""
	for _, command := range availableEndpoints{
		toPrint += "\n	" + command.CommandCall
	}
	return toPrint + "\n"
}
