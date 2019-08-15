package pcc

import (
	"strings"
)

func parseArguments(args []string) (err error) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-g=") {
			userCommand.parameters["g"] = "true"
			hasGroup = true
			group = arg[3:]
			if err != nil {
				return err
			}
		} else if arg == "-j" {
			userCommand.parameters["j"] = "true"
			isJSONOutput = true
		} else if strings.HasPrefix(arg, "-r=") {
			userCommand.parameters["r"] = "true"
			region = arg[3:]
		} else if strings.HasPrefix(arg, "--regionName=") {
			userCommand.parameters["r"] = "true"
			region = arg[13:]
		} else if strings.HasPrefix(arg, "-u=") {
			username = arg[3:]
		} else if strings.HasPrefix(arg, "-p=") {
			password = arg[3:]
		} else if strings.HasPrefix(arg, "-d=") {
			jsonFile = arg[3:]
		} else if strings.HasPrefix(arg, "--data=") {
			jsonFile = arg[7:]
		} else if strings.HasPrefix(arg, "-id=") {
			id = arg[4:]
			userCommand.parameters["id"] = id
		}
	}
	return nil
}
