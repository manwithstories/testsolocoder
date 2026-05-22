package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logger    *log.Logger
	logFile   *os.File
	logLevel  LogLevel = INFO
)

func InitLogger(logPath string) error {
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logFile = file
	logger = log.New(file, "", log.LstdFlags)
	return nil
}

func SetLevel(level LogLevel) {
	logLevel = level
}

func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func logf(level LogLevel, format string, v ...interface{}) {
	if level < logLevel {
		return
	}

	levelStr := [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[level]
	caller := getCallerInfo()
	msg := fmt.Sprintf(format, v...)
	logMsg := fmt.Sprintf("[%s] %s - %s", levelStr, caller, msg)

	if logger != nil {
		logger.Println(logMsg)
	}
	fmt.Println(logMsg)

	if level == FATAL {
		os.Exit(1)
	}
}

func Debug(format string, v ...interface{}) {
	logf(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	logf(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	logf(WARN, format, v...)
}

func Error(format string, v ...interface{}) {
	logf(ERROR, format, v...)
}

func Fatal(format string, v ...interface{}) {
	logf(FATAL, format, v...)
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func GetLogFileName() string {
	return fmt.Sprintf("auction_%s.log", time.Now().Format("2006-01-02"))
}
