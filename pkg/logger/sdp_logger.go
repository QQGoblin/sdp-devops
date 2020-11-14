package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{}
	//log.SetLevel(logrus.InfoLevel)
	log.WithFields(logrus.Fields{
		"sdp-devops": "Kubernetes",
	})
	logrus.Info()
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Infof(format string, args ...interface{}) {
	log.Info(args)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args)
}

func Println(args ...interface{}) {
	log.Println(args)
}
