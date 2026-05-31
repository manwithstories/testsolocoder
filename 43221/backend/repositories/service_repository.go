package repositories

import (
	"errors"
	"time"

	"consultation-platform/models"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{db: utils.GetDB()}
}

func (r *ServiceRepository) Create(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *ServiceRepository) FindByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service
	err := r.db.Preload("Professional").First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Update(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *ServiceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Service{}, "id = ?", id).Error
}

func (r *ServiceRepository) FindByProfessionalID(professionalID uuid.UUID, page, pageSize int) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Service{}).Where("professional_id = ?", professionalID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Preload("Professional").Find(&services).Error
	if err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *ServiceRepository) FindAll(page, pageSize int, serviceType string) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Service{}).Where("status = ?", models.ServiceStatusActive)
	if serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Preload("Professional").Find(&services).Error
	if err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *ServiceRepository) UpdateRating(serviceID uuid.UUID, avgRating float64, reviewCount int) error {
	return r.db.Model(&models.Service{}).Where("id = ?", serviceID).Updates(map[string]interface{}{
		"average_rating": avgRating,
		"review_count":   reviewCount,
	}).Error
}

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{db: utils.GetDB()}
}

func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *ScheduleRepository) CreateBatch(schedules []models.Schedule) error {
	return r.db.Create(&schedules).Error
}

func (r *ScheduleRepository) FindByID(id uuid.UUID) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.First(&schedule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) FindByServiceIDAndDate(serviceID uuid.UUID, date time.Time) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Where("service_id = ? AND date = ?", serviceID, date.Format("2006-01-02")).Order("start_time ASC").Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) FindAvailableByServiceIDAndDate(serviceID uuid.UUID, date time.Time) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Where("service_id = ? AND date = ? AND is_available = ? AND is_booked = ?",
		serviceID, date.Format("2006-01-02"), true, false).Order("start_time ASC").Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) MarkAsBooked(id uuid.UUID) error {
	return r.db.Model(&models.Schedule{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_booked": true,
	}).Error
}

func (r *ScheduleRepository) MarkAsAvailable(id uuid.UUID) error {
	return r.db.Model(&models.Schedule{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_booked": false,
	}).Error
}

func (r *ScheduleRepository) DeleteByServiceIDAndDateRange(serviceID uuid.UUID, startDate, endDate time.Time) error {
	return r.db.Where("service_id = ? AND date >= ? AND date <= ?",
		serviceID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Delete(&models.Schedule{}).Error
}

func (r *ScheduleRepository) WithTx(tx *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: tx}
}

func (r *ServiceRepository) WithTx(tx *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: tx}
}

func (r *ScheduleRepository) FindByIDWithLock(id uuid.UUID) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.Set("gorm:query_option", "FOR UPDATE").First(&schedule, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}
	return &schedule, nil
}
