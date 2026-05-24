package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/repository"
	"luxury-trading-platform/internal/utils"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	userRepo    *repository.UserRepository
	redisClient *cache.RedisClient
	db          *gorm.DB
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, userRepo *repository.UserRepository, redisClient *cache.RedisClient, db *gorm.DB) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		redisClient: redisClient,
		db:          db,
	}
}

type CreateOrderRequest struct {
	ProductID       uint   `json:"product_id" binding:"required"`
	ShippingAddress string `json:"shipping_address" binding:"required"`
	NeedAuth        bool   `json:"need_auth"`
	Remark          string `json:"remark"`
}

type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wechat card"`
}

type ShipOrderRequest struct {
	TrackingNumber string `json:"tracking_number" binding:"required"`
}

func (s *OrderService) CreateOrder(ctx context.Context, buyerID uint, req *CreateOrderRequest) (*model.Order, error) {
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if product.Status != model.ProductStatusOnSale {
		return nil, errors.New("product is not on sale")
	}

	if product.Stock <= 0 {
		return nil, errors.New("product is out of stock")
	}

	if product.SellerID == buyerID {
		return nil, errors.New("cannot buy your own product")
	}

	order := &model.Order{
		OrderNumber:     utils.GenerateOrderNumber(),
		BuyerID:         buyerID,
		SellerID:        product.SellerID,
		ProductID:       req.ProductID,
		Price:           product.Price,
		Status:          model.OrderStatusPending,
		PaymentStatus:   model.PaymentStatusPending,
		ShippingAddress: req.ShippingAddress,
		NeedAuth:        req.NeedAuth,
		Remark:          req.Remark,
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	stockKey := fmt.Sprintf("stock:product:%d", req.ProductID)
	if s.redisClient != nil {
		_, _ = s.redisClient.Decr(ctx, stockKey)
		_ = s.redisClient.Expire(ctx, stockKey, 24*time.Hour)
	}

	err = s.productRepo.DecrementStock(req.ProductID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	order, err := s.orderRepo.FindByOrderNumber(orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) ListOrders(page, pageSize int, buyerID, sellerID *uint, status model.OrderStatus) ([]model.Order, int64, error) {
	page, pageSize = utils.ValidatePage(page, pageSize)
	return s.orderRepo.List(page, pageSize, buyerID, sellerID, status)
}

func (s *OrderService) PayOrder(ctx context.Context, orderID uint, buyerID uint, req *PayOrderRequest) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.BuyerID != buyerID {
		return nil, errors.New("permission denied: this order does not belong to you")
	}

	if order.Status != model.OrderStatusPending {
		return nil, errors.New("order cannot be paid")
	}

	err = s.orderRepo.UpdatePaymentStatus(orderID, model.PaymentStatusSuccess, req.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	err = s.orderRepo.UpdateStatus(orderID, model.OrderStatusPaid)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return s.orderRepo.FindByID(orderID)
}

func (s *OrderService) ShipOrder(ctx context.Context, orderID uint, sellerID uint, req *ShipOrderRequest) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.SellerID != sellerID {
		return nil, errors.New("permission denied: this order does not belong to you")
	}

	if order.Status != model.OrderStatusPaid {
		return nil, errors.New("order cannot be shipped")
	}

	err = s.orderRepo.UpdateShipping(orderID, req.TrackingNumber, order.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to update shipping info: %w", err)
	}

	err = s.orderRepo.UpdateStatus(orderID, model.OrderStatusShipped)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return s.orderRepo.FindByID(orderID)
}

func (s *OrderService) ConfirmDelivery(ctx context.Context, orderID uint, buyerID uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.BuyerID != buyerID {
		return nil, errors.New("permission denied: this order does not belong to you")
	}

	if order.Status != model.OrderStatusShipped {
		return nil, errors.New("order cannot be confirmed")
	}

	err = s.orderRepo.UpdateStatus(orderID, model.OrderStatusDelivered)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	err = s.orderRepo.UpdateStatus(orderID, model.OrderStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	product, _ := s.productRepo.FindByID(order.ProductID)
	if product != nil {
		_ = s.productRepo.UpdateStatus(order.ProductID, model.ProductStatusSold)
	}

	return s.orderRepo.FindByID(orderID)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uint, userID uint, reason string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.BuyerID != userID && order.SellerID != userID {
		return nil, errors.New("permission denied")
	}

	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusPaid {
		return nil, errors.New("order cannot be cancelled")
	}

	err = s.orderRepo.UpdateStatus(orderID, model.OrderStatusCancelled)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	product, _ := s.productRepo.FindByID(order.ProductID)
	if product != nil {
		_ = s.productRepo.IncrementStock(order.ProductID, 1)

		stockKey := fmt.Sprintf("stock:product:%d", order.ProductID)
		if s.redisClient != nil {
			_, _ = s.redisClient.Incr(ctx, stockKey)
		}
	}

	return s.orderRepo.FindByID(orderID)
}

func (s *OrderService) ValidateOrderOwnership(order *model.Order, userID uint) bool {
	return order.BuyerID == userID || order.SellerID == userID
}
