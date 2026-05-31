package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"tea-platform/pkg/logger"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		fields := []interface{}{
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"ip", clientIP,
			"user_agent", userAgent,
		}

		if errors != "" {
			fields = append(fields, "errors", errors)
		}

		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, "user_id", userID)
		}
		if username, exists := c.Get("username"); exists {
			fields = append(fields, "username", username)
		}

		switch {
		case statusCode >= 500:
			logger.Sugar.Errorw("Server Error", fields...)
		case statusCode >= 400:
			logger.Sugar.Warnw("Client Error", fields...)
		default:
			logger.Sugar.Infow("Request", fields...)
		}
	}
}
