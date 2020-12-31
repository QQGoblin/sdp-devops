package deploy

import "github.com/spf13/pflag"

var (
	deployAction string
)

func AddDeployFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&deployAction, "deploy-action", "a", "check", "相关操作名称。")
}
