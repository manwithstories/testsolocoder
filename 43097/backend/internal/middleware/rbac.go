package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hotel-system/internal/pkg/rbac"
)

func RBACMiddleware(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c)
		if role == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "用户角色信息不存在",
			})
			c.Abort()
			return
		}

		if !rbac.CheckPermission(role, resource, action) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足，无法访问该资源",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c)
		if !rbac.HasAdminRole(role) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func FrontDeskRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c)
		if !rbac.HasFrontDeskRole(role) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要前台或管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
