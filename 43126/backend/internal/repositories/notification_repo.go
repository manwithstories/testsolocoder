package repositories

import (
	"time"

	"meeting-room/internal/models"
	"meeting-room/internal/utils"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return utils.DB.Create(notification).Error
}

func (r *NotificationRepository) FindByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	err := utils.DB.First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) Update(notification *models.Notification) error {
	return utils.DB.Save(notification).Error
}

func (r *NotificationRepository) GetPending() ([]models.Notification, error) {
	var notifications []models.Notification
	err := utils.DB.Where("status = ?", models.NotificationStatusPending).
		Where("retry_count < ?", 3).
		Limit(50).
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkSent(id uint) error {
	now := time.Now()
	return utils.DB.Model(&models.Notification{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  models.NotificationStatusSent,
			"sent_at": now,
		}).Error
}

func (r *NotificationRepository) MarkFailed(id uint, errMsg string) error {
	return utils.DB.Model(&models.Notification{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     models.NotificationStatusFailed,
			"error_msg":  errMsg,
			"retry_count": 3,
		}).Error
}
