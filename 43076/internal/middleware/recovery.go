package middleware

import (
	"log"
	"net/http"

	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				utils.Error(c, http.StatusInternalServerError, "Internal server error")
			}
		}()
		c.Next()
	}
}
