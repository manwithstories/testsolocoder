package repository

import (
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"
)

type PaymentRepository struct{}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (r *PaymentRepository) Create(payment *model.Payment) error {
	return database.DB.Create(payment).Error
}

func (r *PaymentRepository) GetByID(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := database.DB.Preload("Order").First(&payment, id).Error
	return &payment, err
}

func (r *PaymentRepository) GetByOrderID(orderID uint) (*model.Payment, error) {
	var payment model.Payment
	err := database.DB.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepository) List(req *dto.PaymentListRequest) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	query := database.DB.Model(&model.Payment{}).Preload("Order")

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&payments).Error
	return payments, total, err
}

func (r *PaymentRepository) Update(payment *model.Payment) error {
	return database.DB.Save(payment).Error
}

func (r *PaymentRepository) GetByDateRange(startDate, endDate string) ([]model.Payment, error) {
	var payments []model.Payment
	err := database.DB.Preload("Order").Where(
		"created_at >= ? AND created_at <= ?",
		startDate+" 00:00:00",
		endDate+" 23:59:59",
	).Order("created_at DESC").Find(&payments).Error
	return payments, err
}
