package config

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	KubeConfigStr     string
	ShellNamespace    string
	ShellDaemonset    string
	ShellNodeName     string
	ShellNodeNameFile string
	HttpTimeOutInSec  int
	CurrentThreadNum  int
	FileDistTimeOut   int
)

func AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&KubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群的config配置文件。")

	flags.StringVar(&ShellNamespace, "namespace", "node-shell", "安装Shell客户端的命名空间名。")
	flags.StringVar(&ShellDaemonset, "ds", "node-shell-tool", "安装Shell客户端的控制器名称。")

	flags.StringVar(&ShellNodeName, "node", "", "在指定宿主机节点执行操作。")
	flags.StringVar(&ShellNodeNameFile, "nodefile", "", "通过文件指定要运行命令的宿主机。")

}

func AddShellFlags(flags *pflag.FlagSet) {
	flags.IntVar(&HttpTimeOutInSec, "kubelet-timeout", 30, "连接Kubelet超时时间。")
	flags.IntVar(&CurrentThreadNum, "thread", 1, "执行shell命令的并发数。")
}

func AddDistFlags(flags *pflag.FlagSet) {
	flags.IntVar(&FileDistTimeOut, "dist-timeout", 15, "单个文件传输超时时间。")
}
