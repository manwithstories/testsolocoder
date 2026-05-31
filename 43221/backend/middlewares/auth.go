package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"consultation-platform/config"
	"consultation-platform/utils"
	"go.uber.org/zap"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, 401, "Missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.ErrorResponse(c, http.StatusUnauthorized, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], secret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			utils.ErrorResponse(c, http.StatusUnauthorized, 401, "Token has expired")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, err := utils.GetUserRoleFromContext(c)
		if err != nil {
			utils.ErrorResponse(c, http.StatusForbidden, 403, "Permission denied")
			c.Abort()
			return
		}

		allowed := false
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorResponse(c, http.StatusForbidden, 403, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				utils.ErrorResponse(c, http.StatusInternalServerError, 500, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimitMiddleware(cfg config.RateLimitConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	requestCounts := make(map[string]int)
	var lastCleanup time.Time

	return func(c *gin.Context) {
		now := time.Now()
		if now.Sub(lastCleanup) > time.Minute {
			requestCounts = make(map[string]int)
			lastCleanup = now
		}

		ip := c.ClientIP()
		requestCounts[ip]++

		if requestCounts[ip] > cfg.RequestsPerMinute {
			utils.ErrorResponse(c, http.StatusTooManyRequests, 429, "Rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
