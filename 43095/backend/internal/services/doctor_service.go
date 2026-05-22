package services

import (
	"errors"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"
	"time"

	"gorm.io/gorm"
)

type DoctorService struct {
	db *gorm.DB
}

func NewDoctorService() *DoctorService {
	return &DoctorService{
		db: database.GetDB(),
	}
}

type DoctorListQuery struct {
	DepartmentID uint   `form:"department_id"`
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"page_size,default=10"`
	SortByRating string `form:"sort_by_rating"`
}

type CreateDoctorRequest struct {
	UserID          uint              `json:"user_id" binding:"required"`
	DepartmentID    uint              `json:"department_id" binding:"required"`
	Title           models.DoctorTitle `json:"title" binding:"required"`
	Specialty       string            `json:"specialty"`
	Introduction    string            `json:"introduction"`
	ConsultationFee float64           `json:"consultation_fee"`
	RegistrationFee float64           `json:"registration_fee"`
}

type UpdateDoctorRequest struct {
	DepartmentID    uint              `json:"department_id"`
	Title           models.DoctorTitle `json:"title"`
	Specialty       string            `json:"specialty"`
	Introduction    string            `json:"introduction"`
	ConsultationFee *float64          `json:"consultation_fee"`
	RegistrationFee *float64          `json:"registration_fee"`
}

type ScheduleRequest struct {
	DayOfWeek   models.DayOfWeek `json:"day_of_week" binding:"required,min=0,max=6"`
	StartTime   string           `json:"start_time" binding:"required"`
	EndTime     string           `json:"end_time" binding:"required"`
	MaxPatients int              `json:"max_patients" binding:"min=1"`
	TimeSlot    int              `json:"time_slot_minutes" binding:"min=5"`
	IsAvailable *bool            `json:"is_available"`
}

type AvailableTimeSlot struct {
	Date     string   `json:"date"`
	TimeSlots []string `json:"time_slots"`
}

func (s *DoctorService) GetDoctorList(query DoctorListQuery) ([]models.Doctor, int64, error) {
	var doctors []models.Doctor
	var total int64

	db := s.db.Model(&models.Doctor{}).Preload("User").Preload("Department")

	if query.DepartmentID > 0 {
		db = db.Where("department_id = ?", query.DepartmentID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.SortByRating == "desc" {
		db = db.Order("average_rating DESC")
	} else if query.SortByRating == "asc" {
		db = db.Order("average_rating ASC")
	} else {
		db = db.Order("id DESC")
	}

	if err := db.Scopes(database.Paginate(query.Page, query.PageSize)).Find(&doctors).Error; err != nil {
		return nil, 0, err
	}

	return doctors, total, nil
}

func (s *DoctorService) GetDoctorByID(id uint) (*models.Doctor, error) {
	var doctor models.Doctor
	err := s.db.Preload("User").
		Preload("Department").
		Preload("Schedules").
		Preload("Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Patient.User").Order("created_at DESC")
		}).
		First(&doctor, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("医生不存在")
		}
		return nil, err
	}
	return &doctor, nil
}

func (s *DoctorService) GetDoctorByUserID(userID uint) (*models.Doctor, error) {
	var doctor models.Doctor
	err := s.db.Where("user_id = ?", userID).First(&doctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("医生不存在")
		}
		return nil, err
	}
	return &doctor, nil
}

func (s *DoctorService) CreateDoctor(req CreateDoctorRequest) (*models.Doctor, error) {
	var count int64
	s.db.Model(&models.Doctor{}).Where("user_id = ?", req.UserID).Count(&count)
	if count > 0 {
		return nil, errors.New("该用户已是医生")
	}

	var user models.User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	var department models.Department
	if err := s.db.First(&department, req.DepartmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("科室不存在")
		}
		return nil, err
	}

	doctor := &models.Doctor{
		UserID:          req.UserID,
		DepartmentID:    req.DepartmentID,
		Title:           req.Title,
		Specialty:       req.Specialty,
		Introduction:    req.Introduction,
		ConsultationFee: req.ConsultationFee,
		RegistrationFee: req.RegistrationFee,
		AverageRating:   0,
		ReviewCount:     0,
	}

	if err := s.db.Create(doctor).Error; err != nil {
		return nil, err
	}

	return s.GetDoctorByID(doctor.ID)
}

func (s *DoctorService) UpdateDoctor(id uint, req UpdateDoctorRequest) (*models.Doctor, error) {
	doctor, err := s.GetDoctorByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.DepartmentID > 0 {
		updates["department_id"] = req.DepartmentID
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Specialty != "" {
		updates["specialty"] = req.Specialty
	}
	if req.Introduction != "" {
		updates["introduction"] = req.Introduction
	}
	if req.ConsultationFee != nil {
		updates["consultation_fee"] = *req.ConsultationFee
	}
	if req.RegistrationFee != nil {
		updates["registration_fee"] = *req.RegistrationFee
	}

	if len(updates) > 0 {
		if err := s.db.Model(doctor).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetDoctorByID(id)
}

func (s *DoctorService) DeleteDoctor(id uint) error {
	doctor, err := s.GetDoctorByID(id)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("doctor_id = ?", id).Delete(&models.Schedule{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(doctor).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *DoctorService) GetSchedules(doctorID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := s.db.Where("doctor_id = ?", doctorID).Order("day_of_week ASC, start_time ASC").Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (s *DoctorService) CreateSchedule(doctorID uint, req ScheduleRequest) (*models.Schedule, error) {
	if _, err := s.GetDoctorByID(doctorID); err != nil {
		return nil, err
	}

	var count int64
	s.db.Model(&models.Schedule{}).
		Where("doctor_id = ? AND day_of_week = ? AND start_time = ?", doctorID, req.DayOfWeek, req.StartTime).
		Count(&count)
	if count > 0 {
		return nil, errors.New("该时间段已有排班")
	}

	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	maxPatients := 20
	if req.MaxPatients > 0 {
		maxPatients = req.MaxPatients
	}

	timeSlot := 15
	if req.TimeSlot > 0 {
		timeSlot = req.TimeSlot
	}

	schedule := &models.Schedule{
		DoctorID:    doctorID,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		MaxPatients: maxPatients,
		TimeSlot:    timeSlot,
		IsAvailable: isAvailable,
	}

	if err := s.db.Create(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *DoctorService) UpdateSchedule(scheduleID uint, req ScheduleRequest) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := s.db.First(&schedule, scheduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("排班不存在")
		}
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.DayOfWeek >= 0 && req.DayOfWeek <= 6 {
		updates["day_of_week"] = req.DayOfWeek
	}
	if req.StartTime != "" {
		updates["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		updates["end_time"] = req.EndTime
	}
	if req.MaxPatients > 0 {
		updates["max_patients"] = req.MaxPatients
	}
	if req.TimeSlot > 0 {
		updates["time_slot"] = req.TimeSlot
	}
	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}

	if len(updates) > 0 {
		if err := s.db.Model(&schedule).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if err := s.db.First(&schedule, scheduleID).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (s *DoctorService) DeleteSchedule(scheduleID uint) error {
	var schedule models.Schedule
	if err := s.db.First(&schedule, scheduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("排班不存在")
		}
		return err
	}
	return s.db.Delete(&schedule).Error
}

func (s *DoctorService) GetAvailableTimeSlots(doctorID uint, date string) (*AvailableTimeSlot, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("日期格式错误，应为 YYYY-MM-DD")
	}

	dayOfWeek := models.DayOfWeek(parsedDate.Weekday())

	var schedules []models.Schedule
	err = s.db.Where("doctor_id = ? AND day_of_week = ? AND is_available = ?", doctorID, dayOfWeek, true).
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	var bookedSlots []string
	err = s.db.Model(&models.Appointment{}).
		Where("doctor_id = ? AND DATE(appointment_date) = ? AND status IN (?, ?)",
			doctorID, parsedDate.Format("2006-01-02"),
			models.AppointmentPending, models.AppointmentConfirmed).
		Pluck("start_time", &bookedSlots).Error
	if err != nil {
		return nil, err
	}

	bookedMap := make(map[string]bool)
	for _, slot := range bookedSlots {
		bookedMap[slot] = true
	}

	var allTimeSlots []string
	for _, schedule := range schedules {
		startTime, err := time.Parse("15:04", schedule.StartTime)
		if err != nil {
			continue
		}
		endTime, err := time.Parse("15:04", schedule.EndTime)
		if err != nil {
			continue
		}

		for current := startTime; current.Before(endTime); current = current.Add(time.Duration(schedule.TimeSlot) * time.Minute) {
			slotStr := current.Format("15:04")
			if !bookedMap[slotStr] {
				allTimeSlots = append(allTimeSlots, slotStr)
			}
		}
	}

	return &AvailableTimeSlot{
		Date:      date,
		TimeSlots: allTimeSlots,
	}, nil
}
