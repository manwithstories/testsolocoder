package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"museum-server/pkg/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.InfoWithCtx("HTTP", "[%s] %s %s %d %v",
			c.Request.Method, path, query, status, latency,
		)
	}
}
