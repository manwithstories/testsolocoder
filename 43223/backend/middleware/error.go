package middleware

import (
	"net/http"

	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			utils.Error(c, http.StatusInternalServerError, err.Error())
		}
	}
}
