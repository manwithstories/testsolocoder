package services

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"meeting-room/internal/config"
	"meeting-room/internal/models"
	"meeting-room/internal/repositories"
	"meeting-room/internal/utils"

	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
	userRepo         *repositories.UserRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepo: repositories.NewNotificationRepository(),
		userRepo:         repositories.NewUserRepository(),
	}
}

func (s *NotificationService) ProcessNotificationQueue() {
	for {
		result, err := utils.RedisBRPop("notification_queue", 30*time.Second)
		if err != nil || len(result) < 2 {
			continue
		}

		notificationID, err := strconv.ParseUint(result[1], 10, 32)
		if err != nil {
			continue
		}

		notification, err := s.notificationRepo.FindByID(uint(notificationID))
		if err != nil {
			continue
		}

		if notification.Status == models.NotificationStatusSent {
			continue
		}

		s.sendNotification(notification)
	}
}

func (s *NotificationService) sendNotification(notification *models.Notification) {
	user, err := s.userRepo.FindByID(notification.UserID)
	if err != nil {
		utils.Logger.Errorf("查找通知用户失败: %v", err)
		s.notificationRepo.MarkFailed(notification.ID, err.Error())
		return
	}

	err = s.sendEmail(user.Email, notification.Subject, notification.Content)
	if err != nil {
		utils.Logger.Errorf("发送邮件失败: %v", err)
		if notification.RetryCount < 2 {
			notification.RetryCount++
			s.notificationRepo.Update(notification)
			utils.RedisLPush("notification_queue", notification.ID)
		} else {
			s.notificationRepo.MarkFailed(notification.ID, err.Error())
		}
		return
	}

	s.notificationRepo.MarkSent(notification.ID)
	utils.Logger.Infof("通知发送成功: %d -> %s", notification.ID, user.Email)
}

func (s *NotificationService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Cfg.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		config.Cfg.Email.SMTPHost,
		config.Cfg.Email.SMTPPort,
		config.Cfg.Email.Username,
		config.Cfg.Email.Password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

func (s *NotificationService) CheckAndSendReminders() {
	bookings, err := repositories.NewBookingRepository().GetUpcomingReminders()
	if err != nil {
		utils.Logger.Errorf("获取即将开始的预订失败: %v", err)
		return
	}

	for _, booking := range bookings {
		_, err := s.userRepo.FindByID(booking.UserID)
		if err != nil {
			continue
		}

		roomName := "会议室"
		if booking.Room != nil {
			roomName = booking.Room.Name
		}

		notification := &models.Notification{
			UserID:    booking.UserID,
			BookingID: &booking.ID,
			Type:      models.NotificationTypeBookingReminder,
			Subject:   fmt.Sprintf("会议即将开始: %s", booking.Title),
			Content:   fmt.Sprintf("温馨提醒：您的会议 \"%s\" 将在1小时后开始。\n会议室: %s\n时间: %s - %s", booking.Title, roomName, booking.StartTime.Format("2006-01-02 15:04"), booking.EndTime.Format("15:04")),
		}
		s.notificationRepo.Create(notification)
		utils.RedisLPush("notification_queue", notification.ID)

		repositories.NewBookingRepository().MarkReminded(booking.ID)
	}
}
