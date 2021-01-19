package main

import (
	"github.com/spf13/cobra"
	"os"
	"sdp-devops/pkg/alert"
	"sdp-devops/pkg/alert/config"
	"sdp-devops/pkg/util"
)

func main() {
	util.InitLogger()
	command := NewAlertServerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewAlertServerCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sdp-alert",
		Short: "SDP Kubernetes 告警网关。",
		Run: func(cmd *cobra.Command, args []string) {
			alert.Main()
		},
	}
	flags := rootCmd.PersistentFlags()
	config.AddFlags(flags)
	return rootCmd
}
