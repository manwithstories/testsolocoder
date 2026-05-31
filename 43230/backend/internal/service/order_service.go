package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"print3d-platform/internal/config"
	"print3d-platform/internal/models"
	"print3d-platform/internal/repository"
	"print3d-platform/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo  *repository.OrderRepository
	modelRepo  *repository.ModelRepository
	userRepo   *repository.UserRepository
	printerRepo *repository.PrinterRepository
	cfg        *config.Config
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	modelRepo *repository.ModelRepository,
	userRepo *repository.UserRepository,
	printerRepo *repository.PrinterRepository,
	cfg *config.Config,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		modelRepo:   modelRepo,
		userRepo:    userRepo,
		printerRepo: printerRepo,
		cfg:         cfg,
	}
}

type CreateOrderRequest struct {
	ModelID       uuid.UUID         `json:"model_id" binding:"required"`
	Quantity      int               `json:"quantity" binding:"required,min=1,max=100"`
	MaterialID    uuid.UUID         `json:"material_id" binding:"required"`
	Color         string            `json:"color" binding:"required"`
	Quality       models.PrintQuality `json:"quality" binding:"required"`
	LayerHeight   float64           `json:"layer_height"`
	InfillPercent int               `json:"infill_percent"`
	Supports      bool              `json:"supports"`
	ShippingAddress string          `json:"shipping_address" binding:"required"`
	Notes         string            `json:"notes"`
}

type PriceEstimateRequest struct {
	ModelID       uuid.UUID         `json:"model_id" binding:"required"`
	Quantity      int               `json:"quantity" binding:"required"`
	MaterialID    uuid.UUID         `json:"material_id" binding:"required"`
	Color         string            `json:"color"`
	Quality       models.PrintQuality `json:"quality" binding:"required"`
	LayerHeight   float64           `json:"layer_height"`
	InfillPercent int               `json:"infill_percent"`
	Supports      bool              `json:"supports"`
}

type PriceEstimateResponse struct {
	BasePrice       float64 `json:"base_price"`
	MaterialCost    float64 `json:"material_cost"`
	ServiceFee      float64 `json:"service_fee"`
	ShippingFee     float64 `json:"shipping_fee"`
	TotalAmount     float64 `json:"total_amount"`
	EstimatedVolume float64 `json:"estimated_volume"`
	EstimatedWeight float64 `json:"estimated_weight"`
	EstimatedPrintTime float64 `json:"estimated_print_time"`
}

func (s *OrderService) EstimatePrice(ctx context.Context, req *PriceEstimateRequest) (*PriceEstimateResponse, error) {
	model, err := s.modelRepo.FindByID(ctx, req.ModelID)
	if err != nil {
		return nil, errors.New("model not found")
	}

	material, err := s.printerRepo.GetMaterial(ctx, req.MaterialID)
	if err != nil {
		return nil, errors.New("material not found")
	}

	qualityMultipliers := map[models.PrintQuality]float64{
		models.QualityDraft:    0.7,
		models.QualityStandard: 1.0,
		models.QualityHigh:    1.4,
		models.QualityUltra:   1.8,
	}

	infillMultiplier := 1.0 + float64(req.InfillPercent)/100.0

	qualityMult := qualityMultipliers[req.Quality]
	if qualityMult == 0 {
		qualityMult = 1.0
	}

	estimatedVolume := model.Volume
	if estimatedVolume <= 0 {
		estimatedVolume = 100.0
	}

	estimatedWeight := estimatedVolume * material.Density * 0.001
	estimatedPrintTime := utils.EstimatePrintTime(estimatedVolume, string(req.Quality))

	basePrice := math.Round(model.Price*100) / 100

	materialCost := math.Round(
		estimatedWeight*float64(req.Quantity)*material.PricePerGram*
			qualityMult*infillMultiplier*100,
	) / 100

	serviceFee := math.Round(
		(basePrice*0.1+estimatedPrintTime*s.cfg.Pricing.TimePriceMultiplier*10)*
			float64(req.Quantity)*100,
	) / 100

	var supportsFee float64
	if req.Supports {
		supportsFee = math.Round(estimatedVolume*0.01*float64(req.Quantity)*100) / 100
	}

	shippingFee := math.Round((5.0+estimatedWeight*0.01)*float64(req.Quantity)*100) / 100

	totalAmount := math.Round((basePrice+materialCost+serviceFee+supportsFee+shippingFee)*100) / 100

	return &PriceEstimateResponse{
		BasePrice:         basePrice,
		MaterialCost:      materialCost + supportsFee,
		ServiceFee:        serviceFee,
		ShippingFee:       shippingFee,
		TotalAmount:       totalAmount,
		EstimatedVolume:   estimatedVolume,
		EstimatedWeight:   estimatedWeight * float64(req.Quantity),
		EstimatedPrintTime: estimatedPrintTime * float64(req.Quantity),
	}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID uuid.UUID, req *CreateOrderRequest) (*models.PrintOrder, error) {
	model, err := s.modelRepo.FindByID(ctx, req.ModelID)
	if err != nil {
		return nil, errors.New("model not found")
	}

	material, err := s.printerRepo.GetMaterial(ctx, req.MaterialID)
	if err != nil {
		return nil, errors.New("material not found")
	}

	customer, err := s.userRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	_ = model
	_ = material

	estimate, err := s.EstimatePrice(ctx, &PriceEstimateRequest{
		ModelID:       req.ModelID,
		Quantity:      req.Quantity,
		MaterialID:    req.MaterialID,
		Quality:       req.Quality,
		LayerHeight:   req.LayerHeight,
		InfillPercent: req.InfillPercent,
		Supports:      req.Supports,
	})
	if err != nil {
		return nil, err
	}

	if customer.Balance < estimate.TotalAmount {
		return nil, errors.New("insufficient balance")
	}

	now := time.Now()
	order := &models.PrintOrder{
		ID:                uuid.New(),
		OrderNo:           utils.GenerateOrderNo(),
		CustomerID:        customerID,
		ModelID:           req.ModelID,
		Quantity:          req.Quantity,
		MaterialID:        req.MaterialID,
		Color:             req.Color,
		Quality:           req.Quality,
		LayerHeight:       req.LayerHeight,
		InfillPercent:     req.InfillPercent,
		Supports:          req.Supports,
		EstimatedVolume:   estimate.EstimatedVolume,
		EstimatedWeight:   estimate.EstimatedWeight,
		EstimatedPrintTime: estimate.EstimatedPrintTime,
		BasePrice:         estimate.BasePrice,
		MaterialCost:      estimate.MaterialCost,
		ServiceFee:        estimate.ServiceFee,
		ShippingFee:       estimate.ShippingFee,
		TotalAmount:       estimate.TotalAmount,
		Status:            models.OrderStatusPending,
		ShippingAddress:   req.ShippingAddress,
		Notes:             req.Notes,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err = s.orderRepo.ExecTx(ctx, func(tx *gorm.DB) error {
		err = tx.WithContext(ctx).Create(order).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(&models.User{}).
			Where("id = ?", customerID).
			Update("balance", gorm.Expr("balance - ?", estimate.TotalAmount)).Error
		if err != nil {
			return err
		}

		customerTx := &models.Transaction{
			ID:            uuid.New(),
			UserID:        customerID,
			OrderID:       order.ID,
			Type:          "expense",
			Amount:        estimate.TotalAmount,
			BalanceAfter:  customer.Balance - estimate.TotalAmount,
			Description:   fmt.Sprintf("Payment for order %s", order.OrderNo),
			TransactionNo: utils.GenerateTransactionNo(),
			Status:        "completed",
			CreatedAt:     now,
		}
		err = tx.WithContext(ctx).Create(customerTx).Error
		if err != nil {
			return err
		}

		history := &models.OrderHistory{
			ID:          uuid.New(),
			OrderID:     order.ID,
			Status:      models.OrderStatusPending,
			Description: "Order created and pending payment verification",
			CreatedAt:   now,
		}
		err = tx.WithContext(ctx).Create(history).Error

		return err
	})

	if err != nil {
		return nil, err
	}

	utils.LogInfo("Order created: %s by customer %s, amount: %.2f", order.OrderNo, customerID, order.TotalAmount)
	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.PrintOrder, error) {
	return s.orderRepo.FindByID(ctx, id)
}

func (s *OrderService) GetOrderByNo(ctx context.Context, orderNo string) (*models.PrintOrder, error) {
	return s.orderRepo.FindByOrderNo(ctx, orderNo)
}

func (s *OrderService) ListCustomerOrders(ctx context.Context, customerID uuid.UUID, page, pageSize int, status *string) ([]models.PrintOrder, int64, error) {
	return s.orderRepo.ListCustomerOrders(ctx, customerID, page, pageSize, status)
}

func (s *OrderService) ListPrinterOrders(ctx context.Context, printerID uuid.UUID, page, pageSize int, status *string) ([]models.PrintOrder, int64, error) {
	return s.orderRepo.ListPrinterOrders(ctx, printerID, page, pageSize, status)
}

func (s *OrderService) AssignPrinter(ctx context.Context, orderID, printerID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusPaid && order.Status != models.OrderStatusPending {
		return errors.New("order cannot be assigned at this status")
	}

	err = s.orderRepo.AssignPrinter(ctx, orderID, printerID)
	if err != nil {
		return err
	}

	return s.addOrderHistory(ctx, orderID, models.OrderStatusPaid,
		fmt.Sprintf("Order assigned to printer %s", printerID))
}

func (s *OrderService) StartPrinting(ctx context.Context, orderID, printerID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.PrinterID != printerID {
		return errors.New("not authorized to print this order")
	}

	if order.Status != models.OrderStatusPaid {
		return errors.New("order cannot be printed at this status")
	}

	err = s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusPrinting)
	if err != nil {
		return err
	}

	return s.addOrderHistory(ctx, orderID, models.OrderStatusPrinting, "Printing started")
}

func (s *OrderService) CompletePrinting(ctx context.Context, orderID, printerID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.PrinterID != printerID {
		return errors.New("not authorized")
	}

	if order.Status != models.OrderStatusPrinting {
		return errors.New("order is not being printed")
	}

	err = s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusQualityCheck)
	if err != nil {
		return err
	}

	return s.addOrderHistory(ctx, orderID, models.OrderStatusQualityCheck, "Printing completed, quality check in progress")
}

func (s *OrderService) ApproveQuality(ctx context.Context, orderID, printerID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.PrinterID != printerID {
		return errors.New("not authorized")
	}

	if order.Status != models.OrderStatusQualityCheck {
		return errors.New("order is not in quality check")
	}

	err = s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusShipped)
	if err != nil {
		return err
	}

	return s.addOrderHistory(ctx, orderID, models.OrderStatusShipped, "Quality check passed, order shipped")
}

func (s *OrderService) ShipOrder(ctx context.Context, orderID uuid.UUID, trackingNumber string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusShipped {
		return errors.New("order is not ready for shipping")
	}

	err = s.orderRepo.UpdateShipping(ctx, orderID, trackingNumber, order.ShippingAddress)
	if err != nil {
		return err
	}

	return s.addOrderHistory(ctx, orderID, models.OrderStatusShipped,
		fmt.Sprintf("Order shipped, tracking: %s", trackingNumber))
}

func (s *OrderService) DeliverOrder(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusShipped {
		return errors.New("order is not shipped")
	}

	err = s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusDelivered)
	if err != nil {
		return err
	}

	return s.addOrderHistory(ctx, orderID, models.OrderStatusDelivered, "Order delivered")
}

func (s *OrderService) CompleteOrder(ctx context.Context, orderID, customerID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.CustomerID != customerID {
		return errors.New("not authorized")
	}

	if order.Status != models.OrderStatusDelivered {
		return errors.New("order is not delivered")
	}

	err = s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusCompleted)
	if err != nil {
		return err
	}

	err = s.addOrderHistory(ctx, orderID, models.OrderStatusCompleted, "Order completed by customer")
	if err != nil {
		return err
	}

	return s.ProcessSettlement(ctx, orderID)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID uuid.UUID, reason string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.CustomerID != userID && order.PrinterID != userID {
		return errors.New("not authorized")
	}

	if order.Status == models.OrderStatusCompleted || order.Status == models.OrderStatusCancelled {
		return errors.New("order cannot be cancelled")
	}

	err = s.orderRepo.ExecTx(ctx, func(tx *gorm.DB) error {
		err = tx.WithContext(ctx).Model(&models.PrintOrder{}).
			Where("id = ?", orderID).
			Updates(map[string]interface{}{
				"status":           models.OrderStatusCancelled,
				"cancelled_reason": reason,
				"cancelled_at":     time.Now(),
			}).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(&models.User{}).
			Where("id = ?", order.CustomerID).
			Update("balance", gorm.Expr("balance + ?", order.TotalAmount)).Error
		if err != nil {
			return err
		}

		refundTx := &models.Transaction{
			ID:            uuid.New(),
			UserID:        order.CustomerID,
			OrderID:       orderID,
			Type:          "refund",
			Amount:        order.TotalAmount,
			Description:   fmt.Sprintf("Refund for cancelled order %s", order.OrderNo),
			TransactionNo: utils.GenerateTransactionNo(),
			Status:        "completed",
			CreatedAt:     time.Now(),
		}
		err = tx.WithContext(ctx).Create(refundTx).Error
		if err != nil {
			return err
		}

		history := &models.OrderHistory{
			ID:          uuid.New(),
			OrderID:     orderID,
			Status:      models.OrderStatusCancelled,
			Description: fmt.Sprintf("Order cancelled: %s", reason),
			CreatedAt:   time.Now(),
		}
		return tx.WithContext(ctx).Create(history).Error
	})

	if err != nil {
		return err
	}

	utils.LogInfo("Order %s cancelled: %s", order.OrderNo, reason)
	return nil
}

func (s *OrderService) GetOrderHistory(ctx context.Context, orderID uuid.UUID) ([]models.OrderHistory, error) {
	return s.orderRepo.GetHistory(ctx, orderID)
}

func (s *OrderService) ProcessSettlement(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	existingSettlement, _ := s.orderRepo.GetSettlementByOrderID(ctx, orderID)
	if existingSettlement != nil {
		return nil
	}

	model, err := s.modelRepo.FindByID(ctx, order.ModelID)
	if err != nil {
		return err
	}

	totalAmount := order.TotalAmount
	platformFee := math.Round(totalAmount*s.cfg.Pricing.PlatformFeeRate*100) / 100
	designerShare := math.Round((totalAmount-platformFee)*s.cfg.Pricing.DesignerFeeRate*100) / 100
	printerShare := math.Round((totalAmount-platformFee)*s.cfg.Pricing.PrinterFeeRate*100) / 100

	err = s.orderRepo.ExecTx(ctx, func(tx *gorm.DB) error {
		settlement := &models.Settlement{
			ID:            uuid.New(),
			OrderID:       orderID,
			TotalAmount:   totalAmount,
			PlatformFee:   platformFee,
			DesignerShare: designerShare,
			PrinterShare:  printerShare,
			Status:        "completed",
			SettledAt:     ptr(time.Now()),
			CreatedAt:     time.Now(),
		}
		err = tx.WithContext(ctx).Create(settlement).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(&models.User{}).
			Where("id = ?", model.DesignerID).
			Update("balance", gorm.Expr("balance + ?", designerShare)).Error
		if err != nil {
			return err
		}

		designerTx := &models.Transaction{
			ID:            uuid.New(),
			UserID:        model.DesignerID,
			OrderID:       orderID,
			Type:          "income",
			Amount:        designerShare,
			Description:   fmt.Sprintf("Designer share for order %s", order.OrderNo),
			TransactionNo: utils.GenerateTransactionNo(),
			Status:        "completed",
			CreatedAt:     time.Now(),
		}
		err = tx.WithContext(ctx).Create(designerTx).Error
		if err != nil {
			return err
		}

		if order.PrinterID != uuid.Nil {
			err = tx.WithContext(ctx).Model(&models.User{}).
				Where("id = ?", order.PrinterID).
				Update("balance", gorm.Expr("balance + ?", printerShare)).Error
			if err != nil {
				return err
			}

			printerTx := &models.Transaction{
				ID:            uuid.New(),
				UserID:        order.PrinterID,
				OrderID:       orderID,
				Type:          "income",
				Amount:        printerShare,
				Description:   fmt.Sprintf("Printer share for order %s", order.OrderNo),
				TransactionNo: utils.GenerateTransactionNo(),
				Status:        "completed",
				CreatedAt:     time.Now(),
			}
			err = tx.WithContext(ctx).Create(printerTx).Error
		}

		return err
	})

	if err != nil {
		return err
	}

	utils.LogInfo("Settlement processed for order %s: platform=%.2f, designer=%.2f, printer=%.2f",
		order.OrderNo, platformFee, designerShare, printerShare)
	return nil
}

func (s *OrderService) GetOrderStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.orderRepo.GetOrderStats(ctx, startDate, endDate)
}

func (s *OrderService) GetPendingOrders(ctx context.Context) ([]models.PrintOrder, error) {
	return s.orderRepo.GetPendingOrders(ctx)
}

func (s *OrderService) addOrderHistory(ctx context.Context, orderID uuid.UUID, status models.OrderStatus, description string) error {
	history := &models.OrderHistory{
		ID:          uuid.New(),
		OrderID:     orderID,
		Status:      status,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return s.orderRepo.AddHistory(ctx, history)
}

func ptr(t time.Time) *time.Time {
	return &t
}
