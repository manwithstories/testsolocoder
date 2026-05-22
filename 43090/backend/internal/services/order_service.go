package services

import (
	"errors"
	"fmt"
	"time"

	"auction-system/internal/dto"
	"auction-system/internal/models"
	"auction-system/internal/utils"
	"auction-system/pkg/logger"
)

type OrderService struct {
	notifyService *NotificationService
}

func NewOrderService() *OrderService {
	return &OrderService{
		notifyService: NewNotificationService(),
	}
}

const (
	OrderStatusPendingPayment = 0
	OrderStatusPaid           = 1
	OrderStatusShipped        = 2
	OrderStatusDelivered      = 3
	OrderStatusCompleted      = 4
	OrderStatusCancelled      = 5
)

const (
	PaymentStatusPending = 0
	PaymentStatusSuccess = 1
	PaymentStatusFailed  = 2
)

func (s *OrderService) CreateOrder(buyerID uint, req *dto.CreateOrderRequest) (*models.Order, error) {
	item, err := NewAuctionItemService().GetItemByID(req.AuctionItemID)
	if err != nil {
		return nil, errors.New("拍卖品不存在")
	}

	var existingBid models.Bid
	if err := models.DB.Where("auction_item_id = ? AND user_id = ? AND is_winning = 1", req.AuctionItemID, buyerID).First(&existingBid).Error; err != nil {
		return nil, errors.New("您不是该拍卖品的最高出价者")
	}

	var existingOrder models.Order
	if models.DB.Where("auction_item_id = ?", req.AuctionItemID).First(&existingOrder).Error == nil {
		return nil, errors.New("该拍卖品已生成订单")
	}

	orderNo := utils.GenerateOrderNo()

	order := &models.Order{
		OrderNo:       orderNo,
		AuctionItemID: req.AuctionItemID,
		BuyerID:       buyerID,
		SellerID:      item.SellerID,
		Price:         item.CurrentPrice,
		Status:        OrderStatusPendingPayment,
		ShippingInfo:  req.ShippingInfo,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tx := models.BeginTransaction()

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to create order: %v", err)
		return nil, errors.New("创建订单失败")
	}

	if err := tx.Model(&models.AuctionItem{}).Where("id = ?", req.AuctionItemID).Update("status", ItemStatusSold).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("更新拍卖品状态失败")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("提交订单失败")
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(orderID uint, userID uint) (*models.Order, error) {
	var order models.Order
	if err := models.DB.Preload("AuctionItem.Images").Preload("Buyer").Preload("Seller").Preload("Payments").First(&order, orderID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.BuyerID != userID && order.SellerID != userID {
		return nil, errors.New("无权查看此订单")
	}

	return &order, nil
}

func (s *OrderService) GetOrderByNo(orderNo string, userID uint) (*models.Order, error) {
	var order models.Order
	if err := models.DB.Preload("AuctionItem.Images").Preload("Buyer").Preload("Seller").Preload("Payments").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.BuyerID != userID && order.SellerID != userID {
		return nil, errors.New("无权查看此订单")
	}

	return &order, nil
}

func (s *OrderService) GetBuyerOrders(buyerID uint, page, pageSize int, status *int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := models.DB.Model(&models.Order{}).Where("buyer_id = ?", buyerID).Preload("AuctionItem.Images")
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (s *OrderService) GetSellerOrders(sellerID uint, page, pageSize int, status *int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := models.DB.Model(&models.Order{}).Where("seller_id = ?", sellerID).Preload("AuctionItem.Images")
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (s *OrderService) PayOrder(orderID uint, buyerID uint, req *dto.PayOrderRequest) (*models.Payment, error) {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.BuyerID != buyerID {
		return nil, errors.New("无权支付此订单")
	}

	if order.Status != OrderStatusPendingPayment {
		return nil, errors.New("订单状态不允许支付")
	}

	user, err := NewUserService().GetUserByID(buyerID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Balance < order.Price {
		return nil, errors.New("余额不足")
	}

	paymentNo := fmt.Sprintf("PAY%s", utils.GenerateOrderNo())

	tx := models.BeginTransaction()

	if err := tx.Model(&models.User{}).Where("id = ?", buyerID).Update("balance", user.Balance-order.Price).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("扣款失败")
	}

	payment := &models.Payment{
		OrderID:       orderID,
		PaymentNo:     paymentNo,
		Amount:        order.Price,
		Method:        req.Method,
		Status:        PaymentStatusSuccess,
		TransactionID: utils.GenerateRandomString(32),
		CreatedAt:     time.Now(),
		PaidAt:        &time.Time{},
	}
	*payment.PaidAt = time.Now()

	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("创建支付记录失败")
	}

	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":       OrderStatusPaid,
		"payment_time": now,
		"updated_at":   now,
	}).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("更新订单状态失败")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("支付失败")
	}

	s.notifyService.NotifyPaymentSuccess(buyerID, order.AuctionItemID, "", order.OrderNo)

	return payment, nil
}

func (s *OrderService) ShipOrder(orderID uint, sellerID uint, trackingNo string) error {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.SellerID != sellerID {
		return errors.New("无权操作此订单")
	}

	if order.Status != OrderStatusPaid {
		return errors.New("订单状态不允许发货")
	}

	return models.DB.Model(&order).Updates(map[string]interface{}{
		"status":      OrderStatusShipped,
		"tracking_no": trackingNo,
		"updated_at":  time.Now(),
	}).Error
}

func (s *OrderService) ConfirmDelivery(orderID uint, buyerID uint) error {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != buyerID {
		return errors.New("无权操作此订单")
	}

	if order.Status != OrderStatusShipped {
		return errors.New("订单状态不允许确认收货")
	}

	return models.DB.Model(&order).Updates(map[string]interface{}{
		"status":     OrderStatusDelivered,
		"updated_at": time.Now(),
	}).Error
}

func (s *OrderService) CompleteOrder(orderID uint) error {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.Status != OrderStatusDelivered {
		return errors.New("订单状态不允许完成")
	}

	tx := models.BeginTransaction()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":     OrderStatusCompleted,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	seller, err := NewUserService().GetUserByID(order.SellerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.User{}).Where("id = ?", order.SellerID).Update("balance", seller.Balance+order.Price).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *OrderService) GetAllOrders(page, pageSize int, status *int, keyword string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := models.DB.Model(&models.Order{}).Preload("AuctionItem").Preload("Buyer").Preload("Seller")
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		query = query.Joins("JOIN auction_items ON auction_items.id = orders.auction_item_id").
			Where("orders.order_no LIKE ? OR auction_items.title LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}
