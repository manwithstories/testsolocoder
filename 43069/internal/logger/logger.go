package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/apitester/apitester/pkg/models"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	level      int
	fileWriter io.Writer
	console    bool
	mu         sync.Mutex
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func Init(level string, logFile string, console bool) error {
	once.Do(func() {
		defaultLogger = &Logger{
			level:   parseLevel(level),
			console: console,
		}

		if logFile != "" {
			f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to open log file: %v\n", err)
			} else {
				defaultLogger.fileWriter = f
			}
		}
	})
	return nil
}

func Get() *Logger {
	if defaultLogger == nil {
		Init("info", "", true)
	}
	return defaultLogger
}

func parseLevel(level string) int {
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

func (l *Logger) log(level int, levelStr string, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelStr, msg)

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.console {
		fmt.Print(logLine)
	}
	if l.fileWriter != nil {
		l.fileWriter.Write([]byte(logLine))
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, "DEBUG", format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, "INFO", format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, "WARN", format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, "ERROR", format, args...)
}

func (l *Logger) LogRequest(entry *models.LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		l.Error("Failed to marshal log entry: %v", err)
		return
	}

	if l.fileWriter != nil {
		l.fileWriter.Write(data)
		l.fileWriter.Write([]byte("\n"))
	}
}

func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if closer, ok := l.fileWriter.(io.Closer); ok {
		closer.Close()
	}
}

func Debug(format string, args ...interface{}) {
	Get().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	Get().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	Get().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	Get().Error(format, args...)
}

func LogRequest(entry *models.LogEntry) {
	Get().LogRequest(entry)
}

func Close() {
	Get().Close()
}
