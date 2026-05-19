package middleware

import (
	"bytes"
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		start := time.Now()
		c.Next()

		userID, _ := c.Get("userID")
		role, _ := c.Get("role")

		var memberID *uint
		var operatorType string

		if role == "member" && userID != nil {
			id := userID.(uint)
			memberID = &id
			operatorType = "member"
		} else if role == "admin" {
			operatorType = "admin"
		} else {
			operatorType = "guest"
		}

		log := models.OperationLog{
			MemberID:     memberID,
			OperatorID:   getUintFromInterface(userID),
			OperatorType: operatorType,
			Action:       c.Request.Method,
			ResourceType: c.FullPath(),
			Detail:       string(bodyBytes),
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
		}

		database.GetDB().Create(&log)

		_ = time.Since(start)
	}
}

func getUintFromInterface(v interface{}) uint {
	if v == nil {
		return 0
	}
	return v.(uint)
}
