package services

import (
	"github.com/notification-center/internal/database"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChannelService struct {
	senderService *SenderService
}

func NewChannelService(senderService *SenderService) *ChannelService {
	return &ChannelService{
		senderService: senderService,
	}
}

func (s *ChannelService) Create(channel *models.Channel) (*models.Channel, error) {
	var existing models.Channel
	if err := database.DB.Where("name = ?", channel.Name).First(&existing).Error; err == nil {
		return nil, errors.DuplicateChannel(channel.Name)
	} else if err != gorm.ErrRecordNotFound {
		logger.Error("check duplicate channel failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	if err := database.DB.Create(channel).Error; err != nil {
		logger.Error("create channel failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("channel created", zap.Uint("id", channel.ID), zap.String("name", channel.Name))
	return channel, nil
}

func (s *ChannelService) GetByID(id uint) (*models.Channel, error) {
	var channel models.Channel
	if err := database.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ChannelNotFound(id)
		}
		logger.Error("get channel failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}
	return &channel, nil
}

func (s *ChannelService) GetByName(name string) (*models.Channel, error) {
	var channel models.Channel
	if err := database.DB.Where("name = ?", name).First(&channel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("channel", err)
		}
		logger.Error("get channel by name failed", zap.Error(err), zap.String("name", name))
		return nil, errors.DatabaseError(err)
	}
	return &channel, nil
}

func (s *ChannelService) List(enabledOnly bool, channelType models.ChannelType, page, pageSize int) ([]models.Channel, int64, error) {
	var channels []models.Channel
	var total int64

	query := database.DB.Model(&models.Channel{})

	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}
	if channelType != "" {
		query = query.Where("type = ?", channelType)
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error("count channels failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("priority DESC, id ASC").Offset(offset).Limit(pageSize).Find(&channels).Error; err != nil {
		logger.Error("list channels failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	return channels, total, nil
}

func (s *ChannelService) Update(id uint, updates map[string]interface{}) (*models.Channel, error) {
	channel, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := database.DB.Model(channel).Updates(updates).Error; err != nil {
		logger.Error("update channel failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}

	if err := database.DB.First(channel, id).Error; err != nil {
		logger.Error("reload channel failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("channel updated", zap.Uint("id", id))
	return channel, nil
}

func (s *ChannelService) Delete(id uint) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(&models.Channel{}, id).Error; err != nil {
		logger.Error("delete channel failed", zap.Error(err), zap.Uint("id", id))
		return errors.DatabaseError(err)
	}

	logger.Info("channel deleted", zap.Uint("id", id))
	return nil
}

func (s *ChannelService) Enable(id uint) error {
	result := database.DB.Model(&models.Channel{}).Where("id = ?", id).Update("enabled", true)
	if result.Error != nil {
		logger.Error("enable channel failed", zap.Error(result.Error), zap.Uint("id", id))
		return errors.DatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ChannelNotFound(id)
	}
	logger.Info("channel enabled", zap.Uint("id", id))
	return nil
}

func (s *ChannelService) Disable(id uint) error {
	result := database.DB.Model(&models.Channel{}).Where("id = ?", id).Update("enabled", false)
	if result.Error != nil {
		logger.Error("disable channel failed", zap.Error(result.Error), zap.Uint("id", id))
		return errors.DatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ChannelNotFound(id)
	}
	logger.Info("channel disabled", zap.Uint("id", id))
	return nil
}

func (s *ChannelService) UpdatePriority(id uint, priority int) error {
	result := database.DB.Model(&models.Channel{}).Where("id = ?", id).Update("priority", priority)
	if result.Error != nil {
		logger.Error("update channel priority failed", zap.Error(result.Error), zap.Uint("id", id))
		return errors.DatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ChannelNotFound(id)
	}
	logger.Info("channel priority updated", zap.Uint("id", id), zap.Int("priority", priority))
	return nil
}

func (s *ChannelService) TestConnection(id uint) error {
	channel, err := s.GetByID(id)
	if err != nil {
		return err
	}

	logger.Info("testing channel connection", zap.Uint("id", id), zap.String("type", string(channel.Type)))

	if s.senderService == nil {
		logger.Error("sender service not initialized")
		return errors.ChannelError("sender service not available", nil)
	}

	if err := s.senderService.TestConnection(channel); err != nil {
		logger.Error("channel connection test failed",
			zap.Uint("id", id),
			zap.String("type", string(channel.Type)),
			zap.Error(err),
		)
		return err
	}

	logger.Info("channel connection test passed", zap.Uint("id", id), zap.String("type", string(channel.Type)))
	return nil
}
