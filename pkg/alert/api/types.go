package api

import (
	"fmt"
	"sdp-devops/pkg/alert/config"
	"strings"
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

func (u Notify) Title() string {

	areaName := config.GlobalAlertConfig.AreaNameMap[u.CommonLabels["AREA"]]
	alertname := config.GlobalAlertConfig.AlertNameMap[u.CommonLabels["alertname"]]
	var msg string
	if strings.EqualFold(u.Status, "resolved") {
		msg = fmt.Sprintf("%s %s", areaName, "异常恢复")
	} else {
		msg = fmt.Sprintf("%s %s", areaName, alertname)
	}

	return msg
}

func (u Notify) Summary() string {
	summary := make([]string, len(u.Alerts))
	for i, s := range u.Alerts {
		summary[i] = s.Annotations.Summary
	}

	return strings.Join(summary, "\n")
}

type FalconPhoneAlert struct {
	Type      string   `json:"type"`
	Content   string   `json: "content"`
	Recievers []string `json: "recievers"`
}
