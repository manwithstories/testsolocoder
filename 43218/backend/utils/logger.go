package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"secondhand-platform/config"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(cfg *config.LogConfig) {
	Logger = logrus.New()

	switch cfg.Level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	if cfg.Format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	var writers []io.Writer

	if cfg.Output == "stdout" || cfg.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	if cfg.Output == "file" || cfg.Output == "both" {
		logDir := "./logs"
		os.MkdirAll(logDir, 0755)
		logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			writers = append(writers, file)
		}
	}

	if len(writers) > 0 {
		Logger.SetOutput(io.MultiWriter(writers...))
	} else {
		Logger.SetOutput(os.Stdout)
	}

	logrus.SetOutput(Logger.Out)
	logrus.SetFormatter(Logger.Formatter)
	logrus.SetLevel(Logger.Level)
}

func LogDebug(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func LogInfo(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func LogWarn(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func LogError(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func LogFatal(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}
