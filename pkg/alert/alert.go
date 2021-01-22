package alert

import (
	restful "github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"sdp-devops/pkg/alert/api"
	"sdp-devops/pkg/alert/config"
)

func Main() {

	config.LoadConfig()
	restful.DefaultContainer.Add(api.WebService())
	logrus.Printf("start listening on 0.0.0.0:%s", config.GlobalAlertConfig.Port)
	logrus.Fatal(http.ListenAndServe("0.0.0.0:"+config.GlobalAlertConfig.Port, nil))
}
