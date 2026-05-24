package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/utils"
	"recruitment-platform/pkg/logger"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v", err)
				utils.InternalError(c, "服务器内部错误")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
