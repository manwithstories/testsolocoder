package repositories

import (
	"consultation-platform/models"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{db: utils.GetDB()}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) FindByID(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) FindByUserID(userID uuid.UUID, page, pageSize int, onlyUnread bool) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	if onlyUnread {
		query = query.Where("is_read = ?", false)
	}
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *NotificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_read": true,
	}).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Updates(map[string]interface{}{
		"is_read": true,
	}).Error
}

func (r *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

type NotificationTemplateRepository struct {
	db *gorm.DB
}

func NewNotificationTemplateRepository() *NotificationTemplateRepository {
	return &NotificationTemplateRepository{db: utils.GetDB()}
}

func (r *NotificationTemplateRepository) Create(template *models.NotificationTemplate) error {
	return r.db.Create(template).Error
}

func (r *NotificationTemplateRepository) FindByID(id uuid.UUID) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *NotificationTemplateRepository) FindByType(templateType string) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.Where("type = ? AND is_active = ?", templateType, true).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *NotificationTemplateRepository) Update(template *models.NotificationTemplate) error {
	return r.db.Save(template).Error
}

func (r *NotificationTemplateRepository) FindAll(page, pageSize int) ([]models.NotificationTemplate, int64, error) {
	var templates []models.NotificationTemplate
	var total int64

	offset, limit := utils.Pagination(page, pageSize)

	query := r.db.Model(&models.NotificationTemplate{})
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}
