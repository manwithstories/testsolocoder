package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"wedding-planner/internal/services"
	"wedding-planner/pkg/database"

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

func OperationLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		module := extractModule(path)
		action := extractAction(method, path)

		if module == "" || action == "" {
			c.Next()
			return
		}

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		userID, _ := c.Get("user_id")
		userIDUint, _ := userID.(uint)

		if userIDUint > 0 {
			detail := string(bodyBytes)
			if len(detail) > 500 {
				detail = detail[:500]
			}

			if detail != "" {
				var data map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &data); err == nil {
					if id, ok := data["id"]; ok {
						if idFloat, ok := id.(float64); ok {
							services.CreateOperationLog(
								database.GetDB(),
								userIDUint,
								module,
								action,
								uint(idFloat),
								detail,
								c.ClientIP(),
							)
							return
						}
					}
				}
			}

			services.CreateOperationLog(
				database.GetDB(),
				userIDUint,
				module,
				action,
				0,
				detail,
				c.ClientIP(),
			)
		}
	}
}

func extractModule(path string) string {
	modules := map[string]string{
		"/api/auth":        "auth",
		"/api/users":       "user",
		"/api/weddings":    "wedding",
		"/api/vendors":     "vendor",
		"/api/guests":      "guest",
		"/api/budget":      "budget",
		"/api/tasks":       "task",
		"/api/documents":   "document",
		"/api/dashboard":   "dashboard",
		"/api/notifications": "notification",
	}

	for prefix, module := range modules {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			return module
		}
	}

	return ""
}

func extractAction(method, path string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}
