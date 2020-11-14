package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sdp-devops/pkg/exporter/collector"
	"sdp-devops/pkg/exporter/config"
	"sdp-devops/pkg/util"
)

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
	if err := http.ListenAndServe("0.0.0.0:"+config.Port, nil); err != nil {
		util.Error.Printf("创建SDPCollector实例失败: %s", err)
	}
}
