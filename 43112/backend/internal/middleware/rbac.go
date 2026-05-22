package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"e-learning-platform/internal/utils"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, "Role not found")
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		hasRole := false
		for _, r := range roles {
			if r == roleStr {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

func RequireInstructor() gin.HandlerFunc {
	return RequireRole("instructor", "admin")
}

func RequireStudent() gin.HandlerFunc {
	return RequireRole("student", "instructor", "admin")
}

func RequireOwnerOrAdmin(idField string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role == "admin" {
			c.Next()
			return
		}

		userID, _ := c.Get("user_id")
		paramID := c.Param(idField)

		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			utils.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		if paramID != userUUID.String() {
			utils.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}
