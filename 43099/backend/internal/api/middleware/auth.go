package middleware

import (
	"net/http"
	"strings"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/auth"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.Error(401, "Authorization header is required"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.Error(401, "Invalid authorization format"))
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := auth.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Error(401, "Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.Error(401, "User not authenticated"))
			c.Abort()
			return
		}

		userRole, ok := role.(model.UserRole)
		if !ok {
			c.JSON(http.StatusForbidden, dto.Error(403, "Invalid user role"))
			c.Abort()
			return
		}

		if userRole != model.RoleAdmin && userRole != model.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, dto.Error(403, "Admin permission required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func SuperAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.Error(401, "User not authenticated"))
			c.Abort()
			return
		}

		userRole, ok := role.(model.UserRole)
		if !ok {
			c.JSON(http.StatusForbidden, dto.Error(403, "Invalid user role"))
			c.Abort()
			return
		}

		if userRole != model.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, dto.Error(403, "Super admin permission required"))
			c.Abort()
			return
		}

		c.Next()
	}
}
