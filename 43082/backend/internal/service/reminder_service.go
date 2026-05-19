package service

import (
	"errors"
	"fmt"
	"gym-management/config"
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"gym-management/internal/pkg/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReminderService interface {
	Create(reminder *models.Reminder) error
	SendReminder(reminderID uint) error
	ProcessPendingReminders() error
	GenerateMembershipExpireReminders() error
	GenerateCourseReminders() error
	SendReminders() error
}

type reminderService struct {
	db *gorm.DB
}

func NewReminderService() ReminderService {
	return &reminderService{
		db: database.GetDB(),
	}
}

func (s *reminderService) Create(reminder *models.Reminder) error {
	return s.db.Create(reminder).Error
}

func (s *reminderService) SendReminder(reminderID uint) error {
	var reminder models.Reminder
	if err := s.db.First(&reminder, reminderID).Error; err != nil {
		return errors.New("提醒不存在")
	}

	if reminder.Status == 2 {
		return nil
	}

	err := s.send(&reminder)
	now := time.Now()

	if err != nil {
		reminder.RetryCount++
		if reminder.RetryCount >= 3 {
			reminder.Status = 3
		}
		logger.Error("Failed to send reminder", zap.Uint("reminder_id", reminderID), zap.Error(err))
	} else {
		reminder.Status = 2
		reminder.SentTime = &now
		logger.Info("Reminder sent successfully", zap.Uint("reminder_id", reminderID))
	}

	return s.db.Save(&reminder).Error
}

func (s *reminderService) send(reminder *models.Reminder) error {
	var member models.Member
	if err := s.db.First(&member, reminder.MemberID).Error; err != nil {
		return err
	}

	logger.Info("Sending reminder",
		zap.Uint("member_id", member.ID),
		zap.String("member_name", member.Name),
		zap.String("title", reminder.Title),
		zap.String("content", reminder.Content),
		zap.String("phone", member.Phone),
		zap.String("email", member.Email),
	)

	return nil
}

func (s *reminderService) ProcessPendingReminders() error {
	var reminders []models.Reminder
	now := time.Now()

	err := s.db.Where("status = 1 AND schedule_time <= ?", now).
		Order("schedule_time ASC").
		Limit(100).
		Find(&reminders).Error
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		go func(r models.Reminder) {
			_ = s.SendReminder(r.ID)
		}(reminder)
	}

	return nil
}

func (s *reminderService) GenerateMembershipExpireReminders() error {
	var memberships []models.Membership
	sevenDaysLater := time.Now().AddDate(0, 0, 7)

	err := s.db.Where("status = 1 AND end_date <= ? AND end_date > ?", sevenDaysLater, time.Now()).
		Preload("Member").
		Find(&memberships).Error
	if err != nil {
		return err
	}

	for _, membership := range memberships {
		var existingReminder models.Reminder
		result := s.db.Where("member_id = ? AND type = 1 AND DATE(schedule_time) = DATE(?)",
			membership.MemberID, membership.EndDate.AddDate(0, 0, -3)).First(&existingReminder)

		if result.Error == gorm.ErrRecordNotFound && membership.Member != nil {
			reminder := &models.Reminder{
				MemberID:   membership.MemberID,
				Type:       1,
				Title:      "会员卡到期提醒",
				Content:    fmt.Sprintf(config.AppConfig.Reminder.TemplateExpire, membership.EndDate.Format("2006-01-02")),
				ScheduleTime: membership.EndDate.AddDate(0, 0, -3),
				Status:     1,
			}
			_ = s.Create(reminder)
		}
	}

	return nil
}

func (s *reminderService) GenerateCourseReminders() error {
	var bookings []models.Booking
	now := time.Now()
	twoDaysLater := now.AddDate(0, 0, 2)
	threeHoursLater := now.Add(3 * time.Hour)

	err := s.db.Where("status = 1").
		Preload("Schedule").
		Preload("Schedule.Course").
		Preload("Member").
		Joins("JOIN course_schedules ON course_schedules.id = bookings.schedule_id").
		Where("course_schedules.start_time > ? AND course_schedules.start_time <= ?", now, twoDaysLater).
		Find(&bookings).Error
	if err != nil {
		return err
	}

	for _, booking := range bookings {
		if booking.Schedule == nil || booking.Member == nil {
			continue
		}

		scheduleTime := booking.Schedule.StartTime
		courseName := booking.Schedule.Course.Name

		oneDayBefore := scheduleTime.AddDate(0, 0, -1)
		if oneDayBefore.After(now) && oneDayBefore.Before(twoDaysLater) {
			s.createCourseReminderIfNotExists(booking.MemberID, booking.ScheduleID, 2,
				oneDayBefore, courseName, scheduleTime.Format("15:04"))
		}

		twoHoursBefore := scheduleTime.Add(-2 * time.Hour)
		if twoHoursBefore.After(now) && twoHoursBefore.Before(threeHoursLater) {
			s.createCourseReminderIfNotExists(booking.MemberID, booking.ScheduleID, 3,
				twoHoursBefore, courseName, scheduleTime.Format("15:04"))
		}
	}

	return nil
}

func (s *reminderService) createCourseReminderIfNotExists(memberID, scheduleID uint, reminderType int,
	scheduleTime time.Time, courseName, startTime string) {
	var existingReminder models.Reminder
	result := s.db.Where("member_id = ? AND type = ? AND schedule_id = ?", memberID, reminderType, scheduleID).First(&existingReminder)

	if result.Error == gorm.ErrRecordNotFound {
		var title, template string
		if reminderType == 2 {
			title = "课程预约提醒（1天前）"
			template = config.AppConfig.Reminder.TemplateCourse1D
		} else {
			title = "课程即将开始提醒"
			template = config.AppConfig.Reminder.TemplateCourse2H
		}

		reminder := &models.Reminder{
			MemberID:     memberID,
			Type:         reminderType,
			Title:        title,
			Content:      fmt.Sprintf(template, courseName, startTime),
			ScheduleTime: scheduleTime,
			Status:       1,
		}
		_ = s.Create(reminder)
	}
}

func (s *reminderService) SendReminders() error {
	if err := s.GenerateMembershipExpireReminders(); err != nil {
		logger.Error("Failed to generate membership expire reminders", zap.Error(err))
	}

	if err := s.GenerateCourseReminders(); err != nil {
		logger.Error("Failed to generate course reminders", zap.Error(err))
	}

	if err := s.ProcessPendingReminders(); err != nil {
		logger.Error("Failed to process pending reminders", zap.Error(err))
	}

	return nil
}
