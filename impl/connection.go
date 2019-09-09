package impl

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ConnectionProvider

// ConnectionProvider interface defines a way to get connection information for a Geode cluster
type ConnectionProvider interface {
	GetConnectionData(commandData *domain.CommandData) error
}
