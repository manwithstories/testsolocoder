package utils

import (
	"os"
	"time"

	"meeting-room/internal/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	lj := &lumberjack.Logger{
		Filename:   config.Cfg.Log.Filename,
		MaxSize:    config.Cfg.Log.MaxSize,
		MaxBackups: config.Cfg.Log.MaxBackups,
		MaxAge:     config.Cfg.Log.MaxAge,
	}

	Logger.SetOutput(lj)

	level, err := logrus.ParseLevel(config.Cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}

func LogDir() {
	os.MkdirAll("./logs", 0755)
	os.MkdirAll(config.Cfg.Upload.Dir, 0755)
}
