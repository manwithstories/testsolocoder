package middleware

import (
	"ticket-system/internal/common/exception"
	"ticket-system/internal/common/response"
	"ticket-system/internal/logger"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Errorf("Panic recovered: %v", err)
				if be, ok := err.(*exception.BusinessException); ok {
					response.Fail(c, be.Code, be.Message)
				} else {
					response.ServerError(c, "服务器内部错误")
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
