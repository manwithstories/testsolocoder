package services

import (
	"errors"
	"time"

	"secondhand-platform/database"
	"secondhand-platform/models"
	"secondhand-platform/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(buyerID, productID uint, receiverName, receiverPhone, receiverAddress string, negotiatedPrice float64) (*models.Order, error) {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return nil, errors.New("商品不存在")
	}

	if product.Status != models.ProductStatusOnSale {
		return nil, errors.New("商品不在售卖状态")
	}

	if product.SellerID == buyerID {
		return nil, errors.New("不能购买自己的商品")
	}

	finalPrice := product.Price
	negotiated := false
	if negotiatedPrice > 0 && negotiatedPrice < product.Price {
		finalPrice = negotiatedPrice
		negotiated = true
	}

	warrantyUntil := time.Now().AddDate(0, 0, product.WarrantyDays)

	order := &models.Order{
		OrderNo:         utils.GenerateOrderNo(),
		BuyerID:         buyerID,
		SellerID:        product.SellerID,
		ProductID:       productID,
		ProductTitle:    product.Title,
		ProductImage:    product.Images,
		OriginalPrice:   product.Price,
		FinalPrice:      finalPrice,
		Negotiated:      negotiated,
		Status:          models.OrderStatusPending,
		ReceiverName:    receiverName,
		ReceiverPhone:   receiverPhone,
		ReceiverAddress: receiverAddress,
		WarrantyDays:    product.WarrantyDays,
		WarrantyUntil:   &warrantyUntil,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		if negotiated {
			negotiation := &models.Negotiation{
				OrderID:      order.ID,
				BuyerID:      buyerID,
				SellerID:     product.SellerID,
				OfferedPrice: negotiatedPrice,
				Status:       models.NegotiationStatusPending,
			}
			return tx.Create(negotiation).Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := database.DB.Preload("Buyer").Preload("Seller").Preload("Product").
		Preload("Negotiations").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetOrderByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	if err := database.DB.Preload("Buyer").Preload("Seller").Preload("Product").
		Preload("Negotiations").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) ListOrders(page, pageSize int, userID uint, userRole string, status int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	db := database.DB.Model(&models.Order{})

	if userRole == models.RoleBuyer {
		db = db.Where("buyer_id = ?", userID)
	} else if userRole == models.RoleSeller {
		db = db.Where("seller_id = ?", userID)
	}

	if status > 0 {
		db = db.Where("status = ?", status)
	}

	db.Count(&total)
	if err := db.Preload("Buyer").Preload("Seller").Preload("Product").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *OrderService) PayOrder(userID uint, orderNo string, paymentMethod string) error {
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("订单状态不允许支付")
	}

	commission := order.FinalPrice * 0.05

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var buyer models.User
		if err := tx.First(&buyer, userID).Error; err != nil {
			return err
		}

		if paymentMethod == "wallet" {
			if buyer.WalletBalance < order.FinalPrice {
				return errors.New("钱包余额不足")
			}
			if err := tx.Model(&buyer).Update("wallet_balance", buyer.WalletBalance-order.FinalPrice).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		order.Status = models.OrderStatusPaid
		order.PaymentMethod = paymentMethod
		order.PaidAt = &now
		order.Commission = commission
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		if paymentMethod == "wallet" {
			buyerLog := models.WalletLog{
				UserID:      userID,
				Type:        models.WalletTypePayment,
				Amount:      -order.FinalPrice,
				Balance:     buyer.WalletBalance - order.FinalPrice,
				OrderNo:     order.OrderNo,
				Description: "购买商品支付",
			}
			if err := tx.Create(&buyerLog).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *OrderService) ShipOrder(sellerID uint, orderNo, trackingNo, trackingCompany string) error {
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.SellerID != sellerID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.OrderStatusPaid {
		return errors.New("订单状态不允许发货")
	}

	now := time.Now()
	order.Status = models.OrderStatusShipped
	order.ShippedAt = &now
	order.TrackingNo = trackingNo
	order.TrackingCompany = trackingCompany

	return database.DB.Save(&order).Error
}

func (s *OrderService) ConfirmDelivery(userID uint, orderNo string) error {
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.OrderStatusShipped {
		return errors.New("订单状态不允许确认收货")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		order.Status = models.OrderStatusCompleted
		order.DeliveredAt = &now
		order.CompletedAt = &now
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		var seller models.User
		if err := tx.First(&seller, order.SellerID).Error; err != nil {
			return err
		}

		income := order.FinalPrice - order.Commission
		seller.WalletBalance += income
		if err := tx.Save(&seller).Error; err != nil {
			return err
		}

		sellerLog := models.WalletLog{
			UserID:      order.SellerID,
			Type:        models.WalletTypeIncome,
			Amount:      income,
			Balance:     seller.WalletBalance,
			OrderNo:     order.OrderNo,
			Description: "商品销售收入",
		}
		if err := tx.Create(&sellerLog).Error; err != nil {
			return err
		}

		return tx.Model(&models.Product{}).Where("id = ?", order.ProductID).
			Updates(map[string]interface{}{
				"status":     models.ProductStatusSoldOut,
				"sold_count": gorm.Expr("sold_count + 1"),
			}).Error
	})
}

func (s *OrderService) CancelOrder(userID uint, orderNo string, userRole string) error {
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if userRole == models.RoleBuyer && order.BuyerID != userID {
		return errors.New("无权操作此订单")
	}
	if userRole == models.RoleSeller && order.SellerID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status > models.OrderStatusPaid {
		return errors.New("订单状态不允许取消")
	}

	now := time.Now()
	order.Status = models.OrderStatusCancelled
	order.CancelledAt = &now

	return database.DB.Save(&order).Error
}

func (s *OrderService) RefundOrder(userID uint, orderNo string, reason string) error {
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.OrderStatusShipped && order.Status != models.OrderStatusDelivered {
		return errors.New("订单状态不允许申请退款")
	}

	order.Status = models.OrderStatusRefunding
	if reason != "" {
		order.Remark = reason
	}

	return database.DB.Save(&order).Error
}

func (s *OrderService) HandleRefund(orderID uint, approved bool) error {
	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.Status != models.OrderStatusRefunding {
		return errors.New("订单状态不允许处理退款")
	}

	if !approved {
		order.Status = models.OrderStatusShipped
		return database.DB.Save(&order).Error
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		order.Status = models.OrderStatusRefunded
		order.RefundedAt = &now
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		if order.PaymentMethod == "wallet" {
			var buyer models.User
			if err := tx.First(&buyer, order.BuyerID).Error; err != nil {
				return err
			}
			buyer.WalletBalance += order.FinalPrice
			if err := tx.Save(&buyer).Error; err != nil {
				return err
			}

			refundLog := models.WalletLog{
				UserID:      order.BuyerID,
				Type:        models.WalletTypeRefund,
				Amount:      order.FinalPrice,
				Balance:     buyer.WalletBalance,
				OrderNo:     order.OrderNo,
				Description: "订单退款",
			}
			if err := tx.Create(&refundLog).Error; err != nil {
				return err
			}
		}

		return tx.Model(&models.Product{}).Where("id = ?", order.ProductID).
			Update("status", models.ProductStatusOnSale).Error
	})
}

func (s *OrderService) NegotiatePrice(userID uint, orderNo string, offeredPrice float64, message string) error {
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID {
		return errors.New("无权操作此订单")
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("订单状态不允许议价")
	}

	var negotiation models.Negotiation
	if err := database.DB.Where("order_id = ?", order.ID).First(&negotiation).Error; err != nil {
		return errors.New("议价记录不存在")
	}

	negotiation.OfferedPrice = offeredPrice
	negotiation.BuyerMessage = message
	negotiation.Status = models.NegotiationStatusPending

	return database.DB.Save(&negotiation).Error
}

func (s *OrderService) HandleNegotiation(sellerID uint, orderNo string, accepted bool, counterPrice float64, message string) error {
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.SellerID != sellerID {
		return errors.New("无权操作此订单")
	}

	var negotiation models.Negotiation
	if err := database.DB.Where("order_id = ?", order.ID).First(&negotiation).Error; err != nil {
		return errors.New("议价记录不存在")
	}

	if accepted {
		order.FinalPrice = negotiation.OfferedPrice
		negotiation.Status = models.NegotiationStatusAccepted
	} else {
		negotiation.Status = models.NegotiationStatusRejected
		negotiation.CounterPrice = counterPrice
		negotiation.SellerMessage = message
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		return tx.Save(&negotiation).Error
	})
}

func (s *OrderService) GetOrderStats(userID uint, userRole string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var pendingCount int64
	var paidCount int64
	var shippedCount int64
	var completedCount int64

	db := database.DB.Model(&models.Order{})
	if userRole == models.RoleBuyer {
		db = db.Where("buyer_id = ?", userID)
	} else if userRole == models.RoleSeller {
		db = db.Where("seller_id = ?", userID)
	}

	db.Where("status = ?", models.OrderStatusPending).Count(&pendingCount)
	db.Where("status = ?", models.OrderStatusPaid).Count(&paidCount)
	db.Where("status = ?", models.OrderStatusShipped).Count(&shippedCount)
	db.Where("status = ?", models.OrderStatusCompleted).Count(&completedCount)

	stats["pending"] = pendingCount
	stats["paid"] = paidCount
	stats["shipped"] = shippedCount
	stats["completed"] = completedCount

	return stats, nil
}

func init() {
	logrus.Info("Order service initialized")
}
