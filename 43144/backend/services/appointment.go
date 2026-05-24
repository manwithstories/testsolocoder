package services

import (
	"errors"
	"fmt"
	"time"

	"pet-adoption-platform/database"
	"pet-adoption-platform/models"
)

func CreateAppointment(req *models.CreateAppointmentRequest, userID uint) (*models.Appointment, error) {
	pet, err := GetPetByID(req.PetID)
	if err != nil {
		return nil, errors.New("pet not found")
	}

	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return nil, errors.New("invalid appointment date")
	}

	var existing models.Appointment
	if err := database.DB.Where(
		"appointment_date = ? AND start_time = ? AND status NOT IN ?",
		appointmentDate, req.StartTime,
		[]models.AppointmentStatus{models.AppointmentStatusCancelled, models.AppointmentStatusCompleted},
	).First(&existing).Error; err == nil {
		return nil, errors.New("this time slot is already booked")
	}

	appointment := &models.Appointment{
		UserID:          userID,
		PetID:           req.PetID,
		RescueID:        pet.RescueID,
		AppointmentType: models.AppointmentType(req.AppointmentType),
		AppointmentDate: appointmentDate,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Status:          models.AppointmentStatusPending,
		Location:        req.Location,
		Notes:           req.Notes,
	}

	if err := database.DB.Create(appointment).Error; err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

func GetAppointmentByID(id uint) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := database.DB.Preload("Pet").Preload("User").First(&appointment, id).Error; err != nil {
		return nil, err
	}
	return &appointment, nil
}

func ListAppointments(query *models.AppointmentListQuery, userID uint, role string, rescueID *uint) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	db := database.DB.Model(&models.Appointment{})

	switch role {
	case "adopter":
		db = db.Where("user_id = ?", userID)
	case "rescue":
		if rescueID != nil {
			db = db.Where("rescue_id = ?", *rescueID)
		}
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.PetID > 0 {
		db = db.Where("pet_id = ?", query.PetID)
	}
	if query.RescueID > 0 {
		db = db.Where("rescue_id = ?", query.RescueID)
	}
	if query.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", query.DateFrom)
		if err == nil {
			db = db.Where("appointment_date >= ?", dateFrom)
		}
	}
	if query.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", query.DateTo)
		if err == nil {
			db = db.Where("appointment_date <= ?", dateTo)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Pet").Preload("User").
		Order("appointment_date ASC, start_time ASC").
		Offset(offset).Limit(query.PageSize).
		Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func UpdateAppointment(id uint, req *models.UpdateAppointmentRequest) (*models.Appointment, error) {
	updates := make(map[string]interface{})

	if req.AppointmentDate != "" {
		date, err := time.Parse("2006-01-02", req.AppointmentDate)
		if err == nil {
			updates["appointment_date"] = date
		}
	}
	if req.StartTime != "" {
		updates["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		updates["end_time"] = req.EndTime
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}

	result := database.DB.Model(&models.Appointment{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	return GetAppointmentByID(id)
}

func CancelAppointment(id uint, reason string) (*models.Appointment, error) {
	appointment, err := GetAppointmentByID(id)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.Status == models.AppointmentStatusCancelled ||
		appointment.Status == models.AppointmentStatusCompleted {
		return nil, errors.New("cannot cancel this appointment")
	}

	appointment.Status = models.AppointmentStatusCancelled
	appointment.CancelReason = reason

	if err := database.DB.Save(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}

func RescheduleAppointment(id uint, req *models.UpdateAppointmentRequest) (*models.Appointment, error) {
	appointment, err := GetAppointmentByID(id)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.Status == models.AppointmentStatusCancelled ||
		appointment.Status == models.AppointmentStatusCompleted {
		return nil, errors.New("cannot reschedule this appointment")
	}

	newAppointment := &models.Appointment{
		UserID:          appointment.UserID,
		PetID:           appointment.PetID,
		RescueID:        appointment.RescueID,
		AppointmentType: appointment.AppointmentType,
		AppointmentDate: appointment.AppointmentDate,
		StartTime:       appointment.StartTime,
		EndTime:         appointment.EndTime,
		Location:        appointment.Location,
		Notes:           appointment.Notes,
		Status:          models.AppointmentStatusRescheduled,
		OriginalID:      &id,
	}

	if req.AppointmentDate != "" {
		date, err := time.Parse("2006-01-02", req.AppointmentDate)
		if err == nil {
			newAppointment.AppointmentDate = date
		}
	}
	if req.StartTime != "" {
		newAppointment.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		newAppointment.EndTime = req.EndTime
	}
	if req.Location != "" {
		newAppointment.Location = req.Location
	}
	if req.Notes != "" {
		newAppointment.Notes = req.Notes
	}

	if err := database.DB.Create(newAppointment).Error; err != nil {
		return nil, err
	}

	appointment.Status = models.AppointmentStatusRescheduled
	database.DB.Save(appointment)

	return newAppointment, nil
}

func ConfirmAppointment(id uint) (*models.Appointment, error) {
	appointment, err := GetAppointmentByID(id)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.Status != models.AppointmentStatusPending {
		return nil, errors.New("can only confirm pending appointments")
	}

	appointment.Status = models.AppointmentStatusConfirmed
	if err := database.DB.Save(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}

func CompleteAppointment(id uint) (*models.Appointment, error) {
	appointment, err := GetAppointmentByID(id)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.Status == models.AppointmentStatusCancelled {
		return nil, errors.New("cannot complete cancelled appointment")
	}

	appointment.Status = models.AppointmentStatusCompleted
	if err := database.DB.Save(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}
