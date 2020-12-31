package config

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	KubeConfigStr string
	ShellToolName string
)

func AddCommonFlags(flags *pflag.FlagSet) {
	flags.StringVar(&KubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群的config配置文件。")
	flags.StringVar(&ShellToolName, "shell-tool-name", "node-shell", "Shell客户端工具名称。")
}
