package requests

import (
	"errors"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
)

// MapUserInputToAvailableEndpoint matches a requested enpoint to available endpoints
func MapUserInputToAvailableEndpoint(commandData *domain.CommandData) error {
	for _, ep := range commandData.AvailableEndpoints {
		if ep.CommandCall == commandData.UserCommand.Command {
			commandData.Endpoint = ep
			return nil
		}
	}
	return errors.New(util.NoEndpointFoundMessage)
}
