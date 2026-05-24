package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"recruitment-platform/pkg/logger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if statusCode >= 500 {
			logger.Error("[%s] %s %d %v IP:%s Errors:%s", method, path, statusCode, latency, clientIP, errors)
		} else if statusCode >= 400 {
			logger.Warn("[%s] %s %d %v IP:%s", method, path, statusCode, latency, clientIP)
		} else {
			logger.Info("[%s] %s %d %v IP:%s", method, path, statusCode, latency, clientIP)
		}
	}
}
