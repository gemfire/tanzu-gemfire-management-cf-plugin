package util

import (
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

// ParseArguments parses command line arguments into the CommandData struct passed by pointer
func ParseArguments(args []string, commandData *domain.CommandData) (err error) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-g=") {
			commandData.UserCommand.Parameters["g"] = "true"
			commandData.HasGroup = true
			commandData.Group = arg[3:]
			if err != nil {
				return
			}
		} else if arg == "-j" {
			commandData.UserCommand.Parameters["j"] = "true"
			commandData.IsJSONOutput = true
		} else if strings.HasPrefix(arg, "-r=") {
			commandData.UserCommand.Parameters["r"] = "true"
			commandData.Region = arg[3:]
		} else if strings.HasPrefix(arg, "--regionName=") {
			commandData.UserCommand.Parameters["r"] = "true"
			commandData.Region = arg[13:]
		} else if strings.HasPrefix(arg, "-u=") {
			commandData.Username = arg[3:]
		} else if strings.HasPrefix(arg, "-p=") {
			commandData.Password = arg[3:]
		} else if strings.HasPrefix(arg, "-d=") {
			commandData.JSONFile = arg[3:]
		} else if strings.HasPrefix(arg, "--data=") {
			commandData.JSONFile = arg[7:]
		} else if strings.HasPrefix(arg, "-id=") {
			commandData.ID = arg[4:]
			commandData.UserCommand.Parameters["id"] = commandData.ID
		}
	}
	return
}
