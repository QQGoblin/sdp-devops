package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/toolkits/nux"
	"os"
	"sdp-devops/pkg/exporter/config"
	"sdp-devops/pkg/logger"
	"strings"
)

type diskCollector struct {
	used     *prometheus.Desc // 已使用空间
	capacity *prometheus.Desc // 磁盘容量
	isMount  *prometheus.Desc // 是否磁盘独立挂载
}

const (
	diskCollectorSubsystem = "disk"
)

func init() {
	registerCollector(diskCollectorSubsystem, NewDiskCollector)
}

// 创建磁盘采集器
func NewDiskCollector() (Collector, error) {
	return &diskCollector{
		used: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "used_size"),
			"目录使用容量",
			[]string{"node", "directory"},
			nil,
		),
		capacity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "capacity_size"),
			"磁盘容量",
			[]string{"node", "directory"},
			nil,
		),
		isMount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "ismount"),
			"磁盘独立挂载",
			[]string{"node", "directory"},
			nil,
		),
	}, nil
}

// 实现采集接口
func (c *diskCollector) Update(ch chan<- prometheus.Metric) error {

	nodename, _ := os.Hostname()
	monitorStr := config.MonitorDirectories
	if strings.EqualFold(monitorStr, "") {
		logrus.Errorf("请指定磁盘监控目录")
		return nil
	}
	directories := strings.Split(monitorStr, ",")
	for _, directory := range directories {
		dir := strings.ReplaceAll(directory[1:], "/", "-")
		isMount := -1
		_, isNotExist := os.Stat(directory)
		if isNotExist != nil {
			logrus.Errorf("%s 目录不存在(%s)。\n", directory, isNotExist.Error())
		} else {
			if L, err := nux.ListMountPoint(); err != nil {
				logrus.Errorf("获取系统挂载信息失败(%s)。", err.Error())
			} else {
				for _, arr := range L {
					if strings.EqualFold(directory, arr[1]) {
						isMount = 1
						deviceUsage, _ := nux.BuildDeviceUsage(arr[0], arr[1], arr[2])
						ch <- prometheus.MustNewConstMetric(c.capacity, prometheus.GaugeValue, float64(deviceUsage.BlocksAll), nodename, dir)
						ch <- prometheus.MustNewConstMetric(c.used, prometheus.GaugeValue, float64(deviceUsage.BlocksUsed), nodename, dir)
						break
					}
				}
			}
		}
		ch <- prometheus.MustNewConstMetric(c.isMount, prometheus.GaugeValue, float64(isMount), nodename, dir)
	}
	return nil
}
