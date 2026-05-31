package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pet-board/internal/models"
	"pet-board/internal/service"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func OperationLogMiddleware(logService *service.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		w := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		latency := time.Since(start)

		var userID *uuid.UUID
		var username string
		if uid, exists := c.Get("user_id"); exists {
			parsedID, err := uuid.Parse(uid.(string))
			if err == nil {
				userID = &parsedID
			}
		}
		if uname, exists := c.Get("username"); exists {
			username = uname.(string)
		}

		params := ""
		if len(bodyBytes) > 0 && len(bodyBytes) < 10000 {
			params = string(bodyBytes)
		}

		result := ""
		if w.body.Len() > 0 && w.body.Len() < 10000 {
			result = w.body.String()
		}

		errorMsg := ""
		if c.Writer.Status() >= 500 {
			var resp struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(w.body.Bytes(), &resp); err == nil {
				errorMsg = resp.Message
			}
		}

		log := &models.OperationLog{
			ID:        uuid.New(),
			UserID:    userID,
			Username:  username,
			Action:    c.FullPath(),
			Method:    c.Request.Method,
			URL:       c.Request.URL.String(),
			IP:        c.ClientIP(),
			Params:    params,
			Result:    result,
			Status:    c.Writer.Status(),
			ExecTime:  latency.Milliseconds(),
			ErrorMsg:  errorMsg,
			CreatedAt: start,
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			_ = logService.Create(log)
		}()
	}
}
