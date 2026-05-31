package audit

import (
	"encoding/json"
	"time"

	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LogAction(c *gin.Context, action string, entityType string, entityID *uuid.UUID, oldValue interface{}, newValue interface{}) {
	var oldValueStr, newValueStr string

	if oldValue != nil {
		oldBytes, _ := json.Marshal(oldValue)
		oldValueStr = string(oldBytes)
	}

	if newValue != nil {
		newBytes, _ := json.Marshal(newValue)
		newValueStr = string(newBytes)
	}

	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if parsed, err := uuid.Parse(uid.(string)); err == nil {
			userID = &parsed
		}
	}

	log := model.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValue:   oldValueStr,
		NewValue:   newValueStr,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		RequestID:  c.GetHeader("X-Request-ID"),
		Status:     "success",
		CreatedAt:  time.Now(),
	}

	if err := database.DB.Create(&log).Error; err != nil {
		logrus.Errorf("Failed to create audit log: %v", err)
	}
}

func LogActionWithStatus(c *gin.Context, action string, entityType string, entityID *uuid.UUID, oldValue interface{}, newValue interface{}, status string) {
	var oldValueStr, newValueStr string

	if oldValue != nil {
		oldBytes, _ := json.Marshal(oldValue)
		oldValueStr = string(oldBytes)
	}

	if newValue != nil {
		newBytes, _ := json.Marshal(newValue)
		newValueStr = string(newBytes)
	}

	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if parsed, err := uuid.Parse(uid.(string)); err == nil {
			userID = &parsed
		}
	}

	log := model.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValue:   oldValueStr,
		NewValue:   newValueStr,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		RequestID:  c.GetHeader("X-Request-ID"),
		Status:     status,
		CreatedAt:  time.Now(),
	}

	if err := database.DB.Create(&log).Error; err != nil {
		logrus.Errorf("Failed to create audit log: %v", err)
	}
}
