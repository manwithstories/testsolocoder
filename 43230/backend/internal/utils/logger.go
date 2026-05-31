package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(level, logDir string, maxSize, maxBackups, maxAge int) error {
	logger := logrus.New()

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	logger.SetLevel(parsedLevel)

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFilePath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(mw)

	Logger = logger
	return nil
}

func LogInfo(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, args...)
	}
}

func LogWarn(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Warnf(format, args...)
	}
}

func LogError(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, args...)
	}
}

func LogDebug(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, args...)
	}
}

func LogWithFields(fields map[string]interface{}, level, message string) {
	if Logger != nil {
		entry := Logger.WithFields(logrus.Fields(fields))
		switch level {
		case "info":
			entry.Info(message)
		case "warn":
			entry.Warn(message)
		case "error":
			entry.Error(message)
		case "debug":
			entry.Debug(message)
		}
	}
}
