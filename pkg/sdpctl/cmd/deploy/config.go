package deploy

import "github.com/spf13/pflag"

var (
	action string
)

func AddNodeShellFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&action, "action", "a", "check", "相关操作名称。")
}
