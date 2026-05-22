package repository

import (
	"time"

	"gorm.io/gorm"
	"ticket-system/internal/database"
	"ticket-system/internal/models"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return database.DB.Create(order).Error
}

func (r *OrderRepository) GetByID(id uint64) (*models.Order, error) {
	var order models.Order
	err := database.DB.Preload("Tickets").Preload("Refund").First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	err := database.DB.Where("order_no = ?", orderNo).Preload("Tickets").Preload("Refund").First(&order).Error
	return &order, err
}

func (r *OrderRepository) List(userID uint64, page, pageSize int, status int, startDate, endDate, keyword string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := database.DB.Model(&models.Order{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("Tickets").Preload("Refund").Offset(offset).Limit(pageSize).Order("created_at desc").Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) Update(order *models.Order) error {
	return database.DB.Save(order).Error
}

func (r *OrderRepository) UpdateStatus(orderID uint64, status int) error {
	return database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *OrderRepository) MarkPaid(orderID uint64, payType int) error {
	now := time.Now()
	return database.DB.Model(&models.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":   models.OrderStatusPaid,
			"pay_type": payType,
			"pay_time": now,
		}).Error
}

func (r *OrderRepository) CreateTicket(ticket *models.Ticket) error {
	return database.DB.Create(ticket).Error
}

func (r *OrderRepository) GetTicketByID(id uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.First(&ticket, id).Error
	return &ticket, err
}

func (r *OrderRepository) GetTicketByTicketNo(ticketNo string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.Where("ticket_no = ?", ticketNo).First(&ticket).Error
	return &ticket, err
}

func (r *OrderRepository) GetTicketsByOrderID(orderID uint64) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Where("order_id = ?", orderID).Find(&tickets).Error
	return tickets, err
}

func (r *OrderRepository) UpdateTicket(ticket *models.Ticket) error {
	return database.DB.Save(ticket).Error
}

func (r *OrderRepository) MarkTicketCheckedIn(ticketID uint64) error {
	now := time.Now()
	return database.DB.Model(&models.Ticket{}).Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"checked_in":   1,
			"checkin_time": now,
		}).Error
}

func (r *OrderRepository) CreateRefund(refund *models.Refund) error {
	return database.DB.Create(refund).Error
}

func (r *OrderRepository) GetRefundByID(id uint64) (*models.Refund, error) {
	var refund models.Refund
	err := database.DB.First(&refund, id).Error
	return &refund, err
}

func (r *OrderRepository) GetRefundByRefundNo(refundNo string) (*models.Refund, error) {
	var refund models.Refund
	err := database.DB.Where("refund_no = ?", refundNo).First(&refund).Error
	return &refund, err
}

func (r *OrderRepository) GetRefundByOrderID(orderID uint64) (*models.Refund, error) {
	var refund models.Refund
	err := database.DB.Where("order_id = ?", orderID).First(&refund).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &refund, err
}

func (r *OrderRepository) UpdateRefund(refund *models.Refund) error {
	return database.DB.Save(refund).Error
}

func (r *OrderRepository) AuditRefund(refundID uint64, status int, remark string) error {
	now := time.Now()
	return database.DB.Model(&models.Refund{}).Where("id = ?", refundID).
		Updates(map[string]interface{}{
			"status":       status,
			"audit_time":   now,
			"audit_remark": remark,
		}).Error
}

func (r *OrderRepository) CreateCheckinLog(log *models.CheckinLog) error {
	return database.DB.Create(log).Error
}

func (r *OrderRepository) CreatePaymentLog(log *models.PaymentLog) error {
	return database.DB.Create(log).Error
}

func (r *OrderRepository) GetPaymentLogsByOrderID(orderID uint64) ([]models.PaymentLog, error) {
	var logs []models.PaymentLog
	err := database.DB.Where("order_id = ?", orderID).Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *OrderRepository) UpdatePaymentLogStatus(logID uint64, status int, response string) error {
	return database.DB.Model(&models.PaymentLog{}).Where("id = ?", logID).
		Updates(map[string]interface{}{
			"status":   status,
			"response": response,
		}).Error
}

func (r *OrderRepository) IncrementPaymentRetry(logID uint64) error {
	return database.DB.Model(&models.PaymentLog{}).Where("id = ?", logID).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

func (r *OrderRepository) GetAllOrdersForExport(startDate, endDate string, status int) ([]models.Order, error) {
	var orders []models.Order
	query := database.DB.Model(&models.Order{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}
	err := query.Preload("Tickets").Order("created_at desc").Find(&orders).Error
	return orders, err
}
