package sdpk8s

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sdp-devops/pkg/sdpctl/config"
	k8stools "sdp-devops/pkg/util/kubernetes"
	systools "sdp-devops/pkg/util/sys"
	"strings"
)

// 返回目标节点（Node List）的shell pod列表
func GetShellPodDict(kubeClientSet *kubernetes.Clientset) map[string]*v1.Pod {

	shellPods, _ := k8stools.GetPodList(kubeClientSet, config.ShellNamespace, "name="+config.ShellDaemonset)

	// 根据参数判断需要传在那些Node执行Shell
	nodeList := make([]string, 0)
	// 所有Node执行Shell
	if strings.EqualFold(config.ShellNodeName, "") && strings.EqualFold(config.ShellNodeNameFile, "") {
		nodes, _ := kubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
		for _, node := range nodes.Items {
			nodeList = append(nodeList, node.Name)
		}
	}

	// 指定Node执行shell
	if !strings.EqualFold(config.ShellNodeName, "") {
		nodeList = append(nodeList, config.ShellNodeName)
	}

	// nodefile文件指定node执行shell
	if strings.EqualFold(config.ShellNodeName, "") && !strings.EqualFold(config.ShellNodeNameFile, "") {
		nodeList = systools.ReadLine(config.ShellNodeNameFile)
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
