package geode

import (
	"errors"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
)

type geodeConnection struct {
}

// NewGeodeConnectionProvider provides a constructor for the Geode standalone implementation of ConnectionProvider
func NewGeodeConnectionProvider() (impl.ConnectionProvider, error) {
	return &geodeConnection{}, nil
}

func (gc *geodeConnection) GetConnectionData(args []string) (domain.ConnectionData, error) {
	connectionData := domain.ConnectionData{}
	for _, value := range args {
		if strings.HasPrefix(value, "-u=") {
			connectionData.Username = value[3:]
		} else if strings.HasPrefix(value, "-p=") {
			connectionData.Password = value[3:]
		} else if strings.HasPrefix(value, "http") {
			connectionData.LocatorAddress = value
		}
	}

	if len(connectionData.LocatorAddress) < 7 {
		return connectionData, errors.New(util.NoEndpointFoundMessage)
	}
	return connectionData, nil
}
