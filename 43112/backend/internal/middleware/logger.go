package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"e-learning-platform/internal/utils"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		logFields := map[string]interface{}{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"latency":     latency.String(),
			"client_ip":   clientIP,
			"user_agent":  c.Request.UserAgent(),
		}

		if raw != "" {
			logFields["query"] = raw
		}

		if len(bodyBytes) > 0 && len(bodyBytes) < 2048 {
			logFields["request_body"] = string(bodyBytes)
		}

		if errors != "" {
			logFields["errors"] = errors
		}

		if statusCode >= 500 {
			utils.Logger.WithFields(logFields).Error("Request failed")
		} else if statusCode >= 400 {
			utils.Logger.WithFields(logFields).Warn("Request warning")
		} else {
			utils.Logger.WithFields(logFields).Info("Request completed")
		}
	}
}
