package node

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdK8sNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "k8s-node",
		Short:                 "k8s节点管理工具",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			k8sNode(cmd, args)
		},
	}
	AddNodeFlags(cmd.Flags())
	return cmd
}

func k8sNode(cmd *cobra.Command, args []string) {

	// k8s节点管理工具
	switch action {
	case "check":
		// 检查节点配置是否正常
		nodeCheck(cmd, args)
		break
	case "init":
		// 生成初始化脚本
		nodeInit(cmd, args)
		break
	default:
		logrus.Error("不支持该操作")
	}
}
