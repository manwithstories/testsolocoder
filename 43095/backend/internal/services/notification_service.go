package services

import (
	"crypto/tls"
	"errors"
	"fmt"
	"medical-platform/internal/config"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"
	"medical-platform/pkg/utils"
	"net/smtp"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		db: database.GetDB(),
	}
}

type SendNotificationRequest struct {
	UserID      uint                      `json:"user_id" binding:"required"`
	Type        models.NotificationType   `json:"type" binding:"required"`
	Title       string                    `json:"title" binding:"required,max=200"`
	Content     string                    `json:"content" binding:"required"`
	Channel     models.NotificationChannel `json:"channel" binding:"required"`
	RelatedID   uint                      `json:"related_id"`
	RelatedType string                    `json:"related_type"`
}

type SendEmailRequest struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func (s *NotificationService) SendInAppNotification(req *SendNotificationRequest) (*models.Notification, error) {
	var user models.User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	notification := &models.Notification{
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		Channel:     models.ChannelInApp,
		IsRead:      false,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) SendEmail(req *SendEmailRequest) error {
	if config.AppConfig.SMTPHost == "" {
		return errors.New("SMTP配置未设置")
	}

	smtpHost := config.AppConfig.SMTPHost
	smtpPort := config.AppConfig.SMTPPort
	smtpUser := config.AppConfig.SMTPUser
	smtpPass := config.AppConfig.SMTPPass
	smtpFrom := config.AppConfig.SMTPFrom

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	if err := client.Mail(smtpFrom); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	if err := client.Rcpt(req.To); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("创建邮件数据写入器失败: %w", err)
	}

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=utf-8\r\n\r\n",
		smtpFrom, req.To, req.Subject)

	if _, err := w.Write([]byte(headers + req.Body)); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("关闭邮件写入器失败: %w", err)
	}

	return nil
}

func (s *NotificationService) GetNotificationList(userID uint, page, pageSize int, isRead *bool) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (s *NotificationService) MarkAsRead(userID uint, notificationID uint) error {
	var notification models.Notification
	if err := s.db.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return err
	}

	if notification.IsRead {
		return nil
	}

	now := time.Now()
	if err := s.db.Model(&notification).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": &now,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	now := time.Now()
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) DeleteNotification(userID uint, notificationID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&models.Notification{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("通知不存在")
	}
	return nil
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *NotificationService) SendAppointmentConfirmation(appointment *models.Appointment) error {
	var patient models.Patient
	if err := s.db.Preload("User").First(&patient, appointment.PatientID).Error; err != nil {
		return fmt.Errorf("获取患者信息失败: %w", err)
	}

	var doctor models.Doctor
	if err := s.db.Preload("User").Preload("Department").First(&doctor, appointment.DoctorID).Error; err != nil {
		return fmt.Errorf("获取医生信息失败: %w", err)
	}

	title := "预约成功确认"
	content := fmt.Sprintf("您的预约已成功！\n就诊时间：%s %s-%s\n就诊科室：%s\n就诊医生：%s %s\n请按时就诊。",
		utils.FormatDate(appointment.AppointmentDate),
		appointment.StartTime,
		appointment.EndTime,
		doctor.Department.Name,
		doctor.User.FullName,
		doctor.Title)

	inAppReq := &SendNotificationRequest{
		UserID:      patient.UserID,
		Type:        models.NotificationAppointmentConfirmation,
		Title:       title,
		Content:     content,
		Channel:     models.ChannelInApp,
		RelatedID:   appointment.ID,
		RelatedType: "appointment",
	}

	if _, err := s.SendInAppNotification(inAppReq); err != nil {
		config.Logger.Error("发送站内通知失败", zap.Error(err))
	}

	if patient.User.Email != "" && config.AppConfig.SMTPHost != "" {
		emailBody := fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #3498db;">预约成功确认</h2>
				<p>尊敬的 %s 患者：</p>
				<p>您的预约已成功！</p>
				<div style="background: #f8f9fa; padding: 15px; border-radius: 8px; margin: 20px 0;">
					<p><strong>就诊时间：</strong>%s %s-%s</p>
					<p><strong>就诊科室：</strong>%s</p>
					<p><strong>就诊医生：</strong>%s %s</p>
				</div>
				<p>请按时就诊，如有特殊情况请及时取消预约。</p>
				<p>感谢您的信任！</p>
			</div>
		`, patient.User.FullName,
			utils.FormatDate(appointment.AppointmentDate),
			appointment.StartTime,
			appointment.EndTime,
			doctor.Department.Name,
			doctor.User.FullName,
			doctor.Title)

		emailReq := &SendEmailRequest{
			To:      patient.User.Email,
			Subject: "预约成功确认",
			Body:    emailBody,
		}

		if err := s.SendEmail(emailReq); err != nil {
			config.Logger.Error("发送邮件通知失败", zap.Error(err))
		}
	}

	return nil
}

func (s *NotificationService) SendAppointmentReminders() error {
	tomorrow := time.Now().AddDate(0, 0, 1)
	tomorrowStart := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	tomorrowEnd := tomorrowStart.AddDate(0, 0, 1)

	var appointments []models.Appointment
	if err := s.db.Where("appointment_date >= ? AND appointment_date < ? AND status IN (?)",
		tomorrowStart, tomorrowEnd,
		[]models.AppointmentStatus{models.AppointmentPending, models.AppointmentConfirmed}).
		Preload("Patient.User").
		Preload("Doctor.User").
		Preload("Doctor.Department").
		Find(&appointments).Error; err != nil {
		return fmt.Errorf("获取预约列表失败: %w", err)
	}

	for _, apt := range appointments {
		if err := s.sendSingleAppointmentReminder(&apt); err != nil {
			config.Logger.Error("发送预约提醒失败", zap.Uint("appointment_id", apt.ID), zap.Error(err))
		}
	}

	return nil
}

func (s *NotificationService) sendSingleAppointmentReminder(appointment *models.Appointment) error {
	patient := appointment.Patient
	doctor := appointment.Doctor

	title := "就诊提醒"
	content := fmt.Sprintf("明天您有一个预约！\n就诊时间：%s %s-%s\n就诊科室：%s\n就诊医生：%s %s\n请按时就诊。",
		utils.FormatDate(appointment.AppointmentDate),
		appointment.StartTime,
		appointment.EndTime,
		doctor.Department.Name,
		doctor.User.FullName,
		doctor.Title)

	inAppReq := &SendNotificationRequest{
		UserID:      patient.UserID,
		Type:        models.NotificationAppointmentReminder,
		Title:       title,
		Content:     content,
		Channel:     models.ChannelInApp,
		RelatedID:   appointment.ID,
		RelatedType: "appointment",
	}

	if _, err := s.SendInAppNotification(inAppReq); err != nil {
		return err
	}

	if patient.User.Email != "" && config.AppConfig.SMTPHost != "" {
		emailBody := fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #e74c3c;">就诊提醒</h2>
				<p>尊敬的 %s 患者：</p>
				<p>明天您有一个预约，请注意就诊时间！</p>
				<div style="background: #f8f9fa; padding: 15px; border-radius: 8px; margin: 20px 0;">
					<p><strong>就诊时间：</strong>%s %s-%s</p>
					<p><strong>就诊科室：</strong>%s</p>
					<p><strong>就诊医生：</strong>%s %s</p>
				</div>
				<p>请按时就诊，如有特殊情况请及时取消预约。</p>
			</div>
		`, patient.User.FullName,
			utils.FormatDate(appointment.AppointmentDate),
			appointment.StartTime,
			appointment.EndTime,
			doctor.Department.Name,
			doctor.User.FullName,
			doctor.Title)

		emailReq := &SendEmailRequest{
			To:      patient.User.Email,
			Subject: "就诊提醒",
			Body:    emailBody,
		}

		if err := s.SendEmail(emailReq); err != nil {
			return err
		}
	}

	return nil
}
