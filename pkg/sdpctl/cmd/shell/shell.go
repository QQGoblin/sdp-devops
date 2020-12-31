package shell

import "github.com/spf13/cobra"

func NewCmdSh() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sh [command]",
		Short:                 "在宿主机的客户端中执行Shell命令,请用“--”分隔开Shell命令，慎用！！！",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if !localNsenter {
				nodeShell(cmd, args)
			} else {
				dockerNet(cmd, args)
			}
		},
	}
	AddShellFlags(cmd.Flags())
	return cmd
}
