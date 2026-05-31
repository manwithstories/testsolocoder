package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"museum-server/pkg/logger"
	"museum-server/pkg/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorWithCtx("Recovery", "panic recovered: %v", err)
				response.Error(c, http.StatusInternalServerError, 500, "Internal Server Error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
