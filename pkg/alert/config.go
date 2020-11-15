package alert

import (
	"github.com/spf13/pflag"
)

var (
	apiToken     string
	falconServer string
	reciever     string

	areaMap = map[string]string{
		"CL":    "长乐环境",
		"AWSCA": "AWS 加利福尼亚环境",
		"AWSBH": "AWS 巴林环境",
		"HK":    "香港环境",
		"WX":    "无锡生产环境",
		"WXPRE": "无锡预生产环境",
		"VOD":   "阿里云VOD环境",
	}
)

const (
	alarmURL = "/v1/api/alarms"
)

func AddFlags(flags *pflag.FlagSet) {

	flags.StringVar(&apiToken, "apitoken", "c119a423a2c05da9703a76ab57fa4570f0fa655951da2", "tomcat日志的根目录。")
	flags.StringVar(&falconServer, "falcon-server", "alarm.falcon.cl.sdp:6063", "falcon服务端地址。")
	flags.StringVar(&reciever, "reciever", "", "电话告警通知。")
}
