package alert

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const PhoneAlertType = "phone"
const PhoneAlertURL = "/v1/alert/phone"

type PhoneAlert struct {
	Type      string   `json:"type"`
	Content   string   `json: "content"`
	Recievers []string `json: "recievers"`
}

type PhoneAlertResponse struct {
	Code string `json:"code"`
}

type AnnotationsMsg struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type AlerStatus struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  AnnotationsMsg    `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type AlertNotifyMsg struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []AlerStatus      `json:"alerts"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
}

func PostPhoneAlert(content, reciever string) {

	url := "http://" + falconServer + alarmURL
	contentType := "application/json"
	alert := PhoneAlert{
		Type:      PhoneAlertType,
		Content:   content,
		Recievers: []string{reciever},
	}
	body, _ := json.Marshal(alert)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		logrus.Errorf("构造告警Request失败（%s）。", err.Error())
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Api-Token", apiToken)

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("发送电话告警信息失败（%s），Respone Code (%d) %s。", content, resp.StatusCode, err.Error())
		return
	}
	respone := PhoneAlertResponse{}
	json.Unmarshal(rbody, &respone)

	if strings.EqualFold(respone.Code, "ok") {
		logrus.Infof("发送电话告警信息:（%s）-- %s。", content, reciever)
	} else {
		logrus.Errorf("发送电话告警信息:（%s）-- %s，服务端异常（%s）。", content, reciever, string(rbody))
	}
	return

}

func PhoneAlertHandler(w http.ResponseWriter, r *http.Request) {

	s, _ := ioutil.ReadAll(r.Body)
	var alertNotifyMsg AlertNotifyMsg
	if err := json.Unmarshal(s, &alertNotifyMsg); err != nil {
		logrus.Println(err.Error())
	}

	logrus.Infof("Alerting 告警信息：%s", alertNotifyMsg)

	area := alertNotifyMsg.CommonLabels["AREA"]
	alertname := alertNotifyMsg.CommonLabels["alertname"]
	var content string
	switch alertname {
	case "HighErrorRate":
		content = areaMap[area] + "，磁盘使用率超过80%"
	default:
		content = areaMap[area] + "，" + alertname
	}
	if !strings.EqualFold(reciever, "") {
		PostPhoneAlert(content, reciever)
	}

}
