package utils

import (
	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"time"
)

func LogOperation(operatorID uint, operatorRole, module, action string, targetID *uint, targetType, content, ip, userAgent string) error {
	log := models.OperationLog{
		OperatorID:   operatorID,
		OperatorRole: operatorRole,
		Module:       module,
		Action:       action,
		TargetID:     targetID,
		TargetType:   targetType,
		Content:      content,
		IP:           ip,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}
	return config.DB.Create(&log).Error
}
