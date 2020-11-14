package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
	"os"
	"path"
	"sdp-devops/pkg/exporter"
	"sdp-devops/pkg/util"
	dockertools "sdp-devops/pkg/util/docker"
	k8stools "sdp-devops/pkg/util/kubernetes"
	systools "sdp-devops/pkg/util/sys"
	"strings"
)

type containerCollector struct {
	containerSize *prometheus.Desc
	logSize       *prometheus.Desc
	userLogSize   *prometheus.Desc
}

const (
	containerCollectorSubsystem = "container"
)

func init() {
	registerCollector(containerCollectorSubsystem, NewContainerCollector)
}

// 创建容器采集器
func NewContainerCollector() (Collector, error) {
	return &containerCollector{
		containerSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, containerCollectorSubsystem, "container_size"),
			"Pod容器运行时（Diff）占用的磁盘空间大小",
			[]string{"pod", "namespace", "node"},
			nil,
		),
		logSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, containerCollectorSubsystem, "container_log_std_size"),
			"Docker标准输出日志占用的磁盘空间大小",
			[]string{"pod", "namespace", "node"},
			nil,
		),
		userLogSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, containerCollectorSubsystem, "container_log_tomcat_size"),
			"用户Tomcat日志占用的磁盘空间大小",
			[]string{"pod", "namespace", "node"},
			nil,
		),
	}, nil
}

// 实现采集接口
func (c *containerCollector) Update(ch chan<- prometheus.Metric) error {
	k8scli := k8stools.KubeClientByConfig(exporter.KubeConfigStr)
	dockercli := dockertools.DockerClient("")
	podDict, err := k8stools.GetPodDict(k8scli, "")
	if err != nil {
		panic(err.Error())
	}
	nodename, _ := os.Hostname()
	pods := podDict[nodename]
	if pods == nil || len(pods) == 0 {
		util.Warning.Printf("%s 不是Kubernetes集群的节点或者该节点没有Pod运行\n", nodename)
		return nil
	}

	for _, pod := range pods {
		if pod.Status.Phase != v1.PodRunning {
			util.Warning.Printf("%s 没有正常运行，状态：%s\n", pod.Name, pod.Status.Phase)
			continue
		}
		var containerSize int64
		var dockerLogSize int64
		for _, container := range pod.Status.ContainerStatuses {
			containerId := strings.Replace(container.ContainerID, "docker://", "", -1)
			dockerLogSize += dockertools.ContainerLogSize(containerId, dockercli)
			containerSize += dockertools.ContainerSize(containerId, dockercli)
		}
		tomcatLogDirPath := path.Join(exporter.TomcatLogDir, pod.Name)

		_, isExist := os.Stat(tomcatLogDirPath)
		var tomcatLogSize int64
		if isExist == nil {
			tomcatLogSize, _ = systools.CalculateDirSize(tomcatLogDirPath)
		}
		ch <- prometheus.MustNewConstMetric(c.containerSize, prometheus.GaugeValue, float64(containerSize), pod.Name, pod.Namespace, nodename)
		ch <- prometheus.MustNewConstMetric(c.logSize, prometheus.GaugeValue, float64(dockerLogSize), pod.Name, pod.Namespace, nodename)
		ch <- prometheus.MustNewConstMetric(c.userLogSize, prometheus.GaugeValue, float64(tomcatLogSize), pod.Name, pod.Namespace, nodename)
	}

	return nil
}
