package api

import (
	"encoding/json"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"sdp-devops/pkg/alert/config"
	"sdp-devops/pkg/util/wxwork"
)

var wxWorkClient *wxwork.WXWorkClient

func WebService() *restful.WebService {
	ws := new(restful.WebService)

	ws.Path("/v1/alert")
	ws.Route(ws.POST("/falcon").To(alertFalcon).Consumes())
	ws.Route(ws.POST("/wx").To(alertWX))
	return ws
}

func parseAlertNotify(request *restful.Request, response *restful.Response) *Notify {
	var alertNotify Notify
	if err := request.ReadEntity(&alertNotify); err != nil {
		response.WriteErrorString(http.StatusBadRequest, errors.Wrap(err, "读取告警信息失败").Error())
		return nil
	}
	alertNotifyRaw, _ := json.Marshal(alertNotify)
	logrus.Infof("Alerting 告警信息：%s", string(alertNotifyRaw))
	return &alertNotify
}

func alertFalcon(request *restful.Request, response *restful.Response) {

	alertNotify := parseAlertNotify(request, response)
	if alertNotify == nil {
		return
	}
	// 通过Falcon 发送电话告警
	client := resty.New()
	resp, err := client.R().
		SetBody(FalconPhoneAlert{
			Type:      "phone",
			Content:   alertNotify.Title(),
			Recievers: config.GlobalAlertConfig.Falcon.Numbers,
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Api-Token", config.GlobalAlertConfig.Falcon.Token).
		Post("http://" + config.GlobalAlertConfig.Falcon.Server + "/v1/api/alarms")

	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, errors.Wrap(err, "发送Falcon电话告警失败").Error())
		return
	} else {
		logrus.Infof("发送Falcon电话告警通知：%s", resp.String())
	}

}

func alertWX(request *restful.Request, response *restful.Response) {

	alertNotify := parseAlertNotify(request, response)
	if alertNotify == nil {
		return
	}
	// 发送微信告警信息
	if wxWorkClient == nil {
		wxWorkClient = wxwork.New(config.GlobalAlertConfig.WXWork.CorpId, config.GlobalAlertConfig.WXWork.CorpSecret)
	}

	textCradMsg := wxwork.TextCardMsg{
		Title:       alertNotify.Title(),
		Description: alertNotify.Summary(),
		Url:         "http://grafana.k8s.101.com",
	}
	wxWorkClient.SendMsg(
		config.GlobalAlertConfig.WXWork.AgentId,
		config.GlobalAlertConfig.WXWork.ToParty,
		textCradMsg,
	)
}
