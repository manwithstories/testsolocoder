package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"coffee-platform/database"
	"coffee-platform/models"

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

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

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

		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)

		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		var uid uint
		if userID != nil {
			uid = userID.(uint)
		}

		var uname string
		if username != nil {
			uname = username.(string)
		}

		action := c.Request.Method + " " + c.FullPath()

		log := models.OperationLog{
			UserID:      uid,
			Username:    uname,
			Action:      action,
			Resource:    getResourceFromPath(c.FullPath()),
			ResourceID:  c.Param("id"),
			Method:      c.Request.Method,
			Path:        c.Request.URL.Path,
			IP:          c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			StatusCode:  c.Writer.Status(),
			RequestData: string(requestBody),
			ResponseData: w.body.String(),
		}

		if c.Writer.Status() >= 400 {
			var resp struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(w.body.Bytes(), &resp); err == nil {
				log.ErrorMsg = resp.Message
			}
		}

		if latency > 5*time.Second {
			log.ErrorMsg = "请求耗时过长: " + latency.String()
		}

		database.DB.Create(&log)
	}
}

func getResourceFromPath(path string) string {
	parts := splitPath(path)
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func splitPath(path string) []string {
	var parts []string
	for _, part := range splitString(path, "/") {
		if part != "" && part[0] != ':' {
			parts = append(parts, part)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep[0] {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
