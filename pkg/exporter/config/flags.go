package config

import "github.com/spf13/pflag"

var (
	Port               string
	KubeConfigStr      string
	TomcatLogDir       string
	MonitorDirectories string
	ExcludingCol       string
	IncludingCol       string
	MetricsURL         string
	DockerRootDir      string
)

func AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&Port, "port", "17000", "metrics服务端口。")
	flags.StringVar(&KubeConfigStr, "kubeconfig", "", "Kubernetes Config配置文件。")
	flags.StringVar(&TomcatLogDir, "tomcat-log-dir", "/data/container_logs", "tomcat 日志根目录。")
	flags.StringVar(&MonitorDirectories, "monitor-dirs", "/data", "要监控的磁盘目录，用逗号分隔。")
	flags.StringVar(&ExcludingCol, "excluding", "", "采集器名单黑名单，使用逗号分隔。")
	flags.StringVar(&IncludingCol, "including", "", "采集器名单白名单，使用逗号分隔。优先于黑名单")
	flags.StringVar(&MetricsURL, "url", "/metrics", "采集器地址。")
	flags.StringVar(&DockerRootDir, "docker-root-dir", "/data/var/lib/docker", "Docker Enginer 根目录。")
}
