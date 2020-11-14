package config

import "github.com/spf13/pflag"

var (
	Port               string
	KubeConfigStr      string
	TomcatLogDir       string
	MonitorDirectories string
)

func AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&Port, "port", "17000", "metrics服务端口。")
	flags.StringVar(&KubeConfigStr, "kubeconfig", "", "Kubernetes Config配置文件。")
	flags.StringVar(&TomcatLogDir, "tomcat-log-dir", "/data/container_logs", "tomcat 日志根目录。")
	flags.StringVar(&MonitorDirectories, "monitor-dirs", "/data", "要监控的磁盘目录，用逗号分隔。")

}
