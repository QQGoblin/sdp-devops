package get

import (
	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get",
		Short:                 "打印信息",
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(NewCmdNode())
	return cmd
}

func NewCmdNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "node",
		Short:                 "打印更加细致的集群信息",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			node(cmd, args)
		},
	}
	return cmd
}
