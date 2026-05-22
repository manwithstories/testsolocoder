package services

import (
	"sync"
	"time"

	"travel-planner/internal/database"
	"travel-planner/internal/logger"
	"travel-planner/internal/models"
)

type ReminderScheduler struct {
	emailService *EmailService
	ticker       *time.Ticker
	stopChan     chan struct{}
	mu           sync.Mutex
	running      bool
}

var (
	schedulerInstance *ReminderScheduler
	schedulerOnce     sync.Once
)

func GetReminderScheduler(emailService *EmailService) *ReminderScheduler {
	schedulerOnce.Do(func() {
		schedulerInstance = &ReminderScheduler{
			emailService: emailService,
			stopChan:     make(chan struct{}),
		}
	})
	return schedulerInstance
}

func (s *ReminderScheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		logger.Info("Reminder scheduler is already running")
		return
	}

	s.ticker = time.NewTicker(1 * time.Minute)
	s.running = true

	go func() {
		logger.Info("Reminder scheduler started")
		s.checkReminders()

		for {
			select {
			case <-s.ticker.C:
				s.checkReminders()
			case <-s.stopChan:
				logger.Info("Reminder scheduler stopped")
				return
			}
		}
	}()
}

func (s *ReminderScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.ticker.Stop()
	close(s.stopChan)
	s.running = false
}

func (s *ReminderScheduler) checkReminders() {
	now := time.Now()

	var reminders []models.Reminder
	if err := database.DB.Where(
		"reminder_time <= ? AND is_sent = ?",
		now,
		false,
	).Preload("User").Preload("Plan").Find(&reminders).Error; err != nil {
		logger.Errorf("Failed to fetch due reminders: %v", err)
		return
	}

	if len(reminders) == 0 {
		return
	}

	logger.Infof("Found %d due reminders to process", len(reminders))

	for _, reminder := range reminders {
		s.processReminder(reminder)
	}
}

func (s *ReminderScheduler) processReminder(reminder models.Reminder) {
	logger.Infof("Processing reminder: %s for user %s", reminder.ID, reminder.UserID)

	planTitle := ""
	if reminder.Plan != nil {
		planTitle = reminder.Plan.Title
	}

	if reminder.Channel == "email" && reminder.User != nil && reminder.User.Email != "" {
		err := s.emailService.SendReminderEmail(
			reminder.User.Email,
			reminder.Title,
			reminder.Description,
			planTitle,
		)
		if err != nil {
			logger.Errorf("Failed to send email reminder %s: %v", reminder.ID, err)
		} else {
			logger.Infof("Email reminder sent successfully: %s", reminder.ID)
		}
	}

	if err := database.DB.Model(&reminder).Updates(map[string]interface{}{
		"is_sent": true,
		"sent_at": time.Now(),
	}).Error; err != nil {
		logger.Errorf("Failed to mark reminder as sent: %v", err)
	} else {
		logger.Infof("Reminder marked as sent: %s", reminder.ID)
	}
}
