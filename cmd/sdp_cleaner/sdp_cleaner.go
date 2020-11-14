package main

import (
	"github.com/spf13/cobra"
	"os"
	"sdp-devops/pkg/cleaner"
)

func main() {

	command := NewCleanerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewCleanerCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sdp-cleaner",
		Short: "SDP Kubernetes 磁盘清理工具。",
		Run: func(cmd *cobra.Command, args []string) {
			cleaner.Main()
		},
	}
	flags := rootCmd.PersistentFlags()
	cleaner.AddFlags(flags)
	return rootCmd
}
