package config

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	KubeConfigStr string
)

func AddCommonFlags(flags *pflag.FlagSet) {
	flags.StringVar(&KubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群的config配置文件。")

}
