package main

import (
	"strings"
)

func parseArguments(args []string)(err error){
	for _, arg := range args {
		if strings.HasPrefix(arg, "-g="){
			APICallStruct.parameters["g"] = "true"
			hasGroup = true
			group = arg[3:]
			if err != nil{
				return err
			}
		} else if arg == "-j"{
			APICallStruct.parameters["j"] ="true"
			isJSONOutput = true
		} else if strings.HasPrefix(arg, "-r="){
			APICallStruct.parameters["r"] = "true"
			region = arg[3:]
		} else if strings.HasPrefix(arg, "--regionName="){
			APICallStruct.parameters["r"] = "true"
			region = arg[13:]
		} else if strings.HasPrefix(arg, "-u="){
			username = arg[3:]
		} else if strings.HasPrefix(arg, "-p="){
			password = arg[3:]
		} else if strings.HasPrefix(arg, "-d="){
			jsonFile = arg[3:]
		} else if strings.HasPrefix(arg, "--data="){
			jsonFile = arg[7:]
		} else if strings.HasPrefix(arg, "-id="){
			id=arg[4:]
			APICallStruct.parameters["id"] = id
		}
	}
	return nil
}
