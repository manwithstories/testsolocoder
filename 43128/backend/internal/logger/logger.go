package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	std    *log.Logger
	level  Level = LevelInfo
	outTag = map[Level]string{
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
	}
)

func Init(lv, dir string) {
	_ = os.MkdirAll(dir, 0755)
	switch strings.ToLower(lv) {
	case "debug":
		level = LevelDebug
	case "info":
		level = LevelInfo
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	default:
		level = LevelInfo
	}

	var writers []io.Writer
	writers = append(writers, os.Stdout)
	f, err := os.OpenFile(filepath.Join(dir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		writers = append(writers, f)
	}
	std = log.New(io.MultiWriter(writers...), "", 0)
}

func write(l Level, format string, args ...interface{}) {
	if std == nil {
		return
	}
	if l < level {
		return
	}
	tag := outTag[l]
	ts := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	std.Printf("[%s] [%s] %s\n", ts, tag, msg)
}

func Debugf(format string, args ...interface{}) { write(LevelDebug, format, args...) }
func Infof(format string, args ...interface{})  { write(LevelInfo, format, args...) }
func Warnf(format string, args ...interface{})  { write(LevelWarn, format, args...) }
func Errorf(format string, args ...interface{}) { write(LevelError, format, args...) }

func GinLog(c *gin.Context) {
	start := time.Now()
	c.Next()
	latency := time.Since(start)
	write(LevelInfo, "method=%s path=%s status=%d latency=%s ip=%s",
		c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency.String(), c.ClientIP())
}
