package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	LogFile *os.File
)

func InitLogger(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	LogFile = f

	gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, f)

	return nil
}

func CloseLogger() {
	if LogFile != nil {
		LogFile.Close()
	}
}

func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf("[INFO] %s %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	fmt.Fprintln(gin.DefaultWriter, msg)
}

func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf("[ERROR] %s %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	fmt.Fprintln(gin.DefaultErrorWriter, msg)
}

func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf("[WARN] %s %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	fmt.Fprintln(gin.DefaultWriter, msg)
}
