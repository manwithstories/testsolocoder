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

type CheckInService interface {
	CreateCheckIn(req *dto.CheckInCreateRequest) (*model.CheckIn, error)
	CheckOut(id uint, req *dto.CheckOutRequest) (*dto.CheckOutResponse, error)
	ExtendStay(id uint, req *dto.ExtendStayRequest) (*model.CheckIn, error)
	AddExtraCharge(id uint, req *dto.ExtraChargeRequest) (*model.CheckIn, error)
	GetCheckInDetail(id uint) (*model.CheckIn, error)
	ListCheckIns(req *dto.CheckInListRequest) ([]model.CheckIn, int64, error)
}

type checkInService struct {
	checkInRepo repository.CheckInRepository
	bookingRepo repository.BookingRepository
	roomRepo    repository.RoomRepository
	db          *gorm.DB
}

func NewCheckInService(checkInRepo repository.CheckInRepository, bookingRepo repository.BookingRepository, roomRepo repository.RoomRepository, db *gorm.DB) CheckInService {
	return &checkInService{
		checkInRepo: checkInRepo,
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		db:          db,
	}
}

func (s *checkInService) CreateCheckIn(req *dto.CheckInCreateRequest) (*model.CheckIn, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, errors.New("办理入住失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var room model.Room
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&room, req.RoomID).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("房间不存在: room_id=%d", req.RoomID)
			return nil, errors.New("房间不存在")
		}
		logger.Errorf("锁定房间失败: %v", err)
		return nil, errors.New("办理入住失败")
	}

	if room.Status != model.RoomStatusAvailable && room.Status != model.RoomStatusReserved {
		tx.Rollback()
		logger.Warnf("房间状态不可用: room_id=%d, status=%s", req.RoomID, room.Status)
		return nil, errors.New("房间状态不可用")
	}

	var activeCheckIn model.CheckIn
	err = tx.Where("room_id = ? AND status = ?", req.RoomID, model.CheckInStatusActive).First(&activeCheckIn).Error
	if err == nil {
		tx.Rollback()
		logger.Warnf("房间已有活跃入住记录: room_id=%d, check_in_id=%d", req.RoomID, activeCheckIn.ID)
		return nil, errors.New("房间已有客人入住")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		logger.Errorf("检查房间入住状态失败: %v", err)
		return nil, errors.New("办理入住失败")
	}

	if req.BookingID != nil {
		var booking model.Booking
		err = tx.First(&booking, *req.BookingID).Error
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Warnf("预订不存在: booking_id=%d", *req.BookingID)
				return nil, errors.New("预订不存在")
			}
			logger.Errorf("查询预订失败: %v", err)
			return nil, errors.New("办理入住失败")
		}

		if booking.Status != model.BookingStatusConfirmed {
			tx.Rollback()
			logger.Warnf("预订状态不允许入住: booking_id=%d, status=%s", *req.BookingID, booking.Status)
			return nil, errors.New("预订状态不允许入住")
		}

		if booking.RoomID != req.RoomID {
			tx.Rollback()
			logger.Warnf("预订房间与入住房间不匹配: booking_id=%d, booking_room_id=%d, req_room_id=%d", *req.BookingID, booking.RoomID, req.RoomID)
			return nil, errors.New("预订房间与入住房间不匹配")
		}

		err = tx.Model(&booking).Update("status", model.BookingStatusCompleted).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("更新预订状态失败: %v", err)
			return nil, errors.New("办理入住失败")
		}
	}

	days := int(req.ExpectedCheckOut.Sub(req.CheckInTime).Hours() / 24)
	if days <= 0 {
		tx.Rollback()
		logger.Warnf("入住天数无效: days=%d", days)
		return nil, errors.New("入住天数必须大于0")
	}

	totalAmount := room.Price * float64(days)

	checkInNo, err := generateCheckInNo()
	if err != nil {
		tx.Rollback()
		logger.Errorf("生成入住单号失败: %v", err)
		return nil, errors.New("办理入住失败")
	}

	checkIn := &model.CheckIn{
		CheckInNo:        checkInNo,
		BookingID:        req.BookingID,
		RoomID:           req.RoomID,
		GuestName:        req.GuestName,
		GuestPhone:       req.GuestPhone,
		GuestIDCard:      req.GuestIDCard,
		CheckInTime:      req.CheckInTime,
		ExpectedCheckOut: req.ExpectedCheckOut,
		Status:           model.CheckInStatusActive,
		Deposit:          req.Deposit,
		ExtraCharges:     0,
		TotalAmount:      totalAmount,
	}

	err = tx.Create(checkIn).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("创建入住记录失败: %v", err)
		return nil, errors.New("办理入住失败")
	}

	err = tx.Model(&room).Update("status", model.RoomStatusOccupied).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("更新房间状态失败: %v", err)
		return nil, errors.New("办理入住失败")
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("提交事务失败: %v", err)
		return nil, errors.New("办理入住失败")
	}

	logger.Infof("入住办理成功: check_in_no=%s, room_id=%d", checkInNo, req.RoomID)
	return checkIn, nil
}

func (s *checkInService) CheckOut(id uint, req *dto.CheckOutRequest) (*dto.CheckOutResponse, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, errors.New("退房失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var checkIn model.CheckIn
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&checkIn, id).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("入住记录不存在: id=%d", id)
			return nil, errors.New("入住记录不存在")
		}
		logger.Errorf("锁定入住记录失败: %v", err)
		return nil, errors.New("退房失败")
	}

	if checkIn.Status != model.CheckInStatusActive {
		tx.Rollback()
		logger.Warnf("入住状态不允许退房: id=%d, status=%s", id, checkIn.Status)
		return nil, errors.New("该入住记录不是活跃状态")
	}

	payAmount := checkIn.TotalAmount - checkIn.Deposit

	now := time.Now()
	checkIn.ActualCheckOut = &now
	checkIn.Status = model.CheckInStatusCheckedOut

	err = tx.Save(&checkIn).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("更新入住记录失败: %v", err)
		return nil, errors.New("退房失败")
	}

	err = tx.Model(&model.Room{}).Where("id = ?", checkIn.RoomID).Update("status", model.RoomStatusAvailable).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("更新房间状态失败: %v", err)
		return nil, errors.New("退房失败")
	}

	paymentNo, err := generateCheckInPaymentNo()
	if err != nil {
		tx.Rollback()
		logger.Errorf("生成支付单号失败: %v", err)
		return nil, errors.New("退房失败")
	}

	payment := &model.Payment{
		PaymentNo:     paymentNo,
		OrderType:     model.OrderTypeCheckIn,
		OrderID:       checkIn.ID,
		Amount:        payAmount,
		PaymentMethod: req.PaymentMethod,
		PaymentType:   model.PaymentTypeExtra,
		Status:        model.PaymentStatusCompleted,
		TransactionID: req.TransactionID,
		Remark:        req.Remark,
	}

	err = tx.Create(payment).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("创建支付记录失败: %v", err)
		return nil, errors.New("退房失败")
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("提交事务失败: %v", err)
		return nil, errors.New("退房失败")
	}

	logger.Infof("退房成功: check_in_no=%s, pay_amount=%.2f", checkIn.CheckInNo, payAmount)
	return &dto.CheckOutResponse{
		CheckIn:      nil,
		TotalAmount:  checkIn.TotalAmount,
		Deposit:      checkIn.Deposit,
		ExtraCharges: checkIn.ExtraCharges,
		PayAmount:    payAmount,
		PaymentNo:    paymentNo,
	}, nil
}

func (s *checkInService) ExtendStay(id uint, req *dto.ExtendStayRequest) (*model.CheckIn, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, errors.New("续住失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var checkIn model.CheckIn
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&checkIn, id).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("入住记录不存在: id=%d", id)
			return nil, errors.New("入住记录不存在")
		}
		logger.Errorf("锁定入住记录失败: %v", err)
		return nil, errors.New("续住失败")
	}

	if checkIn.Status != model.CheckInStatusActive {
		tx.Rollback()
		logger.Warnf("入住状态不允许续住: id=%d, status=%s", id, checkIn.Status)
		return nil, errors.New("该入住记录不是活跃状态")
	}

	newExpectedCheckOut := checkIn.ExpectedCheckOut.AddDate(0, 0, req.ExtendDays)

	var count int64
	query := tx.Model(&model.Booking{}).Where(
		"room_id = ? AND status IN ? AND check_in_date < ? AND check_out_date > ?",
		checkIn.RoomID,
		[]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed},
		newExpectedCheckOut,
		checkIn.ExpectedCheckOut,
	)
	err = query.Count(&count).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("检查房间可用性失败: %v", err)
		return nil, errors.New("续住失败")
	}
	if count > 0 {
		tx.Rollback()
		logger.Warnf("续住期间房间已被预订: room_id=%d", checkIn.RoomID)
		return nil, errors.New("续住期间房间已被预订，请选择其他日期")
	}

	var room model.Room
	err = tx.First(&room, checkIn.RoomID).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("获取房间信息失败: %v", err)
		return nil, errors.New("续住失败")
	}

	extendAmount := room.Price * float64(req.ExtendDays)
	checkIn.ExpectedCheckOut = newExpectedCheckOut
	checkIn.TotalAmount += extendAmount

	err = tx.Save(&checkIn).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("更新入住记录失败: %v", err)
		return nil, errors.New("续住失败")
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("提交事务失败: %v", err)
		return nil, errors.New("续住失败")
	}

	logger.Infof("续住成功: check_in_no=%s, extend_days=%d, extend_amount=%.2f", checkIn.CheckInNo, req.ExtendDays, extendAmount)
	return &checkIn, nil
}

func (s *checkInService) AddExtraCharge(id uint, req *dto.ExtraChargeRequest) (*model.CheckIn, error) {
	checkIn, err := s.checkInRepo.GetByID(id)
	if err != nil {
		logger.Errorf("添加上消费失败，入住记录不存在: id=%d, err=%v", id, err)
		return nil, errors.New("入住记录不存在")
	}

	if checkIn.Status != model.CheckInStatusActive {
		logger.Warnf("入住状态不允许添加消费: id=%d, status=%s", id, checkIn.Status)
		return nil, errors.New("该入住记录不是活跃状态")
	}

	err = s.checkInRepo.AddExtraCharge(id, req.Amount, req.Description)
	if err != nil {
		logger.Errorf("添加额外消费失败: %v", err)
		return nil, errors.New("添加额外消费失败")
	}

	checkIn, err = s.checkInRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取更新后的入住记录失败: %v", err)
		return nil, errors.New("添加额外消费失败")
	}

	logger.Infof("添加额外消费成功: check_in_no=%s, amount=%.2f, description=%s", checkIn.CheckInNo, req.Amount, req.Description)
	return checkIn, nil
}

func (s *checkInService) GetCheckInDetail(id uint) (*model.CheckIn, error) {
	checkIn, err := s.checkInRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取入住详情失败: id=%d, err=%v", id, err)
		return nil, errors.New("入住记录不存在")
	}
	return checkIn, nil
}

func (s *checkInService) ListCheckIns(req *dto.CheckInListRequest) ([]model.CheckIn, int64, error) {
	page := req.GetPage()
	pageSize := req.GetPageSize()

	checkIns, total, err := s.checkInRepo.List(
		page, pageSize,
		req.CheckInNo,
		req.BookingID,
		req.RoomID,
		req.GuestName,
		req.GuestPhone,
		req.Status,
		req.CheckInTime,
		req.CheckOutTime,
	)
	if err != nil {
		logger.Errorf("获取入住列表失败: %v", err)
		return nil, 0, errors.New("获取入住列表失败")
	}

	return checkIns, total, nil
}

func generateCheckInNo() (string, error) {
	now := time.Now().Format("20060102150405")
	random, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("CK%s%03d", now, random.Int64()), nil
}

func generateCheckInPaymentNo() (string, error) {
	now := time.Now().Format("20060102150405")
	random, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PAY%s%04d", now, random.Int64()), nil
}
