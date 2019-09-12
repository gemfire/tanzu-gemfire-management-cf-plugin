package input

import "github.com/gemfire/cloudcache-management-cf-plugin/domain"

// HasOption checks if a option has been passed in on the command line
func HasOption(commandData *domain.CommandData, option string) bool {
	return commandData.UserCommand.Parameters[option] != "" || commandData.Target == option
}
