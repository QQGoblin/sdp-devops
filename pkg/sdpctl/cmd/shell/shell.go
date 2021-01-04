package shell

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdSh() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sh [command]",
		Short:                 "执行Shell命令,请用“--”分隔开Shell命令",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			switch shellMode {
			case "k8s-node":
				nodeShell(cmd, args)
				break
			case "docker-net":
				dockerNet(cmd, args)
				break
			default:
				logrus.Error("不支持该模式执行shell")
			}
		},
	}
	AddShellFlags(cmd.Flags())
	return cmd
}

type OutPut struct {
	Title  string
	StdOut *bytes.Buffer
	StdErr *bytes.Buffer
}
