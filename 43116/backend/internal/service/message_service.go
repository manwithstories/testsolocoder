package service

import (
	"car-rental/internal/config"
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"encoding/json"
	"time"

	cachedb "car-rental/internal/config"

	"gopkg.in/gomail.v2"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	cfg         *config.EmailConfig
}

func NewMessageService(cfg *config.EmailConfig) *MessageService {
	return &MessageService{
		messageRepo: repository.NewMessageRepository(),
		cfg:         cfg,
	}
}

type NotificationMessage struct {
	Type      string `json:"type"`
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	RelatedID *uint  `json:"related_id"`
}

func (s *MessageService) SendMessage(userID uint, msgType model.MessageType, title, content string, relatedID *uint) error {
	message := &model.Message{
		UserID:    userID,
		Type:      msgType,
		Title:     title,
		Content:   content,
		RelatedID: relatedID,
		IsRead:    false,
	}
	return s.messageRepo.Create(message)
}

func (s *MessageService) GetUserMessages(userID uint, page, pageSize int) ([]model.Message, int64, error) {
	return s.messageRepo.FindByUserID(userID, page, pageSize)
}

func (s *MessageService) GetUnreadCount(userID uint) (int64, error) {
	return s.messageRepo.GetUnreadCount(userID)
}

func (s *MessageService) MarkAsRead(id uint, userID uint) error {
	return s.messageRepo.MarkAsRead(id, userID)
}

func (s *MessageService) MarkAllAsRead(userID uint) error {
	return s.messageRepo.MarkAllAsRead(userID)
}

func (s *MessageService) DeleteMessage(id uint, userID uint) error {
	return s.messageRepo.Delete(id, userID)
}

func (s *MessageService) SendEmailNotification(msg *NotificationMessage) error {
	notification := &model.Notification{
		UserID:  msg.UserID,
		Type:    msg.Type,
		Title:   msg.Title,
		Content: msg.Content,
		Channel: "email",
		IsSent:  false,
	}

	_ = s.messageRepo.CreateNotification(notification)

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.User)
	m.SetHeader("To", msg.Email)
	m.SetHeader("Subject", msg.Title)
	m.SetBody("text/html", msg.Content)

	d := gomail.NewDialer(s.cfg.Host, s.cfg.Port, s.cfg.User, s.cfg.Password)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		notification.ErrorMsg = err.Error()
		s.messageRepo.UpdateNotification(notification)
		return err
	}

	now := time.Now()
	notification.IsSent = true
	notification.SentAt = &now
	return s.messageRepo.UpdateNotification(notification)
}

func (s *MessageService) QueueNotification(msg *NotificationMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return cachedb.LPush("notification_queue", string(data))
}

func (s *MessageService) ProcessNotificationQueue() {
	for {
		data, err := cachedb.BLPop("notification_queue", 5*time.Second)
		if err != nil {
			continue
		}

		var msg NotificationMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		_ = s.SendMessage(msg.UserID, model.MessageType(msg.Type), msg.Title, msg.Content, msg.RelatedID)

		if msg.Email != "" {
			_ = s.SendEmailNotification(&msg)
		}
	}
}

func (s *MessageService) SendBookingConfirmation(userID uint, email string, bookingNo string, relatedID uint) {
	content := "<h3>预订确认</h3><p>您的预订已确认，预订号：" + bookingNo + "</p>"
	msg := &NotificationMessage{
		Type:      "booking",
		UserID:    userID,
		Email:     email,
		Title:     "预订确认通知",
		Content:   content,
		RelatedID: &relatedID,
	}
	_ = s.QueueNotification(msg)
}

func (s *MessageService) SendPickupReminder(userID uint, email string, bookingNo string, relatedID uint) {
	content := "<h3>取车提醒</h3><p>您预订的车辆即将到期取车，请按时前往取车门店。预订号：" + bookingNo + "</p>"
	msg := &NotificationMessage{
		Type:      "pickup",
		UserID:    userID,
		Email:     email,
		Title:     "取车提醒通知",
		Content:   content,
		RelatedID: &relatedID,
	}
	_ = s.QueueNotification(msg)
}

func (s *MessageService) SendReturnReminder(userID uint, email string, bookingNo string, relatedID uint) {
	content := "<h3>还车提醒</h3><p>您的车辆即将到期归还，请按时前往还车门店。预订号：" + bookingNo + "</p>"
	msg := &NotificationMessage{
		Type:      "return",
		UserID:    userID,
		Email:     email,
		Title:     "还车提醒通知",
		Content:   content,
		RelatedID: &relatedID,
	}
	_ = s.QueueNotification(msg)
}