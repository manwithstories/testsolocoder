package repository

import (
	"beauty-salon-system/internal/model"

	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) List(page, pageSize int, filters map[string]interface{}) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.Model(&model.AuditLog{}).Preload("User")
	for key, value := range filters {
		switch key {
		case "user_id":
			query = query.Where("user_id = ?", value)
		case "module":
			query = query.Where("module = ?", value)
		case "action":
			query = query.Where("action = ?", value)
		}
	}
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetByUserID(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&notifications).Error
	return notifications, total, err
}

func (r *NotificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkAsRead(userID uint, id uint) error {
	return r.db.Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) GetByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Preload("Customer.User").Preload("Technician.User").
		Preload("Service").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) GetByTechnicianID(technicianID uint, page, pageSize int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	query := r.db.Model(&model.Review{}).Where("technician_id = ?", technicianID).
		Preload("Customer.User").Preload("Service")
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) GetAverageRating(technicianID uint) (float64, int, error) {
	var result struct {
		AvgRating   float64
		TotalCount  int
	}
	err := r.db.Model(&model.Review{}).
		Where("technician_id = ?", technicianID).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as total_count").
		Scan(&result).Error
	return result.AvgRating, result.TotalCount, err
}

func (r *ReviewRepository) GetByAppointmentID(appointmentID uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("appointment_id = ?", appointmentID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}
