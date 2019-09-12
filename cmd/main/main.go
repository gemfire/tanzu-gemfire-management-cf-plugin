package main

import (
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/geode"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/pcc"
)

func main() {
	helper := common.Requester{}
	commonCode, err := common.NewCommandProcessor(&helper)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// figure out who is calling
	if strings.Contains(os.Args[0], ".cf/plugins") {
		basicPlugin, err := pcc.NewBasicPlugin(commonCode)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		plugin.Start(&basicPlugin)
	} else {
		geodeCommand, err := geode.NewGeodeCommand(commonCode)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = geodeCommand.Run(os.Args)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
