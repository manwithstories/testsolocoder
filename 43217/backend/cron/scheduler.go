package cron

import (
	"fmt"
	"log"
	"time"

	"health-platform/config"
	"health-platform/models"
	"health-platform/repository"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron            *cron.Cron
	appointmentRepo *repository.AppointmentRepository
	reminderRepo    *repository.RecheckReminderRepository
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:            cron.New(),
		appointmentRepo: repository.NewAppointmentRepository(),
		reminderRepo:    repository.NewRecheckReminderRepository(),
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("0 9 * * *", s.RemindAppointment)
	
	s.cron.AddFunc("0 0 * * *", s.ExpireAppointments)
	
	s.cron.AddFunc("0 10 * * *", s.SendRecheckReminders)
	
	s.cron.AddFunc("0 0 1 * *", s.ClearExpiredCache)

	s.cron.Start()
	fmt.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) RemindAppointment() {
	fmt.Println("Running appointment reminder task...")
	
	appointments, err := s.appointmentRepo.FindNeedRemindAppointments(7)
	if err != nil {
		log.Printf("Error finding appointments to remind: %v", err)
		return
	}

	for _, appt := range appointments {
		if appt.Employee.User != nil {
			message := fmt.Sprintf("您预约的%s体检将于%s进行，请按时参加。机构：%s，套餐：%s",
				appt.AppointmentDate.Format("2006-01-02"),
				appt.StartTime,
				appt.Agency.Name,
				appt.Package.Name,
			)
			fmt.Printf("Sending reminder to %s: %s\n", appt.Employee.User.Phone, message)
		}
		
		s.appointmentRepo.MarkAsReminded(appt.ID)
	}

	fmt.Printf("Sent reminders for %d appointments\n", len(appointments))
}

func (s *Scheduler) ExpireAppointments() {
	fmt.Println("Running appointment expiration task...")
	
	appointments, err := s.appointmentRepo.FindExpiredAppointments()
	if err != nil {
		log.Printf("Error finding expired appointments: %v", err)
		return
	}

	for _, appt := range appointments {
		s.appointmentRepo.UpdateStatus(appt.ID, models.AppointmentStatusExpired)
		fmt.Printf("Expired appointment: %s\n", appt.AppointmentNo)
	}

	fmt.Printf("Expired %d appointments\n", len(appointments))
}

func (s *Scheduler) SendRecheckReminders() {
	fmt.Println("Running recheck reminder task...")
	
	reminders, err := s.reminderRepo.GetNeedSendReminders()
	if err != nil {
		log.Printf("Error finding reminders to send: %v", err)
		return
	}

	for _, reminder := range reminders {
		fmt.Printf("Sending recheck reminder: %s\n", reminder.Content)
		s.reminderRepo.MarkAsSent(reminder.ID)
	}

	fmt.Printf("Sent %d recheck reminders\n", len(reminders))
}

func (s *Scheduler) ClearExpiredCache() {
	fmt.Println("Running cache cleanup task...")
	
	patterns := []string{
		"company_stats_*",
		"hot_packages_*",
	}

	for _, pattern := range patterns {
		keys, err := config.RedisClient.Keys(config.Ctx, pattern).Result()
		if err != nil {
			continue
		}
		
		for _, key := range keys {
			ttl := config.RedisClient.TTL(config.Ctx, key).Val()
			if ttl == -1 || ttl < time.Hour {
				config.RedisClient.Del(config.Ctx, key)
			}
		}
	}

	fmt.Println("Cache cleanup completed")
}
