package alert

import (
	restful "github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"sdp-devops/pkg/alert/api"
)

func Main() {

	restful.DefaultContainer.Add(api.WebService())
	logrus.Printf("start listening on 0.0.0.0:8080")
	logrus.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
