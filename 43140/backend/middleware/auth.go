package middleware

import (
	"strings"

	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Unauthorized(c, "User role not found")
			c.Abort()
			return
		}

		roleStr := role.(string)
		allowed := false
		for _, r := range roles {
			if roleStr == r {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CompanyMiddleware() gin.HandlerFunc {
	return RoleMiddleware("company", "admin")
}

func JobSeekerMiddleware() gin.HandlerFunc {
	return RoleMiddleware("jobseeker", "admin")
}

func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

func GetUserID(c *gin.Context) uint {
	userId, _ := c.Get("user_id")
	return userId.(uint)
}

func GetUserRole(c *gin.Context) string {
	role, _ := c.Get("role")
	return role.(string)
}
