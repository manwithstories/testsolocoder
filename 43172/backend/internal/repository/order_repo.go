package repository

import (
	"errors"
	"time"

	"luxury-trading-platform/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Buyer").
		Preload("Seller").
		Preload("Product").
		Preload("Authentication").
		Preload("Review").
		First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByOrderNumber(orderNumber string) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("order_number = ?", orderNumber).
		Preload("Buyer").
		Preload("Seller").
		Preload("Product").
		Preload("Authentication").
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) List(page, pageSize int, buyerID, sellerID *uint, status model.OrderStatus) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{})
	if buyerID != nil {
		query = query.Where("buyer_id = ?", *buyerID)
	}
	if sellerID != nil {
		query = query.Where("seller_id = ?", *sellerID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Buyer").
		Preload("Seller").
		Preload("Product").
		Preload("Authentication").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *OrderRepository) UpdateStatus(id uint, status model.OrderStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	now := time.Now()
	switch status {
	case model.OrderStatusPaid:
		updates["payment_time"] = now
	case model.OrderStatusShipped:
		updates["shipped_at"] = now
	case model.OrderStatusDelivered:
		updates["delivered_at"] = now
	case model.OrderStatusCompleted:
		updates["completed_at"] = now
	case model.OrderStatusCancelled:
		updates["cancelled_at"] = now
	}
	return r.db.Model(&model.Order{}).Where("id = ?", id).Updates(updates).Error
}

func (r *OrderRepository) UpdatePaymentStatus(id uint, status model.PaymentStatus, method string) error {
	return r.db.Model(&model.Order{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status": status,
			"payment_method": method,
			"payment_time":   time.Now(),
		}).Error
}

func (r *OrderRepository) UpdateShipping(id uint, trackingNumber, address string) error {
	return r.db.Model(&model.Order{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"tracking_number": trackingNumber,
			"shipping_address": address,
			"shipped_at":      time.Now(),
		}).Error
}
