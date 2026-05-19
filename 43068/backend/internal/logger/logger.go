package logger

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	AccessLog    *log.Logger
)

func InitLogger() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFile, err := os.OpenFile(logDir+"/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	accessFile, err := os.OpenFile(logDir+"/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open access log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	accessWriter := io.MultiWriter(os.Stdout, accessFile)

	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	AccessLog = log.New(accessWriter, "ACCESS: ", log.Ldate|log.Ltime)
}

func LogInfo(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

func LogError(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

func LogAccess(method, path, clientIP string, statusCode int, duration time.Duration) {
	AccessLog.Printf("%s %s %s %d %v", clientIP, method, path, statusCode, duration)
}

func LogOperation(userID uint, operation, details string) {
	InfoLogger.Printf("User: %d, Operation: %s, Details: %s", userID, operation, details)
}
