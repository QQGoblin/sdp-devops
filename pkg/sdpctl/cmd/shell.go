package cmd

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"os"
	"sdp-devops/pkg/sdpctl/config"
	"sdp-devops/pkg/sdpctl/sdpk8s"
	k8stools "sdp-devops/pkg/util/kubernetes"
	systools "sdp-devops/pkg/util/sys"
	"strconv"
	"strings"
)

func RunShell(cmd *cobra.Command, args []string) {

	kubeClientSet, kubeClientConfig := k8stools.KubeClientAndConfig(config.KubeConfigStr)

	cmdStr := strings.Join(args, " ")
	// 返回所有需要运行运行的Node列表
	shellPodTargets := sdpk8s.GetShellPodDict(kubeClientSet)
	i := 0
	threadNum := 0
	total := len(shellPodTargets)
	tChan := make(chan int, len(shellPodTargets))
	outPutBuffers := make([]*bytes.Buffer, len(shellPodTargets))
	defer close(tChan)

	for n, pod := range shellPodTargets {
		outPutBuffer := bytes.NewBufferString("------------------------------> No." + strconv.Itoa(i) + " Shell on node: " + n + " <------------------------------\n")
		outPutBuffers[i] = outPutBuffer
		if pod != nil {
			shExecOps := k8stools.ExecOptions{
				Command:       cmdStr,
				ContainerName: "",
				In:            nil,
				Out:           outPutBuffer,
				Err:           os.Stderr,
				Istty:         false,
				TimeOut:       config.HttpTimeOutInSec,
			}
			go execCmdParallel(kubeClientSet, kubeClientConfig, pod, shExecOps, tChan)
			threadNum += 1
		} else {
			outPutBuffer.WriteString("Can't find shell pod on " + n + "\n")
		}
		i += 1
		if threadNum == config.CurrentThreadNum || total == i {
			systools.WaitAllThreadFinish(threadNum, tChan, config.HttpTimeOutInSec)
			threadNum = 0
		}
	}

	for _, output := range outPutBuffers {
		fmt.Print(output.String())
	}
}

func execCmdParallel(kubeClientSet *kubernetes.Clientset, kubeClientConfig *restclient.Config, pod *v1.Pod, execOptions k8stools.ExecOptions, tChan chan int) {
	err := k8stools.ExecCmd(kubeClientSet, kubeClientConfig, pod, execOptions)
	if err != nil {
		logrus.Println("请求 API Service 返回异常：", pod.Status.HostIP)
		//panic(err.Error())
	}
	tChan <- 1
}

func NewCmdSh() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sh [command]",
		Short:                 "在宿主机的客户端中执行Shell命令,不支持管道、&、重定向等Shell操作符，慎用！！！",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunShell(cmd, args)
		},
	}
	config.AddShellFlags(cmd.Flags())
	return cmd
}
