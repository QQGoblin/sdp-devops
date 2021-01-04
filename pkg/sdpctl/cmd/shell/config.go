package shell

import "github.com/spf13/pflag"

var (
	httpTimeOutInSec int
	currentThreadNum int
	targetNode       string
	targetNodeFile   string
	shellMode        string
	toolName         string
)

func AddShellFlags(flags *pflag.FlagSet) {
	flags.IntVar(&httpTimeOutInSec, "kubelet-timeout", 30, "连接Kubelet超时时间。")
	flags.IntVar(&currentThreadNum, "thread", 1, "执行shell命令的并发数。")
	flags.StringVarP(&targetNode, "node", "n", "", "在指定宿主机节点执行操作。")
	flags.StringVar(&targetNodeFile, "nodefile", "", "通过文件指定要运行命令的宿主机。")
	flags.StringVar(&shellMode, "shellMode", "k8s-node", "执行shell的模式：k8s-node，docker-net")
	flags.StringVar(&toolName, "shell-tool-name", "node-shell", "Shell客户端工具名称。")
}
