package services

import (
	"errors"
	"time"

	"consultation-platform/models"
	"consultation-platform/repositories"
	"consultation-platform/utils"

	"github.com/google/uuid"
)

type ServiceService struct {
	serviceRepo  *repositories.ServiceRepository
	scheduleRepo *repositories.ScheduleRepository
	userRepo     *repositories.UserRepository
}

func NewServiceService() *ServiceService {
	return &ServiceService{
		serviceRepo:  repositories.NewServiceRepository(),
		scheduleRepo: repositories.NewScheduleRepository(),
		userRepo:     repositories.NewUserRepository(),
	}
}

type CreateServiceRequest struct {
	Title           string  `json:"title" binding:"required"`
	Description     string  `json:"description"`
	ServiceType     string  `json:"service_type" binding:"required"`
	Price           float64 `json:"price" binding:"required,gt=0"`
	DurationMinutes int     `json:"duration_minutes" binding:"required,gt=0"`
	Tags            string  `json:"tags"`
}

type CreateScheduleRequest struct {
	ServiceID uuid.UUID `json:"service_id" binding:"required"`
	Date      string    `json:"date" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
}

type BatchCreateScheduleRequest struct {
	ServiceID   uuid.UUID `json:"service_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     string    `json:"end_date" binding:"required"`
	TimeSlots   []TimeSlot `json:"time_slots" binding:"required"`
}

type TimeSlot struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type ServiceServiceInterface interface {
	CreateService(professionalID uuid.UUID, req *CreateServiceRequest) (*models.Service, error)
	GetServiceByID(id uuid.UUID) (*models.Service, error)
	UpdateService(id uuid.UUID, professionalID uuid.UUID, updates map[string]interface{}) (*models.Service, error)
	DeleteService(id uuid.UUID, professionalID uuid.UUID) error
	GetServicesByProfessional(professionalID uuid.UUID, page, pageSize int) ([]models.Service, int64, error)
	GetAllServices(page, pageSize int, serviceType string) ([]models.Service, int64, error)
	CreateSchedule(professionalID uuid.UUID, req *CreateScheduleRequest) (*models.Schedule, error)
	BatchCreateSchedules(professionalID uuid.UUID, req *BatchCreateScheduleRequest) error
	GetSchedulesByServiceIDAndDate(serviceID uuid.UUID, date string) ([]models.Schedule, error)
	GetAvailableSchedules(serviceID uuid.UUID, date string) ([]models.Schedule, error)
	DeleteSchedulesByDateRange(professionalID uuid.UUID, serviceID uuid.UUID, startDate, endDate string) error
}

func (s *ServiceService) CreateService(professionalID uuid.UUID, req *CreateServiceRequest) (*models.Service, error) {
	professional, err := s.userRepo.FindActiveProfessionalByID(professionalID)
	if err != nil {
		return nil, errors.New("professional not found or not active")
	}
	_ = professional

	service := &models.Service{
		ProfessionalID:  professionalID,
		Title:           req.Title,
		Description:     req.Description,
		ServiceType:     models.ServiceType(req.ServiceType),
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		Status:          models.ServiceStatusActive,
		Tags:            req.Tags,
	}

	if err := s.serviceRepo.Create(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	return s.serviceRepo.FindByID(id)
}

func (s *ServiceService) UpdateService(id uuid.UUID, professionalID uuid.UUID, updates map[string]interface{}) (*models.Service, error) {
	service, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if service.ProfessionalID != professionalID {
		return nil, errors.New("you do not have permission to update this service")
	}

	if title, ok := updates["title"].(string); ok {
		service.Title = title
	}
	if desc, ok := updates["description"].(string); ok {
		service.Description = desc
	}
	if serviceType, ok := updates["service_type"].(string); ok {
		service.ServiceType = models.ServiceType(serviceType)
	}
	if price, ok := updates["price"].(float64); ok {
		service.Price = price
	}
	if duration, ok := updates["duration_minutes"].(int); ok {
		service.DurationMinutes = duration
	}
	if tags, ok := updates["tags"].(string); ok {
		service.Tags = tags
	}
	if status, ok := updates["status"].(string); ok {
		service.Status = models.ServiceStatus(status)
	}

	if err := s.serviceRepo.Update(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) DeleteService(id uuid.UUID, professionalID uuid.UUID) error {
	service, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return err
	}

	if service.ProfessionalID != professionalID {
		return errors.New("you do not have permission to delete this service")
	}

	return s.serviceRepo.Delete(id)
}

func (s *ServiceService) GetServicesByProfessional(professionalID uuid.UUID, page, pageSize int) ([]models.Service, int64, error) {
	return s.serviceRepo.FindByProfessionalID(professionalID, page, pageSize)
}

func (s *ServiceService) GetAllServices(page, pageSize int, serviceType string) ([]models.Service, int64, error) {
	return s.serviceRepo.FindAll(page, pageSize, serviceType)
}

func (s *ServiceService) CreateSchedule(professionalID uuid.UUID, req *CreateScheduleRequest) (*models.Schedule, error) {
	service, err := s.serviceRepo.FindByID(req.ServiceID)
	if err != nil {
		return nil, err
	}

	if service.ProfessionalID != professionalID {
		return nil, errors.New("you do not have permission to add schedules for this service")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("cannot create schedule for past dates")
	}

	schedule := &models.Schedule{
		ServiceID:   req.ServiceID,
		Date:        date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsBooked:    false,
		IsAvailable: true,
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *ServiceService) BatchCreateSchedules(professionalID uuid.UUID, req *BatchCreateScheduleRequest) error {
	service, err := s.serviceRepo.FindByID(req.ServiceID)
	if err != nil {
		return err
	}

	if service.ProfessionalID != professionalID {
		return errors.New("you do not have permission to add schedules for this service")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return errors.New("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return errors.New("invalid end date format")
	}

	if startDate.After(endDate) {
		return errors.New("start date must be before end date")
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("cannot create schedule for past dates")
	}

	var schedules []models.Schedule
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		for _, slot := range req.TimeSlots {
			schedules = append(schedules, models.Schedule{
				ServiceID:   req.ServiceID,
				Date:        d,
				StartTime:   slot.StartTime,
				EndTime:     slot.EndTime,
				IsBooked:    false,
				IsAvailable: true,
			})
		}
	}

	if len(schedules) > 0 {
		return s.scheduleRepo.CreateBatch(schedules)
	}

	return nil
}

func (s *ServiceService) GetSchedulesByServiceIDAndDate(serviceID uuid.UUID, dateStr string) ([]models.Schedule, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format")
	}
	return s.scheduleRepo.FindByServiceIDAndDate(serviceID, date)
}

func (s *ServiceService) GetAvailableSchedules(serviceID uuid.UUID, dateStr string) ([]models.Schedule, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format")
	}
	return s.scheduleRepo.FindAvailableByServiceIDAndDate(serviceID, date)
}

func (s *ServiceService) DeleteSchedulesByDateRange(professionalID uuid.UUID, serviceID uuid.UUID, startDateStr, endDateStr string) error {
	service, err := s.serviceRepo.FindByID(serviceID)
	if err != nil {
		return err
	}

	if service.ProfessionalID != professionalID {
		return errors.New("you do not have permission to delete schedules for this service")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return errors.New("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return errors.New("invalid end date format")
	}

	return s.scheduleRepo.DeleteByServiceIDAndDateRange(serviceID, startDate, endDate)
}
