package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"net/http"
	"sdp-devops/pkg/exporter/collector"
	"sdp-devops/pkg/util"
)

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

// 启动 exporter server
func Main() {
	sc, _ := collector.NewNodeCollector()
	r := prometheus.NewRegistry()
	if err := r.Register(sc); err != nil {
		util.Error.Printf("创建SDPCollector实例失败: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{r},
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		},
	)
	http.Handle("/metrics", handler)
	if err := http.ListenAndServe("0.0.0.0:"+Port, nil); err != nil {
		util.Error.Printf("创建SDPCollector实例失败: %s", err)
	}
}
