package api

import (
	"encoding/json"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"sdp-devops/pkg/alert/config"
)

func WebService() *restful.WebService {
	ws := new(restful.WebService)

	ws.Path("/v1/alert").Consumes(restful.MIME_XML, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.POST("/falcon").To(alertFalcon))
	ws.Route(ws.POST("/wx").To(alertFalcon))

	return ws
}

func alertFalcon(request *restful.Request, response *restful.Response) {

	var alertNotify Notify
	if err := request.ReadEntity(&alertNotify); err != nil {
		response.WriteErrorString(http.StatusBadRequest, errors.Wrap(err, "读取告警信息失败").Error())
		return
	}
	alertNotifyRaw, _ := json.Marshal(alertNotify)
	logrus.Infof("Alerting 告警信息：%s", string(alertNotifyRaw))

	// 通过Falcon 发送电话告警
	client := resty.New()
	resp, err := client.R().
		SetBody(FalconPhoneAlert{
			Type:      "phone",
			Content:   alertNotify.AlertNotifyMsg(),
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
	// TODO：发送微信告警信息
}
