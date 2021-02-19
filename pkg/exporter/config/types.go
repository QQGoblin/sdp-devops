package config

type ExporterConfig struct {
	Port        string    `yaml:"port"`
	MetricsPath string    `yaml:"metricsPath"`
	Collector   Collector `yaml:"collector"`
}

type Collector struct {
	Exclude []string        `yaml:"exclude"`
	Include []string        `yaml:"Include"`
	Config  CollectorConfig `yaml:"config"`
}

type CollectorConfig struct {
	DiskUseCheck        DiskUseCheck        `yaml:"diskUseCheck"`
	ProbeHttpStatusCode ProbeHttpStatusCode `yaml:"probeHttpStatusCode"`
}

type DiskUseCheck struct {
	Monitor []string `yaml:"monitor"`
}

type X509KeyPair struct {
	Name     string `yaml:"name"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type TlsConfig struct {
	X509KeyPair        []X509KeyPair `yaml:"x509KeyPair"`
	InsecureSkipVerify bool          `yaml:"insecureSkipVerify"`
}

type Service struct {
	TargetURL    string   `yaml:"targetURL"`
	Name         string   `yaml:"name"`
	NodeSelector []string `yaml:"nodeSelector"`
}

type ProbeHttpStatusCode struct {
	TimeOutSec int       `yaml:"timeOutSec"`
	TlsConfig  TlsConfig `yaml:"tlsConfig"`
	Service    []Service `yaml:"service"`
}
