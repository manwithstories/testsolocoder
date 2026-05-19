package services

import (
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	apperrors "github.com/notification-center/internal/errors"
	"go.uber.org/zap"
)

type ChannelAdapter interface {
	Send(message *models.Message, config string) error
	TestConnection(config string) error
}

type SenderService struct {
	adapters map[models.ChannelType]ChannelAdapter
}

func NewSenderService() *SenderService {
	return &SenderService{
		adapters: make(map[models.ChannelType]ChannelAdapter),
	}
}

func (ss *SenderService) RegisterAdapter(channelType models.ChannelType, adapter ChannelAdapter) {
	ss.adapters[channelType] = adapter
	logger.Info("channel adapter registered", zap.String("type", string(channelType)))
}

func (ss *SenderService) Send(message *models.Message, channel *models.Channel) error {
	adapter, exists := ss.adapters[channel.Type]
	if !exists {
		return apperrors.ChannelError("no adapter registered for channel type", nil)
	}

	if !channel.Enabled {
		return apperrors.ChannelError("channel is disabled", nil)
	}

	if err := adapter.Send(message, channel.Config); err != nil {
		logger.Error("channel send failed",
			zap.String("channel_type", string(channel.Type)),
			zap.String("message_id", message.MessageID),
			zap.Error(err),
		)
		return apperrors.SendError("failed to send message via channel", err)
	}

	return nil
}

func (ss *SenderService) TestConnection(channel *models.Channel) error {
	adapter, exists := ss.adapters[channel.Type]
	if !exists {
		return apperrors.ChannelError("no adapter registered for channel type", nil)
	}

	return adapter.TestConnection(channel.Config)
}
