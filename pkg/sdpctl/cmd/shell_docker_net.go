package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/toolkits/sys"
	dockertools "sdp-devops/pkg/util/docker"
	"strconv"
	"strings"
)

func RunShellDockerNet(cmd *cobra.Command, args []string) {

	cmdStr := strings.Join(args, " ")
	cli := dockertools.DockerClient("")
	podBriefInfoList := dockertools.GetPodInfo(cli)

	for _, podBriefInfo := range podBriefInfoList {
		fmt.Println("============================= POD: " + podBriefInfo.Name + "=============================")
		if outStr, err := sys.CmdOut("nsenter", "-t", strconv.Itoa(podBriefInfo.PID), "-n", cmdStr); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(outStr)
		}
	}

}

func NewCmdShellDockerNet() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sh-docker-net",
		Short:                 "在当前节点Pod容器的网络空间中执行Shell命令, 基于 nsenter -t ${NET_ID} -n ${CMD}",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunShellDockerNet(cmd, args)
		},
	}
	return cmd
}
