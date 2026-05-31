package middleware

import (
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("role")
		if !ok {
			utils.Forbidden(c, "role not found")
			c.Abort()
			return
		}
		role := models.Role(v.(string))
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		utils.Forbidden(c, "insufficient role")
		c.Abort()
	}
}
