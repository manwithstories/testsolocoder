package repository

import (
	"car-rental/internal/model"
	cachedb "car-rental/internal/config"
	"time"

	"gorm.io/gorm"
)

type MaintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository() *MaintenanceRepository {
	return &MaintenanceRepository{db: cachedb.DB}
}

func (r *MaintenanceRepository) Create(plan *model.MaintenancePlan) error {
	return r.db.Create(plan).Error
}

func (r *MaintenanceRepository) FindByID(id uint) (*model.MaintenancePlan, error) {
	var plan model.MaintenancePlan
	err := r.db.Preload("Car").First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *MaintenanceRepository) FindAll(page, pageSize int, carID uint, status string) ([]model.MaintenancePlan, int64, error) {
	var plans []model.MaintenancePlan
	var total int64

	query := r.db.Model(&model.MaintenancePlan{})
	if carID > 0 {
		query = query.Where("car_id = ?", carID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Car").Offset((page - 1) * pageSize).Limit(pageSize).Order("start_date DESC").Find(&plans).Error
	return plans, total, err
}

func (r *MaintenanceRepository) FindByCarID(carID uint, page, pageSize int) ([]model.MaintenancePlan, int64, error) {
	var plans []model.MaintenancePlan
	var total int64

	query := r.db.Model(&model.MaintenancePlan{}).Where("car_id = ?", carID)
	query.Count(&total)
	err := query.Preload("Car").Offset((page - 1) * pageSize).Limit(pageSize).Order("start_date DESC").Find(&plans).Error
	return plans, total, err
}

func (r *MaintenanceRepository) FindUpcoming() ([]model.MaintenancePlan, error) {
	var plans []model.MaintenancePlan
	now := time.Now()
	err := r.db.Preload("Car").Where("status = ? AND start_date <= ?", model.MaintenanceStatusScheduled, now.AddDate(0, 0, 7)).Order("start_date ASC").Find(&plans).Error
	return plans, err
}

func (r *MaintenanceRepository) Update(plan *model.MaintenancePlan) error {
	return r.db.Save(plan).Error
}

func (r *MaintenanceRepository) Delete(id uint) error {
	return r.db.Delete(&model.MaintenancePlan{}, id).Error
}

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{db: cachedb.DB}
}

func (r *MessageRepository) Create(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) FindByID(id uint) (*model.Message, error) {
	var message model.Message
	err := r.db.First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) FindByUserID(userID uint, page, pageSize int) ([]model.Message, int64, error) {
	var messages []model.Message
	var total int64

	query := r.db.Model(&model.Message{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&messages).Error
	return messages, total, err
}

func (r *MessageRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Message{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *MessageRepository) MarkAsRead(id uint, userID uint) error {
	return r.db.Model(&model.Message{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}

func (r *MessageRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&model.Message{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *MessageRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Message{}).Error
}

func (r *MessageRepository) CreateNotification(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *MessageRepository) UpdateNotification(notification *model.Notification) error {
	return r.db.Save(notification).Error
}

func (r *CarRepository) CountCars() (int64, error) {
	var count int64
	err := r.db.Model(&model.Car{}).Count(&count).Error
	return count, err
}

func (r *CarRepository) CountCarsByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Car{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *CarRepository) GetPopularCars(limit int) ([]CarStat, error) {
	var stats []CarStat
	err := r.db.Table("cars").
		Select("cars.id, cars.brand, cars.model, cars.rating, cars.review_count, COUNT(bookings.id) as booking_count").
		Joins("LEFT JOIN bookings ON bookings.car_id = cars.id AND bookings.status IN ('pending', 'confirmed', 'completed')").
		Group("cars.id, cars.brand, cars.model, cars.rating, cars.review_count").
		Order("rating DESC, review_count DESC, booking_count DESC").
		Limit(limit).
		Scan(&stats).Error
	return stats, err
}

func (r *CarRepository) GetStatusBreakdown() ([]StatusCount, error) {
	var breakdown []StatusCount
	err := r.db.Model(&model.Car{}).Select("status as status, count(*) as count").Group("status").Scan(&breakdown).Error
	return breakdown, err
}

type CarStat struct {
	ID          uint    `json:"id"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`
	BookingCount int64  `json:"booking_count"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

func (r *OrderRepository) GetTotalRevenue(startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := r.db.Model(&model.Order{}).Where("payment_status = ?", model.PaymentStatusPaid)
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	err := query.Select("COALESCE(SUM(final_amount), 0)").Scan(&total).Error
	return total, err
}
