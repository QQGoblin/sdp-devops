package shell_tool

import "github.com/spf13/pflag"

var (
	action string
)

func AddNodeShellFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&action, "action", "a", "install", "相关操作名称。")
}
