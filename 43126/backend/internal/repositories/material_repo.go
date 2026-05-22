package repositories

import (
	"time"

	"meeting-room/internal/models"
	"meeting-room/internal/utils"
)

type MaterialRepository struct{}

func NewMaterialRepository() *MaterialRepository {
	return &MaterialRepository{}
}

func (r *MaterialRepository) Create(material *models.MeetingMaterial) error {
	return utils.DB.Create(material).Error
}

func (r *MaterialRepository) FindByID(id uint) (*models.MeetingMaterial, error) {
	var material models.MeetingMaterial
	err := utils.DB.First(&material, id).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *MaterialRepository) Delete(id uint) error {
	return utils.DB.Delete(&models.MeetingMaterial{}, id).Error
}

func (r *MaterialRepository) GetByBooking(bookingID uint) ([]models.MeetingMaterial, error) {
	var materials []models.MeetingMaterial
	err := utils.DB.Where("booking_id = ?", bookingID).Find(&materials).Error
	return materials, err
}

func (r *MaterialRepository) GetExpired() ([]models.MeetingMaterial, error) {
	var materials []models.MeetingMaterial
	now := time.Now()
	err := utils.DB.Where("expire_at IS NOT NULL AND expire_at < ?", now).Find(&materials).Error
	return materials, err
}
