package ops

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdOps() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ops",
		Short:                 "管理SDP相关依赖",
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(NewCmdNodeShell())
	return cmd
}

func NewCmdNodeShell() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "node-shell",
		Short:                 "部署Shell Pod工具",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			switch action {
			case "install":
				install(cmd, args)
				break
			case "remove":
				clean(cmd, args)
				break
			default:
				logrus.Error("只支持 install / remove 操作")
			}
		},
	}
	AddNodeShellFlags(cmd.Flags())
	return cmd
}
