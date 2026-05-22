package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/repository"
	"math/big"
	"time"

	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(req *dto.PaymentCreateRequest) (*model.Payment, error)
	RefundPayment(req *dto.PaymentRefundRequest) (*model.Payment, error)
	GeneratePaymentVoucher(paymentID uint) (*dto.PaymentVoucherResponse, error)
	GetOrderPayments(orderType model.OrderType, orderID uint) (*dto.OrderPaymentsResponse, error)
	GetPaymentByID(id uint) (*model.Payment, error)
	ListPayments(req *dto.PaymentListRequest) ([]model.Payment, int64, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	db          *gorm.DB
}

func NewPaymentService(paymentRepo repository.PaymentRepository, db *gorm.DB) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		db:          db,
	}
}

func (s *paymentService) CreatePayment(req *dto.PaymentCreateRequest) (*model.Payment, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, errors.New("创建支付记录失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	paymentNo, err := generatePaymentNo()
	if err != nil {
		tx.Rollback()
		logger.Errorf("生成支付凭证号失败: %v", err)
		return nil, errors.New("创建支付记录失败")
	}

	payment := &model.Payment{
		PaymentNo:     paymentNo,
		OrderType:     req.OrderType,
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentType:   req.PaymentType,
		Status:        model.PaymentStatusCompleted,
		TransactionID: req.TransactionID,
		Remark:        req.Remark,
	}

	var memberID *uint
	if req.OrderType == model.OrderTypeBooking {
		var booking model.Booking
		err = tx.Set("gorm:query_option", "FOR UPDATE").First(&booking, req.OrderID).Error
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("预订不存在")
			}
			logger.Errorf("查询预订失败: %v", err)
			return nil, errors.New("创建支付记录失败")
		}
		memberID = booking.MemberID
		booking.PaidAmount += req.Amount
		err = tx.Save(&booking).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("更新预订支付金额失败: %v", err)
			return nil, errors.New("创建支付记录失败")
		}
	} else if req.OrderType == model.OrderTypeCheckIn {
		var checkIn model.CheckIn
		err = tx.Set("gorm:query_option", "FOR UPDATE").First(&checkIn, req.OrderID).Error
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("入住记录不存在")
			}
			logger.Errorf("查询入住记录失败: %v", err)
			return nil, errors.New("创建支付记录失败")
		}
		memberID = s.getMemberIDFromCheckIn(&checkIn)
		checkIn.TotalAmount += req.Amount
		err = tx.Save(&checkIn).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("更新入住记录支付金额失败: %v", err)
			return nil, errors.New("创建支付记录失败")
		}
	}

	err = tx.Create(payment).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("创建支付记录失败: %v", err)
		return nil, errors.New("创建支付记录失败")
	}

	if memberID != nil && req.PaymentType != model.PaymentTypeRefund {
		err = s.calculateMemberPoints(tx, *memberID, req.Amount)
		if err != nil {
			logger.Warnf("计算会员积分失败: member_id=%d, err=%v", *memberID, err)
		}
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("提交事务失败: %v", err)
		return nil, errors.New("创建支付记录失败")
	}

	logger.Infof("支付记录创建成功: payment_no=%s, order_type=%s, order_id=%d, amount=%.2f",
		paymentNo, req.OrderType, req.OrderID, req.Amount)
	return payment, nil
}

func (s *paymentService) RefundPayment(req *dto.PaymentRefundRequest) (*model.Payment, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, errors.New("退款处理失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	originalPayment, err := s.paymentRepo.GetByID(req.PaymentID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("原支付记录不存在")
		}
		logger.Errorf("查询原支付记录失败: %v", err)
		return nil, errors.New("退款处理失败")
	}

	if originalPayment.Status == model.PaymentStatusRefunded {
		tx.Rollback()
		return nil, errors.New("该支付记录已退款")
	}

	if req.RefundAmount > originalPayment.Amount {
		tx.Rollback()
		return nil, errors.New("退款金额不能超过原支付金额")
	}

	refundNo, err := generatePaymentNo()
	if err != nil {
		tx.Rollback()
		logger.Errorf("生成退款凭证号失败: %v", err)
		return nil, errors.New("退款处理失败")
	}

	refundPayment := &model.Payment{
		PaymentNo:     refundNo,
		OrderType:     originalPayment.OrderType,
		OrderID:       originalPayment.OrderID,
		Amount:        req.RefundAmount,
		PaymentMethod: req.PaymentMethod,
		PaymentType:   model.PaymentTypeRefund,
		Status:        model.PaymentStatusCompleted,
		TransactionID: req.TransactionID,
		Remark:        fmt.Sprintf("退款，原支付号: %s，原因: %s", originalPayment.PaymentNo, req.Reason),
	}

	var memberID *uint
	if originalPayment.OrderType == model.OrderTypeBooking {
		var booking model.Booking
		err = tx.Set("gorm:query_option", "FOR UPDATE").First(&booking, originalPayment.OrderID).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("查询预订失败: %v", err)
			return nil, errors.New("退款处理失败")
		}
		memberID = booking.MemberID
		booking.PaidAmount -= req.RefundAmount
		if booking.PaidAmount < 0 {
			booking.PaidAmount = 0
		}
		err = tx.Save(&booking).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("更新预订支付金额失败: %v", err)
			return nil, errors.New("退款处理失败")
		}
	} else if originalPayment.OrderType == model.OrderTypeCheckIn {
		var checkIn model.CheckIn
		err = tx.Set("gorm:query_option", "FOR UPDATE").First(&checkIn, originalPayment.OrderID).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("查询入住记录失败: %v", err)
			return nil, errors.New("退款处理失败")
		}
		memberID = s.getMemberIDFromCheckIn(&checkIn)
		checkIn.TotalAmount -= req.RefundAmount
		if checkIn.TotalAmount < 0 {
			checkIn.TotalAmount = 0
		}
		err = tx.Save(&checkIn).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("更新入住记录支付金额失败: %v", err)
			return nil, errors.New("退款处理失败")
		}
	}

	err = tx.Create(refundPayment).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("创建退款记录失败: %v", err)
		return nil, errors.New("退款处理失败")
	}

	if req.RefundAmount >= originalPayment.Amount {
		err = tx.Model(&model.Payment{}).Where("id = ?", originalPayment.ID).Update("status", model.PaymentStatusRefunded).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("更新原支付状态失败: %v", err)
			return nil, errors.New("退款处理失败")
		}
	}

	if memberID != nil {
		err = s.refundMemberPoints(tx, *memberID, req.RefundAmount)
		if err != nil {
			logger.Warnf("回退会员积分失败: member_id=%d, err=%v", *memberID, err)
		}
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("提交事务失败: %v", err)
		return nil, errors.New("退款处理失败")
	}

	logger.Infof("退款处理成功: refund_no=%s, original_payment_no=%s, refund_amount=%.2f",
		refundNo, originalPayment.PaymentNo, req.RefundAmount)
	return refundPayment, nil
}

func (s *paymentService) GeneratePaymentVoucher(paymentID uint) (*dto.PaymentVoucherResponse, error) {
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		logger.Errorf("查询支付记录失败: id=%d, err=%v", paymentID, err)
		return nil, errors.New("支付记录不存在")
	}

	voucher := &dto.PaymentVoucherResponse{
		PaymentNo:     payment.PaymentNo,
		OrderType:     payment.OrderType,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		PaymentType:   payment.PaymentType,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
		Remark:        payment.Remark,
		CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if payment.OrderType == model.OrderTypeBooking {
		var booking model.Booking
		err = s.db.Preload("Room").First(&booking, payment.OrderID).Error
		if err == nil {
			voucher.GuestName = booking.GuestName
			if booking.Room != nil {
				voucher.RoomNo = booking.Room.RoomNo
			}
		}
	} else if payment.OrderType == model.OrderTypeCheckIn {
		var checkIn model.CheckIn
		err = s.db.Preload("Room").First(&checkIn, payment.OrderID).Error
		if err == nil {
			voucher.GuestName = checkIn.GuestName
			if checkIn.Room != nil {
				voucher.RoomNo = checkIn.Room.RoomNo
			}
		}
	}

	return voucher, nil
}

func (s *paymentService) GetOrderPayments(orderType model.OrderType, orderID uint) (*dto.OrderPaymentsResponse, error) {
	payments, err := s.paymentRepo.GetPaymentsByOrder(orderType, orderID)
	if err != nil {
		logger.Errorf("获取订单支付记录失败: order_type=%s, order_id=%d, err=%v", orderType, orderID, err)
		return nil, errors.New("获取订单支付记录失败")
	}

	totalPaid, err := s.paymentRepo.GetTotalPaidByOrder(orderType, orderID)
	if err != nil {
		logger.Errorf("计算订单已支付总额失败: order_type=%s, order_id=%d, err=%v", orderType, orderID, err)
		return nil, errors.New("计算订单已支付总额失败")
	}

	totalRefund := 0.0
	for _, p := range payments {
		if p.PaymentType == model.PaymentTypeRefund && p.Status == model.PaymentStatusCompleted {
			totalRefund += p.Amount
		}
	}

	paymentResponses := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		paymentResponses = append(paymentResponses, convertToPaymentResponse(&p))
	}

	return &dto.OrderPaymentsResponse{
		Payments:    paymentResponses,
		TotalPaid:   totalPaid,
		TotalRefund: totalRefund,
	}, nil
}

func (s *paymentService) GetPaymentByID(id uint) (*model.Payment, error) {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取支付记录失败: id=%d, err=%v", id, err)
		return nil, errors.New("支付记录不存在")
	}
	return payment, nil
}

func (s *paymentService) ListPayments(req *dto.PaymentListRequest) ([]model.Payment, int64, error) {
	page := req.GetPage()
	pageSize := req.GetPageSize()

	payments, total, err := s.paymentRepo.List(
		page, pageSize,
		req.PaymentNo,
		req.OrderType,
		req.OrderID,
		req.PaymentMethod,
		req.PaymentType,
		req.Status,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		logger.Errorf("获取支付记录列表失败: %v", err)
		return nil, 0, errors.New("获取支付记录列表失败")
	}

	return payments, total, nil
}

func generatePaymentNo() (string, error) {
	now := time.Now().Format("20060102150405")
	random, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PY%s%05d", now, random.Int64()), nil
}

func (s *paymentService) getMemberIDFromCheckIn(checkIn *model.CheckIn) *uint {
	if checkIn.BookingID != nil {
		var booking model.Booking
		err := s.db.Select("member_id").First(&booking, *checkIn.BookingID).Error
		if err == nil {
			return booking.MemberID
		}
	}
	return nil
}

func (s *paymentService) calculateMemberPoints(tx *gorm.DB, memberID uint, amount float64) error {
	var member model.Member
	err := tx.Preload("Level").First(&member, memberID).Error
	if err != nil {
		return err
	}

	if member.Status != model.MemberStatusActive || member.Level == nil {
		return nil
	}

	pointsRate := member.Level.PointsRate
	if pointsRate <= 0 {
		pointsRate = 1.0
	}

	points := int(amount * pointsRate)
	if points > 0 {
		member.Points += points
		err = tx.Save(&member).Error
		if err != nil {
			return err
		}
		logger.Infof("会员积分增加: member_id=%d, points=+%d, total=%d", memberID, points, member.Points)
	}

	return nil
}

func (s *paymentService) refundMemberPoints(tx *gorm.DB, memberID uint, amount float64) error {
	var member model.Member
	err := tx.Preload("Level").First(&member, memberID).Error
	if err != nil {
		return err
	}

	if member.Status != model.MemberStatusActive || member.Level == nil {
		return nil
	}

	pointsRate := member.Level.PointsRate
	if pointsRate <= 0 {
		pointsRate = 1.0
	}

	points := int(amount * pointsRate)
	if points > 0 {
		member.Points -= points
		if member.Points < 0 {
			member.Points = 0
		}
		err = tx.Save(&member).Error
		if err != nil {
			return err
		}
		logger.Infof("会员积分回退: member_id=%d, points=-%d, total=%d", memberID, points, member.Points)
	}

	return nil
}

func convertToPaymentResponse(payment *model.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:            payment.ID,
		PaymentNo:     payment.PaymentNo,
		OrderType:     payment.OrderType,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		PaymentType:   payment.PaymentType,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
		Remark:        payment.Remark,
		CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
