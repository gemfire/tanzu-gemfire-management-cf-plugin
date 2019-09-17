package common

import (
	"os"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

// GetTargetAndClusterCommand extracts the target and command from the args and environment variables
func GetTargetAndClusterCommand(args []string) (target string, userCommand domain.UserCommand) {
	if len(args) < 2 {
		return
	}
	target = os.Getenv("GEODE_TARGET")
	commandStart := 2
	if target == "" && !strings.HasPrefix(args[1], "-") {
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
				userCommand.Parameters[option] = ""
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
		userCommand.Parameters[option] = ""
	}
	return
}

// HasOption checks if a option has been passed in on the command line
func HasOption(parameters map[string]string, options []string) bool {
	for _, option := range options {
		_, available := parameters[option]
		if available {
			return true
		}
	}
	return false
}

// GetOption retrieves entries from the map of parameters by name
func GetOption(parameters map[string]string, options []string) string {
	for _, option := range options {
		value := parameters[option]
		if value != "" {
			return value
		}
	}
	return ""
}
