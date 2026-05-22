package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"travel-planner/internal/database"
	"travel-planner/internal/models"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		userID := GetCurrentUserID(c)
		action := getActionFromMethod(c.Request.Method)
		resource := c.Request.URL.Path

		oldValue := string(requestBody)
		newValue := w.body.String()

		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		auditLog := models.AuditLog{
			UserID:    userID,
			Action:    action,
			Resource:  resource,
			OldValue:  truncateString(oldValue, 1000),
			NewValue:  truncateString(newValue, 1000),
			IPAddress: ip,
			UserAgent: userAgent,
			CreatedAt: time.Now(),
		}

		if planID := c.GetString("plan_id_context"); planID != "" {
			auditLog.PlanID, _ = parseUUID(planID)
		}

		database.DB.Create(&auditLog)
	}
}

func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "other"
	}
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

func parseUUID(s string) (uuid.UUID, error) {
	var id uuid.UUID
	err := json.Unmarshal([]byte(`"`+s+`"`), &id)
	return id, err
}
