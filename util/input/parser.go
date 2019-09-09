package input

import (
	"os"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

// GetTargetAndClusterCommand extracts the target and command from the args and environment variables
func GetTargetAndClusterCommand(args []string) (target string, userCommand domain.UserCommand) {
	target = os.Getenv("CFPCC")
	if len(args) < 2 {
		return
	}

	commandStart := 2
	if target == "" {
		target = args[1]
	} else if target != args[1] {
		commandStart = 1
	}

	userCommand.Parameters = make(map[string]string)
	// find the command name before the options
	var option = ""
	for i := commandStart; i < len(args); i++ {
		token := args[i]
		if strings.HasPrefix(token, "-") {
			if option != "" {
				userCommand.Parameters[option] = "true"
			}
			option = token
		} else if option == "" {
			userCommand.Command += token + " "
		} else {
			userCommand.Parameters[option] = token
			option = ""
		}
	}
	userCommand.Command = strings.Trim(userCommand.Command, " ")
	if option != "" {
		userCommand.Parameters[option] = "true"
	}
	return
}
