package services

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"consultation-platform/config"
	"consultation-platform/models"
	"consultation-platform/repositories"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
	templateRepo     *repositories.NotificationTemplateRepository
	userRepo         *repositories.UserRepository
	cfg              *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		notificationRepo: repositories.NewNotificationRepository(),
		templateRepo:     repositories.NewNotificationTemplateRepository(),
		userRepo:         repositories.NewUserRepository(),
		cfg:              cfg,
	}
}

type NotificationServiceInterface interface {
	SendNotification(userID uuid.UUID, notificationType models.NotificationType, title string, content string, data map[string]interface{}) error
	SendEmailNotification(userID uuid.UUID, notificationType models.NotificationType, data map[string]interface{}) error
	GetNotifications(userID uuid.UUID, page, pageSize int, onlyUnread bool) ([]models.Notification, int64, error)
	MarkAsRead(userID, notificationID uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	GetUnreadCount(userID uuid.UUID) (int64, error)
	GetTemplates(page, pageSize int) ([]models.NotificationTemplate, int64, error)
	UpdateTemplate(id uuid.UUID, title, content string) error
}

func (s *NotificationService) SendNotification(userID uuid.UUID, notificationType models.NotificationType, title string, content string, data map[string]interface{}) error {
	dataStr := fmt.Sprintf("%v", data)
	notification := &models.Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Content: content,
		Data:    dataStr,
		IsRead:  false,
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) SendEmailNotification(userID uuid.UUID, notificationType models.NotificationType, data map[string]interface{}) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	template, err := s.templateRepo.FindByType(string(notificationType))
	if err != nil {
		return nil
	}

	emailBody := s.renderTemplate(template.Content, data)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.cfg.Email.FromName, s.cfg.Email.Username))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", s.renderTemplate(template.Title, data))
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(
		s.cfg.Email.SMTPHost,
		s.cfg.Email.SMTPPort,
		s.cfg.Email.Username,
		s.cfg.Email.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (s *NotificationService) GetNotifications(userID uuid.UUID, page, pageSize int, onlyUnread bool) ([]models.Notification, int64, error) {
	return s.notificationRepo.FindByUserID(userID, page, pageSize, onlyUnread)
}

func (s *NotificationService) MarkAsRead(userID, notificationID uuid.UUID) error {
	notification, err := s.notificationRepo.FindByID(notificationID)
	if err != nil {
		return err
	}

	if notification.UserID != userID {
		return fmt.Errorf("notification does not belong to user")
	}

	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

func (s *NotificationService) GetTemplates(page, pageSize int) ([]models.NotificationTemplate, int64, error) {
	return s.templateRepo.FindAll(page, pageSize)
}

func (s *NotificationService) UpdateTemplate(id uuid.UUID, title, content string) error {
	template, err := s.templateRepo.FindByID(id)
	if err != nil {
		return err
	}

	template.Title = title
	template.Content = content

	return s.templateRepo.Update(template)
}

func (s *NotificationService) renderTemplate(tmplStr string, data map[string]interface{}) string {
	tmpl, err := template.New("notification").Parse(tmplStr)
	if err != nil {
		return tmplStr
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return tmplStr
	}

	return buf.String()
}

func (s *NotificationService) NotifyAppointmentSuccess(userID uuid.UUID, appointment *models.Appointment) {
	title := "预约成功"
	content := fmt.Sprintf("您的预约已成功，时间：%s %s-%s",
		appointment.Schedule.Date.Format("2006-01-02"),
		appointment.Schedule.StartTime,
		appointment.Schedule.EndTime,
	)

	data := map[string]interface{}{
		"appointment_id": appointment.ID.String(),
		"date":           appointment.Schedule.Date.Format("2006-01-02"),
		"start_time":     appointment.Schedule.StartTime,
		"end_time":       appointment.Schedule.EndTime,
		"professional":   appointment.Professional.FullName,
		"service":        appointment.Service.Title,
		"amount":         appointment.Payment.Amount,
	}

	_ = s.SendNotification(userID, models.NotificationAppointmentSuccess, title, content, data)
	_ = s.SendEmailNotification(userID, models.NotificationAppointmentSuccess, data)
}

func (s *NotificationService) NotifyAppointmentCancel(userID uuid.UUID, appointment *models.Appointment, reason string) {
	title := "预约已取消"
	content := fmt.Sprintf("您的预约已取消，原因：%s", reason)

	data := map[string]interface{}{
		"appointment_id": appointment.ID.String(),
		"reason":         reason,
	}

	_ = s.SendNotification(userID, models.NotificationAppointmentCancel, title, content, data)
	_ = s.SendEmailNotification(userID, models.NotificationAppointmentCancel, data)
}

func (s *NotificationService) NotifyAppointmentRemind(userID uuid.UUID, appointment *models.Appointment) {
	title := "预约提醒"
	content := fmt.Sprintf("您有一个预约即将开始：%s %s-%s",
		appointment.Schedule.Date.Format("2006-01-02"),
		appointment.Schedule.StartTime,
		appointment.Schedule.EndTime,
	)

	data := map[string]interface{}{
		"appointment_id": appointment.ID.String(),
		"date":           appointment.Schedule.Date.Format("2006-01-02"),
		"start_time":     appointment.Schedule.StartTime,
		"end_time":       appointment.Schedule.EndTime,
	}

	_ = s.SendNotification(userID, models.NotificationAppointmentRemind, title, content, data)
}

func (s *NotificationService) NotifyPaymentSuccess(userID uuid.UUID, payment *models.Payment) {
	title := "支付成功"
	content := fmt.Sprintf("您的付款已成功，金额：￥%.2f", payment.Amount)

	data := map[string]interface{}{
		"payment_id": payment.ID.String(),
		"amount":     payment.Amount,
	}

	_ = s.SendNotification(userID, models.NotificationPaymentSuccess, title, content, data)
	_ = s.SendEmailNotification(userID, models.NotificationPaymentSuccess, data)
}

func (s *NotificationService) NotifyPaymentRefund(userID uuid.UUID, payment *models.Payment, reason string) {
	title := "退款通知"
	content := fmt.Sprintf("您的退款已处理，金额：￥%.2f，原因：%s", payment.Amount, reason)

	data := map[string]interface{}{
		"payment_id": payment.ID.String(),
		"amount":     payment.Amount,
		"reason":     reason,
	}

	_ = s.SendNotification(userID, models.NotificationPaymentRefund, title, content, data)
	_ = s.SendEmailNotification(userID, models.NotificationPaymentRefund, data)
}

var _ = time.Now
