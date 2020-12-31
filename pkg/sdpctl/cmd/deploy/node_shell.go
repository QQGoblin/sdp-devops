package deploy

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sdp-devops/pkg/sdpctl/config"
	k8stools "sdp-devops/pkg/util/kubernetes"
	"strings"
)

func install(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := k8stools.KubeClientAndConfig(config.KubeConfigStr)

	if _, err := kubeClientSet.AppsV1().DaemonSets(config.ShellToolName).Get(config.ShellToolName, metav1.GetOptions{
		TypeMeta:        metav1.TypeMeta{},
		ResourceVersion: "",
	}); err == nil {
		logrus.Info("服务已经部署")
		return
	}

	shellNS := v1.Namespace{
		TypeMeta:   metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: config.ShellToolName},
	}
	if _, err := kubeClientSet.CoreV1().Namespaces().Create(&shellNS); err != nil {
		logrus.Error("创建命名空间", config.ShellToolName, "失败")
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
			Name:      config.ShellToolName,
			Namespace: config.ShellToolName,
		},
		Spec: appv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": config.ShellToolName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": config.ShellToolName,
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
	if _, err := kubeClientSet.AppsV1().DaemonSets(config.ShellToolName).Create(&nodeShellDSDefine); err != nil {
		logrus.Error("创建DaemonSet", nodeShellDSDefine.Name, "失败")
		panic(err.Error())
	}
}

func clean(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := k8stools.KubeClientAndConfig(config.KubeConfigStr)

	if _, err := kubeClientSet.CoreV1().Namespaces().Get(config.ShellToolName, metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
	}); err != nil {
		logrus.Info("清理命名空间", config.ShellToolName, "成功！")
		return
	}
	if err := kubeClientSet.AppsV1().DaemonSets(config.ShellToolName).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{}); err != nil {
		logrus.Error("清理命名空间", config.ShellToolName, "失败！")
		panic(err.Error())
	}

	if err := kubeClientSet.CoreV1().Namespaces().Delete(config.ShellToolName, &metav1.DeleteOptions{}); err != nil {
		logrus.Error("清理命名空间", config.ShellToolName, "失败！")
		panic(err.Error())
	}
	logrus.Info("清理命名空间", config.ShellToolName, "成功！")
}
