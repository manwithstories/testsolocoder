package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"skillshare/pkg/auth"
	"skillshare/pkg/response"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1], secret)
		if err != nil {
			response.Unauthorized(c, "无效或过期的认证令牌")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_nickname", claims.Nickname)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if statusCode >= 500 && len(c.Errors) > 0 {
			_ = c.Errors.Last()
		}

		_ = latency
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	var requestCounts = make(map[string]int)
	var lastCleanup = time.Now()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if time.Since(lastCleanup) > time.Minute {
			requestCounts = make(map[string]int)
			lastCleanup = time.Now()
		}

		requestCounts[clientIP]++

		if requestCounts[clientIP] > 60 {
			c.JSON(429, gin.H{"error": "请求过于频繁"})
			c.Abort()
			return
		}

		c.Next()
	}
}
