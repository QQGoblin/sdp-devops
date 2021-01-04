package sdpk8s

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8stools "sdp-devops/pkg/util/kubernetes"
	systools "sdp-devops/pkg/util/sys"
	"strconv"
	"strings"
)

// 返回目标节点（Node List）的shell pod列表
func GetShellPodDict(kubeClientSet *kubernetes.Clientset, shellNodeName, shellNodeNameFile, toolName string) map[string]*v1.Pod {

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
		nodeList = systools.ReadLine(shellNodeNameFile)
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

// 根据 K8S 容器的命名规则从 docker 容器名称解析对应的Pod信息
// 映射规则如下：
// k8s_POD_${POD_FULL_NAME}_${NAMESPACES}_${POD_ID}_${Restart_Count}
// k8s_${Container_NAME}_${POD_FULL_NAME}_${NAMESPACES}_${POD_ID}_${Restart_Count}
//
func GetPodInfo(cli *client.Client) []PodBriefInfo {

	podBriefInfoList := make([]PodBriefInfo, 0)
	if containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{}); err != nil {
		logrus.Error(err)
	} else {
		for _, container := range containers {
			if strings.EqualFold(container.State, "running") && strings.HasPrefix(container.Names[0], "/k8s_POD_") {
				infoStr := strings.Replace(container.Names[0], "/k8s_POD_", "", -1)
				info := strings.Split(infoStr, "_")
				podBriefInfo := PodBriefInfo{
					Name:        info[0],
					UID:         info[2],
					NameSpace:   info[1],
					Status:      container.State,
					Node:        "local",
					NetworkMode: container.HostConfig.NetworkMode,
					Restart:     0,
					PID:         1,
				}
				if restartCount, err := strconv.Atoi(info[3]); err != nil {
					logrus.Error(err)
					continue
				} else {
					podBriefInfo.Restart = restartCount
				}
				if containerInspect, err := cli.ContainerInspect(context.Background(), container.ID); err == nil {
					podBriefInfo.PID = containerInspect.State.Pid
				} else {
					logrus.Error(err)
					continue
				}
				podBriefInfoList = append(podBriefInfoList, podBriefInfo)
			}
		}
	}
	return podBriefInfoList
}
