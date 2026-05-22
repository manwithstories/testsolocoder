package service

import (
	"fmt"
	"time"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
)

type TechnicianService struct {
	techRepo        *repository.TechnicianRepository
	userRepo        *repository.UserRepository
	appointmentRepo *repository.AppointmentRepository
	notificationSvc *NotificationService
}

func NewTechnicianService(
	techRepo *repository.TechnicianRepository,
	userRepo *repository.UserRepository,
	appointmentRepo *repository.AppointmentRepository,
	notificationSvc *NotificationService,
) *TechnicianService {
	return &TechnicianService{
		techRepo:        techRepo,
		userRepo:        userRepo,
		appointmentRepo: appointmentRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *TechnicianService) SetNotificationService(svc *NotificationService) {
	s.notificationSvc = svc
}

type CreateTechnicianRequest struct {
	UserID        uint   `json:"user_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Title         string `json:"title"`
	Avatar        string `json:"avatar"`
	Specialties   string `json:"specialties"`
	Description   string `json:"description"`
	WorkStartTime string `json:"work_start_time"`
	WorkEndTime   string `json:"work_end_time"`
	WorkDays      string `json:"work_days"`
}

type UpdateTechnicianRequest struct {
	Name          string `json:"name"`
	Title         string `json:"title"`
	Avatar        string `json:"avatar"`
	Specialties   string `json:"specialties"`
	Description   string `json:"description"`
	WorkStartTime string `json:"work_start_time"`
	WorkEndTime   string `json:"work_end_time"`
	WorkDays      string `json:"work_days"`
	Status        int    `json:"status"`
}

func (s *TechnicianService) Create(req *CreateTechnicianRequest) (*model.Technician, error) {
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	existing, _ := s.techRepo.GetByUserID(req.UserID)
	if existing != nil {
		return nil, fmt.Errorf("technician profile already exists")
	}

	tech := &model.Technician{
		UserID:        req.UserID,
		Name:          req.Name,
		Title:         req.Title,
		Avatar:        req.Avatar,
		Specialties:   req.Specialties,
		Description:   req.Description,
		Rating:        5.0,
		WorkStartTime: req.WorkStartTime,
		WorkEndTime:   req.WorkEndTime,
		WorkDays:      req.WorkDays,
		Status:        1,
		User:          user,
	}

	if tech.WorkStartTime == "" {
		tech.WorkStartTime = "09:00"
	}
	if tech.WorkEndTime == "" {
		tech.WorkEndTime = "21:00"
	}
	if tech.WorkDays == "" {
		tech.WorkDays = "1,2,3,4,5,6"
	}

	if err := s.techRepo.Create(tech); err != nil {
		return nil, fmt.Errorf("create technician: %w", err)
	}

	return tech, nil
}

func (s *TechnicianService) GetByID(id uint) (*model.Technician, error) {
	return s.techRepo.GetByID(id)
}

func (s *TechnicianService) GetByUserID(userID uint) (*model.Technician, error) {
	return s.techRepo.GetByUserID(userID)
}

func (s *TechnicianService) Update(id uint, req *UpdateTechnicianRequest) (*model.Technician, error) {
	tech, err := s.techRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("technician not found: %w", err)
	}

	if req.Name != "" {
		tech.Name = req.Name
	}
	if req.Title != "" {
		tech.Title = req.Title
	}
	if req.Avatar != "" {
		tech.Avatar = req.Avatar
	}
	if req.Specialties != "" {
		tech.Specialties = req.Specialties
	}
	if req.Description != "" {
		tech.Description = req.Description
	}
	if req.WorkStartTime != "" {
		tech.WorkStartTime = req.WorkStartTime
	}
	if req.WorkEndTime != "" {
		tech.WorkEndTime = req.WorkEndTime
	}
	if req.WorkDays != "" {
		tech.WorkDays = req.WorkDays
	}
	if req.Status > 0 {
		tech.Status = req.Status
	}

	if err := s.techRepo.Update(tech); err != nil {
		return nil, fmt.Errorf("update technician: %w", err)
	}

	return tech, nil
}

func (s *TechnicianService) List(page, pageSize int, status int) ([]model.Technician, int64, error) {
	return s.techRepo.List(page, pageSize, status)
}

func (s *TechnicianService) ListAll() ([]model.Technician, error) {
	return s.techRepo.ListAll()
}

func (s *TechnicianService) AddLeave(technicianID uint, dateStr string, reason string) error {
	t, _ := time.Parse("2006-01-02", dateStr)

	if err := s.techRepo.AddLeave(&model.TechnicianLeave{
		TechnicianID: technicianID,
		LeaveDate:    t,
		Reason:       reason,
		Approved:     true,
	}); err != nil {
		return err
	}

	appointments, err := s.appointmentRepo.GetByTechnicianAndDate(technicianID, t)
	if err != nil {
		return nil
	}

	for _, appt := range appointments {
		if appt.Status == "confirmed" {
			appt.Status = "rescheduled"
			appt.CancelReason = fmt.Sprintf("技师请假：%s", reason)
			if err := s.appointmentRepo.Update(&appt); err != nil {
				continue
			}

			fullAppt, _ := s.appointmentRepo.GetByID(appt.ID)
			if fullAppt != nil && s.notificationSvc != nil {
				s.notificationSvc.SendAppointmentRescheduleNotification(fullAppt)
			}
		}
	}

	return nil
}

func (s *TechnicianService) GetLeaves(technicianID uint, month int) ([]model.TechnicianLeave, error) {
	return s.techRepo.GetLeaves(technicianID, month)
}
