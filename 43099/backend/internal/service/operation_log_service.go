package service

import (
	"encoding/json"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"

	"github.com/gin-gonic/gin"
)

type OperationLogService struct{}

func NewOperationLogService() *OperationLogService {
	return &OperationLogService{}
}

func (s *OperationLogService) Log(c *gin.Context, userID uint, action, module string, details interface{}) error {
	var detailsStr string
	if details != nil {
		detailsBytes, _ := json.Marshal(details)
		detailsStr = string(detailsBytes)
	}

	log := &model.OperationLog{
		UserID:    userID,
		Action:    action,
		Module:    module,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Details:   detailsStr,
	}

	return database.DB.Create(log).Error
}
