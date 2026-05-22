package repository

import (
	"hotel-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	GetByID(id uint) (*model.Payment, error)
	GetByPaymentNo(paymentNo string) (*model.Payment, error)
	Update(payment *model.Payment) error
	List(page, pageSize int, paymentNo string, orderType model.OrderType, orderID uint, paymentMethod model.PaymentMethod, paymentType model.PaymentType, status model.PaymentStatus, startDate, endDate string) ([]model.Payment, int64, error)
	GetPaymentsByOrder(orderType model.OrderType, orderID uint) ([]model.Payment, error)
	GetTotalPaidByOrder(orderType model.OrderType, orderID uint) (float64, error)
	UpdateStatus(id uint, status model.PaymentStatus) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByID(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByPaymentNo(paymentNo string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("payment_no = ?", paymentNo).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) List(page, pageSize int, paymentNo string, orderType model.OrderType, orderID uint, paymentMethod model.PaymentMethod, paymentType model.PaymentType, status model.PaymentStatus, startDate, endDate string) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	query := r.db.Model(&model.Payment{})

	if paymentNo != "" {
		query = query.Where("payment_no LIKE ?", "%"+paymentNo+"%")
	}
	if orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}
	if orderID > 0 {
		query = query.Where("order_id = ?", orderID)
	}
	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}
	if paymentType != "" {
		query = query.Where("payment_type = ?", paymentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}
	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("created_at < ?", end)
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *paymentRepository) GetPaymentsByOrder(orderType model.OrderType, orderID uint) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Where("order_type = ? AND order_id = ?", orderType, orderID).
		Order("created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) GetTotalPaidByOrder(orderType model.OrderType, orderID uint) (float64, error) {
	type Result struct {
		PaymentType model.PaymentType
		TotalAmount float64
	}

	var results []Result
	err := r.db.Model(&model.Payment{}).
		Select("payment_type, COALESCE(SUM(amount), 0) as total_amount").
		Where("order_type = ? AND order_id = ? AND status = ?", orderType, orderID, model.PaymentStatusCompleted).
		Group("payment_type").
		Scan(&results).Error
	if err != nil {
		return 0, err
	}

	totalPaid := 0.0
	for _, r := range results {
		if r.PaymentType == model.PaymentTypeRefund {
			totalPaid -= r.TotalAmount
		} else {
			totalPaid += r.TotalAmount
		}
	}

	return totalPaid, nil
}

func (r *paymentRepository) UpdateStatus(id uint, status model.PaymentStatus) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", id).Update("status", status).Error
}
