package deploy

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sdp-devops/pkg/sdpctl/config"
)

func NewCmdDeploy() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "deploy",
		Short:                 "在集群部署基于Daemonset的Shell客户端",
		DisableFlagsInUseLine: true,
	}
	config.AddDeployFlags(cmd.Flags())
	cmd.AddCommand(NewCmdNodeShell())
	return cmd
}

func NewCmdNodeShell() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "install",
		Short:                 "安装Shell Pod服务",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			switch config.DeployAction {
			case "install":
				install(cmd, args)
				break
			case "remove":
				clean(cmd, args)
				break
			default:
				logrus.Error("只支持 install / remove 操作")
			}
		},
	}
	return cmd
}
