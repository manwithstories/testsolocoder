package repositories

import (
	"time"

	"consultation-platform/models"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository() *AppointmentRepository {
	return &AppointmentRepository{db: utils.GetDB()}
}

func (r *AppointmentRepository) Create(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *AppointmentRepository) FindByID(id uuid.UUID) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.Preload("Client").
		Preload("Professional").
		Preload("Service").
		Preload("Schedule").
		Preload("Payment").
		First(&appointment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) Update(appointment *models.Appointment) error {
	return r.db.Save(appointment).Error
}

func (r *AppointmentRepository) UpdateStatus(id uuid.UUID, status models.AppointmentStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == models.AppointmentCancelled {
		now := time.Now()
		updates["cancelled_at"] = &now
	}
	if status == models.AppointmentCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}
	return r.db.Model(&models.Appointment{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AppointmentRepository) FindByClientID(clientID uuid.UUID, page, pageSize int, status string) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Appointment{}).Where("client_id = ?", clientID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").
		Preload("Professional").
		Preload("Service").
		Preload("Schedule").
		Preload("Payment").
		Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (r *AppointmentRepository) FindByProfessionalID(professionalID uuid.UUID, page, pageSize int, status string) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Appointment{}).Where("professional_id = ?", professionalID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").
		Preload("Client").
		Preload("Service").
		Preload("Schedule").
		Preload("Payment").
		Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (r *AppointmentRepository) FindByScheduleID(scheduleID uuid.UUID) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.Where("schedule_id = ?", scheduleID).First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) FindExpiredPendingAppointments(expireDuration time.Duration) ([]models.Appointment, error) {
	var appointments []models.Appointment
	cutoffTime := time.Now().Add(-expireDuration)
	err := r.db.Where("status = ? AND created_at < ?", models.AppointmentPending, cutoffTime).Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

func (r *AppointmentRepository) WithTx(tx *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: tx}
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: utils.GetDB()}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) FindByAppointmentID(appointmentID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("appointment_id = ?", appointmentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) UpdateStatus(id uuid.UUID, status models.PaymentStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == models.PaymentPaid {
		now := time.Now()
		updates["paid_at"] = &now
	}
	if status == models.PaymentRefunded {
		now := time.Now()
		updates["refunded_at"] = &now
	}
	return r.db.Model(&models.Payment{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PaymentRepository) FindExpiredPayments() ([]models.Payment, error) {
	var payments []models.Payment
	now := time.Now()
	err := r.db.Where("status = ? AND expires_at < ?", models.PaymentPending, now).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindByUserID(userID uuid.UUID, page, pageSize int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Payment{}).Where("client_id = ? OR professional_id = ?", userID, userID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *PaymentRepository) WithTx(tx *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: tx}
}
