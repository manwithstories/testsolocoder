package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	level    string
	filePath string
}

var AppLogger *Logger

func InitLogger(level, output string) {
	AppLogger = &Logger{
		level:    level,
		filePath: output,
	}

	if output != "stdout" {
		logFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("无法创建日志文件: %v, 使用stdout", err)
		} else {
			log.SetOutput(logFile)
		}
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level == "debug" || l.level == "info" {
		log.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level == "debug" {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	log.Printf("[WARNING] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Audit(action, targetType string, targetID uint, operatorID uint, status string, details string) {
	log.Printf("[AUDIT] time=%s action=%s targetType=%s targetID=%d operatorID=%d status=%s details=%s",
		time.Now().Format(time.RFC3339),
		action, targetType, targetID, operatorID, status, details,
	)
}

func LogInfo(format string, v ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", v...)
	if AppLogger != nil {
		AppLogger.Info(format, v...)
	}
}

func LogError(format string, v ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", v...)
	if AppLogger != nil {
		AppLogger.Error(format, v...)
	}
}

func LogWarning(format string, v ...interface{}) {
	fmt.Printf("[WARNING] "+format+"\n", v...)
	if AppLogger != nil {
		AppLogger.Warning(format, v...)
	}
}
