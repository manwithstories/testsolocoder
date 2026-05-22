package middleware

import (
	"bytes"
	"io"
	"time"

	"car-rental/internal/model"
	cachedb "car-rental/internal/config"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.URL.Path == "/api/login" || c.Request.URL.Path == "/api/register" {
			c.Next()
			return
		}

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		start := time.Now()

		c.Next()

		latency := time.Since(start)
		user := GetUserContext(c)

		module := "other"
		action := c.Request.Method
		path := c.Request.URL.Path

		if len(path) > 4 && path[:4] == "/api" {
			parts := splitPath(path[4:])
			if len(parts) > 0 {
				module = parts[0]
			}
		}

		log := model.OperationLog{
			Module:      module,
			Action:      action,
			Description: string(requestBody),
			IP:          c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			Method:      c.Request.Method,
			Path:        path,
			Status:      c.Writer.Status(),
		}

		if user != nil {
			log.UserID = user.UserID
		}

		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		cachedb.DB.Create(&log)
	}
}

func splitPath(path string) []string {
	if path[0] == '/' {
		path = path[1:]
	}
	return splitString(path, '/')
}

func splitString(s string, sep rune) []string {
	var parts []string
	current := ""
	for _, r := range s {
		if r == sep {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}