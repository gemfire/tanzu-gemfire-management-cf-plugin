package pcc

import (
	"errors"
)

func mapUserInputToAvailableEndpoint() (IndividualEndpoint, error) {
	for _, ep := range availableEndpoints {
		if ep.CommandCall == userCommand.command {
			return ep, nil
		}
	}
	return IndividualEndpoint{}, errors.New(NoEndpointFoundMessage)
}
