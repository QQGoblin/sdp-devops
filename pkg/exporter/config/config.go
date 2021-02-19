package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	ConfigPath           string
	GlobalExporterConfig ExporterConfig
)

func AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&ConfigPath, "config", "/etc/sdp/exporter.conf", "sdp-exporter服务配置文件地址。")
}

func LoadConfig() {
	yamlFile, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		logrus.Error(errors.Wrapf(err, "读取配置文件失败：%s", ConfigPath))
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &GlobalExporterConfig)
	if err != nil {
		logrus.Error(errors.Wrapf(err, "配置文件格式异常：%s", ConfigPath))
		panic(err)
	}
}

func GetDiskUseCheck() *DiskUseCheck {
	return &GlobalExporterConfig.Collector.Config.DiskUseCheck
}

func GetProbeHttpStatusCode() *ProbeHttpStatusCode {
	return &GlobalExporterConfig.Collector.Config.ProbeHttpStatusCode
}
