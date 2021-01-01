package node

type NodeConfig struct {
	Env     int           `json:"env" bson:"env"`
	Area    string        `json:"area" bson:"area"`
	Skydns  []string      `json:"skydns" bson:"skydns"`
	Salt    SaltConfig    `json:"salt" bson:"salt"`
	Falcon  FalconConfig  `json:"falcon" bson:"falcon"`
	NetWork NetWorkConfig `json:"net_work" bson:"net_work"`
	Kubelet KubeletConfig `json:"kubelet" bson:"kubelet"`
}

type SaltConfig struct {
	Master       string `json:"master" bson:"master"`
	MinionSuffix string `json:"minion_suffix" bson:"minion_suffix"`
}

type FalconConfig struct {
	EndpointPrefix string `json:"endpoint_prefix" bson:"endpoint_prefix"`
}
type NetWorkConfig struct {
	IdcInnerSegmet string `json:"idc_inner_segmet" bson:"idc_inner_segmet"`
}
type KubeletConfig struct {
	K8sPodSubnet     string `json:"k8s_pod_subnet" bson:"k8s_pod_subnet"`
	K8sServiceSubnet string `json:"k8s_service_subnet" bson:"k8s_service_subnet"`
}
