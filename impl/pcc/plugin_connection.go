package pcc

import (
	"encoding/json"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
)

type pluginConnection struct {
	cliConnection plugin.CliConnection
}

// NewPluginConnectionProvider provides a constructor for the PCC implementation of ConnectionProvider
func NewPluginConnectionProvider(connection plugin.CliConnection) (impl.ConnectionProvider, error) {
	return &pluginConnection{cliConnection: connection}, nil
}

// GetConnectionData provides the connection data from a PCC cluster using the CF CLI
func (pc *pluginConnection) GetConnectionData(commandData *domain.CommandData) error {
	commandData.ConnnectionData = domain.ConnectionData{}
	serviceKey, err := pc.getServiceKey(commandData.Target)
	if err != nil {
		return err
	}

	return pc.getServiceKeyDetails(commandData, serviceKey)

}

func (pc *pluginConnection) getServiceKey(target string) (serviceKey string, err error) {
	results, err := pc.cliConnection.CliCommandWithoutTerminalOutput("service-keys", target)
	if err != nil {
		return "", err
	}
	hasKey := false
	if strings.Contains(results[1], "No service key for service instance") {
		return "", fmt.Errorf(common.NoServiceKeyMessage, target, target)
	}
	for _, value := range results {
		line := strings.Fields(value)
		if len(line) > 0 {
			if hasKey {
				serviceKey = line[0]
				return
			} else if line[0] == "name" {
				hasKey = true
			}
		}
	}
	if serviceKey == "" {
		err = fmt.Errorf(common.NoServiceKeyMessage, target, target)
	}
	return
}

func (pc *pluginConnection) getServiceKeyDetails(commandData *domain.CommandData, serviceKey string) (err error) {
	keyInfo, err := pc.cliConnection.CliCommandWithoutTerminalOutput("service-key", commandData.Target, serviceKey)
	if err != nil {
		return err
	}

	if len(keyInfo) < 2 {
		return errors.New(common.InvalidServiceKeyResponse)
	}
	keyInfo = keyInfo[2:] //take out first two lines of cf service-key ... output
	joinKeyInfo := strings.Join(keyInfo, "\n")
	serviceKeyStruct := domain.ServiceKey{}

	err = json.Unmarshal([]byte(joinKeyInfo), &serviceKeyStruct)
	if err != nil {
		return err
	}

	commandData.ConnnectionData.LocatorAddress = strings.TrimSuffix(serviceKeyStruct.Urls.Gfsh, "/gemfire/v1")
	for _, user := range serviceKeyStruct.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			commandData.ConnnectionData.Username = user.Username
			commandData.ConnnectionData.Password = user.Password
		}
	}

	if commandData.ConnnectionData.Username == "" || commandData.ConnnectionData.Password == "" {
		return errors.New("Unable to retrieve username/password from the servicekey.")
	}
	return
}
