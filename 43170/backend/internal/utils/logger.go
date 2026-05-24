package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Logger *AppLogger

type AppLogger struct {
	logFile   *os.File
	infoLog   *log.Logger
	errorLog  *log.Logger
	accessLog *log.Logger
}

func InitLogger() {
	logDir := "logs"
	os.MkdirAll(logDir, 0755)

	date := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", date))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	Logger = &AppLogger{
		logFile:   logFile,
		infoLog:   log.New(multiWriter, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog:  log.New(multiWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		accessLog: log.New(multiWriter, "[ACCESS] ", log.Ldate|log.Ltime),
	}
}

func (l *AppLogger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

func (l *AppLogger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

func (l *AppLogger) Access(method, path string, status int, duration time.Duration) {
	l.accessLog.Printf("%s %s - %d - %v", method, path, status, duration)
}

func (l *AppLogger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}
