package geode

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
)

type geodeConnection struct {
}

// NewGeodeConnectionProvider provides a constructor for the Geode standalone implementation of ConnectionProvider
func NewGeodeConnectionProvider() (impl.ConnectionProvider, error) {
	return &geodeConnection{}, nil
}

func (gc *geodeConnection) GetConnectionData(args ...string) (domain.ConnectionData, error) {
	return domain.ConnectionData{}, nil
}
