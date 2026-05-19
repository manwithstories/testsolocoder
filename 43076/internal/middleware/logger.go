package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var requestBody string
		if len(bodyBytes) > 0 {
			requestBody = string(bodyBytes)
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		ticketNo, _ := c.Get("ticket_no")

		log.Printf(
			"[API] %s | %d | %v | %s | %s | UserID: %v | Username: %v | TicketNo: %v | Request: %s | Response: %s",
			c.Request.Method,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.URL.Path,
			userID,
			username,
			ticketNo,
			requestBody,
			w.body.String(),
		)
	}
}

func SetRequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID, exists := c.Get("user_id"); exists {
			fields := make(map[string]interface{})
			fields["user_id"] = userID
			fields["path"] = c.Request.URL.Path
			fields["method"] = c.Request.Method
			fields["ip"] = c.ClientIP()
			c.Set("request_context", fields)
		}
		c.Next()
	}
}
