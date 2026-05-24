package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "User not found",
			})
			c.Abort()
			return
		}

		if user.Status != models.UserStatusActive {
			c.JSON(http.StatusForbidden, models.Response{
				Code:    403,
				Message: "Account is not active",
			})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Set("user", &user)

		now := time.Now()
		database.DB.Model(&user).Updates(map[string]interface{}{
			"last_login_at": now,
		})

		c.Next()
	}
}

func RoleAuth(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(models.UserRole)
		if !ok {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Invalid role type",
			})
			c.Abort()
			return
		}

		authorized := false
		for _, r := range roles {
			if role == r {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, models.Response{
				Code:    403,
				Message: "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: err.Error(),
			})
		}
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

		logEntry := map[string]interface{}{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"latency":     latency.String(),
			"timestamp":   time.Now().Format(time.RFC3339),
		}

		if statusCode >= 400 {
			gin.DefaultErrorWriter.Write([]byte(
				logEntry["timestamp"].(string) + " " +
					logEntry["method"].(string) + " " +
					logEntry["path"].(string) + " " +
					fmt.Sprintf("%d", logEntry["status_code"].(int)) + " " +
					logEntry["latency"].(string) + "\n",
			))
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func ParseUUIDParam(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param(paramName)
		if idStr == "" {
			c.Next()
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Code:    400,
				Message: "Invalid UUID format",
			})
			c.Abort()
			return
		}

		c.Set(paramName+"_uuid", id)
		c.Next()
	}
}
