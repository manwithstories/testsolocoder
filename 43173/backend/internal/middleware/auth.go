package middleware

import (
	"time"

	"music-platform/pkg/jwt"
	"music-platform/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, "未提供Token")
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "无效的Token")
			c.Abort()
			return
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			response.Unauthorized(c, "Token已过期")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func JWTAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := jwt.ParseToken(token)
		if err == nil && time.Now().Unix() <= claims.ExpiresAt.Unix() {
			c.Set("claims", claims)
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}

func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		for _, role := range roles {
			if role == userRole {
				c.Next()
				return
			}
		}
		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}
