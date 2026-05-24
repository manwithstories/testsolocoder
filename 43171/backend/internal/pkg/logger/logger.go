package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"drone-rental/internal/config"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	file   *os.File
	writer io.Writer
}

var L *Logger

func New(cfg config.LogConfig) (*Logger, error) {
	dir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	l := &Logger{
		file:   f,
		writer: io.MultiWriter(os.Stdout, f),
	}
	L = l
	return l, nil
}

func (l *Logger) Writer() io.Writer {
	return l.writer
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) Info(msg string) {
	l.writer.Write([]byte("[INFO] " + time.Now().Format("2006-01-02 15:04:05") + " " + msg + "\n"))
}

func (l *Logger) Error(msg string) {
	l.writer.Write([]byte("[ERROR] " + time.Now().Format("2006-01-02 15:04:05") + " " + msg + "\n"))
}

func (l *Logger) Warn(msg string) {
	l.writer.Write([]byte("[WARN] " + time.Now().Format("2006-01-02 15:04:05") + " " + msg + "\n"))
}

func GinLogger() gin.HandlerFunc {
	return gin.LoggerWithWriter(L.Writer())
}
