package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"luxury-trading-platform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func AuditLogMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		startTime := time.Now()

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		params := ""
		if len(bodyBytes) > 0 {
			params = string(bodyBytes)
		}

		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		var uid *uint
		if uidVal, ok := userID.(uint); ok {
			uid = &uidVal
		}

		usernameStr := ""
		if u, ok := username.(string); ok {
			usernameStr = u
		}

		result := "success"
		if c.Writer.Status() >= http.StatusBadRequest {
			result = "failed"
		}

		auditLog := &model.AuditLog{
			UserID:    uid,
			Username:  usernameStr,
			Action:    c.Request.Method + " " + c.FullPath(),
			Module:    extractModule(c.FullPath()),
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Method:    c.Request.Method,
			Path:      c.FullPath(),
			Params:    params,
			Result:    result,
			Detail:    w.body.String(),
		}

		if db := model.GetDB(); db != nil {
			if err := db.Create(auditLog).Error; err != nil {
				log.Errorf("Failed to create audit log: %v", err)
			}
		}

		log.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.FullPath(),
			"status":   c.Writer.Status(),
			"duration": time.Since(startTime),
			"ip":       c.ClientIP(),
			"user_id":  uid,
		}).Info("Request processed")
	}
}

func extractModule(path string) string {
	parts := splitPath(path)
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}

func splitPath(path string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}
	return parts
}

func RetryMiddleware(maxRetries int, delay time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		var lastErr error
		for i := 0; i < maxRetries; i++ {
			c.Next()

			if c.Writer.Status() < http.StatusInternalServerError {
				return
			}

			lastErr = getErrorFromContext(c)

			if i < maxRetries-1 {
				time.Sleep(delay)
				c.Writer = &responseBodyWriter{
					ResponseWriter: c.Writer,
					body:           &bytes.Buffer{},
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		if lastErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "request failed after retries",
				"details": lastErr.Error(),
			})
		}
	}
}

func getErrorFromContext(c *gin.Context) error {
	if err := c.Errors.Last(); err != nil {
		return err.Err
	}
	return nil
}

var bodyBytes []byte

func ErrorHandlerMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			log.Errorf("Errors: %v", c.Errors.String())
		}
	}
}

func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			contentType := c.GetHeader("Content-Type")
			if contentType == "application/json" {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
					c.Abort()
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				bodyBytes = bodyBytes
			}
		}
		c.Next()
	}
}

func ValidateJSON(schema map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
				c.Abort()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var data map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
				c.Abort()
				return
			}

			for key, val := range schema {
				if _, exists := data[key]; !exists {
					if required, ok := val.(bool); ok && required {
						c.JSON(http.StatusBadRequest, gin.H{"error": "missing required field: " + key})
						c.Abort()
						return
					}
				}
			}
		}
		c.Next()
	}
}
