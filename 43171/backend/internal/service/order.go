package service

import (
	"drone-rental/internal/config"
	"drone-rental/internal/dto"
	"drone-rental/internal/model"
	"drone-rental/internal/pkg/utils"
	"drone-rental/internal/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo   *repository.OrderRepo
	paymentRepo *repository.PaymentRepo
	droneRepo   *repository.DroneRepo
	userRepo    *repository.UserRepo
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:   repository.NewOrderRepo(),
		paymentRepo: repository.NewPaymentRepo(),
		droneRepo:   repository.NewDroneRepo(),
		userRepo:    repository.NewUserRepo(),
	}
}

func (s *OrderService) Create(userID uint, req *dto.CreateOrderReq) (*model.RentalOrder, error) {
	drone, err := s.droneRepo.GetByID(req.DroneID)
	if err != nil {
		return nil, errors.New("设备不存在")
	}
	if drone.Status != model.DroneStatusOnline {
		return nil, errors.New("设备当前不可租")
	}
	startDate, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
	if err != nil {
		return nil, errors.New("开始日期格式错误")
	}
	endDate, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
	if err != nil {
		return nil, errors.New("结束日期格式错误")
	}
	if !endDate.After(startDate) {
		return nil, errors.New("结束日期必须晚于开始日期")
	}
	totalDays := utils.DaysBetween(startDate, endDate)
	rentalFee := float64(totalDays) * drone.PricePerDay
	deposit := drone.Deposit
	if deposit <= 0 {
		deposit = rentalFee * config.Cfg.Insurance.DepositRate
	}
	insuranceFee := rentalFee * config.Cfg.Insurance.InsuranceRate
	totalAmount := rentalFee + deposit + insuranceFee

	order := &model.RentalOrder{
		OrderNo:      utils.GenerateOrderNo("RN"),
		UserID:       userID,
		DroneID:      req.DroneID,
		StartDate:    startDate,
		EndDate:      endDate,
		Region:       req.Region,
		Address:      req.Address,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		TotalDays:    totalDays,
		PricePerDay:  drone.PricePerDay,
		RentalFee:    rentalFee,
		Deposit:      deposit,
		InsuranceFee: insuranceFee,
		LateFee:      0,
		TotalAmount:  totalAmount,
		PaidAmount:   0,
		Status:       model.OrderStatusPending,
		Remark:       req.Remark,
	}
	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetByID(id uint) (*model.RentalOrder, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) List(page, pageSize int, userID, droneID uint, status model.OrderStatus) ([]model.RentalOrder, int64, error) {
	return s.orderRepo.List(page, pageSize, userID, droneID, status)
}

func (s *OrderService) PayOrder(userID uint, req *dto.PayOrderReq) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		order, err := s.orderRepo.GetByID(req.OrderID)
		if err != nil {
			return errors.New("订单不存在")
		}
		if order.UserID != userID {
			return errors.New("无权操作该订单")
		}
		if order.Status != model.OrderStatusPending {
			return errors.New("订单状态不允许支付")
		}
		user, err := s.userRepo.GetByID(userID)
		if err != nil {
			return errors.New("用户不存在")
		}
		if req.PayType == "balance" {
			if user.Balance < order.TotalAmount {
				return errors.New("余额不足")
			}
			user.Balance -= order.TotalAmount
		}
		user.Deposit += order.Deposit
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		order.Status = model.OrderStatusPaid
		order.PaidAmount = order.TotalAmount
		if err := tx.Save(order).Error; err != nil {
			return err
		}
		now := time.Now()
		payment := &model.Payment{
			PaymentNo: utils.GenerateOrderNo("PM"),
			OrderID:   order.ID,
			UserID:    userID,
			Amount:    order.TotalAmount,
			PayType:   req.PayType,
			Status:    model.PaymentStatusSuccess,
			PaidAt:    &now,
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		return tx.Model(&model.Drone{}).Where("id = ?", order.DroneID).
			Update("status", model.DroneStatusRented).Error
	})
}

func (s *OrderService) CancelOrder(userID uint, req *dto.CancelOrderReq) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		order, err := s.orderRepo.GetByID(req.OrderID)
		if err != nil {
			return errors.New("订单不存在")
		}
		if order.UserID != userID {
			return errors.New("无权操作该订单")
		}
		if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusPaid {
			return errors.New("订单状态不允许取消")
		}
		if order.Status == model.OrderStatusPaid {
			user, _ := s.userRepo.GetByID(userID)
			user.Balance += order.TotalAmount
			user.Deposit -= order.Deposit
			tx.Save(user)
		}
		order.Status = model.OrderStatusCancelled
		order.CancelReason = req.CancelReason
		tx.Save(order)
		return tx.Model(&model.Drone{}).Where("id = ?", order.DroneID).
			Update("status", model.DroneStatusOnline).Error
	})
}

func (s *OrderService) PickupOrder(userID, orderID uint) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.UserID != userID {
		return errors.New("无权操作该订单")
	}
	if order.Status != model.OrderStatusPaid {
		return errors.New("订单状态不允许取机")
	}
	order.Status = model.OrderStatusPicked
	return s.orderRepo.Update(order)
}

func (s *OrderService) ConfirmReturn(userID uint, req *dto.ConfirmReturnReq) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		order, err := s.orderRepo.GetByID(req.OrderID)
		if err != nil {
			return errors.New("订单不存在")
		}
		if order.UserID != userID {
			return errors.New("无权操作该订单")
		}
		if order.Status != model.OrderStatusPicked {
			return errors.New("订单状态不允许归还")
		}
		returnDate := time.Now()
		if req.ReturnDate != nil {
			returnDate = *req.ReturnDate
		}
		order.ReturnDate = &returnDate
		if returnDate.After(order.EndDate) {
			lateDays := utils.DaysBetween(order.EndDate, returnDate) - 1
			order.LateFee = float64(lateDays) * order.PricePerDay * config.Cfg.Insurance.LateFeeRate
		}
		user, _ := s.userRepo.GetByID(userID)
		refundDeposit := order.Deposit - order.LateFee
		if refundDeposit > 0 {
			user.Deposit -= refundDeposit
			user.Balance += refundDeposit
			order.RefundAmount = refundDeposit
			tx.Save(user)
		}
		order.Status = model.OrderStatusReturned
		tx.Save(order)
		return tx.Model(&model.Drone{}).Where("id = ?", order.DroneID).
			Update("status", model.DroneStatusOnline).Error
	})
}

func (s *OrderService) CompleteOrder(orderID uint) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != model.OrderStatusReturned {
		return errors.New("订单状态不允许完成")
	}
	order.Status = model.OrderStatusCompleted
	return s.orderRepo.Update(order)
}

func (s *OrderService) ProcessLateFees() error {
	orders, err := s.orderRepo.GetOverdueOrders()
	if err != nil {
		return err
	}
	for _, order := range orders {
		if order.ReturnDate == nil {
			lateDays := utils.DaysBetween(order.EndDate, time.Now()) - 1
			order.LateFee = float64(lateDays) * order.PricePerDay * config.Cfg.Insurance.LateFeeRate
			s.orderRepo.Update(&order)
		}
	}
	return nil
}
