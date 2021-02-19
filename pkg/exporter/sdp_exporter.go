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

	sc, _ := collector.NewCollector()
	registry := prometheus.NewRegistry()
	registry.MustRegister(sc)

	http.HandleFunc(config.GlobalExporterConfig.MetricsPath, func(w http.ResponseWriter, r *http.Request) {
		sc.Params = r.URL.Query()
		h := promhttp.HandlerFor(
			prometheus.Gatherers{registry},
			promhttp.HandlerOpts{
				ErrorHandling: promhttp.ContinueOnError,
			},
		)
		h.ServeHTTP(w, r)
	})
	if err := http.ListenAndServe("0.0.0.0:"+config.GlobalExporterConfig.Port, nil); err != nil {
		logrus.Errorf("创建SDPCollector实例失败: %s", err)
	}
}
