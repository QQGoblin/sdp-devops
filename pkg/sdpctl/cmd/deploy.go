package cmd

import (
	"github.com/modood/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sdp-devops/pkg/sdpctl/config"
	"sdp-devops/pkg/sdpctl/sdpk8s"
	k8stools "sdp-devops/pkg/util/kubernetes"
	"strings"
	"time"
)

func RunInstall(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := k8stools.KubeClientAndConfig(config.KubeConfigStr)

	if _, err := kubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).Get(config.ShellDaemonset, metav1.GetOptions{
		TypeMeta:        metav1.TypeMeta{},
		ResourceVersion: "",
	}); err == nil {
		PrintPodSimpleInfo(kubeClientSet, config.ShellNamespace, "name="+config.ShellDaemonset)
		return
	}

	shellNS := v1.Namespace{
		TypeMeta:   metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: config.ShellNamespace},
	}
	if _, err := kubeClientSet.CoreV1().Namespaces().Create(&shellNS); err != nil {
		logrus.Error("创建命名空间", config.ShellNamespace, "失败")
		panic(err.Error())
	}

	isPrivileged := true
	var priority int32 = 0
	nodeShellDSDefine := appv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "apps/v1",
			APIVersion: "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.ShellDaemonset,
			Namespace: config.ShellNamespace,
		},
		Spec: appv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": config.ShellDaemonset,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": config.ShellDaemonset,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "shell",
							Image:           "docker.io/alpine:3.9",
							ImagePullPolicy: "IfNotPresent",
							Command:         strings.Fields("nsenter"),
							Args:            strings.Fields("-t 1 -m -u -i -n sleep inf"),
							SecurityContext: &v1.SecurityContext{
								Privileged: &isPrivileged,
							},
							WorkingDir: "/root",
						},
					},
					DNSPolicy:     v1.DNSClusterFirst,
					HostIPC:       true,
					HostNetwork:   true,
					HostPID:       true,
					Priority:      &priority,
					RestartPolicy: v1.RestartPolicyAlways,
					Tolerations: []v1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: v1.TolerationOpExists,
							Effect:   v1.TaintEffectNoSchedule,
						},
						{
							Key:      "build",
							Value:    "type",
							Operator: v1.TolerationOpEqual,
							Effect:   v1.TaintEffectNoExecute,
						},
					},
				},
			},
		},
	}
	if _, err := kubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).Create(&nodeShellDSDefine); err != nil {
		logrus.Error("创建DaemonSet", nodeShellDSDefine.Name, "失败")
		panic(err.Error())
	}
	time.Sleep(5 * 1e9)
	PrintPodSimpleInfo(kubeClientSet, config.ShellNamespace, "name="+config.ShellDaemonset)
}

func RunClean(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := k8stools.KubeClientAndConfig(config.KubeConfigStr)

	if _, err := kubeClientSet.CoreV1().Namespaces().Get(config.ShellNamespace, metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
	}); err != nil {
		logrus.Info("清理命名空间", config.ShellNamespace, "成功！")
		return
	}
	if err := kubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{}); err != nil {
		logrus.Error("清理命名空间", config.ShellNamespace, "失败！")
		panic(err.Error())
	}

	if err := kubeClientSet.CoreV1().Namespaces().Delete(config.ShellNamespace, &metav1.DeleteOptions{}); err != nil {
		logrus.Error("清理命名空间", config.ShellNamespace, "失败！")
		panic(err.Error())
	}
	logrus.Info("清理命名空间", config.ShellNamespace, "成功！")
}

func PrintPodSimpleInfo(kubeClientSet *kubernetes.Clientset, namespace, lableSelector string) {

	pods, err := kubeClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: lableSelector,
	})
	if err != nil {
		panic(err.Error())
	}
	podInfoList := make([]sdpk8s.PodBriefInfo, len(pods.Items))
	for i := 0; i < len(pods.Items); i++ {
		podInfo := sdpk8s.PodBriefInfo{
			Name:      pods.Items[i].Name,
			NameSpace: config.ShellNamespace,
			Status:    string(pods.Items[i].Status.Phase),
			Node:      pods.Items[i].Status.HostIP,
		}
		podInfoList[i] = podInfo

	}
	table.Output(podInfoList)
}

func NewCmdInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "install",
		Short:                 "安装Shell Pod服务",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunInstall(cmd, args)
		},
	}
	return cmd
}

func NewCmdClean() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "clean",
		Short:                 "清理Shell Pod服务",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunClean(cmd, args)
		},
	}
	return cmd
}

func NewCmdDeploy() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "deploy",
		Short:                 "在集群部署基于Daemonset的Shell客户端",
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(NewCmdInstall())
	cmd.AddCommand(NewCmdClean())
	return cmd
}
