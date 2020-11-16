package cmd

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"os"
	"sdp-devops/pkg/sdpctl/config"
	k8stools "sdp-devops/pkg/util/kubernetes"
	systools "sdp-devops/pkg/util/sys"
	"strconv"
	"strings"
)

var (
	httpTimeOutInSec int
	currentThreadNum int
)

func RunShell(cmd *cobra.Command, args []string) {

	kubeClientSet, kubeClientConfig := k8stools.KubeClientAndConfig(config.KubeConfigStr)

	cmdStr := strings.Join(args, " ")
	// 返回所有需要运行运行的Node列表
	shellPodTargets := k8stools.GetShellPodDict(kubeClientSet)
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
				TimeOut:       httpTimeOutInSec,
			}
			go execCmdParallel(kubeClientSet, kubeClientConfig, pod, shExecOps, tChan)
			threadNum += 1
		} else {
			outPutBuffer.WriteString("Can't find shell pod on " + n + "\n")
		}
		i += 1
		if threadNum == currentThreadNum || total == i {
			systools.WaitAllThreadFinish(threadNum, tChan, httpTimeOutInSec)
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

func addShFlag(flags *pflag.FlagSet) {

	flags.IntVar(&httpTimeOutInSec, "timeout", 30, "连接Kubelet超时时间。")
	flags.IntVar(&currentThreadNum, "thread", 1, "执行shell命令的并发数。")
}

func NewCmdSh() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sh [command]",
		Short:                 "在宿主机的客户端中执行Shell命令",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunShell(cmd, args)
		},
	}
	addShFlag(cmd.Flags())
	return cmd
}
