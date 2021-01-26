package cleaner

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sdp-devops/pkg/util"
	"strings"
	"time"
)

type CleanNotification struct {
	Log    string    `json:"log"`
	Stream string    `json:"stream"`
	Time   time.Time `json:"time"`
}

func echo(filepath string) {

	cleanN := CleanNotification{
		Log:    CLEAN_MSG,
		Stream: "stderr",
		Time:   time.Now(),
	}
	cleanNStr, _ := json.Marshal(cleanN)
	cmd := exec.Command("/bin/sh", "-c", "echo "+strings.ReplaceAll(string(cleanNStr), "\"", "\\\"")+" > "+filepath)
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("清空文件失败：%s (%s)", filepath, err.Error())
	}
}

/**
清理目录，返回值表示该目录是否需要删除
*/
func cleanDir(dirpath string) bool {

	lastEditTime := time.Now().AddDate(0, 0, -logEditTimeLimit).Local()

	if err := os.Chdir(dirpath); err != nil {
		logrus.Errorf("目录不存在或者目录不可用：%s (%s)", dirpath, err.Error())
		return false
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		logrus.Errorf("查询目录信息失败：%s (%s)", dirpath, err.Error())
		return false
	}
	if len(files) == 0 {
		// 该目录是空目录可以删除
		return true
	}
	for _, file := range files {
		if file.Mode().IsRegular() {
			fileExt := path.Ext(file.Name())
			if file.ModTime().Before(lastEditTime) && IsDeleteFiles(file.Name()) {
				logrus.Infof("超出编辑时间限制，删除文件：%s (%s) ", path.Join(dirpath, file.Name()), file.ModTime())
				if err = os.Remove(path.Join(dirpath, file.Name())); err != nil {
					logrus.Errorf("删除文件失败：%s (%s)", file.Name(), err.Error())
				}
				continue
			}
			if file.Size() > logSizeLimit && logFileExt.Contains(fileExt) {
				logrus.Infof("文件超出大小限制，清除数据：%s (%s) ", path.Join(dirpath, file.Name()), util.FormatByte(file.Size()))
				echo(path.Join(dirpath, file.Name()))
				continue
			}
		}
		if file.Mode().IsDir() {
			isEmpty := cleanDir(path.Join(dirpath, file.Name()))

			if isEmpty && file.ModTime().Before(lastEditTime) {
				logrus.Infof("超出编辑时间限制，删除空白目录：%s (%s) ", path.Join(dirpath, file.Name()), file.ModTime())
				if err = os.Remove(path.Join(dirpath, file.Name())); err != nil {
					logrus.Errorf("删除目录失败：%s (%s)", file.Name(), err.Error())
				}
			}
		}
	}
	return false
}

/**
判断字符串是否符合正则表达试：

*/
func IsDeleteFiles(fname string) bool {

	fext := fname
	if indx := strings.Index(fname, "."); indx >= 0 {
		fext = fname[indx:]
	}

	if deleteFileExt.Contains(fext) {
		return true
	}

	isMatch := false
	for _, p := range regFileExt {
		isMatch, _ = regexp.MatchString(p, fext)
		if isMatch {
			break
		}
	}
	return isMatch

}

func cleanCrontab(crontabStr string) {

	crontab := cron.New()
	task := func() {
		logrus.Println("#######################################################################################")
		logrus.Printf("Start Clean Job:%s \n", time.Now())
		logrus.Println("#######################################################################################")
		cleanDir(tomcatLogDir)
		cleanContainerStdLog()

	}
	if _, err := crontab.AddFunc(crontabStr, task); err != nil {
		logrus.Errorf(err.Error())
		return
	}
	crontab.Start()
	select {}
}

func cleanContainerStdLog() {
	containerBase := path.Join(dockerBase, "containers")

	if err := os.Chdir(containerBase); err != nil {
		logrus.Errorf("目录不存在或者目录不可用：%s (%s)", containerBase, err.Error())
		return
	}
	containers, _ := ioutil.ReadDir(".")

	for _, c := range containers {
		if c.IsDir() {
			files, err := ioutil.ReadDir(c.Name())
			if err != nil {
				logrus.Errorf("读取容器(%s)日志目录失败：%s", c.Name(), err.Error())
				continue
			}
			for _, file := range files {
				if strings.EqualFold(file.Name(), c.Name()+"-json.log") && file.Size() > logSizeLimit {
					stdLogPath := path.Join(containerBase, c.Name(), file.Name())
					logrus.Infof("文件超出大小限制，清除数据：%s (%s) ", stdLogPath, util.FormatByte(file.Size()))
					echo(stdLogPath)
					break
				}
			}

		}
	}

}

/**
   清理当前节点的以下目录：
	1. tomcat 日志目录，默认为： /data/container_logs/
		- 大于 指定大小 的文件（echo方式清理）
		- 时间长于一定时间没有编辑文件 （直接删除）
	2. 容器标准输出日志，默认为： /data/var/lib/docker/containers/{container-id} （docker配置了回滚操作时，可以关闭该配置）
		- 大于 指定大小 的文件（echo方式清理）
*/
func Main() {
	if !isServer {
		cleanDir(tomcatLogDir)
		cleanContainerStdLog()
	} else {
		cleanCrontab(cronStr)
	}
}
