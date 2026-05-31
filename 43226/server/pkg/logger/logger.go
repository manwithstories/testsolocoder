package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
)

func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, fmt.Sprintf("museum_%s.log", time.Now().Format("2006-01-02")))

	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	infoLogger = log.New(multiWriter, "[INFO] ", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(multiWriter, "[ERROR] ", log.LstdFlags|log.Lshortfile)

	return nil
}

func Info(format string, v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}

func InfoWithCtx(module string, format string, v ...interface{}) {
	msg := fmt.Sprintf("[%s] %s", module, fmt.Sprintf(format, v...))
	infoLogger.Output(2, msg)
}

func ErrorWithCtx(module string, format string, v ...interface{}) {
	msg := fmt.Sprintf("[%s] %s", module, fmt.Sprintf(format, v...))
	errorLogger.Output(2, msg)
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
