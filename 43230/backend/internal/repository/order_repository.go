package repository

import (
	"context"
	"errors"
	"time"

	"print3d-platform/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.PrintOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PrintOrder, error) {
	var order models.PrintOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Printer").
		Preload("Model").
		Preload("Material").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*models.PrintOrder, error) {
	var order models.PrintOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Printer").
		Preload("Model").
		Preload("Material").
		First(&order, "order_no = ?", orderNo).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.PrintOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	now := time.Now()
	switch status {
	case models.OrderStatusPrinting:
		updates["print_started_at"] = now
	case models.OrderStatusShipped:
		updates["shipped_at"] = now
	case models.OrderStatusDelivered:
		updates["delivered_at"] = now
	case models.OrderStatusCompleted:
		updates["completed_at"] = now
	case models.OrderStatusCancelled:
		updates["cancelled_at"] = now
	}

	return r.db.WithContext(ctx).Model(&models.PrintOrder{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *OrderRepository) AddHistory(ctx context.Context, history *models.OrderHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *OrderRepository) GetHistory(ctx context.Context, orderID uuid.UUID) ([]models.OrderHistory, error) {
	var histories []models.OrderHistory
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&histories).Error
	return histories, err
}

func (r *OrderRepository) ListCustomerOrders(ctx context.Context, customerID uuid.UUID, page, pageSize int, status *string) ([]models.PrintOrder, int64, error) {
	var orders []models.PrintOrder
	var total int64

	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Printer").
		Preload("Material").
		Where("customer_id = ?", customerID)

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) ListPrinterOrders(ctx context.Context, printerID uuid.UUID, page, pageSize int, status *string) ([]models.PrintOrder, int64, error) {
	var orders []models.PrintOrder
	var total int64

	query := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Model").
		Preload("Material").
		Where("printer_id = ?", printerID)

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) GetPendingOrders(ctx context.Context) ([]models.PrintOrder, error) {
	var orders []models.PrintOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Model").
		Preload("Material").
		Where("status = ?", models.OrderStatusPaid).
		Order("created_at ASC").
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) AssignPrinter(ctx context.Context, orderID, printerID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.PrintOrder{}).
		Where("id = ?", orderID).
		Update("printer_id", printerID).Error
}

func (r *OrderRepository) UpdateShipping(ctx context.Context, id uuid.UUID, trackingNumber, address string) error {
	return r.db.WithContext(ctx).Model(&models.PrintOrder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"tracking_number": trackingNumber,
			"shipping_address": address,
		}).Error
}

func (r *OrderRepository) CreateSettlement(ctx context.Context, settlement *models.Settlement) error {
	return r.db.WithContext(ctx).Create(settlement).Error
}

func (r *OrderRepository) GetSettlementByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Settlement, error) {
	var settlement models.Settlement
	err := r.db.WithContext(ctx).
		Preload("Order").
		First(&settlement, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &settlement, nil
}

func (r *OrderRepository) UpdateSettlement(ctx context.Context, settlement *models.Settlement) error {
	return r.db.WithContext(ctx).Save(settlement).Error
}

func (r *OrderRepository) GetCompletedOrdersForSettlement(ctx context.Context) ([]models.PrintOrder, error) {
	var orders []models.PrintOrder
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Printer").
		Preload("Model").
		Preload("Material").
		Joins("LEFT JOIN settlements ON settlements.order_id = print_orders.id").
		Where("print_orders.status = ? AND settlements.id IS NULL", models.OrderStatusCompleted).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetOrderStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalOrders int64
	r.db.WithContext(ctx).Model(&models.PrintOrder{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalOrders)
	stats["total_orders"] = totalOrders

	var totalRevenue float64
	r.db.WithContext(ctx).Model(&models.PrintOrder{}).
		Where("created_at BETWEEN ? AND ? AND status NOT IN ?", startDate, endDate,
			[]models.OrderStatus{models.OrderStatusCancelled, models.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)
	stats["total_revenue"] = totalRevenue

	statusStats := make(map[string]int64)
	type StatusCount struct {
		Status models.OrderStatus
		Count  int64
	}
	var results []StatusCount
	r.db.WithContext(ctx).Model(&models.PrintOrder{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results)
	for _, r := range results {
		statusStats[string(r.Status)] = r.Count
	}
	stats["status_stats"] = statusStats

	return stats, nil
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.PrintOrder{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *OrderRepository) ExecTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (r *OrderRepository) LockForUpdate(ctx context.Context, id uuid.UUID) (*models.PrintOrder, error) {
	var order models.PrintOrder
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
