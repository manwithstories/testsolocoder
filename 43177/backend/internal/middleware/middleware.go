package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"repair-platform/internal/utils"
	"repair-platform/pkg/logger"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "未提供认证令牌",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "认证令牌无效或已过期",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, utils.Response{
				Code:    403,
				Message: "无权限访问",
			})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, utils.Response{
			Code:    403,
			Message: "无权限访问",
		})
		c.Abort()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logFields := map[string]interface{}{
			"method":   method,
			"path":     path,
			"status":   statusCode,
			"latency":  latency.String(),
			"client_ip": c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		if statusCode >= 500 {
			logger.Errorf("Server error: %+v", logFields)
		} else if statusCode >= 400 {
			logger.Warnf("Client error: %+v", logFields)
		} else {
			logger.Debugf("Request: %+v", logFields)
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v, Path: %s", err, c.Request.URL.Path)
				c.JSON(http.StatusInternalServerError, utils.Response{
					Code:    500,
					Message: "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimitMiddleware(maxRequests int, windowSeconds int) gin.HandlerFunc {
	requestCounts := make(map[string]int)
	windowStart := time.Now()

	return func(c *gin.Context) {
		now := time.Now()
		if now.Sub(windowStart).Seconds() > float64(windowSeconds) {
			requestCounts = make(map[string]int)
			windowStart = now
		}

		clientIP := c.ClientIP()
		requestCounts[clientIP]++

		if requestCounts[clientIP] > maxRequests {
			c.JSON(http.StatusTooManyRequests, utils.Response{
				Code:    429,
				Message: "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
