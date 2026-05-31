package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"tea-platform/internal/models"
	"tea-platform/pkg/database"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var sensitiveOperations = map[string]string{
	"DELETE": "删除",
	"PUT":    "更新",
	"PATCH":  "部分更新",
	"POST":   "创建",
}

var sensitivePaths = []string{
	"/api/users",
	"/api/teas",
	"/api/orders",
	"/api/tasting-records",
	"/api/activities",
	"/api/posts",
	"/api/appraisals",
}

func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		if !shouldLog(method, path) {
			c.Next()
			return
		}

		operation := sensitiveOperations[method]
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		var requestBody string
		if method == "POST" || method == "PUT" || method == "PATCH" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		detail := map[string]interface{}{
			"method":        method,
			"path":          path,
			"latency":       latency.String(),
			"response_body": truncateString(blw.body.String(), 1000),
		}
		if requestBody != "" {
			detail["request_body"] = truncateString(requestBody, 1000)
		}

		detailJSON, _ := json.Marshal(detail)

		log := models.OperationLog{
			UserID:    toUint(userID),
			Username:  toString(username),
			Operation: operation,
			Method:    method,
			Path:      path,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Status:    c.Writer.Status(),
			Detail:    string(detailJSON),
		}

		go func() {
			db := database.GetDB()
			if db != nil {
				db.Create(&log)
			}
		}()
	}
}

func shouldLog(method, path string) bool {
	if _, ok := sensitiveOperations[method]; !ok {
		return false
	}
	for _, p := range sensitivePaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func toUint(v interface{}) uint {
	if u, ok := v.(uint); ok {
		return u
	}
	return 0
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
