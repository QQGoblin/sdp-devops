package shell

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

func dockerNet(cmd *cobra.Command, args []string) {

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
