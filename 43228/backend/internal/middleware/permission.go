package middleware

import (
	"github.com/gin-gonic/gin"

	"tea-platform/internal/utils"
)

func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Fail(c, utils.CodeUnauthorized, "未获取到用户角色")
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			utils.Fail(c, utils.CodeUnauthorized, "用户角色类型错误")
			c.Abort()
			return
		}

		allowed := false
		for _, r := range roles {
			if userRole == r {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.Fail(c, utils.CodeForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
