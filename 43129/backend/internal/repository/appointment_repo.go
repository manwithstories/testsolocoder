package repository

import (
	"beauty-salon-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(appointment *model.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *AppointmentRepository) GetByID(id uint) (*model.Appointment, error) {
	var appt model.Appointment
	err := r.db.Preload("Customer.User").Preload("Technician.User").
		Preload("Service").First(&appt, id).Error
	if err != nil {
		return nil, err
	}
	return &appt, nil
}

func (r *AppointmentRepository) Update(appointment *model.Appointment) error {
	return r.db.Save(appointment).Error
}

func (r *AppointmentRepository) List(page, pageSize int, filters map[string]interface{}) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	query := r.db.Model(&model.Appointment{}).
		Preload("Customer.User").
		Preload("Technician.User").
		Preload("Service")

	for key, value := range filters {
		switch key {
		case "customer_id":
			query = query.Where("customer_id = ?", value)
		case "technician_id":
			query = query.Where("technician_id = ?", value)
		case "service_id":
			query = query.Where("service_id = ?", value)
		case "status":
			query = query.Where("status = ?", value)
		case "start_date":
			query = query.Where("appointment_date >= ?", value)
		case "end_date":
			query = query.Where("appointment_date <= ?", value)
		}
	}

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("appointment_date DESC, start_time ASC").Find(&appointments).Error
	return appointments, total, err
}

func (r *AppointmentRepository) GetByTechnicianAndDate(technicianID uint, date time.Time) ([]model.Appointment, error) {
	var appointments []model.Appointment
	err := r.db.Where("technician_id = ? AND DATE(appointment_date) = ? AND status != ?",
		technicianID, date.Format("2006-01-02"), "cancelled").
		Preload("Customer.User").Preload("Service").
		Order("start_time ASC").Find(&appointments).Error
	return appointments, err
}

func (r *AppointmentRepository) GetByCustomer(customerID uint, page, pageSize int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	query := r.db.Model(&model.Appointment{}).Where("customer_id = ?", customerID).
		Preload("Technician.User").Preload("Service")
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("appointment_date DESC, start_time ASC").Find(&appointments).Error
	return appointments, total, err
}

func (r *AppointmentRepository) CheckTimeConflict(technicianID uint, date time.Time, startTime, endTime string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Appointment{}).
		Where("technician_id = ? AND DATE(appointment_date) = ? AND status != ?",
			technicianID, date.Format("2006-01-02"), "cancelled").
		Where("(start_time < ? AND end_time > ?)", endTime, startTime)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *AppointmentRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Appointment{}).Where("id = ?", id).
		Update("status", status).Error
}
