package main

import (
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/pcc"
)

func main() {
	// figure out who is calling
	if strings.Contains(os.Args[0], ".cf/plugins") {
		plugin.Start(new(pcc.BasicPlugin))
	} else {
		fmt.Println("Standalone mode is a future feature")
	}
}
