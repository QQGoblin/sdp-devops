package sdpctl

import (
	"github.com/spf13/cobra"
	"os"
	"sdp-devops/pkg/sdpctl/cmd/get"
	"sdp-devops/pkg/sdpctl/cmd/ops"
	"sdp-devops/pkg/sdpctl/cmd/shell"
	"sdp-devops/pkg/sdpctl/config"
)

func Main() {
	var rootCmd = &cobra.Command{
		Use:   "sdpctl",
		Short: "ND Kubernetes 运维工具",
		Run:   runHelp,
	}

	flags := rootCmd.PersistentFlags()

	config.AddCommonFlags(flags)
	rootCmd.AddCommand(get.NewCmdGet())
	rootCmd.AddCommand(ops.NewCmdOps())
	rootCmd.AddCommand(shell.NewCmdSh())
	if err := execute(rootCmd); err != nil {
		os.Exit(1)
	}
}

func execute(cmd *cobra.Command) error {
	err := cmd.Execute()
	return err
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
