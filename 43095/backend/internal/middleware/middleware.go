package middleware

import (
	"medical-platform/internal/config"
	"medical-platform/internal/models"
	"medical-platform/pkg/utils"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.Duration("cost", cost),
		}

		if statusCode >= 500 {
			config.Logger.Error("Server error", fields...)
		} else if statusCode >= 400 {
			config.Logger.Warn("Client error", fields...)
		} else {
			config.Logger.Info("Request", fields...)
		}
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未提供认证令牌")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "认证令牌格式错误")
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "认证令牌无效或已过期")
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetCurrentUser(c)
		if user == nil {
			utils.Unauthorized(c, "请先登录")
			return
		}

		for _, role := range roles {
			if string(user.Role) == role {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "权限不足")
	}
}

func AdminRequired() gin.HandlerFunc {
	return RoleRequired(string(models.RoleAdmin))
}

func DoctorRequired() gin.HandlerFunc {
	return RoleRequired(string(models.RoleDoctor), string(models.RoleAdmin))
}

func PatientRequired() gin.HandlerFunc {
	return RoleRequired(string(models.RolePatient), string(models.RoleAdmin))
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		config.Logger.Error("Panic recovered", zap.Any("error", err))
		utils.InternalError(c, nil)
	})
}
