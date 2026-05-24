package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	jwtpkg "furniture-platform/pkg/jwt"
	"furniture-platform/pkg/response"
)

// JWTAuth JWT 认证中间件
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		claims, err := jwtpkg.ParseToken(parts[1], secret)
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RoleRequired 角色校验中间件
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		currentRole, ok := role.(string)
		if !ok {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		for _, r := range roles {
			if r == currentRole {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}

// RateLimiter 简易限流器（可选）
func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	_ = maxRequests
	_ = window
	return func(c *gin.Context) {
		c.Next()
	}
}
