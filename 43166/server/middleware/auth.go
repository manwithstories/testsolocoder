package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"business-registration-platform/utils"
	"business-registration-platform/models"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := utils.GetTokenFromHeader(authHeader)
		if tokenString == "" {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			utils.Unauthorized(c, "User role not found")
			c.Abort()
			return
		}

		role, ok := userRole.(models.UserRole)
		if !ok {
			utils.Unauthorized(c, "Invalid user role type")
			c.Abort()
			return
		}

		allowed := false
		for _, r := range roles {
			if role == r {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.Forbidden(c, "Access denied: insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
