package service

import (
	"errors"
	"fmt"
	"ticket-system/internal/dto"
	"ticket-system/internal/logger"
	"ticket-system/internal/models"
	"ticket-system/internal/redis"
	"ticket-system/internal/util"
	"time"
)

type OrderService struct {
	ticketTypeService *TicketTypeService
	couponService     *CouponService
}

func NewOrderService() *OrderService {
	return &OrderService{
		ticketTypeService: NewTicketTypeService(),
		couponService:     NewCouponService(),
	}
}

func (s *OrderService) Create(userID uint, req *dto.OrderCreateRequest) (*models.Order, error) {
	orderNo := util.GenerateOrderNo()

	locked, err := redis.LockOrder(orderNo, 30*time.Second)
	if err != nil || !locked {
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	defer redis.UnlockOrder(orderNo)

	var activity models.Activity
	if err := models.DB.First(&activity, req.ActivityID).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	if activity.Status != models.ActivityStatusPublished {
		return nil, errors.New("活动未发布，无法购票")
	}

	ticketMap := make(map[uint]int)
	for _, t := range req.Tickets {
		ticketMap[t.TicketTypeID] += t.Quantity
	}

	var totalAmount float64
	var orderItems []models.OrderItem
	var ticketTypes []*models.TicketType

	for ticketTypeID, quantity := range ticketMap {
		ticketType, err := s.ticketTypeService.GetByID(ticketTypeID)
		if err != nil {
			return nil, errors.New("票型不存在")
		}

		if ticketType.ActivityID != req.ActivityID {
			return nil, errors.New("票型不属于该活动")
		}

		if ticketType.Status != models.TicketStatusOnSale {
			return nil, fmt.Errorf("票型[%s]已下架或售罄", ticketType.Name)
		}

		stock, err := redis.GetStock(ticketTypeID)
		if err != nil {
			stock = int64(ticketType.Stock)
		}

		if stock < int64(quantity) {
			return nil, fmt.Errorf("票型[%s]库存不足", ticketType.Name)
		}

		remaining, err := redis.DecrementStock(ticketTypeID, quantity)
		if err != nil {
			return nil, errors.New("扣减库存失败")
		}

		if remaining < 0 {
			redis.IncrementStock(ticketTypeID, quantity)
			return nil, fmt.Errorf("票型[%s]库存不足", ticketType.Name)
		}

		subtotal := float64(quantity) * ticketType.Price
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			TicketTypeID: ticketTypeID,
			Quantity:     quantity,
			UnitPrice:    ticketType.Price,
			Subtotal:     subtotal,
		})

		ticketTypes = append(ticketTypes, ticketType)
	}

	var discount float64
	var couponID *uint

	if req.CouponCode != "" {
		coupon, err := s.couponService.GetByCode(req.CouponCode)
		if err != nil {
			s.rollbackStock(ticketMap)
			return nil, errors.New("优惠券不存在")
		}

		discount, err = s.couponService.Validate(coupon, totalAmount)
		if err != nil {
			s.rollbackStock(ticketMap)
			return nil, err
		}

		couponID = &coupon.ID
	}

	payAmount := totalAmount - discount
	if payAmount < 0 {
		payAmount = 0
	}

	tx := models.DB.Begin()
	if tx.Error != nil {
		s.rollbackStock(ticketMap)
		return nil, errors.New("创建订单失败")
	}

	order := &models.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		ActivityID:  req.ActivityID,
		CouponID:    couponID,
		TotalAmount: totalAmount,
		Discount:    discount,
		PayAmount:   payAmount,
		Status:      models.OrderStatusPending,
		Remark:      req.Remark,
		OrderItems:  orderItems,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		s.rollbackStock(ticketMap)
		logger.Log.Errorf("Create order failed: %v", err)
		return nil, errors.New("创建订单失败")
	}

	var checkIns []models.CheckIn
	for _, item := range order.OrderItems {
		for i := 0; i < item.Quantity; i++ {
			qrCode := util.GenerateQRCodeContent(orderNo, item.ID)
			checkIns = append(checkIns, models.CheckIn{
				OrderID:      order.ID,
				OrderItemID:  item.ID,
				UserID:       userID,
				ActivityID:   req.ActivityID,
				TicketTypeID: item.TicketTypeID,
				QrCode:       qrCode,
				CheckedIn:    false,
			})
		}
	}

	if len(checkIns) > 0 {
		if err := tx.Create(&checkIns).Error; err != nil {
			tx.Rollback()
			s.rollbackStock(ticketMap)
			logger.Log.Errorf("Create checkin failed: %v", err)
			return nil, errors.New("创建订单失败")
		}
	}

	if couponID != nil {
		if err := s.couponService.UseCoupon(*couponID); err != nil {
			tx.Rollback()
			s.rollbackStock(ticketMap)
			return nil, err
		}
	}

	for _, tt := range ticketTypes {
		tt.SoldCount += ticketMap[tt.ID]
		tx.Save(tt)
		s.ticketTypeService.UpdateSoldStatus(tt.ID)
	}

	if err := tx.Commit().Error; err != nil {
		s.rollbackStock(ticketMap)
		return nil, errors.New("创建订单失败")
	}

	logger.Log.Infof("Order created: %s, user: %d, amount: %.2f", orderNo, userID, payAmount)
	return order, nil
}

func (s *OrderService) rollbackStock(ticketMap map[uint]int) {
	for id, qty := range ticketMap {
		redis.IncrementStock(id, qty)
	}
}

func (s *OrderService) GetList(req *dto.OrderListRequest, userID uint, isAdmin bool) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := models.DB.Model(&models.Order{}).Preload("Activity").Preload("User").Preload("OrderItems.TicketType")

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if req.ActivityID > 0 {
		query = query.Where("activity_id = ?", req.ActivityID)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.OrderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *OrderService) GetByID(id uint, userID uint, isAdmin bool) (*models.Order, error) {
	var order models.Order
	query := models.DB.Preload("Activity").Preload("User").Preload("Coupon").Preload("OrderItems.TicketType").Preload("OrderItems.CheckIns")

	if isAdmin {
		query = query.First(&order, id)
	} else {
		query = query.Where("id = ? AND user_id = ?", id, userID).First(&order)
	}

	if query.Error != nil {
		return nil, query.Error
	}

	return &order, nil
}

func (s *OrderService) Pay(id uint, userID uint, isAdmin bool) (*models.Order, error) {
	order, err := s.GetByID(id, userID, isAdmin)
	if err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.Status != models.OrderStatusPending {
		return nil, errors.New("订单状态不正确，无法支付")
	}

	now := time.Now()
	order.Status = models.OrderStatusPaid
	order.PaidAt = &now

	if err := models.DB.Save(order).Error; err != nil {
		return nil, err
	}

	logger.Log.Infof("Order paid: %s", order.OrderNo)
	return order, nil
}

func (s *OrderService) Cancel(id uint, userID uint, isAdmin bool) (*models.Order, error) {
	order, err := s.GetByID(id, userID, isAdmin)
	if err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.Status != models.OrderStatusPending {
		return nil, errors.New("订单状态不正确，无法取消")
	}

	tx := models.DB.Begin()

	order.Status = models.OrderStatusCanceled
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, item := range order.OrderItems {
		if err := redis.IncrementStock(item.TicketTypeID, item.Quantity); err != nil {
			logger.Log.Warnf("Rollback stock failed: %d", item.TicketTypeID)
		}

		var ticketType models.TicketType
		tx.First(&ticketType, item.TicketTypeID)
		ticketType.SoldCount -= item.Quantity
		if ticketType.Status == models.TicketStatusSoldOut {
			ticketType.Status = models.TicketStatusOnSale
		}
		tx.Save(&ticketType)
	}

	if order.CouponID != nil {
		var coupon models.Coupon
		if err := tx.First(&coupon, *order.CouponID).Error; err == nil {
			coupon.UsedCount--
			if coupon.UsedCount < 0 {
				coupon.UsedCount = 0
			}
			if coupon.Status == models.CouponStatusUsed && coupon.UsedCount < coupon.TotalCount {
				coupon.Status = models.CouponStatusActive
			}
			tx.Save(&coupon)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	logger.Log.Infof("Order canceled: %s", order.OrderNo)
	return order, nil
}
