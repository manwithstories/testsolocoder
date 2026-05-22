package services

import (
	"errors"
	"medical-platform/internal/config"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppointmentService struct {
	db *gorm.DB
}

func NewAppointmentService() *AppointmentService {
	return &AppointmentService{
		db: database.GetDB(),
	}
}

type CreateAppointmentRequest struct {
	DoctorID        uint      `json:"doctor_id" binding:"required"`
	AppointmentDate time.Time `json:"appointment_date" binding:"required"`
	StartTime       string    `json:"start_time" binding:"required"`
	EndTime         string    `json:"end_time" binding:"required"`
	Symptoms        string    `json:"symptoms"`
	Notes           string    `json:"notes"`
}

type UpdateAppointmentRequest struct {
	AppointmentDate time.Time `json:"appointment_date"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	Symptoms        string    `json:"symptoms"`
	Notes           string    `json:"notes"`
}

type CancelAppointmentRequest struct {
	CancelReason string `json:"cancel_reason" binding:"required"`
}

type TimeSlotAvailability struct {
	IsAvailable bool `json:"is_available"`
	Reason      string `json:"reason,omitempty"`
}

func (s *AppointmentService) CreateAppointment(patientUserID uint, req *CreateAppointmentRequest) (*models.Appointment, error) {
	var appointment *models.Appointment

	err := database.WithTransaction(func(tx *gorm.DB) error {
		var patient models.Patient
		if err := tx.Where("user_id = ?", patientUserID).First(&patient).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				patient = models.Patient{
					UserID: patientUserID,
				}
				if err := tx.Create(&patient).Error; err != nil {
					return errors.New("创建患者信息失败")
				}

				healthRecord := &models.HealthRecord{
					PatientID: patient.ID,
				}
				if err := tx.Create(healthRecord).Error; err != nil {
					return errors.New("创建健康档案失败")
				}
			} else {
				return errors.New("患者信息不存在")
			}
		}

		var doctor models.Doctor
		if err := tx.Preload("User").Preload("Department").First(&doctor, req.DoctorID).Error; err != nil {
			return errors.New("医生信息不存在")
		}

		date := req.AppointmentDate
		dayOfWeek := models.DayOfWeek(date.Weekday())

		var schedule models.Schedule
		if err := tx.Where("doctor_id = ? AND day_of_week = ? AND is_available = ?",
			req.DoctorID, dayOfWeek, true).First(&schedule).Error; err != nil {
			return errors.New("该医生在该日期没有排班")
		}

		if req.StartTime < schedule.StartTime || req.EndTime > schedule.EndTime {
			return errors.New("预约时间不在医生排班时间范围内")
		}

		var count int64
		if err := tx.Model(&models.Appointment{}).
			Where("doctor_id = ? AND DATE(appointment_date) = DATE(?) AND status IN (?)",
				req.DoctorID, req.AppointmentDate,
				[]models.AppointmentStatus{models.AppointmentPending, models.AppointmentConfirmed}).
			Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(schedule.MaxPatients) {
			return errors.New("该医生当日预约已满")
		}

		var conflictCount int64
		if err := tx.Model(&models.Appointment{}).
			Where(`doctor_id = ? AND DATE(appointment_date) = DATE(?) AND status IN (?) AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?))`,
				req.DoctorID, req.AppointmentDate,
				[]models.AppointmentStatus{models.AppointmentPending, models.AppointmentConfirmed},
				req.EndTime, req.StartTime,
				req.EndTime, req.StartTime,
				req.StartTime, req.EndTime).
			Count(&conflictCount).Error; err != nil {
			return err
		}
		if conflictCount > 0 {
			return errors.New("该时间段已被预约")
		}

		appointment = &models.Appointment{
			PatientID:       patient.ID,
			DoctorID:        req.DoctorID,
			AppointmentDate: req.AppointmentDate,
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			Status:          models.AppointmentPending,
			Symptoms:        req.Symptoms,
			Notes:           req.Notes,
		}

		if err := tx.Create(appointment).Error; err != nil {
			return err
		}

		appointment.Patient = patient
		appointment.Doctor = doctor

		return nil
	})

	if err != nil {
		return nil, err
	}

	go func() {
		notificationService := NewNotificationService()
		if err := notificationService.SendAppointmentConfirmation(appointment); err != nil {
			config.Logger.Error("发送预约确认通知失败", zap.Error(err))
		}
	}()

	return appointment, nil
}

func (s *AppointmentService) GetAppointmentList(userID uint, role models.UserRole, page, pageSize int, status string) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := s.db.Model(&models.Appointment{}).Preload("Patient.User").Preload("Doctor.User")

	switch role {
	case models.RolePatient:
		var patient models.Patient
		if err := s.db.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			return nil, 0, errors.New("患者信息不存在")
		}
		query = query.Where("patient_id = ?", patient.ID)
	case models.RoleDoctor:
		var doctor models.Doctor
		if err := s.db.Where("user_id = ?", userID).First(&doctor).Error; err != nil {
			return nil, 0, errors.New("医生信息不存在")
		}
		query = query.Where("doctor_id = ?", doctor.ID)
	case models.RoleAdmin:
	default:
		return nil, 0, errors.New("无效的用户角色")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).
		Order("appointment_date DESC, start_time ASC").
		Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (s *AppointmentService) GetAppointmentDetail(id uint, userID uint, role models.UserRole) (*models.Appointment, error) {
	var appointment models.Appointment
	query := s.db.Preload("Patient.User").Preload("Doctor.User").Preload("Consultation.Prescription.Items").Preload("Consultation.Reports")

	if err := query.First(&appointment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预约不存在")
		}
		return nil, err
	}

	switch role {
	case models.RolePatient:
		var patient models.Patient
		if err := s.db.Where("user_id = ?", userID).First(&patient).Error; err != nil {
			return nil, errors.New("患者信息不存在")
		}
		if appointment.PatientID != patient.ID {
			return nil, errors.New("无权查看该预约")
		}
	case models.RoleDoctor:
		var doctor models.Doctor
		if err := s.db.Where("user_id = ?", userID).First(&doctor).Error; err != nil {
			return nil, errors.New("医生信息不存在")
		}
		if appointment.DoctorID != doctor.ID {
			return nil, errors.New("无权查看该预约")
		}
	case models.RoleAdmin:
	default:
		return nil, errors.New("无效的用户角色")
	}

	return &appointment, nil
}

func (s *AppointmentService) CancelAppointment(id uint, userID uint, role models.UserRole, reason string) error {
	return database.WithTransaction(func(tx *gorm.DB) error {
		var appointment models.Appointment
		if err := tx.First(&appointment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("预约不存在")
			}
			return err
		}

		if appointment.Status == models.AppointmentCancelled || appointment.Status == models.AppointmentCompleted {
			return errors.New("该预约状态不允许取消")
		}

		switch role {
		case models.RolePatient:
			var patient models.Patient
			if err := tx.Where("user_id = ?", userID).First(&patient).Error; err != nil {
				return errors.New("患者信息不存在")
			}
			if appointment.PatientID != patient.ID {
				return errors.New("无权取消该预约")
			}
		case models.RoleAdmin:
		default:
			return errors.New("无权取消该预约")
		}

		if err := tx.Model(&appointment).Updates(map[string]interface{}{
			"status":        models.AppointmentCancelled,
			"cancel_reason": reason,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *AppointmentService) RescheduleAppointment(id uint, userID uint, role models.UserRole, req *UpdateAppointmentRequest) (*models.Appointment, error) {
	var updatedAppointment *models.Appointment

	err := database.WithTransaction(func(tx *gorm.DB) error {
		var appointment models.Appointment
		if err := tx.First(&appointment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("预约不存在")
			}
			return err
		}

		if appointment.Status != models.AppointmentPending && appointment.Status != models.AppointmentConfirmed {
			return errors.New("该预约状态不允许改签")
		}

		switch role {
		case models.RolePatient:
			var patient models.Patient
			if err := tx.Where("user_id = ?", userID).First(&patient).Error; err != nil {
				return errors.New("患者信息不存在")
			}
			if appointment.PatientID != patient.ID {
				return errors.New("无权改签该预约")
			}
		case models.RoleAdmin:
		default:
			return errors.New("无权改签该预约")
		}

		if req.AppointmentDate.IsZero() {
			req.AppointmentDate = appointment.AppointmentDate
		}
		if req.StartTime == "" {
			req.StartTime = appointment.StartTime
		}
		if req.EndTime == "" {
			req.EndTime = appointment.EndTime
		}

		date := req.AppointmentDate
		dayOfWeek := models.DayOfWeek(date.Weekday())

		var schedule models.Schedule
		if err := tx.Where("doctor_id = ? AND day_of_week = ? AND is_available = ?",
			appointment.DoctorID, dayOfWeek, true).First(&schedule).Error; err != nil {
			return errors.New("该医生在该日期没有排班")
		}

		if req.StartTime < schedule.StartTime || req.EndTime > schedule.EndTime {
			return errors.New("预约时间不在医生排班时间范围内")
		}

		var count int64
		if err := tx.Model(&models.Appointment{}).
			Where("doctor_id = ? AND DATE(appointment_date) = DATE(?) AND status IN (?) AND id != ?",
				appointment.DoctorID, req.AppointmentDate,
				[]models.AppointmentStatus{models.AppointmentPending, models.AppointmentConfirmed},
				id).
			Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(schedule.MaxPatients) {
			return errors.New("该医生当日预约已满")
		}

		var conflictCount int64
		if err := tx.Model(&models.Appointment{}).
			Where(`doctor_id = ? AND DATE(appointment_date) = DATE(?) AND status IN (?) AND id != ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?))`,
				appointment.DoctorID, req.AppointmentDate,
				[]models.AppointmentStatus{models.AppointmentPending, models.AppointmentConfirmed},
				id,
				req.EndTime, req.StartTime,
				req.EndTime, req.StartTime,
				req.StartTime, req.EndTime).
			Count(&conflictCount).Error; err != nil {
			return err
		}
		if conflictCount > 0 {
			return errors.New("该时间段已被预约")
		}

		updates := map[string]interface{}{
			"appointment_date": req.AppointmentDate,
			"start_time":       req.StartTime,
			"end_time":         req.EndTime,
			"status":          models.AppointmentPending,
		}
		if req.Symptoms != "" {
			updates["symptoms"] = req.Symptoms
		}
		if req.Notes != "" {
			updates["notes"] = req.Notes
		}

		if err := tx.Model(&appointment).Updates(updates).Error; err != nil {
			return err
		}

		updatedAppointment = &appointment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedAppointment, nil
}

func (s *AppointmentService) ConfirmAppointment(id uint, doctorUserID uint) error {
	return database.WithTransaction(func(tx *gorm.DB) error {
		var doctor models.Doctor
		if err := tx.Where("user_id = ?", doctorUserID).First(&doctor).Error; err != nil {
			return errors.New("医生信息不存在")
		}

		var appointment models.Appointment
		if err := tx.First(&appointment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("预约不存在")
			}
			return err
		}

		if appointment.DoctorID != doctor.ID {
			return errors.New("无权确认该预约")
		}

		if appointment.Status != models.AppointmentPending {
			return errors.New("该预约状态不允许确认")
		}

		if err := tx.Model(&appointment).Update("status", models.AppointmentConfirmed).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *AppointmentService) CompleteAppointment(id uint, doctorUserID uint) error {
	return database.WithTransaction(func(tx *gorm.DB) error {
		var doctor models.Doctor
		if err := tx.Where("user_id = ?", doctorUserID).First(&doctor).Error; err != nil {
			return errors.New("医生信息不存在")
		}

		var appointment models.Appointment
		if err := tx.First(&appointment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("预约不存在")
			}
			return err
		}

		if appointment.DoctorID != doctor.ID {
			return errors.New("无权完成该预约")
		}

		if appointment.Status != models.AppointmentConfirmed {
			return errors.New("该预约状态不允许完成")
		}

		if err := tx.Model(&appointment).Update("status", models.AppointmentCompleted).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *AppointmentService) CheckTimeSlotAvailability(doctorID uint, appointmentDate time.Time, startTime, endTime string) (*TimeSlotAvailability, error) {
	var doctor models.Doctor
	if err := s.db.First(&doctor, doctorID).Error; err != nil {
		return nil, errors.New("医生信息不存在")
	}

	date := appointmentDate
	dayOfWeek := models.DayOfWeek(date.Weekday())

	var schedule models.Schedule
	if err := s.db.Where("doctor_id = ? AND day_of_week = ? AND is_available = ?",
		doctorID, dayOfWeek, true).First(&schedule).Error; err != nil {
		return &TimeSlotAvailability{
			IsAvailable: false,
			Reason:      "该医生在该日期没有排班",
		}, nil
	}

	if startTime < schedule.StartTime || endTime > schedule.EndTime {
		return &TimeSlotAvailability{
			IsAvailable: false,
			Reason:      "预约时间不在医生排班时间范围内",
		}, nil
	}

	var count int64
	if err := s.db.Model(&models.Appointment{}).
		Where("doctor_id = ? AND DATE(appointment_date) = DATE(?) AND status IN (?)",
			doctorID, appointmentDate,
			[]models.AppointmentStatus{models.AppointmentPending, models.AppointmentConfirmed}).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count >= int64(schedule.MaxPatients) {
		return &TimeSlotAvailability{
			IsAvailable: false,
			Reason:      "该医生当日预约已满",
		}, nil
	}

	var conflictCount int64
	if err := s.db.Model(&models.Appointment{}).
		Where(`doctor_id = ? AND DATE(appointment_date) = DATE(?) AND status IN (?) AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?))`,
			doctorID, appointmentDate,
			[]models.AppointmentStatus{models.AppointmentPending, models.AppointmentConfirmed},
			endTime, startTime,
			endTime, startTime,
			startTime, endTime).
		Count(&conflictCount).Error; err != nil {
		return nil, err
	}
	if conflictCount > 0 {
		return &TimeSlotAvailability{
			IsAvailable: false,
			Reason:      "该时间段已被预约",
		}, nil
	}

	return &TimeSlotAvailability{
		IsAvailable: true,
	}, nil
}
