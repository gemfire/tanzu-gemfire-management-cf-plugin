package util

import (
	"regexp"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

var paramRegex *regexp.Regexp

func init() {
	paramRegex = regexp.MustCompile(`-+(\w*)=?(@?\w*)`)
}

// ParseArguments parses command line arguments into the CommandData struct passed by pointer
func ParseArguments(args []string, commandData *domain.CommandData) {
	if commandData.UserCommand.Parameters == nil {
		commandData.UserCommand.Parameters = make(map[string]interface{})
	}

	for _, arg := range args {
		parseResults := paramRegex.FindStringSubmatch(arg)

		if len(parseResults) == 3 {
			if len(parseResults[2]) == 0 {
				commandData.UserCommand.Parameters[parseResults[1]] = true
			} else if ok, _ := regexp.MatchString(`^true$`, strings.ToLower(parseResults[2])); ok {
				commandData.UserCommand.Parameters[parseResults[1]] = true
			} else if ok, _ := regexp.MatchString(`^false$`, strings.ToLower(parseResults[2])); ok {
				commandData.UserCommand.Parameters[parseResults[1]] = false
			} else {
				commandData.UserCommand.Parameters[parseResults[1]] = parseResults[2]
			}
		}
	}
	return
}
