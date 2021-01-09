package util

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// 初始化日志文件
func InitLogger() {
	logrus.SetOutput(os.Stdout)
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.WithFields(logrus.Fields{
		"sdp-devops": "Kubernetes",
	})

}

// 计算目录的大小
func CalDirSize(dirpath string) (dirsize int64, err error) {
	err = os.Chdir(dirpath)
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		dirsize += file.Size()
		if file.Mode().IsDir() {
			subDirSize, err2 := CalDirSize(path.Join(dirpath, file.Name()))
			if err2 == nil {
				dirsize += subDirSize
			}
		}
	}
	return
}

// 格式化字节单位到字符串
var (
	KB int64 = 1024
	MB       = KB * 1024
	GB       = MB * 1024
	TB       = GB * 1024
)

func FormatByte(b int64) string {

	if b >= TB {
		return strconv.FormatFloat(float64(b)/float64(TB), 'f', 2, 64) + " Ti"
	}
	if b >= GB {
		return strconv.FormatFloat(float64(b)/float64(GB), 'f', 2, 64) + " Gi"
	}
	if b >= MB {
		return strconv.FormatFloat(float64(b)/float64(MB), 'f', 2, 64) + " Mi"
	}
	if b >= KB {
		return strconv.FormatFloat(float64(b)/float64(KB), 'f', 2, 64) + " Ki"
	}

	return strconv.FormatFloat(float64(b), 'f', 2, 64)
}

// 格式化百分比字符串
func FormatPercentage(usage int64, total int64) string {
	percentage := float64(usage) / float64(total)
	var colorStr string
	if percentage > 0.8 {
		colorStr = color.HiRedString(strconv.FormatFloat(percentage*100, 'f', 2, 64) + "%")
	} else if percentage > 0.6 {
		colorStr = color.HiYellowString(strconv.FormatFloat(percentage*100, 'f', 2, 64) + "%")
	} else {
		colorStr = color.HiGreenString(strconv.FormatFloat(percentage*100, 'f', 2, 64) + "%")
	}
	return colorStr
}
