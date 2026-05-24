package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"business-registration-platform/database"
	"business-registration-platform/models"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		w := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		c.Next()

		userID, _ := c.Get("userID")
		var userIDUint uint
		if uid, ok := userID.(uint); ok {
			userIDUint = uid
		}

		operationLog := models.OperationLog{
			UserID:    userIDUint,
			Module:    c.Request.URL.Path,
			Action:    c.Request.Method,
			Content:   c.Request.URL.Path,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Result:    "success",
		}

		if c.Writer.Status() >= 400 {
			operationLog.Result = "failed"
		}

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if len(bodyBytes) > 0 {
			operationLog.Content = string(bodyBytes)
		}

		database.DB.Create(&operationLog)

		_ = startTime
	}
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, gin.H{
					"code":    500,
					"message": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
