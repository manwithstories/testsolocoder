package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/logger"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		requestID, _ := c.Get("request_id")

		fields := []zap.Field{
			zap.String("request_id", requestID.(string)),
			zap.Int("status", statusCode),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		if len(bodyBytes) > 0 {
			fields = append(fields, zap.String("request_body", string(bodyBytes)))
		}

		if len(w.body.Bytes()) > 0 && statusCode >= 400 {
			fields = append(fields, zap.String("response_body", w.body.String()))
		}

		if statusCode >= 500 {
			logger.Error("request failed", fields...)
		} else if statusCode >= 400 {
			logger.Warn("request completed with warning", fields...)
		} else {
			logger.Info("request completed", fields...)
		}
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *errors.AppError
			if e, ok := err.(*errors.AppError); ok {
				appErr = e
			} else {
				appErr = errors.InternalServerError("internal server error", err)
			}

			requestID, _ := c.Get("request_id")
			c.JSON(appErr.HTTPStatus(), gin.H{
				"code":       appErr.Code,
				"message":    appErr.Message,
				"details":    appErr.Details,
				"request_id": requestID,
				"timestamp":  time.Now().Format(time.RFC3339),
			})
			c.Abort()
			return
		}
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RateLimit(limiter interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		requestID, _ := c.Get("request_id")
		logger.Error("panic recovered",
			zap.Any("error", err),
			zap.String("request_id", requestID.(string)),
			zap.Stack("stack"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":       errors.CodeInternalServerError,
			"message":    "internal server error",
			"request_id": requestID,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		c.Abort()
	})
}
