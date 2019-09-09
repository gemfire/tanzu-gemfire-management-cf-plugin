package impl

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RequestHelper

// RequestHelper interface provides a way to get request related items
type RequestHelper interface {
	ExecuteCommand(endpointURL string, httpAction string, commandData *domain.CommandData) (urlResponse string, err error)
	GetEndPoints(commandData *domain.CommandData) error
}
