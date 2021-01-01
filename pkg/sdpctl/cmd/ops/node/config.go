package node

import "github.com/spf13/pflag"

var (
	action string
	area   string
)

func AddNodeFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&action, "action", "a", "install", "相关操作名称。")
}
