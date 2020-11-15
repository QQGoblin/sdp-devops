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
	logrus.Println("method = ", r.Method) //请求方法
	logrus.Println("URL = ", r.URL)       // 浏览器发送请求文件路径
	logrus.Println("header = ", r.Header) // 请求头
	s, _ := ioutil.ReadAll(r.Body)
	logrus.Println("body = ", string(s)) // 请求包体
	logrus.Println(r.RemoteAddr, "连接成功") //客户端网络地址
	//PostPhoneAlert("恭喜发财", "15860837730")

}
