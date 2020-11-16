package kubernetes

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sdp-devops/pkg/sdpctl/config"
	systools "sdp-devops/pkg/util/sys"
	"strings"
)

// 根据Node返回当前Pod运行字典
func GetPodDict(kubeClientSet *kubernetes.Clientset, lableSelector string) (podDist map[string][]v1.Pod, err error) {
	listOptions := metav1.ListOptions{}
	if !strings.EqualFold(lableSelector, "") {
		listOptions = metav1.ListOptions{
			TypeMeta:      metav1.TypeMeta{},
			LabelSelector: lableSelector,
		}
	}

	pods, err := kubeClientSet.CoreV1().Pods("").List(listOptions)
	podDist = make(map[string][]v1.Pod)
	for _, pod := range pods.Items {
		key := pod.Spec.NodeName
		podListOnNode := podDist[key]
		if podListOnNode == nil {
			podListOnNode = make([]v1.Pod, 0)
		}
		podListOnNode = append(podListOnNode, pod)
		podDist[key] = podListOnNode
	}
	return
}

// 返回指定Pod列表
func GetPodList(kubeClientSet *kubernetes.Clientset, namespaceStr, lableSelector string) (pods *v1.PodList, err error) {

	listOptions := metav1.ListOptions{}
	if !strings.EqualFold(lableSelector, "") {
		listOptions = metav1.ListOptions{
			TypeMeta:      metav1.TypeMeta{},
			LabelSelector: lableSelector,
		}
	}
	pods, err = kubeClientSet.CoreV1().Pods(namespaceStr).List(listOptions)
	if err != nil {
		logrus.Error("获取Pod列表异常")
	}
	return
}

// 返回目标节点（Node List）的shell pod列表
func GetShellPodDict(kubeClientSet *kubernetes.Clientset) map[string]*v1.Pod {

	shellPods, _ := GetPodList(kubeClientSet, config.ShellNamespace, "name="+config.ShellDaemonset)

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
