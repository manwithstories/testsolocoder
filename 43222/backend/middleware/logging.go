package middleware

import (
	"bytes"
	"garden-planner/database"
	"garden-planner/models"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		c.Next()

		userID := GetUserID(c)
		logEntry := models.OperationLog{
			ID:             uuid.New(),
			Action:         c.Request.Method + " " + c.FullPath(),
			Resource:       c.FullPath(),
			Method:         c.Request.Method,
			Path:           c.Request.URL.Path,
			IPAddress:      c.ClientIP(),
			UserAgent:      c.Request.UserAgent(),
			RequestBody:    string(requestBody),
			ResponseStatus: c.Writer.Status(),
			CreatedAt:      startTime,
		}

		if userID != uuid.Nil {
			logEntry.UserID = &userID
		}

		go func() {
			if err := database.DB.Create(&logEntry).Error; err != nil {
				log.Printf("Failed to create operation log: %v", err)
			}
		}()

		log.Printf("[%s] %s %s %d %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			time.Since(startTime),
		)
	}
}
