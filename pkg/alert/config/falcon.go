package config

import (
	"github.com/spf13/pflag"
	"strings"
)

var (
	FalconAPIToken         string
	FalconServer           string
	FalconNotifyNumbers    []string
	FalconNotifyNumbersStr string
)

func AddFalconFlags(flags *pflag.FlagSet) {

	flags.StringVar(&FalconAPIToken, "falcon-api-token", "c119a423a2c05da9703a76ab57fa4570f0fa655951da2", "Falcon 告警服务API Token。")
	flags.StringVar(&FalconServer, "falcon-server", "alarm.falcon.cl.sdp:6063", "falcon服务端地址。")
	flags.StringVar(&FalconNotifyNumbersStr, "falcon-notify-numbers", "", "告警通知的电话号码，用逗号分割。")
	FalconNotifyNumbers = strings.Split(FalconNotifyNumbersStr, ",")

}
