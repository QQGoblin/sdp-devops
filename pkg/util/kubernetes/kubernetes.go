package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

// 获取Kube Client 客户端
func KubeClientByConfig(configStr string) *kubernetes.Clientset {
	var err error
	var config *restclient.Config
	if strings.EqualFold(configStr, "") {
		config, err = restclient.InClusterConfig()
	} else {
		// 获取kube配置对象
		config, err = clientcmd.BuildConfigFromFlags("", configStr)
	}
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset

}

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
