package main

import (
	"github.com/spf13/cobra"
	"os"
	"sdp-devops/pkg/exporter"
	"sdp-devops/pkg/exporter/config"
	sdpLogger "sdp-devops/pkg/logger"
)

func main() {
	sdpLogger.InitLogger()
	command := NewExporterCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewExporterCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sdp-exporter",
		Short: "SDP Kubernetes 指标采集服务。",
		Run: func(cmd *cobra.Command, args []string) {
			exporter.Main()
		},
	}
	flags := rootCmd.PersistentFlags()
	config.AddFlags(flags)
	return rootCmd
}
