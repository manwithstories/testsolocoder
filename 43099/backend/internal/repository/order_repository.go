package repository

import (
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return database.DB.Create(order).Error
}

func (r *OrderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	err := database.DB.Preload("User").First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := database.DB.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (r *OrderRepository) List(req *dto.OrderListRequest, userID *uint) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := database.DB.Model(&model.Order{}).Preload("User")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) Update(order *model.Order) error {
	return database.DB.Save(order).Error
}

func (r *OrderRepository) GetVenueBookingsForDate(venueID uint, date string) ([]model.Order, error) {
	var orders []model.Order
	startDate := date + " 00:00:00"
	endDate := date + " 23:59:59"

	err := database.DB.Where("type = ? AND item_id = ? AND status IN (?, ?, ?) AND start_time >= ? AND end_time <= ?",
		model.OrderTypeVenue, venueID,
		model.OrderStatusPending, model.OrderStatusConfirmed, model.OrderStatusPaid,
		startDate, endDate).Find(&orders).Error

	return orders, err
}

func (r *OrderRepository) GetDeviceBookingsForDate(deviceID uint, date string) ([]model.Order, error) {
	var orders []model.Order
	startDate := date + " 00:00:00"
	endDate := date + " 23:59:59"

	err := database.DB.Where("type = ? AND item_id = ? AND status IN (?, ?, ?) AND start_time >= ? AND end_time <= ?",
		model.OrderTypeDevice, deviceID,
		model.OrderStatusPending, model.OrderStatusConfirmed, model.OrderStatusPaid,
		startDate, endDate).Find(&orders).Error

	return orders, err
}

func (r *OrderRepository) CheckVenueConflict(venueID uint, startTime, endTime string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Order{}).Where(
		"type = ? AND item_id = ? AND status IN (?, ?, ?) AND start_time < ? AND end_time > ?",
		model.OrderTypeVenue, venueID,
		model.OrderStatusPending, model.OrderStatusConfirmed, model.OrderStatusPaid,
		endTime, startTime,
	).Count(&count).Error

	return count > 0, err
}

func (r *OrderRepository) GetCalendarBookings(startDate, endDate string) ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Where(
		"status IN (?, ?, ?) AND start_time >= ? AND end_time <= ?",
		model.OrderStatusPending, model.OrderStatusConfirmed, model.OrderStatusPaid,
		startDate, endDate,
	).Find(&orders).Error

	return orders, err
}
