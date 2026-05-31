package repository

import (
	"context"
	"errors"
	"time"

	"print3d-platform/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrinterRepository struct {
	db *gorm.DB
}

func NewPrinterRepository(db *gorm.DB) *PrinterRepository {
	return &PrinterRepository{db: db}
}

func (r *PrinterRepository) CreateDevice(ctx context.Context, device *models.PrinterDevice) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *PrinterRepository) GetDevice(ctx context.Context, id uuid.UUID) (*models.PrinterDevice, error) {
	var device models.PrinterDevice
	err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *PrinterRepository) GetPrinterDevices(ctx context.Context, printerID uuid.UUID) ([]models.PrinterDevice, error) {
	var devices []models.PrinterDevice
	err := r.db.WithContext(ctx).
		Where("printer_id = ?", printerID).
		Order("created_at DESC").
		Find(&devices).Error
	return devices, err
}

func (r *PrinterRepository) UpdateDevice(ctx context.Context, device *models.PrinterDevice) error {
	return r.db.WithContext(ctx).Save(device).Error
}

func (r *PrinterRepository) UpdateDeviceStatus(ctx context.Context, id uuid.UUID, status models.PrinterStatus) error {
	return r.db.WithContext(ctx).Model(&models.PrinterDevice{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *PrinterRepository) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.PrinterDevice{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("device not found")
	}
	return nil
}

func (r *PrinterRepository) GetIdleDevices(ctx context.Context, printerID uuid.UUID) ([]models.PrinterDevice, error) {
	var devices []models.PrinterDevice
	err := r.db.WithContext(ctx).
		Where("printer_id = ? AND status = ?", printerID, models.PrinterStatusIdle).
		Find(&devices).Error
	return devices, err
}

func (r *PrinterRepository) CreateInventory(ctx context.Context, inventory *models.MaterialInventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

func (r *PrinterRepository) GetPrinterInventory(ctx context.Context, printerID uuid.UUID) ([]models.MaterialInventory, error) {
	var inventory []models.MaterialInventory
	err := r.db.WithContext(ctx).
		Preload("Material").
		Where("printer_id = ?", printerID).
		Order("created_at DESC").
		Find(&inventory).Error
	return inventory, err
}

func (r *PrinterRepository) UpdateInventory(ctx context.Context, inventory *models.MaterialInventory) error {
	return r.db.WithContext(ctx).Save(inventory).Error
}

func (r *PrinterRepository) UpdateInventoryQuantity(ctx context.Context, id uuid.UUID, quantityGrams float64) error {
	return r.db.WithContext(ctx).Model(&models.MaterialInventory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"quantity_grams": quantityGrams,
			"last_updated":   time.Now(),
		}).Error
}

func (r *PrinterRepository) DeleteInventory(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.MaterialInventory{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("inventory item not found")
	}
	return nil
}

func (r *PrinterRepository) CreateSchedule(ctx context.Context, schedule *models.PrintSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *PrinterRepository) GetPrinterSchedules(ctx context.Context, printerID uuid.UUID, startDate, endDate *time.Time) ([]models.PrintSchedule, error) {
	var schedules []models.PrintSchedule
	query := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Device").
		Where("printer_id = ?", printerID)

	if startDate != nil {
		query = query.Where("scheduled_start >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("scheduled_end <= ?", *endDate)
	}

	err := query.Order("scheduled_start ASC").Find(&schedules).Error
	return schedules, err
}

func (r *PrinterRepository) UpdateSchedule(ctx context.Context, schedule *models.PrintSchedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

func (r *PrinterRepository) UpdateScheduleStatus(ctx context.Context, id uuid.UUID, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == "printing" {
		updates["actual_start"] = time.Now()
	} else if status == "completed" {
		updates["actual_end"] = time.Now()
	}
	return r.db.WithContext(ctx).Model(&models.PrintSchedule{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *PrinterRepository) CreateReview(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *PrinterRepository) GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Model").
		Preload("Printer").
		First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *PrinterRepository) GetReviewByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Model").
		Preload("Printer").
		First(&review, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *PrinterRepository) GetModelReviews(ctx context.Context, modelID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.WithContext(ctx).Preload("Customer").Where("model_id = ?", modelID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, total, err
}

func (r *PrinterRepository) GetPrinterReviews(ctx context.Context, printerID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.WithContext(ctx).Preload("Customer").Preload("Model").Where("printer_id = ?", printerID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, total, err
}

func (r *PrinterRepository) UpdateReview(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *PrinterRepository) ListMaterials(ctx context.Context) ([]models.Material, error) {
	var materials []models.Material
	err := r.db.WithContext(ctx).
		Preload("ColorOptions").
		Where("is_available = ?", true).
		Find(&materials).Error
	return materials, err
}

func (r *PrinterRepository) GetMaterial(ctx context.Context, id uuid.UUID) (*models.Material, error) {
	var material models.Material
	err := r.db.WithContext(ctx).
		Preload("ColorOptions").
		First(&material, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *PrinterRepository) GetMaterialStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	type MaterialUsage struct {
		MaterialID uuid.UUID
		MaterialName string
		TotalWeight float64
		OrderCount int64
	}

	var results []MaterialUsage
	r.db.WithContext(ctx).
		Table("print_orders").
		Select("material_id, materials.name as material_name, SUM(estimated_weight) as total_weight, COUNT(*) as order_count").
		Joins("JOIN materials ON materials.id = print_orders.material_id").
		Where("print_orders.created_at BETWEEN ? AND ? AND print_orders.status NOT IN ?", startDate, endDate,
			[]models.OrderStatus{models.OrderStatusCancelled, models.OrderStatusRefunded}).
		Group("material_id, materials.name").
		Scan(&results)
	stats["material_usage"] = results

	return stats, nil
}
