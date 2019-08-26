package impl

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

//go:generate counterfeiter . ConnectionProvider

// ConnectionProvider interface defines a way to get connection information for a Geode cluster
type ConnectionProvider interface {
	GetConnectionData(args ...string) (domain.ConnectionData, error)
}
