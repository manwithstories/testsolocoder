package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"property-management/internal/utils"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "Invalid token format",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusForbidden, utils.Response{
				Code:    403,
				Message: "No permission",
			})
			c.Abort()
			return
		}

		userClaims := claims.(map[string]interface{})
		userRole := userClaims["role"].(string)

		hasPermission := false
		for _, role := range roles {
			if role == userRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, utils.Response{
				Code:    403,
				Message: "No permission",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
