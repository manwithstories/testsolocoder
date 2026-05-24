package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"business-registration-platform/config"
	"business-registration-platform/database"
	"business-registration-platform/models"

	"gopkg.in/gomail.v2"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

type SendNotificationRequest struct {
	UserID        uint                      `json:"userId"`
	ApplicationID *uint                     `json:"applicationId"`
	Type          models.NotificationType   `json:"type"`
	TemplateCode  string                    `json:"templateCode"`
	Variables     map[string]string         `json:"variables"`
}

func (s *NotificationService) SendNotification(req *SendNotificationRequest) error {
	var template models.NotificationTemplate
	if err := database.DB.Where("code = ? AND is_active = ?", req.TemplateCode, true).First(&template).Error; err != nil {
		return errors.New("notification template not found")
	}

	title := template.Title
	content := template.Content

	for key, value := range req.Variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		title = strings.ReplaceAll(title, placeholder, value)
		content = strings.ReplaceAll(content, placeholder, value)
	}

	notification := &models.Notification{
		UserID:        req.UserID,
		ApplicationID: req.ApplicationID,
		Type:          req.Type,
		Title:         title,
		Content:       content,
		Status:        models.NotificationStatusPending,
	}

	if err := database.DB.Create(notification).Error; err != nil {
		return err
	}

	go s.sendNotificationAsync(notification)

	return nil
}

func (s *NotificationService) sendNotificationAsync(notification *models.Notification) {
	now := time.Now()

	switch notification.Type {
	case models.NotificationTypeEmail:
		var user models.User
		if err := database.DB.First(&user, notification.UserID).Error; err == nil && user.Email != "" {
			if err := sendEmail(user.Email, notification.Title, notification.Content); err == nil {
				notification.Status = models.NotificationStatusSent
				notification.SentAt = &now
			} else {
				notification.Status = models.NotificationStatusFailed
			}
		}
	case models.NotificationTypeSMS:
		var user models.User
		if err := database.DB.First(&user, notification.UserID).Error; err == nil && user.Phone != "" {
			if err := sendSMS(user.Phone, notification.Content); err == nil {
				notification.Status = models.NotificationStatusSent
				notification.SentAt = &now
			} else {
				notification.Status = models.NotificationStatusFailed
			}
		}
	default:
		notification.Status = models.NotificationStatusSent
		notification.SentAt = &now
	}

	database.DB.Save(notification)
}

func sendEmail(to, subject, body string) error {
	cfg := config.AppConfig.Email
	if cfg.Host == "" {
		return errors.New("email not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	return d.DialAndSend(m)
}

func sendSMS(phone, content string) error {
	return nil
}

func (s *NotificationService) GetUserNotifications(userID uint, page, pageSize int, isRead *bool) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	db := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID)

	if isRead != nil {
		db = db.Where("is_read = ?", *isRead)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Application").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (s *NotificationService) MarkAsRead(userID, notificationID uint) error {
	now := time.Now()
	return database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
			"status":  models.NotificationStatusRead,
		}).Error
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	now := time.Now()
	return database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
			"status":  models.NotificationStatusRead,
		}).Error
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (s *NotificationService) GetNotificationTemplates() ([]models.NotificationTemplate, error) {
	var templates []models.NotificationTemplate
	if err := database.DB.Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (s *NotificationService) CreateNotificationTemplate(template *models.NotificationTemplate) error {
	return database.DB.Create(template).Error
}

func (s *NotificationService) UpdateNotificationTemplate(id uint, data map[string]interface{}) error {
	return database.DB.Model(&models.NotificationTemplate{}).Where("id = ?", id).Updates(data).Error
}

func (s *NotificationService) DeleteNotificationTemplate(id uint) error {
	return database.DB.Delete(&models.NotificationTemplate{}, id).Error
}

func (s *NotificationService) SendApplicationStatusNotification(applicationID uint, status string) error {
	var application models.Application
	if err := database.DB.Preload("Entrepreneur").First(&application, applicationID).Error; err != nil {
		return err
	}

	if application.Entrepreneur == nil {
		return nil
	}

	statusText := map[string]string{
		"payment_pending": "待支付",
		"pending_review":  "待审核",
		"reviewing":       "审核中",
		"processing":      "处理中",
		"completed":       "已完成",
		"rejected":        "已驳回",
		"cancelled":       "已取消",
	}

	variables := map[string]string{
		"company_name":    application.CompanyName,
		"application_no":  application.ApplicationNo,
		"status":          statusText[status],
	}

	return s.SendNotification(&SendNotificationRequest{
		UserID:        application.EntrepreneurID,
		ApplicationID: &applicationID,
		Type:          models.NotificationTypeSystem,
		TemplateCode:  "APPLICATION_STATUS_CHANGE",
		Variables:     variables,
	})
}

func (s *NotificationService) SendStepCompleteNotification(applicationID uint, stepName string) error {
	var application models.Application
	if err := database.DB.Preload("Entrepreneur").First(&application, applicationID).Error; err != nil {
		return err
	}

	if application.Entrepreneur == nil {
		return nil
	}

	variables := map[string]string{
		"company_name":   application.CompanyName,
		"application_no": application.ApplicationNo,
		"step_name":      stepName,
	}

	return s.SendNotification(&SendNotificationRequest{
		UserID:        application.EntrepreneurID,
		ApplicationID: &applicationID,
		Type:          models.NotificationTypeSystem,
		TemplateCode:  "STEP_COMPLETED",
		Variables:     variables,
	})
}
