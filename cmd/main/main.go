package main

import (
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/geode"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/pcc"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/requests"
)

func main() {
	helper := requests.Helper{}
	common, err := common.NewCommon(&helper)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// figure out who is calling
	if strings.Contains(os.Args[0], ".cf/plugins") {
		basicPlugin, err := pcc.NewBasicPlugin(common)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		plugin.Start(&basicPlugin)
	} else {
		geodeCommand, err := geode.NewGeodeCommand(common)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		geodeCommand.Run(os.Args)
	}
}
