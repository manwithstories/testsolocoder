package utils

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(level, logDir string) {
	Logger = logrus.New()

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	Logger.SetLevel(parsedLevel)

	os.MkdirAll(logDir, 0755)
	logFile, err := os.OpenFile(
		logDir+"/app_"+time.Now().Format("20060102")+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}

	Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
