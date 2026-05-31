package middlewares

import (
	"strings"
	"time"

	"secondhand-platform/cache"
	"secondhand-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func JWTAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorWithStatus(c, 401, 401, "未提供认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorWithStatus(c, 401, 401, "认证令牌格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.ErrorWithStatus(c, 401, 401, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		exists, _ := cache.Exists(c.Request.Context(), "token:blacklist:"+parts[1])
		if exists {
			utils.ErrorWithStatus(c, 401, 401, "认证令牌已被注销")
			c.Abort()
			return
		}

		if len(roles) > 0 {
			hasRole := false
			for _, role := range roles {
				if role == claims.Role {
					hasRole = true
					break
				}
			}
			if !hasRole {
				utils.ErrorWithStatus(c, 403, 403, "权限不足")
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		fields := logrus.Fields{
			"status_code": statusCode,
			"latency":     latency.String(),
			"client_ip":   c.ClientIP(),
			"method":      method,
			"path":        path,
			"user_agent":  c.Request.UserAgent(),
		}

		if statusCode >= 500 {
			logrus.WithFields(fields).Error("Server error")
		} else if statusCode >= 400 {
			logrus.WithFields(fields).Warn("Client error")
		} else {
			logrus.WithFields(fields).Info("Request processed")
		}
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("Panic recovered: %v", err)
				utils.ErrorWithStatus(c, 500, 500, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimitMiddleware(maxRequests int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		count, err := cache.Incr(c.Request.Context(), key)
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			cache.Expire(c.Request.Context(), key, time.Minute)
		}

		if count > int64(maxRequests) {
			utils.ErrorWithStatus(c, 429, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
