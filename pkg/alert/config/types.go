package config

type AlertConfig struct {
	AlertNameMap map[string]string `yaml:"alertNameMap"`
	AreaNameMap  map[string]string `yaml:"areaNameMap"`
	Falcon       FalconConfig      `yaml:"falcon"`
	Port         string            `yaml:"port"`
	WXWork       WXWorkConfig      `yaml:"wxWorkConfig"`
}

type FalconConfig struct {
	Server  string   `yaml:"server"`
	Token   string   `yaml:"token"`
	Numbers []string `yaml:"numbers"`
}

type WXWorkConfig struct {
	CorpId     string `yaml:"corpId"`
	CorpSecret string `yaml:"corpSecret"`
	AgentId    int    `yaml:"agentId"`
	ToParty    string `yaml:"toParty"`
}
