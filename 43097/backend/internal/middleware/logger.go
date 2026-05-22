package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"hotel-system/internal/pkg/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)
		statusCode := c.Writer.Status()

		logFields := map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"user_agent": c.Request.UserAgent(),
		}

		if userID := GetUserID(c); userID > 0 {
			logFields["user_id"] = userID
		}

		if len(c.Errors) > 0 {
			logger.WithFields(logFields).Errorf("Request error: %v", c.Errors.String())
			return
		}

		if statusCode >= 500 {
			logger.WithFields(logFields).Error("Server error")
		} else if statusCode >= 400 {
			logger.WithFields(logFields).Warn("Client error")
		} else {
			logger.WithFields(logFields).Info("Request completed")
		}
	}
}
