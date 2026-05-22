package services

import (
	"errors"
	"time"

	"auction-system/internal/dto"
	"auction-system/internal/models"
	"auction-system/pkg/logger"
)

type AuctionSessionService struct{}

func NewAuctionSessionService() *AuctionSessionService {
	return &AuctionSessionService{}
}

const (
	SessionStatusPending  = 0
	SessionStatusActive   = 1
	SessionStatusEnded    = 2
	SessionStatusCancelled = 3
)

func (s *AuctionSessionService) CreateSession(req *dto.CreateSessionRequest) (*models.AuctionSession, error) {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
	if err != nil {
		return nil, errors.New("开始时间格式错误，请使用 YYYY-MM-DD HH:MM:SS")
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local)
	if err != nil {
		return nil, errors.New("结束时间格式错误，请使用 YYYY-MM-DD HH:MM:SS")
	}

	if endTime.Before(startTime) {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	session := &models.AuctionSession{
		Name:         req.Name,
		Description:  req.Description,
		StartTime:    startTime,
		EndTime:      endTime,
		MinIncrement: req.MinIncrement,
		ExtendTime:   req.ExtendTime,
		Status:       SessionStatusPending,
	}

	if err := models.DB.Create(session).Error; err != nil {
		logger.Error("Failed to create auction session: %v", err)
		return nil, err
	}

	return session, nil
}

func (s *AuctionSessionService) UpdateSession(sessionID uint, req *dto.UpdateSessionRequest) error {
	var session models.AuctionSession
	if err := models.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("拍卖会不存在")
	}

	if session.Status != SessionStatusPending {
		return errors.New("只能修改未开始的拍卖会")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.StartTime != "" {
		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
		if err != nil {
			return errors.New("开始时间格式错误")
		}
		updates["start_time"] = startTime
	}
	if req.EndTime != "" {
		endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local)
		if err != nil {
			return errors.New("结束时间格式错误")
		}
		updates["end_time"] = endTime
	}
	if req.MinIncrement > 0 {
		updates["min_increment"] = req.MinIncrement
	}
	if req.ExtendTime > 0 {
		updates["extend_time"] = req.ExtendTime
	}
	updates["updated_at"] = time.Now()

	return models.DB.Model(&session).Updates(updates).Error
}

func (s *AuctionSessionService) GetSessionByID(sessionID uint) (*models.AuctionSession, error) {
	var session models.AuctionSession
	if err := models.DB.Preload("AuctionItems.AuctionItem.Images").Preload("AuctionItems.AuctionItem.Category").First(&session, sessionID).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *AuctionSessionService) GetSessionList(page, pageSize int, status *int, keyword string) ([]models.AuctionSession, int64, error) {
	var sessions []models.AuctionSession
	var total int64

	query := models.DB.Model(&models.AuctionSession{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&sessions).Error
	return sessions, total, err
}

func (s *AuctionSessionService) AddItemsToSession(sessionID uint, itemIDs []uint) error {
	var session models.AuctionSession
	if err := models.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("拍卖会不存在")
	}

	if session.Status != SessionStatusPending {
		return errors.New("只能向未开始的拍卖会添加拍卖品")
	}

	for _, itemID := range itemIDs {
		var item models.AuctionItem
		if err := models.DB.First(&item, itemID).Error; err != nil {
			continue
		}

		var existing models.AuctionItemSession
		if models.DB.Where("session_id = ? AND auction_item_id = ?", sessionID, itemID).First(&existing).Error == nil {
			continue
		}

		itemSession := models.AuctionItemSession{
			SessionID:     sessionID,
			AuctionItemID: itemID,
			StartTime:     session.StartTime,
			EndTime:       session.EndTime,
			Status:        0,
		}
		models.DB.Create(&itemSession)
	}

	return nil
}

func (s *AuctionSessionService) RemoveItemFromSession(sessionID uint, itemID uint) error {
	return models.DB.Where("session_id = ? AND auction_item_id = ?", sessionID, itemID).Delete(&models.AuctionItemSession{}).Error
}

func (s *AuctionSessionService) StartSession(sessionID uint) error {
	var session models.AuctionSession
	if err := models.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("拍卖会不存在")
	}

	if session.Status != SessionStatusPending {
		return errors.New("只能开始未开始的拍卖会")
	}

	return models.DB.Model(&session).Update("status", SessionStatusActive).Error
}

func (s *AuctionSessionService) EndSession(sessionID uint) error {
	var session models.AuctionSession
	if err := models.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("拍卖会不存在")
	}

	if session.Status != SessionStatusActive {
		return errors.New("只能结束进行中的拍卖会")
	}

	return models.DB.Model(&session).Update("status", SessionStatusEnded).Error
}

func (s *AuctionSessionService) CancelSession(sessionID uint) error {
	var session models.AuctionSession
	if err := models.DB.First(&session, sessionID).Error; err != nil {
		return errors.New("拍卖会不存在")
	}

	if session.Status == SessionStatusEnded {
		return errors.New("已结束的拍卖会不能取消")
	}

	return models.DB.Model(&session).Update("status", SessionStatusCancelled).Error
}

func (s *AuctionSessionService) GetActiveSessions() ([]models.AuctionSession, error) {
	var sessions []models.AuctionSession
	now := time.Now()
	err := models.DB.Where("status = ? AND start_time <= ? AND end_time >= ?", SessionStatusActive, now, now).Find(&sessions).Error
	return sessions, err
}
