package collector

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sdp-devops/pkg/exporter/config"
	"sync"
	"time"
)

type probeCollector struct {
	code *prometheus.Desc // 返回的状态码
}

const (
	probeCollectorSubsystem = "probe"
)

func init() {
	registerCollector(probeCollectorSubsystem, NewProbeCollector)
}

// 创建磁盘采集器
func NewProbeCollector() (Collector, error) {
	return &probeCollector{
		code: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, probeCollectorSubsystem, "probe_http_status_code"),
			"拨测状态码",
			[]string{"node", "target", "type"},
			nil,
		),
	}, nil
}

// 创建 HTTP Client
func InitHttpClient() *http.Client {

	tlsCerts := make([]tls.Certificate, 0)
	for _, t := range config.GetProbeHttpStatusCode().TlsConfig.X509KeyPair {
		tlsCert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
		if err != nil {
			logrus.Error(errors.Wrapf(err, "加载tls证书异常：%s（%s，%s）", t.Name, t.CertFile, t.KeyFile))
			continue
		}
		tlsCerts = append(tlsCerts, tlsCert)
	}
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       tlsCerts,
		},
	}
	httpClient := &http.Client{
		Timeout:   time.Duration(config.GetProbeHttpStatusCode().TimeOutSec) * time.Second,
		Transport: tr,
	}
	return httpClient
}

// 实现采集接口
func (c *probeCollector) Update(ch chan<- prometheus.Metric) error {

	nodename, _ := os.Hostname()

	httpClient := InitHttpClient()
	wg := sync.WaitGroup{}
	wg.Add(len(config.GetProbeHttpStatusCode().Service))
	for _, t := range config.GetProbeHttpStatusCode().Service {
		go func(service config.Service) {
			c.ProbeHTTP(httpClient, nodename, service, ch)
			wg.Done()
		}(t)
	}
	wg.Wait()
	return nil
}

func (c *probeCollector) ProbeHTTP(httpClient *http.Client, nodename string, service config.Service, ch chan<- prometheus.Metric) {

	req, _ := http.NewRequest("GET", service.TargetURL, nil)
	resp, _ := httpClient.Do(req)
	if resp != nil {
		ch <- prometheus.MustNewConstMetric(c.code, prometheus.GaugeValue, float64(resp.StatusCode), nodename, service.TargetURL, service.Name)
	} else {
		ch <- prometheus.MustNewConstMetric(c.code, prometheus.GaugeValue, float64(0), nodename, service.TargetURL, service.Name)
	}
}
