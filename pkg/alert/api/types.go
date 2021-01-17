package api

import "fmt"

type AnnotationsMsg struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type Status struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  AnnotationsMsg    `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type Notify struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []Status          `json:"alerts"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
}

var (
	AreaNameMap = map[string]string{
		"CL":    "长乐环境",
		"AWSCA": "AWS 加利福尼亚环境",
		"AWSBH": "AWS 巴林环境",
		"HK":    "香港环境",
		"WX":    "无锡生产环境",
		"WXPRE": "无锡预生产环境",
		"VOD":   "VOD演练环境",
		"EGPRE": "VOD埃及预生产环境",
	}

	AlertNameMap = map[string]string{
		"HighErrorRate": "磁盘空间不足",
	}
)

func (u Notify) AlertNotifyMsg() string {

	areaName := AreaNameMap[u.CommonLabels["AREA"]]
	alertname := AlertNameMap[u.CommonLabels["alertname"]]
	msg := fmt.Sprintf("告警：%s %s ，请关注", areaName, alertname)
	return msg
}

type FalconPhoneAlert struct {
	Type      string   `json:"type"`
	Content   string   `json: "content"`
	Recievers []string `json: "recievers"`
}
