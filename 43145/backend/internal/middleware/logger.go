package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %d %v | %s", method, path, statusCode, latency, c.ClientIP())
	}
}
