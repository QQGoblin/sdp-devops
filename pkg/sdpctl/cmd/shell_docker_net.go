package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sdp-devops/pkg/sdpctl/sdpk8s"
	dockertools "sdp-devops/pkg/util/docker"
	systools "sdp-devops/pkg/util/sys"
	"strconv"
	"strings"
)

func RunShellDockerNet(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		logrus.Error("请输入要执行的命令")
	}

	cmdStr := strings.Join(args, " ")
	cli := dockertools.DockerClient("")
	podBriefInfoList := sdpk8s.GetPodInfo(cli)

	for i, podBriefInfo := range podBriefInfoList {
		fmt.Println("============================= No.", i, "POD:", podBriefInfo.Name, "=============================")
		if outStr, errStr, err := systools.CmdOutErr("/usr/bin/nsenter", "-t", strconv.Itoa(podBriefInfo.PID), "-n", "/bin/sh", "-c", cmdStr); err != nil {
			fmt.Println(errStr)
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
