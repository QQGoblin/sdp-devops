package shell

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"sdp-devops/pkg/sdpctl/config"
	"sdp-devops/pkg/util/goblin"
	k8stools "sdp-devops/pkg/util/kubernetes"
	"strings"
)

func nodeShell(cmd *cobra.Command, args []string) {

	kubeClientSet, kubeClientConfig := k8stools.KubeClientAndConfig(config.KubeConfigStr)

	cmdStr := strings.Join(args, " ")
	// 返回所有需要运行运行的Node列表
	shellPodTargets := getShellPodDict(kubeClientSet, targetNode, targetNodeFile, toolName)
	i := 0
	threadNum := 0
	total := len(shellPodTargets)
	tChan := make(chan int, len(shellPodTargets))
	outPuts := make([]OutPut, len(shellPodTargets))
	defer close(tChan)

	for n, pod := range shellPodTargets {
		outPut := OutPut{
			NodeName: n,
			StdOut:   bytes.NewBufferString(""),
			StdErr:   bytes.NewBufferString(""),
		}

		outPuts[i] = outPut
		if pod != nil {
			shExecOps := k8stools.ExecOptions{
				Command:       cmdStr,
				ContainerName: "",
				In:            nil,
				Out:           outPut.StdOut,
				Err:           outPut.StdErr,
				Istty:         false,
				TimeOut:       httpTimeOutInSec,
			}
			go execCmdParallel(kubeClientSet, kubeClientConfig, pod, shExecOps, tChan)
			threadNum += 1
		} else {
			outPut.StdErr.WriteString("Can't find shell pod on " + n)
		}
		i += 1
		if threadNum == currentThreadNum || total == i {
			goblin.WaitAllThreadFinish(threadNum, tChan, httpTimeOutInSec)
			threadNum = 0
		}
	}
	printOutput(outPuts)
}

func execCmdParallel(kubeClientSet *kubernetes.Clientset, kubeClientConfig *restclient.Config, pod *v1.Pod, execOptions k8stools.ExecOptions, tChan chan int) {
	err := k8stools.ExecCmd(kubeClientSet, kubeClientConfig, pod, execOptions)
	if err != nil {
		logrus.Println("请求 API Service 返回异常：", pod.Status.HostIP)
		//panic(err.Error())
	}
	tChan <- 1
}

func printOutput(outPuts []OutPut) {
	for i, output := range outPuts {
		switch format {
		case "raw":
			color.Blue("------------------------------> No.%d  Shell on node: %s <------------------------------", i, output.NodeName)
			fmt.Printf(output.StdOut.String())
			color.HiRed(output.StdErr.String())
			break
		case "prefix":
			color.Blue("------------------------------> No.%d  Shell on node: %s <------------------------------", i, output.NodeName)
			prefixStr := color.BlueString("[%s]", output.NodeName)
			for {
				line, err := output.StdOut.ReadString('\n')
				if err != nil || io.EOF == err {
					break
				}
				fmt.Printf("%s %s", prefixStr, line)
			}
			for {
				line, err := output.StdErr.ReadString('\n')
				if err != nil || io.EOF == err {
					break
				}
				fmt.Printf("%s %s", prefixStr, color.RedString(line))
			}
			break
		default:
			logrus.Error("不支持该格式输出")

		}

	}
}

// 返回目标节点（Node List）的shell pod列表
func getShellPodDict(kubeClientSet *kubernetes.Clientset, shellNodeName, shellNodeNameFile, toolName string) map[string]*v1.Pod {

	shellPods, _ := k8stools.GetPodList(kubeClientSet, toolName, "name="+toolName)

	// 根据参数判断需要传在那些Node执行Shell
	nodeList := make([]string, 0)
	// 所有Node执行Shell
	if strings.EqualFold(shellNodeName, "") && strings.EqualFold(shellNodeNameFile, "") {
		nodes, _ := kubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
		for _, node := range nodes.Items {
			nodeList = append(nodeList, node.Name)
		}
	}

	// 指定Node执行shell
	if !strings.EqualFold(shellNodeName, "") {
		nodeList = append(nodeList, shellNodeName)
	}

	// nodefile文件指定node执行shell
	if strings.EqualFold(shellNodeName, "") && !strings.EqualFold(shellNodeNameFile, "") {
		nodeList = goblin.ReadLine(shellNodeNameFile)
	}
	if len(nodeList) == 0 {
		panic("选择节点异常")
	}

	// 获取shellPod
	podTarges := make(map[string]*v1.Pod, 0)
	for _, nodename := range nodeList {
		podTarges[nodename] = nil
	}

	// 填充Shellpod
	for i, shellPod := range shellPods.Items {
		_, isOk := podTarges[shellPod.Spec.NodeName]
		if isOk {
			podTarges[shellPod.Spec.NodeName] = &shellPods.Items[i]
		}
	}
	return podTarges
}
