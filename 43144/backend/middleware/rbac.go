package middleware

import (
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, "role not found")
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "insufficient permissions")
		c.Abort()
	}
}
