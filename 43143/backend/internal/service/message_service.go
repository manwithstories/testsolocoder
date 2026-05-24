package service

import (
	"errors"

	"github.com/google/uuid"
	"skillshare/internal/models"
	"skillshare/internal/repository"
	"skillshare/pkg/encryption"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	encryptKey  string
}

func NewMessageService(messageRepo *repository.MessageRepository, encryptKey string) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		encryptKey:  encryptKey,
	}
}

type SendMessageInput struct {
	ReceiverID uuid.UUID          `json:"receiver_id"`
	Type       models.MessageType `json:"type"`
	Content    string             `json:"content"`
	FileURL    string             `json:"file_url"`
	FileName   string             `json:"file_name"`
	FileSize   int64              `json:"file_size"`
}

func (s *MessageService) SendMessage(senderID uuid.UUID, input *SendMessageInput) (*models.Message, error) {
	content := input.Content
	if input.Type == models.MessageTypeText {
		encrypted, err := encryption.EncryptMessage(input.Content, s.encryptKey)
		if err != nil {
			return nil, errors.New("消息加密失败")
		}
		content = encrypted
	}

	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: input.ReceiverID,
		Type:       input.Type,
		Content:    content,
		Encrypted:  input.Type == models.MessageTypeText,
		FileURL:    input.FileURL,
		FileName:   input.FileName,
		FileSize:   input.FileSize,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, errors.New("发送消息失败")
	}

	return message, nil
}

func (s *MessageService) GetMessages(senderID, receiverID uuid.UUID, page, pageSize int) ([]*models.Message, int64, error) {
	messages, total, err := s.messageRepo.GetMessages(senderID, receiverID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	for _, msg := range messages {
		if msg.Encrypted && msg.Type == models.MessageTypeText {
			decrypted, err := encryption.DecryptMessage(msg.Content, s.encryptKey)
			if err == nil {
				msg.Content = decrypted
			}
		}
	}

	s.messageRepo.MarkAsRead(receiverID, senderID)

	return messages, total, nil
}

func (s *MessageService) GetConversations(userID uuid.UUID) ([]*models.Message, error) {
	return s.messageRepo.GetConversations(userID)
}

func (s *MessageService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.messageRepo.GetUnreadCount(userID)
}

func (s *MessageService) MarkAsRead(senderID, receiverID uuid.UUID) error {
	return s.messageRepo.MarkAsRead(senderID, receiverID)
}
