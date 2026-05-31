package repository

import (
	"health-platform/models"
	"time"

	"gorm.io/gorm"
)

type TimeSlotRepository struct {
	*BaseRepository
}

func NewTimeSlotRepository() *TimeSlotRepository {
	return &TimeSlotRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *TimeSlotRepository) FindByPackageID(packageID uint, startDate, endDate *time.Time) ([]models.TimeSlot, error) {
	var timeSlots []models.TimeSlot
	
	query := r.DB.Where("package_id = ?", packageID)
	
	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}
	
	err := query.Order("date ASC, start_time ASC").Find(&timeSlots).Error
	return timeSlots, err
}

func (r *TimeSlotRepository) FindAvailableByPackageID(packageID uint) ([]models.TimeSlot, error) {
	var timeSlots []models.TimeSlot
	err := r.DB.Where("package_id = ? AND date >= ? AND booked < total AND status = 1",
		packageID, time.Now().Truncate(24*time.Hour)).
		Order("date ASC, start_time ASC").
		Find(&timeSlots).Error
	return timeSlots, err
}

func (r *TimeSlotRepository) IncrementBooked(timeSlotID uint) error {
	return r.DB.Model(&models.TimeSlot{}).Where("id = ?", timeSlotID).
		Update("booked", gorm.Expr("booked + 1")).Error
}

func (r *TimeSlotRepository) DecrementBooked(timeSlotID uint) error {
	return r.DB.Model(&models.TimeSlot{}).Where("id = ?", timeSlotID).
		Update("booked", gorm.Expr("booked - 1")).Error
}

func (r *TimeSlotRepository) BatchCreate(timeSlots []models.TimeSlot) error {
	return r.DB.Create(&timeSlots).Error
}
