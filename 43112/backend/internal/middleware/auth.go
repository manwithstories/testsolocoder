package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"e-learning-platform/internal/utils"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RefreshTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := c.GetHeader("X-Refresh-Token")
		if refreshToken == "" {
			utils.Unauthorized(c, "Refresh token required")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(refreshToken)
		if err != nil {
			utils.Unauthorized(c, "Invalid refresh token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func TokenExpireCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				claims, err := utils.ParseToken(parts[1])
				if err == nil {
					remaining := time.Until(claims.ExpiresAt.Time)
					if remaining < 1*time.Hour && remaining > 0 {
						c.Writer.Header().Set("X-Token-Will-Expire", "true")
					}
				}
			}
		}
		c.Next()
	}
}
