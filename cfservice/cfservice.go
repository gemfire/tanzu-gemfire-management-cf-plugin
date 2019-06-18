package cfservice

import (
	"bytes"
	"errors"
	"os/exec"
)

const invalidPCCInstanceMessage string = `You entered %s which not a deployed PCC instance.
To deploy this as an instance, enter: 

	cf create-service p-cloudcache <region_plan> %s

For help see: cf create-service --help

`

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . CfService

type CfService interface {
	Cmd(name string, options ...string) (string, error)
}

type Cf struct {}

func (c *Cf) Cmd(name string, options ...string) (string, error){
	options = append([]string{name}, options...)
	cmd := exec.Command("cf", options...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", errors.New(invalidPCCInstanceMessage)
	}
	return (&out).String(), nil
}
