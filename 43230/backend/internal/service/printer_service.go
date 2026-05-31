package service

import (
	"context"
	"errors"
	"time"

	"print3d-platform/internal/models"
	"print3d-platform/internal/repository"
	"print3d-platform/internal/utils"

	"github.com/google/uuid"
)

type PrinterService struct {
	printerRepo *repository.PrinterRepository
	userRepo    *repository.UserRepository
	orderRepo   *repository.OrderRepository
	modelRepo   *repository.ModelRepository
}

func NewPrinterService(
	printerRepo *repository.PrinterRepository,
	userRepo *repository.UserRepository,
	orderRepo *repository.OrderRepository,
	modelRepo *repository.ModelRepository,
) *PrinterService {
	return &PrinterService{
		printerRepo: printerRepo,
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		modelRepo:   modelRepo,
	}
}

type CreateDeviceRequest struct {
	Name               string                `json:"name" binding:"required"`
	Model              string                `json:"model"`
	Manufacturer       string                `json:"manufacturer"`
	MaxPrintSize       string                `json:"max_print_size"`
	MaxPrintVolume     float64               `json:"max_print_volume"`
	SupportedMaterials []string              `json:"supported_materials"`
	SupportedQualities []models.PrintQuality `json:"supported_qualities"`
	IPAddress          string                `json:"ip_address"`
}

type UpdateDeviceRequest struct {
	Name               string                `json:"name"`
	Status             models.PrinterStatus  `json:"status"`
	IPAddress          string                `json:"ip_address"`
	FirmwareVersion    string                `json:"firmware_version"`
}

type CreateInventoryRequest struct {
	MaterialID   uuid.UUID `json:"material_id" binding:"required"`
	Color        string    `json:"color" binding:"required"`
	QuantityGrams float64  `json:"quantity_grams" binding:"required,min=0"`
	ReorderLevel float64   `json:"reorder_level"`
}

type CreateReviewRequest struct {
	OrderID      uuid.UUID `json:"order_id" binding:"required"`
	ModelRating  int       `json:"model_rating" binding:"required,min=1,max=5"`
	PrintRating  int       `json:"print_rating" binding:"required,min=1,max=5"`
	ModelComment string    `json:"model_comment"`
	PrintComment string    `json:"print_comment"`
	Images       []string  `json:"images"`
	IsAnonymous  bool      `json:"is_anonymous"`
}

func (s *PrinterService) CreateDevice(ctx context.Context, printerID uuid.UUID, req *CreateDeviceRequest) (*models.PrinterDevice, error) {
	now := time.Now()
	device := &models.PrinterDevice{
		ID:                uuid.New(),
		PrinterID:         printerID,
		Name:              req.Name,
		Model:             req.Model,
		Manufacturer:      req.Manufacturer,
		MaxPrintSize:      req.MaxPrintSize,
		MaxPrintVolume:    req.MaxPrintVolume,
		SupportedMaterials: req.SupportedMaterials,
		SupportedQualities: req.SupportedQualities,
		Status:            models.PrinterStatusIdle,
		IPAddress:         req.IPAddress,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err := s.printerRepo.CreateDevice(ctx, device)
	if err != nil {
		return nil, err
	}

	utils.LogInfo("New device created: %s by printer %s", device.Name, printerID)
	return device, nil
}

func (s *PrinterService) GetDevice(ctx context.Context, id uuid.UUID) (*models.PrinterDevice, error) {
	return s.printerRepo.GetDevice(ctx, id)
}

func (s *PrinterService) GetPrinterDevices(ctx context.Context, printerID uuid.UUID) ([]models.PrinterDevice, error) {
	return s.printerRepo.GetPrinterDevices(ctx, printerID)
}

func (s *PrinterService) UpdateDevice(ctx context.Context, id, printerID uuid.UUID, req *UpdateDeviceRequest) (*models.PrinterDevice, error) {
	device, err := s.printerRepo.GetDevice(ctx, id)
	if err != nil {
		return nil, err
	}

	if device.PrinterID != printerID {
		return nil, errors.New("not authorized")
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Status != "" {
		device.Status = req.Status
	}
	if req.IPAddress != "" {
		device.IPAddress = req.IPAddress
	}
	if req.FirmwareVersion != "" {
		device.FirmwareVersion = req.FirmwareVersion
	}

	device.UpdatedAt = time.Now()
	err = s.printerRepo.UpdateDevice(ctx, device)
	return device, err
}

func (s *PrinterService) DeleteDevice(ctx context.Context, id, printerID uuid.UUID) error {
	device, err := s.printerRepo.GetDevice(ctx, id)
	if err != nil {
		return err
	}

	if device.PrinterID != printerID {
		return errors.New("not authorized")
	}

	return s.printerRepo.DeleteDevice(ctx, id)
}

func (s *PrinterService) CreateInventory(ctx context.Context, printerID uuid.UUID, req *CreateInventoryRequest) (*models.MaterialInventory, error) {
	now := time.Now()
	inventory := &models.MaterialInventory{
		ID:             uuid.New(),
		PrinterID:      printerID,
		MaterialID:     req.MaterialID,
		Color:          req.Color,
		QuantityGrams:  req.QuantityGrams,
		ReorderLevel:   req.ReorderLevel,
		LastUpdated:    now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := s.printerRepo.CreateInventory(ctx, inventory)
	if err != nil {
		return nil, err
	}

	utils.LogInfo("Inventory created for material %s by printer %s", req.MaterialID, printerID)
	return inventory, nil
}

func (s *PrinterService) GetPrinterInventory(ctx context.Context, printerID uuid.UUID) ([]models.MaterialInventory, error) {
	return s.printerRepo.GetPrinterInventory(ctx, printerID)
}

func (s *PrinterService) UpdateInventoryQuantity(ctx context.Context, id, printerID uuid.UUID, quantityGrams float64) error {
	inventory, err := s.printerRepo.GetPrinterInventory(ctx, printerID)
	if err != nil {
		return err
	}

	for _, inv := range inventory {
		if inv.ID == id {
			return s.printerRepo.UpdateInventoryQuantity(ctx, id, quantityGrams)
		}
	}

	return errors.New("inventory not found")
}

func (s *PrinterService) DeleteInventory(ctx context.Context, id, printerID uuid.UUID) error {
	inventory, err := s.printerRepo.GetPrinterInventory(ctx, printerID)
	if err != nil {
		return err
	}

	for _, inv := range inventory {
		if inv.ID == id {
			return s.printerRepo.DeleteInventory(ctx, id)
		}
	}

	return errors.New("inventory not found")
}

func (s *PrinterService) CreateSchedule(ctx context.Context, printerID uuid.UUID, deviceID, orderID uuid.UUID, scheduledStart, scheduledEnd *time.Time, priority int) (*models.PrintSchedule, error) {
	now := time.Now()
	schedule := &models.PrintSchedule{
		ID:             uuid.New(),
		PrinterID:      printerID,
		DeviceID:       deviceID,
		OrderID:        orderID,
		ScheduledStart: scheduledStart,
		ScheduledEnd:   scheduledEnd,
		Status:         "scheduled",
		Priority:       priority,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := s.printerRepo.CreateSchedule(ctx, schedule)
	if err != nil {
		return nil, err
	}

	utils.LogInfo("Schedule created for order %s on device %s", orderID, deviceID)
	return schedule, nil
}

func (s *PrinterService) GetPrinterSchedules(ctx context.Context, printerID uuid.UUID, startDate, endDate *time.Time) ([]models.PrintSchedule, error) {
	return s.printerRepo.GetPrinterSchedules(ctx, printerID, startDate, endDate)
}

func (s *PrinterService) CreateReview(ctx context.Context, customerID uuid.UUID, req *CreateReviewRequest) (*models.Review, error) {
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.CustomerID != customerID {
		return nil, errors.New("not authorized to review this order")
	}

	if order.Status != models.OrderStatusCompleted {
		return nil, errors.New("can only review completed orders")
	}

	existingReview, _ := s.printerRepo.GetReviewByOrderID(ctx, req.OrderID)
	if existingReview != nil {
		return nil, errors.New("review already exists for this order")
	}

	now := time.Now()
	review := &models.Review{
		ID:           uuid.New(),
		OrderID:      req.OrderID,
		CustomerID:   customerID,
		ModelID:      order.ModelID,
		PrinterID:    order.PrinterID,
		ModelRating:  req.ModelRating,
		PrintRating:  req.PrintRating,
		ModelComment: req.ModelComment,
		PrintComment: req.PrintComment,
		Images:       req.Images,
		IsAnonymous:  req.IsAnonymous,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.printerRepo.CreateReview(ctx, review)
	if err != nil {
		return nil, err
	}

	avgModelRating := float64(req.ModelRating)
	avgPrintRating := float64(req.PrintRating)

	_ = s.modelRepo.UpdateRating(ctx, order.ModelID, avgModelRating)
	_ = s.userRepo.UpdateDesignerRating(ctx, order.ModelID, avgModelRating)
	if order.PrinterID != uuid.Nil {
		_ = s.userRepo.UpdatePrinterRating(ctx, order.PrinterID, avgPrintRating)
	}

	utils.LogInfo("Review created for order %s: model=%d, print=%d", req.OrderID, req.ModelRating, req.PrintRating)
	return review, nil
}

func (s *PrinterService) GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	return s.printerRepo.GetReview(ctx, id)
}

func (s *PrinterService) GetModelReviews(ctx context.Context, modelID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	return s.printerRepo.GetModelReviews(ctx, modelID, page, pageSize)
}

func (s *PrinterService) GetPrinterReviews(ctx context.Context, printerID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	return s.printerRepo.GetPrinterReviews(ctx, printerID, page, pageSize)
}

func (s *PrinterService) GetMaterials(ctx context.Context) ([]models.Material, error) {
	return s.printerRepo.ListMaterials(ctx)
}

func (s *PrinterService) GetMaterial(ctx context.Context, id uuid.UUID) (*models.Material, error) {
	return s.printerRepo.GetMaterial(ctx, id)
}

func (s *PrinterService) GetMaterialStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.printerRepo.GetMaterialStats(ctx, startDate, endDate)
}

func (s *PrinterService) GetIdleDevices(ctx context.Context, printerID uuid.UUID) ([]models.PrinterDevice, error) {
	return s.printerRepo.GetIdleDevices(ctx, printerID)
}
