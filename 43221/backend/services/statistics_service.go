package services

import (
	"fmt"
	"time"

	"consultation-platform/models"
	"consultation-platform/repositories"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"github.com/tealeg/xlsx/v3"
	"gorm.io/gorm"
)

type StatisticsService struct {
	appointmentRepo *repositories.AppointmentRepository
	paymentRepo     *repositories.PaymentRepository
	reviewRepo      *repositories.ReviewRepository
	userRepo        *repositories.UserRepository
}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{
		appointmentRepo: repositories.NewAppointmentRepository(),
		paymentRepo:     repositories.NewPaymentRepository(),
		reviewRepo:      repositories.NewReviewRepository(),
		userRepo:        repositories.NewUserRepository(),
	}
}

type AppointmentStats struct {
	Total       int64 `json:"total"`
	Pending     int64 `json:"pending"`
	Confirmed   int64 `json:"confirmed"`
	Completed   int64 `json:"completed"`
	Cancelled   int64 `json:"cancelled"`
	Refunded    int64 `json:"refunded"`
}

type RevenueStats struct {
	TotalRevenue  float64 `json:"total_revenue"`
	PaidCount     int64   `json:"paid_count"`
	RefundedAmount float64 `json:"refunded_amount"`
}

type ReviewStats struct {
	AverageRating float64 `json:"average_rating"`
	TotalReviews  int     `json:"total_reviews"`
}

type ProfessionalStats struct {
	Appointments AppointmentStats `json:"appointments"`
	Revenue      RevenueStats     `json:"revenue"`
	Reviews      ReviewStats      `json:"reviews"`
}

type AdminStats struct {
	TotalUsers      int64            `json:"total_users"`
	TotalClients    int64            `json:"total_clients"`
	TotalProfessionals int64         `json:"total_professionals"`
	TotalServices   int64            `json:"total_services"`
	Appointments    AppointmentStats `json:"appointments"`
	Revenue         RevenueStats     `json:"revenue"`
}

type StatisticsServiceInterface interface {
	GetProfessionalStats(professionalID uuid.UUID, startDate, endDate string) (*ProfessionalStats, error)
	GetAdminStats(startDate, endDate string) (*AdminStats, error)
	ExportAppointments(userID uuid.UUID, userRole string, startDate, endDate string, status string) ([]byte, error)
	ExportRevenue(userID uuid.UUID, userRole string, startDate, endDate string) ([]byte, error)
}

func (s *StatisticsService) GetProfessionalStats(professionalID uuid.UUID, startDate, endDate string) (*ProfessionalStats, error) {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		startTime = time.Now().AddDate(0, -1, 0)
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		endTime = time.Now()
	}

	var stats ProfessionalStats

	db := utils.GetDB()

	db.Model(&models.Appointment{}).
		Where("professional_id = ? AND created_at >= ? AND created_at <= ?", professionalID, startTime, endTime).
		Count(&stats.Appointments.Total)

	db.Model(&models.Appointment{}).
		Where("professional_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", professionalID, models.AppointmentPending, startTime, endTime).
		Count(&stats.Appointments.Pending)

	db.Model(&models.Appointment{}).
		Where("professional_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", professionalID, models.AppointmentConfirmed, startTime, endTime).
		Count(&stats.Appointments.Confirmed)

	db.Model(&models.Appointment{}).
		Where("professional_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", professionalID, models.AppointmentCompleted, startTime, endTime).
		Count(&stats.Appointments.Completed)

	db.Model(&models.Appointment{}).
		Where("professional_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", professionalID, models.AppointmentCancelled, startTime, endTime).
		Count(&stats.Appointments.Cancelled)

	db.Model(&models.Appointment{}).
		Where("professional_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", professionalID, models.AppointmentRefunded, startTime, endTime).
		Count(&stats.Appointments.Refunded)

	var revenueResult struct {
		TotalRevenue  float64
		PaidCount     int64
		RefundedAmount float64
	}

	db.Model(&models.Payment{}).
		Where("professional_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", professionalID, models.PaymentPaid, startTime, endTime).
		Select("COALESCE(SUM(amount), 0) as total_revenue, COUNT(*) as paid_count").
		Scan(&revenueResult)

	stats.Revenue.TotalRevenue = revenueResult.TotalRevenue
	stats.Revenue.PaidCount = revenueResult.PaidCount

	db.Model(&models.Payment{}).
		Where("professional_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", professionalID, models.PaymentRefunded, startTime, endTime).
		Select("COALESCE(SUM(amount), 0) as refunded_amount").
		Scan(&revenueResult)

	stats.Revenue.RefundedAmount = revenueResult.RefundedAmount

	avgRating, reviewCount, err := s.reviewRepo.GetAverageRatingByProfessionalID(professionalID)
	if err != nil {
		return nil, err
	}

	stats.Reviews.AverageRating = avgRating
	stats.Reviews.TotalReviews = reviewCount

	return &stats, nil
}

func (s *StatisticsService) GetAdminStats(startDate, endDate string) (*AdminStats, error) {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		startTime = time.Now().AddDate(0, -1, 0)
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		endTime = time.Now()
	}

	var stats AdminStats

	db := utils.GetDB()

	db.Model(&models.User{}).Count(&stats.TotalUsers)
	db.Model(&models.User{}).Where("role = ?", models.RoleClient).Count(&stats.TotalClients)
	db.Model(&models.User{}).Where("role = ?", models.RoleProfessional).Count(&stats.TotalProfessionals)
	db.Model(&models.Service{}).Count(&stats.TotalServices)

	db.Model(&models.Appointment{}).
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Count(&stats.Appointments.Total)

	db.Model(&models.Appointment{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?", models.AppointmentPending, startTime, endTime).
		Count(&stats.Appointments.Pending)

	db.Model(&models.Appointment{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?", models.AppointmentConfirmed, startTime, endTime).
		Count(&stats.Appointments.Confirmed)

	db.Model(&models.Appointment{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?", models.AppointmentCompleted, startTime, endTime).
		Count(&stats.Appointments.Completed)

	db.Model(&models.Appointment{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?", models.AppointmentCancelled, startTime, endTime).
		Count(&stats.Appointments.Cancelled)

	db.Model(&models.Appointment{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?", models.AppointmentRefunded, startTime, endTime).
		Count(&stats.Appointments.Refunded)

	var revenueResult struct {
		TotalRevenue  float64
		PaidCount     int64
		RefundedAmount float64
	}

	db.Model(&models.Payment{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?", models.PaymentPaid, startTime, endTime).
		Select("COALESCE(SUM(amount), 0) as total_revenue, COUNT(*) as paid_count").
		Scan(&revenueResult)

	stats.Revenue.TotalRevenue = revenueResult.TotalRevenue
	stats.Revenue.PaidCount = revenueResult.PaidCount

	db.Model(&models.Payment{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?", models.PaymentRefunded, startTime, endTime).
		Select("COALESCE(SUM(amount), 0) as refunded_amount").
		Scan(&revenueResult)

	stats.Revenue.RefundedAmount = revenueResult.RefundedAmount

	return &stats, nil
}

func (s *StatisticsService) ExportAppointments(userID uuid.UUID, userRole string, startDate, endDate string, status string) ([]byte, error) {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	var appointments []models.Appointment
	db := utils.GetDB()

	query := db.Preload("Client").Preload("Professional").Preload("Service").Preload("Schedule")

	if userRole == "professional" {
		query = query.Where("professional_id = ?", userID)
	} else if userRole == "client" {
		query = query.Where("client_id = ?", userID)
	}

	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}
	if status != "" {
		query = query.Where("status = ?", models.AppointmentStatus(status))
	}

	if err := query.Order("created_at DESC").Find(&appointments).Error; err != nil {
		return nil, err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Appointments")
	if err != nil {
		return nil, err
	}

	headers := []string{"ID", "Client Name", "Professional Name", "Service", "Date", "Time", "Status", "Amount", "Notes", "Created At"}
	headerRow := sheet.AddRow()
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	for _, apt := range appointments {
		row := sheet.AddRow()
		row.AddCell().Value = apt.ID.String()
		row.AddCell().Value = apt.Client.FullName
		row.AddCell().Value = apt.Professional.FullName
		row.AddCell().Value = apt.Service.Title
		row.AddCell().Value = apt.Schedule.Date.Format("2006-01-02")
		row.AddCell().Value = fmt.Sprintf("%s-%s", apt.Schedule.StartTime, apt.Schedule.EndTime)
		row.AddCell().Value = string(apt.Status)
		if apt.Payment != nil {
			row.AddCell().Value = fmt.Sprintf("%.2f", apt.Payment.Amount)
		} else {
			row.AddCell().Value = "0.00"
		}
		row.AddCell().Value = apt.Notes
		row.AddCell().Value = apt.CreatedAt.Format("2006-01-02 15:04:05")
	}

	var buf []byte
	w := &byteBuffer{buf: &buf}
	if err := file.Write(w); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *StatisticsService) ExportRevenue(userID uuid.UUID, userRole string, startDate, endDate string) ([]byte, error) {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	var payments []models.Payment
	db := utils.GetDB()

	query := db.Preload("Appointment").Preload("Appointment.Client").Preload("Appointment.Professional")

	if userRole == "professional" {
		query = query.Where("professional_id = ?", userID)
	} else if userRole == "client" {
		query = query.Where("client_id = ?", userID)
	}

	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	if err := query.Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Revenue")
	if err != nil {
		return nil, err
	}

	headers := []string{"ID", "Client Name", "Professional Name", "Amount", "Method", "Status", "Transaction ID", "Paid At", "Refunded At", "Created At"}
	headerRow := sheet.AddRow()
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	for _, payment := range payments {
		row := sheet.AddRow()
		row.AddCell().Value = payment.ID.String()
		row.AddCell().Value = payment.Appointment.Client.FullName
		row.AddCell().Value = payment.Appointment.Professional.FullName
		row.AddCell().Value = fmt.Sprintf("%.2f", payment.Amount)
		row.AddCell().Value = string(payment.PaymentMethod)
		row.AddCell().Value = string(payment.Status)
		row.AddCell().Value = payment.TransactionID
		if payment.PaidAt != nil {
			row.AddCell().Value = payment.PaidAt.Format("2006-01-02 15:04:05")
		} else {
			row.AddCell().Value = ""
		}
		if payment.RefundedAt != nil {
			row.AddCell().Value = payment.RefundedAt.Format("2006-01-02 15:04:05")
		} else {
			row.AddCell().Value = ""
		}
		row.AddCell().Value = payment.CreatedAt.Format("2006-01-02 15:04:05")
	}

	var buf []byte
	w := &byteBuffer{buf: &buf}
	if err := file.Write(w); err != nil {
		return nil, err
	}

	return buf, nil
}

type byteBuffer struct {
	buf *[]byte
}

func (b *byteBuffer) Write(p []byte) (n int, err error) {
	*b.buf = append(*b.buf, p...)
	return len(p), nil
}

var _ = gorm.ErrRecordNotFound
