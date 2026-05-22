package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"ticket-system/internal/dto"
	"ticket-system/internal/logging"
	"ticket-system/internal/models"
	"ticket-system/internal/redis"
	"ticket-system/internal/repository"
	"ticket-system/internal/util"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
	showRepo  *repository.ShowRepository
	userRepo  *repository.UserRepository
	userSvc   *UserService
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: repository.NewOrderRepository(),
		showRepo:  repository.NewShowRepository(),
		userRepo:  repository.NewUserRepository(),
		userSvc:   NewUserService(),
	}
}

func (s *OrderService) CreateOrder(userID uint64, req *dto.OrderCreateRequest) (*models.Order, error) {
	if !util.ValidateRealName(req.RealName) {
		return nil, errors.New("真实姓名格式不正确")
	}
	if !util.ValidateIDCard(req.IDCard) {
		return nil, errors.New("身份证号格式不正确")
	}
	if !util.ValidatePhone(req.Phone) {
		return nil, errors.New("手机号格式不正确")
	}

	session, err := s.showRepo.GetSessionByID(req.SessionID)
	if err != nil {
		return nil, errors.New("场次不存在")
	}

	areas, err := s.showRepo.GetSeatAreasBySessionID(req.SessionID)
	if err != nil {
		return nil, err
	}
	areaMap := make(map[uint64]models.SeatArea)
	for _, area := range areas {
		areaMap[area.ID] = area
	}

	seats, err := s.showRepo.GetSeatsByIDs(req.SessionID, req.SeatIDs)
	if err != nil || len(seats) != len(req.SeatIDs) {
		return nil, errors.New("部分座位不存在")
	}

	for _, seat := range seats {
		lockUser, err := redis.GetSeatLockUser(req.SessionID, seat.ID)
		if err != nil || lockUser != userID {
			return nil, fmt.Errorf("座位 %s 未被锁定或已被他人锁定", seat.SeatNo)
		}
		if seat.Status == models.SeatStatusSold {
			return nil, fmt.Errorf("座位 %s 已售出", seat.SeatNo)
		}
	}

	var totalAmount float64
	for _, seat := range seats {
		if area, ok := areaMap[seat.AreaID]; ok {
			totalAmount += area.Price
		}
	}

	discount, _ := s.userSvc.GetMemberDiscount(userID)
	var discountAmount float64
	if req.CouponCode != "" {
		coupon, err := s.userRepo.GetCouponByCode(req.CouponCode)
		if err == nil && coupon != nil && coupon.UserID == userID && coupon.Status == 1 && time.Now().Before(coupon.ExpireAt) {
			discountAmount += coupon.Value
		}
	}

	discountAmount += totalAmount * (1 - discount)
	payAmount := totalAmount - discountAmount
	if payAmount < 0 {
		payAmount = 0
	}

	order := &models.Order{
		OrderNo:     util.GenerateOrderNo(),
		UserID:      userID,
		ShowID:      session.ShowID,
		SessionID:   req.SessionID,
		TotalAmount: util.RoundFloat(totalAmount, 2),
		Discount:    util.RoundFloat(discountAmount, 2),
		PayAmount:   util.RoundFloat(payAmount, 2),
		Status:      models.OrderStatusPending,
		RealName:    req.RealName,
		IDCard:      req.IDCard,
		Phone:       req.Phone,
		Email:       req.Email,
		Remark:      req.Remark,
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	for _, seat := range seats {
		area := areaMap[seat.AreaID]
		ticket := &models.Ticket{
			OrderID:   order.ID,
			UserID:    userID,
			ShowID:    session.ShowID,
			SessionID: req.SessionID,
			SeatID:    seat.ID,
			AreaID:    seat.AreaID,
			TicketNo:  util.GenerateTicketNo(),
			QrCode:    util.GenerateRandomString(32),
			Price:     area.Price,
			SeatInfo:  area.Name + " " + seat.SeatNo,
			Status:    models.TicketStatusValid,
			RealName:  req.RealName,
			IDCard:    req.IDCard,
		}
		_ = s.orderRepo.CreateTicket(ticket)
	}

	if req.CouponCode != "" {
		coupon, _ := s.userRepo.GetCouponByCode(req.CouponCode)
		if coupon != nil {
			coupon.Status = 0
			_ = s.userRepo.UpdateCoupon(coupon)
		}
	}

	logging.Infof("Order created: %s, user=%d", order.OrderNo, userID)
	return order, nil
}

func (s *OrderService) GetOrder(orderNo string) (*models.Order, error) {
	return s.orderRepo.GetByOrderNo(orderNo)
}

func (s *OrderService) GetOrderByID(id uint64) (*models.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) ListOrders(userID uint64, req *dto.OrderListRequest) (*dto.PaginatedResponse, error) {
	orders, total, err := s.orderRepo.List(userID, req.Page, req.PageSize, req.Status, req.StartDate, req.EndDate, req.Keyword)
	if err != nil {
		return nil, err
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	return &dto.PaginatedResponse{
		List: orders,
		Pagination: dto.Pagination{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			Pages:    pages,
		},
	}, nil
}

func (s *OrderService) PayOrder(orderNo string, payType int, userID uint64) (*models.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.UserID != userID {
		return nil, errors.New("无权限支付此订单")
	}

	if order.Status == models.OrderStatusPaid {
		return nil, errors.New("订单已支付")
	}

	if order.Status != models.OrderStatusPending {
		return nil, errors.New("订单状态不正确")
	}

	maxRetries := 3
	retryDelay := 500 * time.Millisecond
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		paymentLog := &models.PaymentLog{
			OrderID: order.ID,
			PayType: payType,
			Amount:  order.PayAmount,
			Status:  0,
		}
		_ = s.orderRepo.CreatePaymentLog(paymentLog)

		success := s.simulatePayment(payType, order.PayAmount)
		if success {
			_ = s.orderRepo.UpdatePaymentLogStatus(paymentLog.ID, 1, "支付成功")
			_ = s.orderRepo.MarkPaid(order.ID, payType)

			tickets, _ := s.orderRepo.GetTicketsByOrderID(order.ID)
			for _, ticket := range tickets {
				_ = s.showRepo.UpdateSeatsStatus(order.SessionID, []uint64{ticket.SeatID}, models.SeatStatusSold)
				_ = redis.SetSeatStatus(order.SessionID, ticket.SeatID, models.SeatStatusSold)
				_ = redis.UnlockSeat(order.SessionID, ticket.SeatID)
				_ = s.showRepo.UpdateSessionSoldSeats(order.SessionID, 1)
				_ = s.showRepo.UpdateAreaSoldSeats(ticket.AreaID, 1)
			}

			points := int(order.PayAmount)
			_ = s.userSvc.AddPoints(userID, points)

			s.sendTickets(order)

			logging.Infof("Order paid: %s, pay_type=%d, attempt=%d", orderNo, payType, attempt)
			return order, nil
		}

		lastErr = fmt.Errorf("支付失败，尝试 %d/%d", attempt, maxRetries)
		_ = s.orderRepo.UpdatePaymentLogStatus(paymentLog.ID, 2, lastErr.Error())
		logging.Warnf("Payment attempt %d failed for order %s: %v", attempt, orderNo, lastErr)

		if attempt < maxRetries {
			time.Sleep(retryDelay)
			retryDelay *= 2
		}
	}

	return nil, fmt.Errorf("支付失败，已重试%d次，请稍后重试或联系客服", maxRetries)
}

func (s *OrderService) simulatePayment(payType int, amount float64) bool {
	time.Sleep(500 * time.Millisecond)
	return util.GenerateRandomInt(0, 100) < 95
}

func (s *OrderService) sendTickets(order *models.Order) {
	if order.Email != "" {
		logging.Infof("Sending ticket email to %s for order %s", order.Email, order.OrderNo)
	}
	if order.Phone != "" {
		logging.Infof("Sending ticket SMS to %s for order %s", order.Phone, order.OrderNo)
	}
}

func (s *OrderService) CancelOrder(orderNo string, userID uint64) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.UserID != userID {
		return errors.New("无权限取消此订单")
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("只能取消待支付订单")
	}

	err = s.orderRepo.UpdateStatus(order.ID, models.OrderStatusCanceled)
	if err != nil {
		return err
	}

	tickets, _ := s.orderRepo.GetTicketsByOrderID(order.ID)
	for _, ticket := range tickets {
		_ = redis.UnlockSeat(order.SessionID, ticket.SeatID)
		_ = redis.SetSeatStatus(order.SessionID, ticket.SeatID, models.SeatStatusAvailable)
	}

	logging.Infof("Order canceled: %s", orderNo)
	return nil
}

func (s *OrderService) RequestRefund(orderNo string, userID uint64, reason string) (*models.Refund, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, errors.New("订单不存在")
	}

	if order.UserID != userID {
		return nil, errors.New("无权限申请退款")
	}

	if order.Status != models.OrderStatusPaid {
		return nil, errors.New("只能对已支付订单申请退款")
	}

	existingRefund, _ := s.orderRepo.GetRefundByOrderID(order.ID)
	if existingRefund != nil {
		return nil, errors.New("退款申请已存在")
	}

	refund := &models.Refund{
		OrderID:      order.ID,
		UserID:       userID,
		RefundNo:     util.GenerateRefundNo(),
		RefundAmount: order.PayAmount,
		Reason:       reason,
		Status:       models.RefundStatusPending,
	}

	err = s.orderRepo.CreateRefund(refund)
	if err != nil {
		return nil, err
	}

	_ = s.orderRepo.UpdateStatus(order.ID, models.OrderStatusRefunding)

	logging.Infof("Refund requested: %s for order %s", refund.RefundNo, orderNo)
	return refund, nil
}

func (s *OrderService) AuditRefund(refundNo string, status int, remark string) error {
	refund, err := s.orderRepo.GetRefundByRefundNo(refundNo)
	if err != nil {
		return errors.New("退款申请不存在")
	}

	if refund.Status != models.RefundStatusPending {
		return errors.New("退款申请已审核")
	}

	err = s.orderRepo.AuditRefund(refund.ID, status, remark)
	if err != nil {
		return err
	}

	if status == models.RefundStatusApproved {
		_ = s.orderRepo.UpdateStatus(refund.OrderID, models.OrderStatusRefunded)

		order, _ := s.orderRepo.GetByID(refund.OrderID)
		if order != nil {
			tickets, _ := s.orderRepo.GetTicketsByOrderID(order.ID)
			for _, ticket := range tickets {
				_ = s.orderRepo.UpdateTicket(&models.Ticket{BaseModel: ticket.BaseModel, Status: models.TicketStatusRefunded})
				_ = s.showRepo.UpdateSeatsStatus(order.SessionID, []uint64{ticket.SeatID}, models.SeatStatusAvailable)
				_ = redis.SetSeatStatus(order.SessionID, ticket.SeatID, models.SeatStatusAvailable)
			}
		}
	} else {
		_ = s.orderRepo.UpdateStatus(refund.OrderID, models.OrderStatusPaid)
	}

	logging.Infof("Refund audited: %s, status=%d", refundNo, status)
	return nil
}

func (s *OrderService) ExportOrders(startDate, endDate string, status int) (*excelize.File, error) {
	orders, err := s.orderRepo.GetAllOrdersForExport(startDate, endDate, status)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetName := "订单列表"
	index, _ := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	headers := []string{"订单号", "下单时间", "用户ID", "实名人", "手机号", "订单金额", "优惠金额", "支付金额", "订单状态", "支付方式", "支付时间"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, order := range orders {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), order.OrderNo)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), order.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), order.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), order.RealName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), order.Phone)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), order.TotalAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), order.Discount)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), order.PayAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), s.getOrderStatusText(order.Status))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), s.getPayTypeText(order.PayType))
		if order.PayTime != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), order.PayTime.Format("2006-01-02 15:04:05"))
		}
	}

	f.SetActiveSheet(index)
	return f, nil
}

func (s *OrderService) getOrderStatusText(status int) string {
	switch status {
	case models.OrderStatusPending:
		return "待支付"
	case models.OrderStatusPaid:
		return "已支付"
	case models.OrderStatusCanceled:
		return "已取消"
	case models.OrderStatusRefunding:
		return "退款中"
	case models.OrderStatusRefunded:
		return "已退款"
	default:
		return "未知"
	}
}

func (s *OrderService) getPayTypeText(payType int) string {
	switch payType {
	case 1:
		return "支付宝"
	case 2:
		return "微信支付"
	default:
		return "未支付"
	}
}
