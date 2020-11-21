package sys

import (
	"bytes"
	"os/exec"
)

func CmdOutErr(name string, arg ...string) (string, string, error) {
	cmd := exec.Command(name, arg...)
	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	err := cmd.Run()
	return cmdOut.String(), cmdErr.String(), err
}
