package config

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	KubeConfigStr string
	ShellToolName string
	DeployAction  string
)

func AddCommonFlags(flags *pflag.FlagSet) {
	flags.StringVar(&KubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群的config配置文件。")
	flags.StringVar(&ShellToolName, "shell-tool-name", "node-shell", "Shell客户端工具名称。")
}

func AddDeployFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&DeployAction, "deploy-action", "a", "check", "相关操作名称。")
}
