package service

import (
	"car-rental/internal/repository"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	bookingRepo       *repository.BookingRepository
	carRepo           *repository.CarRepository
	userRepo          *repository.UserRepository
	maintenanceRepo   *repository.MaintenanceRepository
	messageService    *MessageService
	cron              *cron.Cron
}

func NewSchedulerService(messageService *MessageService) *SchedulerService {
	return &SchedulerService{
		bookingRepo:     repository.NewBookingRepository(),
		carRepo:         repository.NewCarRepository(),
		userRepo:        repository.NewUserRepository(),
		maintenanceRepo: repository.NewMaintenanceRepository(),
		messageService:  messageService,
		cron:            cron.New(),
	}
}

func (s *SchedulerService) Start() {
	_, _ = s.cron.AddFunc("0 * * * *", s.checkPickupReminders)

	_, _ = s.cron.AddFunc("30 * * * *", s.checkReturnReminders)

	_, _ = s.cron.AddFunc("0 0 * * *", s.checkMaintenancePlans)

	s.cron.Start()
	log.Println("定时任务调度器已启动")
}

func (s *SchedulerService) Stop() {
	s.cron.Stop()
}

func (s *SchedulerService) checkPickupReminders() {
	now := time.Now()
	reminderTime := now.Add(2 * time.Hour)

	bookings, err := s.bookingRepo.FindUpcomingPickups(now, reminderTime)
	if err != nil {
		log.Printf("查询即将取车的预订失败: %v", err)
		return
	}

	for _, booking := range bookings {
		user, err := s.userRepo.FindByID(booking.UserID)
		if err != nil {
			continue
		}

		if booking.PickupReminderSent {
			continue
		}

		s.messageService.SendPickupReminder(booking.UserID, user.Email, booking.BookingNo, booking.ID)

		_ = s.bookingRepo.MarkPickupReminderSent(booking.ID)
	}
}

func (s *SchedulerService) checkReturnReminders() {
	now := time.Now()
	reminderTime := now.Add(2 * time.Hour)

	bookings, err := s.bookingRepo.FindUpcomingReturns(now, reminderTime)
	if err != nil {
		log.Printf("查询即将还车的预订失败: %v", err)
		return
	}

	for _, booking := range bookings {
		user, err := s.userRepo.FindByID(booking.UserID)
		if err != nil {
			continue
		}

		if booking.ReturnReminderSent {
			continue
		}

		s.messageService.SendReturnReminder(booking.UserID, user.Email, booking.BookingNo, booking.ID)

		_ = s.bookingRepo.MarkReturnReminderSent(booking.ID)
	}
}

func (s *SchedulerService) checkMaintenancePlans() {
	plans, err := s.maintenanceRepo.FindUpcoming()
	if err != nil {
		log.Printf("查询即将到期的维护计划失败: %v", err)
		return
	}

	for _, plan := range plans {
		if plan.Car != nil {
			car := plan.Car
			if car.Status == "available" {
				_ = s.carRepo.UpdateStatus(car.ID, "maintenance")
				log.Printf("车辆 %s 已自动设置为维护状态", car.LicensePlate)
			}
		}
	}
}
