package repository

import (
	"car-rental/internal/model"
	cachedb "car-rental/internal/config"
	"time"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{db: cachedb.DB}
}

func (r *BookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) FindByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("User").Preload("Car").Preload("PickupStore.City").Preload("ReturnStore.City").Preload("PromoCode").First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindByBookingNo(bookingNo string) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("User").Preload("Car").Preload("PickupStore.City").Preload("ReturnStore.City").Preload("PromoCode").Where("booking_no = ?", bookingNo).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindAll(page, pageSize int, userID uint, status string, carID uint, startDate, endDate *time.Time) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64

	query := r.db.Model(&model.Booking{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if carID > 0 {
		query = query.Where("car_id = ?", carID)
	}
	if startDate != nil {
		query = query.Where("pickup_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("pickup_time <= ?", *endDate)
	}

	query.Count(&total)
	err := query.Preload("User").Preload("Car").Preload("PickupStore.City").Preload("ReturnStore.City").Preload("PromoCode").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepository) UpdateStatus(id uint, status model.BookingStatus) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BookingRepository) Cancel(id uint, reason string) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        model.BookingStatusCancelled,
		"cancel_reason": reason,
		"cancelled_at":  gorm.Expr("NOW()"),
	}).Error
}

func (r *BookingRepository) IsCarAvailable(carID uint, pickupTime, returnTime time.Time) bool {
	var count int64
	r.db.Model(&model.Booking{}).Where(
		"car_id = ? AND status IN ? AND pickup_time < ? AND return_time > ?",
		carID,
		[]string{string(model.BookingStatusPending), string(model.BookingStatusConfirmed)},
		returnTime,
		pickupTime,
	).Count(&count)
	return count == 0
}

func (r *BookingRepository) Delete(id uint) error {
	return r.db.Delete(&model.Booking{}, id).Error
}

func (r *BookingRepository) FindUpcomingPickups(startTime, endTime time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Preload("User").Preload("Car").
		Where("status IN ? AND pickup_time >= ? AND pickup_time <= ? AND pickup_reminder_sent = ?",
			[]string{string(model.BookingStatusPending), string(model.BookingStatusConfirmed)},
			startTime, endTime, false).
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) FindUpcomingReturns(startTime, endTime time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Preload("User").Preload("Car").
		Where("status = ? AND return_time >= ? AND return_time <= ? AND return_reminder_sent = ?",
			string(model.BookingStatusConfirmed),
			startTime, endTime, false).
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) MarkPickupReminderSent(id uint) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Update("pickup_reminder_sent", true).Error
}

func (r *BookingRepository) MarkReturnReminderSent(id uint) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Update("return_reminder_sent", true).Error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: cachedb.DB}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("User").Preload("Car").Preload("Booking.PickupStore.City").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("User").Preload("Car").Preload("Booking.PickupStore.City").Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindAll(page, pageSize int, userID uint, status string, startDate, endDate *time.Time) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	query.Count(&total)
	err := query.Preload("User").Preload("Car").Preload("Booking.PickupStore.City").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) UpdateStatus(id uint, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) UpdatePaymentStatus(id uint, paymentStatus model.PaymentStatus, paidAmount float64) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"payment_status": paymentStatus,
		"paid_amount":    paidAmount,
		"paid_at":        gorm.Expr("NOW()"),
	}).Error
}

func (r *OrderRepository) FindAllForExport(status string, startDate, endDate *time.Time) ([]model.Order, error) {
	var orders []model.Order
	query := r.db.Model(&model.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}
	err := query.Preload("User").Preload("Car").Preload("Booking.PickupStore.City").Order("created_at DESC").Find(&orders).Error
	return orders, err
}