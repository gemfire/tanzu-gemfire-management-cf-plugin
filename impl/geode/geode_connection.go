package geode

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"os"
	"strings"
)

type geodeConnection struct {
}

// NewGeodeConnectionProvider provides a constructor for the Geode standalone implementation of ConnectionProvider
func NewGeodeConnectionProvider() (impl.ConnectionProvider, error) {
	return &geodeConnection{}, nil
}

func (gc *geodeConnection) GetConnectionData(commandData *domain.CommandData) error {
	commandData.ConnnectionData = domain.ConnectionData{}

	// LocatorAddress, Username and Password may be provided as environment variables
	// but can be overridden on the command line
	commandData.ConnnectionData.LocatorAddress = strings.TrimSuffix(commandData.Target, "/")
	commandData.ConnnectionData.Username = common.GetOption(commandData.UserCommand.Parameters, []string{"--user", "-u"})
	if commandData.ConnnectionData.Username == "" {
		commandData.ConnnectionData.Username = os.Getenv("GEODE_USERNAME")
	}
	commandData.ConnnectionData.Password = common.GetOption(commandData.UserCommand.Parameters, []string{"--password", "-p"})
	if commandData.ConnnectionData.Password == "" {
		commandData.ConnnectionData.Password = os.Getenv("GEODE_PASSWORD")
	}

	return nil
}
