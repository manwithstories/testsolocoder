package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"museum-server/internal/config"
	"museum-server/pkg/response"
	"museum-server/pkg/utils"
)

func JWTAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 401, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(cfg, parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, 401, "User role not found")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, 401, "Invalid user role")
			c.Abort()
			return
		}

		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, 403, "Access denied")
		c.Abort()
	}
}
