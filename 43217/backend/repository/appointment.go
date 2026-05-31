package repository

import (
	"health-platform/models"
	"time"

	"gorm.io/gorm"
)

type AppointmentRepository struct {
	*BaseRepository
}

func NewAppointmentRepository() *AppointmentRepository {
	return &AppointmentRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *AppointmentRepository) FindByEmployeeID(employeeID uint, page, pageSize int) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := r.DB.Model(&models.Appointment{}).Where("employee_id = ?", employeeID)
	query.Count(&total)

	err := query.Preload("Package").Preload("Agency").Preload("TimeSlot").
		Order("appointment_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&appointments).Error
	return appointments, total, err
}

func (r *AppointmentRepository) FindByCompanyID(companyID uint, page, pageSize int) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := r.DB.Model(&models.Appointment{}).Where("company_id = ?", companyID)
	query.Count(&total)

	err := query.Preload("Employee").Preload("Package").Preload("Agency").
		Order("appointment_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&appointments).Error
	return appointments, total, err
}

func (r *AppointmentRepository) FindByCompanyIDAndDateRange(companyID uint, startDate, endDate time.Time, page, pageSize int) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := r.DB.Model(&models.Appointment{}).
		Where("company_id = ? AND appointment_date >= ? AND appointment_date <= ?", companyID, startDate, endDate)
	query.Count(&total)

	err := query.Preload("Employee").Preload("Package").Preload("Agency").
		Order("appointment_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&appointments).Error
	return appointments, total, err
}

func (r *AppointmentRepository) FindByCompanyIDAndAgency(companyID, agencyID uint, page, pageSize int) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := r.DB.Model(&models.Appointment{}).
		Where("company_id = ? AND agency_id = ?", companyID, agencyID)
	query.Count(&total)

	err := query.Preload("Employee").Preload("Package").Preload("Agency").
		Order("appointment_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&appointments).Error
	return appointments, total, err
}

func (r *AppointmentRepository) FindByAgencyID(agencyID uint, page, pageSize int) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := r.DB.Model(&models.Appointment{}).Where("agency_id = ?", agencyID)
	query.Count(&total)

	err := query.Preload("Employee").Preload("Package").Preload("Company").
		Order("appointment_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&appointments).Error
	return appointments, total, err
}

func (r *AppointmentRepository) GetWithReport(appointmentID uint) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.DB.Preload("Report").Preload("Report.Items").
		Preload("Employee").Preload("Package").
		First(&appointment, appointmentID).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) UpdateStatus(appointmentID uint, status models.AppointmentStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == models.AppointmentStatusCancelled {
		updates["cancelled_at"] = time.Now()
	}
	if status == models.AppointmentStatusCompleted {
		updates["completed_at"] = time.Now()
	}
	return r.DB.Model(&models.Appointment{}).Where("id = ?", appointmentID).Updates(updates).Error
}

func (r *AppointmentRepository) CancelAppointment(appointmentID uint, reason string) error {
	return r.DB.Model(&models.Appointment{}).Where("id = ?", appointmentID).
		Updates(map[string]interface{}{
			"status":        models.AppointmentStatusCancelled,
			"cancelled_at":  time.Now(),
			"cancel_reason": reason,
		}).Error
}

func (r *AppointmentRepository) RescheduleAppointment(appointmentID uint, timeSlotID uint, appointmentDate time.Time, startTime, endTime string) error {
	return r.DB.Model(&models.Appointment{}).Where("id = ?", appointmentID).
		Updates(map[string]interface{}{
			"time_slot_id":     timeSlotID,
			"appointment_date": appointmentDate,
			"start_time":       startTime,
			"end_time":         endTime,
		}).Error
}

func (r *AppointmentRepository) FindExpiredAppointments() ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Where("status = ? AND appointment_date < ?", 
		models.AppointmentStatusPending, time.Now()).
		Find(&appointments).Error
	return appointments, err
}

func (r *AppointmentRepository) FindNeedRemindAppointments(days int) ([]models.Appointment, error) {
	remindDate := time.Now().AddDate(0, 0, days)
	var appointments []models.Appointment
	err := r.DB.Where("status = ? AND is_reminded = ? AND appointment_date <= ?",
		models.AppointmentStatusPending, false, remindDate).
		Preload("Employee").Preload("Employee.User").
		Find(&appointments).Error
	return appointments, err
}

func (r *AppointmentRepository) MarkAsReminded(appointmentID uint) error {
	return r.DB.Model(&models.Appointment{}).Where("id = ?", appointmentID).
		Updates(map[string]interface{}{
			"is_reminded": true,
			"reminded_at": time.Now(),
		}).Error
}

func (r *AppointmentRepository) GetAppointmentByNo(appointmentNo string) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.DB.Where("appointment_no = ?", appointmentNo).First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) GetEmployeeAppointmentCount(employeeID uint, year int) (int64, error) {
	var count int64
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)
	
	err := r.DB.Model(&models.Appointment{}).
		Where("employee_id = ? AND status IN ? AND appointment_date >= ? AND appointment_date < ?",
			employeeID,
			[]models.AppointmentStatus{
				models.AppointmentStatusPending,
				models.AppointmentStatusConfirmed,
				models.AppointmentStatusCompleted,
			},
			startDate,
			endDate,
		).Count(&count).Error
	return count, err
}

func (r *AppointmentRepository) GetDepartmentAppointmentCount(departmentID uint, year int) (int64, error) {
	var count int64
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)
	
	err := r.DB.Model(&models.Appointment{}).
		Joins("JOIN employees ON employees.id = appointments.employee_id").
		Where("employees.department_id = ? AND appointments.status IN ? AND appointments.appointment_date >= ? AND appointments.appointment_date < ?",
			departmentID,
			[]models.AppointmentStatus{
				models.AppointmentStatusPending,
				models.AppointmentStatusConfirmed,
				models.AppointmentStatusCompleted,
			},
			startDate,
			endDate,
		).Count(&count).Error
	return count, err
}

func (r *AppointmentRepository) GetAgencyRanking() ([]map[string]interface{}, error) {
	type AgencyRanking struct {
		AgencyID   uint    `json:"agency_id"`
		AgencyName string  `json:"agency_name"`
		Count      int64   `json:"count"`
		Rating     float64 `json:"rating"`
	}

	var results []AgencyRanking
	err := r.DB.Table("appointments").
		Select("appointments.agency_id, agencies.name as agency_name, COUNT(*) as count, agencies.rating").
		Joins("JOIN agencies ON agencies.id = appointments.agency_id").
		Where("appointments.status = ?", models.AppointmentStatusCompleted).
		Group("appointments.agency_id, agencies.name, agencies.rating").
		Order("count DESC").
		Limit(10).
		Scan(&results).Error

	var rankings []map[string]interface{}
	for _, r := range results {
		rankings = append(rankings, map[string]interface{}{
			"agency_id":   r.AgencyID,
			"agency_name": r.AgencyName,
			"count":       r.Count,
			"rating":      r.Rating,
		})
	}
	return rankings, err
}
