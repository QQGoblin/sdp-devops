package goblin

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func CmdOutErr(name string, arg ...string) (string, string, error) {
	cmd := exec.Command(name, arg...)
	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	err := cmd.Run()
	return cmdOut.String(), cmdErr.String(), err
}

func AddPrefix(prefix, s string) string {
	lines := strings.Split(s, "\n")
	var build strings.Builder
	for _, line := range lines {
		build.WriteString(fmt.Sprintf("%s %s\n", prefix, line))
	}
	return build.String()
}
