package repository

import (
	"drone-rental/internal/config"
	"drone-rental/internal/model"
	"time"
)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

func (r *OrderRepo) Create(order *model.RentalOrder) error {
	return config.DB.Create(order).Error
}

func (r *OrderRepo) GetByID(id uint) (*model.RentalOrder, error) {
	var order model.RentalOrder
	err := config.DB.Preload("User").Preload("Drone").First(&order, id).Error
	return &order, err
}

func (r *OrderRepo) GetByOrderNo(orderNo string) (*model.RentalOrder, error) {
	var order model.RentalOrder
	err := config.DB.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (r *OrderRepo) Update(order *model.RentalOrder) error {
	return config.DB.Save(order).Error
}

func (r *OrderRepo) UpdateStatus(id uint, status model.OrderStatus) error {
	return config.DB.Model(&model.RentalOrder{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *OrderRepo) List(page, pageSize int, userID, droneID uint, status model.OrderStatus) ([]model.RentalOrder, int64, error) {
	var orders []model.RentalOrder
	var total int64
	db := config.DB.Model(&model.RentalOrder{}).Preload("Drone").Preload("User")
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	if droneID != 0 {
		db = db.Where("drone_id = ?", droneID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)
	err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

type PaymentRepo struct{}

func NewPaymentRepo() *PaymentRepo {
	return &PaymentRepo{}
}

func (r *PaymentRepo) Create(payment *model.Payment) error {
	return config.DB.Create(payment).Error
}

func (r *PaymentRepo) GetByID(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := config.DB.First(&payment, id).Error
	return &payment, err
}

func (r *PaymentRepo) GetByPaymentNo(paymentNo string) (*model.Payment, error) {
	var payment model.Payment
	err := config.DB.Where("payment_no = ?", paymentNo).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepo) Update(payment *model.Payment) error {
	return config.DB.Save(payment).Error
}

func (r *PaymentRepo) GetByOrderID(orderID uint) (*model.Payment, error) {
	var payment model.Payment
	err := config.DB.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepo) ListByUserID(userID uint, page, pageSize int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64
	db := config.DB.Model(&model.Payment{}).Where("user_id = ?", userID)
	db.Count(&total)
	err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&payments).Error
	return payments, total, err
}

func (r *OrderRepo) GetOverdueOrders() ([]model.RentalOrder, error) {
	var orders []model.RentalOrder
	err := config.DB.Where("status = ? AND end_date < ?", model.OrderStatusPicked, time.Now()).
		Find(&orders).Error
	return orders, err
}
