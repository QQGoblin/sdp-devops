package api

import (
	"fmt"
	"sdp-devops/pkg/alert/config"
)

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

func (u Notify) AlertNotifyMsg() string {

	areaName := config.GlobalAlertConfig.AreaNameMap[u.CommonLabels["AREA"]]
	alertname := config.GlobalAlertConfig.AlertNameMap[u.CommonLabels["alertname"]]
	msg := fmt.Sprintf("告警：%s %s ，请关注", areaName, alertname)
	return msg
}

type FalconPhoneAlert struct {
	Type      string   `json:"type"`
	Content   string   `json: "content"`
	Recievers []string `json: "recievers"`
}