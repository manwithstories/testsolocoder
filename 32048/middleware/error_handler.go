package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err.Error())

			status := http.StatusInternalServerError
			if c.Writer.Status() != 200 {
				status = c.Writer.Status()
			}

			c.JSON(status, gin.H{
				"error": err.Error(),
			})
		}
	}
}
