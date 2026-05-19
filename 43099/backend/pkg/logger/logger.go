package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}
