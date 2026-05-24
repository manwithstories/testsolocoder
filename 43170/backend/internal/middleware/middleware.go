package middleware

import (
	"net/http"
	"photo-rental/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var JWTManager *utils.JWTManager

func InitJWT(secret string, expireHour int) {
	JWTManager = utils.NewJWTManager(secret, expireHour)
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := JWTManager.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("verified", claims.Verified)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Role not found")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusForbidden, "Invalid role format")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if roleStr == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		verified, exists := c.Get("verified")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Verification status not found")
			c.Abort()
			return
		}

		isVerified, ok := verified.(bool)
		if !ok || !isVerified {
			utils.ErrorResponse(c, http.StatusForbidden, "Account not verified")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		utils.Logger.Access(method, path, statusCode, latency)
	}
}
