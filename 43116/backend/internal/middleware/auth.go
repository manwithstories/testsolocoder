package middleware

import (
	"car-rental/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserContext struct {
	UserID   uint
	Username string
	RoleID   uint
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "需要认证")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			utils.Unauthorized(c, "Token已过期")
			c.Abort()
			return
		}

		c.Set("user", &UserContext{
			UserID:   claims.UserID,
			Username: claims.Username,
			RoleID:   claims.RoleID,
		})

		c.Next()
	}
}

func GetUserContext(c *gin.Context) *UserContext {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*UserContext)
}

func RequireRole(roleID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUserContext(c)
		if user == nil {
			utils.Unauthorized(c, "需要认证")
			c.Abort()
			return
		}

		if user.RoleID != roleID && user.RoleID != 1 {
			utils.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUserContext(c)
		if user == nil {
			utils.Unauthorized(c, "需要认证")
			c.Abort()
			return
		}

		if user.RoleID != 1 {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}