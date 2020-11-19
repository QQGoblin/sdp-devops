package cleaner

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/pflag"
	unit "sdp-devops/pkg/util/metrics"
)

var (
	tomcatLogDir     string
	dockerBase       string
	logEditTimeLimit int   // 上次编辑时间，单位天
	logSizeLimit     int64 // 文件大小限制，单位byte
	isServer         bool
	cronStr          string

	logFileExt = mapset.NewSet(
		".log",
		".out",
	)
	deleteFileExt, regFileExt = mapset.NewSet(
		".gz",
		".log",
		".out",
	), []string{
		"^\\.[0-9]{4}-(((0[13578]|(10|12))-(0[1-9]|[1-2][0-9]|3[0-1]))|(02-(0[1-9]|[1-2][0-9]))|((0[469]|11)-(0[1-9]|[1-2][0-9]|30)))$",
	}
)

const (
	CLEAN_MSG = "===================== Clean by SDP-Cleaner! This file is too larger! =====================\n"
)

func AddFlags(flags *pflag.FlagSet) {

	flags.StringVar(&tomcatLogDir, "tomcat-log-dir", "/data/container_logs", "tomcat日志的根目录。")
	flags.StringVar(&dockerBase, "container-base", "/data/var/lib/docker", "Docker服务的根目录。")
	flags.IntVar(&logEditTimeLimit, "last-edit-time", 30, "最晚文件编辑时间。")
	flags.Int64Var(&logSizeLimit, "max-size", 1*unit.GB, "最晚文件编辑时间。")
	flags.BoolVar(&isServer, "server", false, "启动定时清理服务。")
	flags.StringVar(&cronStr, "cron", "0 * * * *", "定时清理crontab配置。")
}
