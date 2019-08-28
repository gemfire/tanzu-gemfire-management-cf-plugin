package pcc

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
)

type pluginConnection struct {
	cliConnection plugin.CliConnection
}

// NewPluginConnectionProvider provides a constructor for the PCC implementation of ConnectionProvider
func NewPluginConnectionProvider(connection plugin.CliConnection) (impl.ConnectionProvider, error) {
	return &pluginConnection{cliConnection: connection}, nil
}

// GetConnectionData provides the connection data from a PCC cluster using the CF CLI
func (pc *pluginConnection) GetConnectionData(args []string) (domain.ConnectionData, error) {
	serviceKey, err := pc.getServiceKey(args[0])
	if err != nil {
		return domain.ConnectionData{}, err
	}

	return pc.getServiceKeyDetails(args[0], serviceKey)

}

func (pc *pluginConnection) getServiceKey(target string) (serviceKey string, err error) {
	results, err := pc.cliConnection.CliCommandWithoutTerminalOutput("service-keys", target)
	if err != nil {
		return "", err
	}
	hasKey := false
	if strings.Contains(results[1], "No service key for service instance") {
		return "", fmt.Errorf(util.NoServiceKeyMessage, target, target)
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
		err = fmt.Errorf(util.NoServiceKeyMessage, target, target)
	}
	return
}

func (pc *pluginConnection) getServiceKeyDetails(target string, serviceKey string) (connectionData domain.ConnectionData, err error) {
	connectionData = domain.ConnectionData{}
	keyInfo, err := pc.cliConnection.CliCommandWithoutTerminalOutput("service-key", target, serviceKey)
	if err != nil {
		return connectionData, err
	}

	if len(keyInfo) < 2 {
		return connectionData, errors.New(util.InvalidServiceKeyResponse)
	}
	keyInfo = keyInfo[2:] //take out first two lines of cf service-key ... output
	joinKeyInfo := strings.Join(keyInfo, "\n")
	serviceKeyStruct := domain.ServiceKey{}

	err = json.Unmarshal([]byte(joinKeyInfo), &serviceKeyStruct)
	if err != nil {
		return connectionData, err
	}
	connectionData.LocatorAddress = serviceKeyStruct.Urls.Management
	if connectionData.LocatorAddress == "" {
		connectionData.LocatorAddress = strings.TrimSuffix(serviceKeyStruct.Urls.Gfsh, "/gemfire/v1")
	}
	for _, user := range serviceKeyStruct.Users {
		if strings.HasPrefix(user.Username, "cluster_operator") {
			connectionData.Username = user.Username
			connectionData.Password = user.Password
		}
	}
	return
}
