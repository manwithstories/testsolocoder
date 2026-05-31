package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"
	"translation-platform/internal/database"
	"translation-platform/internal/models"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		w := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		log := models.OperationLog{
			Action:    c.Request.Method + " " + c.FullPath(),
			Module:    extractModule(c.FullPath()),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}

		if userID, exists := c.Get("user_id"); exists {
			uid := userID.(uint)
			log.UserID = &uid
		}

		detail := map[string]interface{}{
			"status":       c.Writer.Status(),
			"duration":     time.Since(start).String(),
			"request_body": truncateString(requestBody, 500),
			"response":     truncateString(w.body.String(), 500),
		}

		if detailJSON, err := json.Marshal(detail); err == nil {
			log.Detail = string(detailJSON)
		}

		database.DB.Create(&log)
	}
}

func extractModule(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 1 {
		return parts[0]
	}
	return "unknown"
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
