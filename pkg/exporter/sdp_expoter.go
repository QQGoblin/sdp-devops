package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"sdp-devops/pkg/exporter/collector"
	"sdp-devops/pkg/exporter/config"
)

// 启动 exporter server
func Main() {
	sc, _ := collector.NewNodeCollector()
	r := prometheus.NewRegistry()
	if err := r.Register(sc); err != nil {
		logrus.Errorf("创建SDPCollector实例失败: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{r},
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		},
	)
	http.Handle(config.MetricsURL, handler)
	if err := http.ListenAndServe("0.0.0.0:"+config.Port, nil); err != nil {
		logrus.Errorf("创建SDPCollector实例失败: %s", err)
	}
}
