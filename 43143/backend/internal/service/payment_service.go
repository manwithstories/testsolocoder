package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"skillshare/internal/models"
	"skillshare/internal/repository"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	bookingRepo *repository.BookingRepository
	escrowDays  int
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, bookingRepo *repository.BookingRepository, escrowDays int) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		bookingRepo: bookingRepo,
		escrowDays:  escrowDays,
	}
}

type CreatePaymentInput struct {
	BookingID uuid.UUID            `json:"booking_id"`
	Method    models.PaymentMethod `json:"method"`
}

func (s *PaymentService) CreatePayment(payerID uuid.UUID, input *CreatePaymentInput) (*models.Payment, error) {
	booking, err := s.bookingRepo.FindByID(input.BookingID)
	if err != nil {
		return nil, errors.New("预约不存在")
	}

	if booking.StudentID != payerID {
		return nil, errors.New("无权限支付")
	}

	if booking.PaymentStatus == models.PaymentStatusPaid || booking.PaymentStatus == models.PaymentStatusHeld {
		return nil, errors.New("已支付")
	}

	platformFee := booking.Price * 0.05
	netAmount := booking.Price - platformFee

	payment := &models.Payment{
		BookingID:      booking.ID,
		PayerID:        payerID,
		PayeeID:        booking.TeacherID,
		Amount:         booking.Price,
		Method:         input.Method,
		Type:           models.TransactionTypePayment,
		Status:         models.TransactionStatusCompleted,
		PlatformFee:    platformFee,
		NetAmount:      netAmount,
		EscrowReleaseAt: func() *time.Time {
			t := time.Now().AddDate(0, 0, s.escrowDays)
			return &t
		}(),
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, errors.New("创建支付失败")
	}

	booking.PaymentStatus = models.PaymentStatusHeld
	s.bookingRepo.Update(booking)

	return payment, nil
}

func (s *PaymentService) GetPayment(id uuid.UUID) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("支付记录不存在")
	}
	return payment, nil
}

func (s *PaymentService) GetUserPayments(userID uuid.UUID, page, pageSize int) ([]*models.Payment, int64, error) {
	return s.paymentRepo.GetUserPayments(userID, page, pageSize)
}

func (s *PaymentService) ReleaseEscrowPayments() error {
	return s.paymentRepo.ReleaseEscrowPayments()
}

func (s *PaymentService) GetWallet(userID uuid.UUID) (*models.Wallet, error) {
	wallet, err := s.paymentRepo.GetUserWallet(userID)
	if err != nil {
		wallet = &models.Wallet{
			UserID:   userID,
			Balance:  0,
			Currency: "CNY",
		}
		if err := s.paymentRepo.CreateWallet(wallet); err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

func (s *PaymentService) Withdraw(userID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("提现金额必须大于0")
	}

	wallet, err := s.GetWallet(userID)
	if err != nil {
		return err
	}

	if wallet.Balance < amount {
		return errors.New("余额不足")
	}

	return s.paymentRepo.UpdateWallet(userID, -amount)
}

type ScheduleService struct {
	scheduleRepo *repository.ScheduleRepository
}

func NewScheduleService(scheduleRepo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{scheduleRepo: scheduleRepo}
}

type CreateScheduleInput struct {
	Type         models.ScheduleType `json:"type"`
	DayOfWeek    models.DayOfWeek    `json:"day_of_week"`
	SpecificDate *time.Time          `json:"specific_date"`
	StartTime    string              `json:"start_time"`
	EndTime      string              `json:"end_time"`
	IsRecurring  bool                `json:"is_recurring"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
}

func (s *ScheduleService) CreateSchedule(userID uuid.UUID, input *CreateScheduleInput) (*models.Schedule, error) {
	schedule := &models.Schedule{
		UserID:       userID,
		Type:         input.Type,
		DayOfWeek:    input.DayOfWeek,
		SpecificDate: input.SpecificDate,
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		IsRecurring:  input.IsRecurring,
		Title:        input.Title,
		Description:  input.Description,
		IsActive:     true,
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, errors.New("创建日程失败")
	}
	return schedule, nil
}

func (s *ScheduleService) UpdateSchedule(userID, id uuid.UUID, input *CreateScheduleInput) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("日程不存在")
	}

	if schedule.UserID != userID {
		return nil, errors.New("无权限操作")
	}

	schedule.Type = input.Type
	schedule.DayOfWeek = input.DayOfWeek
	schedule.SpecificDate = input.SpecificDate
	schedule.StartTime = input.StartTime
	schedule.EndTime = input.EndTime
	schedule.IsRecurring = input.IsRecurring
	schedule.Title = input.Title
	schedule.Description = input.Description

	if err := s.scheduleRepo.Update(schedule); err != nil {
		return nil, errors.New("更新日程失败")
	}
	return schedule, nil
}

func (s *ScheduleService) DeleteSchedule(userID, id uuid.UUID) error {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return errors.New("日程不存在")
	}

	if schedule.UserID != userID {
		return errors.New("无权限操作")
	}

	return s.scheduleRepo.Delete(id)
}

func (s *ScheduleService) GetUserSchedules(userID uuid.UUID) ([]*models.Schedule, error) {
	return s.scheduleRepo.GetUserSchedules(userID)
}

func (s *ScheduleService) GetUserAvailability(userID uuid.UUID, dayOfWeek models.DayOfWeek) ([]*models.Schedule, error) {
	return s.scheduleRepo.GetUserAvailability(userID, dayOfWeek)
}

type StatsService struct {
	bookingRepo *repository.BookingRepository
	paymentRepo *repository.PaymentRepository
}

func NewStatsService(bookingRepo *repository.BookingRepository, paymentRepo *repository.PaymentRepository) *StatsService {
	return &StatsService{
		bookingRepo: bookingRepo,
		paymentRepo: paymentRepo,
	}
}

type TeacherStats struct {
	TeachingHours  float64 `json:"teaching_hours"`
	StudentCount   int64   `json:"student_count"`
	TotalIncome    float64 `json:"total_income"`
	AvgRating      float64 `json:"avg_rating"`
}

type MonthlyReport struct {
	Month         string  `json:"month"`
	TeachingHours float64 `json:"teaching_hours"`
	StudentCount  int64   `json:"student_count"`
	Income        float64 `json:"income"`
	Bookings      int64   `json:"bookings"`
}

func (s *StatsService) GetTeacherStats(teacherID uuid.UUID, startDate, endDate time.Time) (*TeacherStats, error) {
	bookings, _, err := s.bookingRepo.ListByUser(teacherID, "teacher", 1, 10000, nil)
	if err != nil {
		return nil, err
	}

	var stats TeacherStats
	for _, booking := range bookings {
		if booking.Status == models.BookingStatusCompleted {
			stats.TeachingHours += booking.ScheduledEnd.Sub(booking.ScheduledStart).Hours()
			stats.TotalIncome += booking.TeacherEarnings
			stats.StudentCount++
		}
	}

	return &stats, nil
}

func (s *StatsService) GetMonthlyReport(teacherID uuid.UUID, year int, month int) (*MonthlyReport, error) {
	bookings, _, err := s.bookingRepo.ListByUser(teacherID, "teacher", 1, 10000, nil)
	if err != nil {
		return nil, err
	}

	var report MonthlyReport
	report.Month = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")

	for _, booking := range bookings {
		if booking.Status == models.BookingStatusCompleted &&
			booking.ScheduledStart.Year() == year &&
			int(booking.ScheduledStart.Month()) == month {
			report.TeachingHours += booking.ScheduledEnd.Sub(booking.ScheduledStart).Hours()
			report.Income += booking.TeacherEarnings
			report.Bookings++
		}
	}

	return &report, nil
}
