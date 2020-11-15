package alert

import (
	"github.com/spf13/pflag"
)

var (
	apiToken     string
	falconServer string
)

const (
	alarmURL = "/v1/api/alarms"
)

func AddFlags(flags *pflag.FlagSet) {

	flags.StringVar(&apiToken, "apitoken", "c119a423a2c05da9703a76ab57fa4570f0fa655951da2", "tomcat日志的根目录。")
	flags.StringVar(&falconServer, "falcon-server", "alarm.falcon.cl.sdp:6063", "falcon服务端地址。")

}
