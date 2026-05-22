package repository

import (
	"beauty-salon-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *model.Payment, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(payment).Error
	}
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) GetByID(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Preload("Appointment").First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetByAppointmentID(appointmentID uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("appointment_id = ?", appointmentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) List(page, pageSize int, filters map[string]interface{}) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	query := r.db.Model(&model.Payment{}).Preload("Appointment.Customer.User")
	for key, value := range filters {
		switch key {
		case "pay_method":
			query = query.Where("pay_method = ?", value)
		case "start_date":
			query = query.Where("DATE(created_at) >= ?", value)
		case "end_date":
			query = query.Where("DATE(created_at) <= ?", value)
		}
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&payments).Error
	return payments, total, err
}

func (r *PaymentRepository) GetRevenueByDateRange(startDate, endDate time.Time) (float64, error) {
	var totalAmount struct {
		Total float64
	}
	err := r.db.Model(&model.Payment{}).
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ? AND status = ?",
			startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), "success").
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&totalAmount).Error
	return totalAmount.Total, err
}

func (r *PaymentRepository) GetRevenueByTechnician(technicianID uint, startDate, endDate time.Time) (float64, error) {
	var totalAmount struct {
		Total float64
	}
	err := r.db.Model(&model.Payment{}).
		Joins("JOIN appointments ON appointments.id = payments.appointment_id").
		Where("appointments.technician_id = ? AND DATE(payments.created_at) >= ? AND DATE(payments.created_at) <= ? AND payments.status = ?",
			technicianID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), "success").
		Select("COALESCE(SUM(payments.amount), 0) as total").
		Scan(&totalAmount).Error
	return totalAmount.Total, err
}

func (r *PaymentRepository) GetRevenueByService(serviceID uint, startDate, endDate time.Time) (float64, error) {
	var totalAmount struct {
		Total float64
	}
	err := r.db.Model(&model.Payment{}).
		Joins("JOIN appointments ON appointments.id = payments.appointment_id").
		Where("appointments.service_id = ? AND DATE(payments.created_at) >= ? AND DATE(payments.created_at) <= ? AND payments.status = ?",
			serviceID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), "success").
		Select("COALESCE(SUM(payments.amount), 0) as total").
		Scan(&totalAmount).Error
	return totalAmount.Total, err
}

func (r *PaymentRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
