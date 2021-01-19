package config

type AlertConfig struct {
	AlertNameMap map[string]string `yaml:"alertNameMap"`
	AreaNameMap  map[string]string `yaml:"areaNameMap"`
	Falcon       FalconConfig      `yaml:"falcon"`
	Port         string            `yaml:"port"`
}

type FalconConfig struct {
	Server  string   `yaml:"server"`
	Token   string   `yaml:"token"`
	Numbers []string `yaml:"numbers"`
}
