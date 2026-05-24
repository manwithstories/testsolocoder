package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"printshop/internal/auth"
	"printshop/internal/config"
	"printshop/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type contextKey string

const (
	CtxUserID  contextKey = "userID"
	CtxUsername contextKey = "username"
	CtxRoleID  contextKey = "roleID"
	CtxRole    contextKey = "role"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		claims, err := auth.ParseToken(cfg, parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.Set(string(CtxUserID), claims.UserID)
		c.Set(string(CtxUsername), claims.Username)
		c.Set(string(CtxRoleID), claims.RoleID)
		c.Set(string(CtxRole), claims.Role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(string(CtxRole))
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
			return
		}
		roleStr := role.(string)
		for _, r := range allowedRoles {
			if r == roleStr {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		userID, _ := c.Get(string(CtxUserID))
		uid, _ := userID.(uint)
		action := c.Request.Method
		resource := c.FullPath()
		status := c.Writer.Status()
		log := models.AuditLog{
			UserID:     uid,
			Action:     action,
			Resource:   resource,
			Detail:     strconv.Itoa(status),
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}
		go db.Create(&log)
		_ = start
	}
}
