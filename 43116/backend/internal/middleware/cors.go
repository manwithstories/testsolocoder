package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RateLimit() gin.HandlerFunc {
	requestCount := make(map[string]int)
	var lastClean time.Time

	return func(c *gin.Context) {
		now := time.Now()
		if now.Sub(lastClean) > time.Hour {
			requestCount = make(map[string]int)
			lastClean = now
		}

		ip := c.ClientIP()
		requestCount[ip]++

		if requestCount[ip] > 1000 {
			c.JSON(429, gin.H{"error": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, gin.H{"error": "服务器内部错误"})
				c.Abort()
			}
		}()
		c.Next()
	}
}