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
	config.LoadConfig()
	http.HandleFunc(config.GlobalExporterConfig.MetricsPath, handlerFunc)
	if err := http.ListenAndServe("0.0.0.0:"+config.GlobalExporterConfig.Port, nil); err != nil {
		logrus.Errorf("创建SDPCollector实例失败: %s", err)
	}
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {

	sc, _ := collector.NewCollector(r)
	registry := prometheus.NewRegistry()
	if err := registry.Register(sc); err != nil {
		logrus.Errorf("创建SDPCollector实例失败: %s", err)
	}

	h := promhttp.HandlerFor(
		prometheus.Gatherers{registry},
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		},
	)
	h.ServeHTTP(w, r)
}
