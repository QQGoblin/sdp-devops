package ops

import (
	"github.com/spf13/cobra"
	"sdp-devops/pkg/sdpctl/cmd/ops/shell_tool"
)

func NewCmdOps() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ops",
		Short:                 "管理SDP相关依赖",
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(shell_tool.NewCmdShellTool())
	return cmd
}
