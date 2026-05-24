package logger

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"watchplatform/internal/config"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	AccessLog io.Writer
	AppLog    *lumberjack.Logger
)

func Init() {
	_ = os.MkdirAll(config.Cfg.LogDir, 0o755)
	AppLog = &lumberjack.Logger{
		Filename:   filepath.Join(config.Cfg.LogDir, "app.log"),
		MaxSize:    50,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	AccessLog = io.MultiWriter(AppLog, os.Stdout)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		latency := time.Since(start)
		line := time.Now().Format("2006-01-02 15:04:05") + " " +
			c.Request.Method + " " + path + " " + query +
			" status=" + strconv.Itoa(c.Writer.Status()) +
			" latency=" + latency.String() + "\n"
		_, _ = AccessLog.Write([]byte(line))
	}
}

func Close() {
	if AppLog != nil {
		_ = AppLog.Close()
	}
}
