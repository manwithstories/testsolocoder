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

const (
	BookingStatusNoShow model.BookingStatus = "no_show"
)

type BookingService interface {
	CreateBooking(req *dto.BookingCreateRequest) (*model.Booking, error)
	ConfirmBooking(id uint, req *dto.BookingConfirmRequest) (*model.Booking, error)
	CancelBooking(id uint, req *dto.BookingCancelRequest) (*model.Booking, error)
	UpdateBooking(id uint, req *dto.BookingUpdateRequest) (*model.Booking, error)
	CalculatePrice(req *dto.BookingPriceCalculationRequest) (*dto.BookingPriceCalculationResponse, error)
	GetBookingDetail(id uint) (*model.Booking, error)
	ListBookings(req *dto.BookingListRequest) ([]model.Booking, int64, error)
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	roomRepo    repository.RoomRepository
	db          *gorm.DB
}

func NewBookingService(bookingRepo repository.BookingRepository, roomRepo repository.RoomRepository, db *gorm.DB) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		db:          db,
	}
}

func (s *bookingService) CreateBooking(req *dto.BookingCreateRequest) (*model.Booking, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, errors.New("创建预订失败")
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
		return nil, errors.New("创建预订失败")
	}

	if room.Status != model.RoomStatusAvailable {
		tx.Rollback()
		logger.Warnf("房间状态不可用: room_id=%d, status=%s", req.RoomID, room.Status)
		return nil, errors.New("房间状态不可用")
	}

	var count int64
	query := tx.Model(&model.Booking{}).Where(
		"room_id = ? AND status IN ? AND check_in_date < ? AND check_out_date > ?",
		req.RoomID,
		[]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed},
		req.CheckOutDate,
		req.CheckInDate,
	)
	err = query.Count(&count).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("检查房间可用性失败: %v", err)
		return nil, errors.New("创建预订失败")
	}
	if count > 0 {
		tx.Rollback()
		logger.Warnf("房间在指定日期已被预订: room_id=%d", req.RoomID)
		return nil, errors.New("房间在指定日期已被预订")
	}

	days := int(req.CheckOutDate.Sub(req.CheckInDate).Hours() / 24)
	if days <= 0 {
		tx.Rollback()
		logger.Warnf("入住天数无效: days=%d", days)
		return nil, errors.New("入住天数必须大于0")
	}

	originalPrice := room.Price * float64(days)
	discountRate := 1.0
	totalPrice := originalPrice

	if req.MemberID != nil {
		var member model.Member
		err = tx.Preload("Level").First(&member, *req.MemberID).Error
		if err == nil && member.Status == model.MemberStatusActive && member.Level != nil {
			discountRate = member.Level.DiscountRate
			totalPrice = originalPrice * discountRate
		}
	}

	bookingNo, err := generateBookingNo()
	if err != nil {
		tx.Rollback()
		logger.Errorf("生成预订号失败: %v", err)
		return nil, errors.New("创建预订失败")
	}

	cancelDeadline := req.CheckInDate.Add(-24 * time.Hour)

	booking := &model.Booking{
		BookingNo:      bookingNo,
		RoomID:         req.RoomID,
		MemberID:       req.MemberID,
		GuestName:      req.GuestName,
		GuestPhone:     req.GuestPhone,
		GuestIDCard:    req.GuestIDCard,
		CheckInDate:    req.CheckInDate,
		CheckOutDate:   req.CheckOutDate,
		Days:           days,
		TotalPrice:     totalPrice,
		Status:         model.BookingStatusPending,
		PaidAmount:     0,
		Remarks:        req.Remarks,
		CancelDeadline: &cancelDeadline,
	}

	err = tx.Create(booking).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("创建预订失败: %v", err)
		return nil, errors.New("创建预订失败")
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("提交事务失败: %v", err)
		return nil, errors.New("创建预订失败")
	}

	logger.Infof("预订创建成功: booking_no=%s, room_id=%d", bookingNo, req.RoomID)
	return booking, nil
}

func (s *bookingService) ConfirmBooking(id uint, req *dto.BookingConfirmRequest) (*model.Booking, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		logger.Errorf("确认预订失败，预订不存在: id=%d, err=%v", id, err)
		return nil, errors.New("预订不存在")
	}

	if booking.Status != model.BookingStatusPending {
		logger.Warnf("预订状态不允许确认: id=%d, status=%s", id, booking.Status)
		return nil, errors.New("只有待确认状态的预订才能确认")
	}

	booking.Status = model.BookingStatusConfirmed
	if req.PaidAmount > 0 {
		booking.PaidAmount = req.PaidAmount
	}

	err = s.bookingRepo.Update(booking)
	if err != nil {
		logger.Errorf("确认预订失败: %v", err)
		return nil, errors.New("确认预订失败")
	}

	logger.Infof("预订确认成功: booking_no=%s", booking.BookingNo)
	return booking, nil
}

func (s *bookingService) CancelBooking(id uint, req *dto.BookingCancelRequest) (*model.Booking, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		logger.Errorf("取消预订失败，预订不存在: id=%d, err=%v", id, err)
		return nil, errors.New("预订不存在")
	}

	if booking.Status == model.BookingStatusCancelled {
		logger.Warnf("预订已取消: id=%d", id)
		return booking, nil
	}

	if booking.Status != model.BookingStatusPending && booking.Status != model.BookingStatusConfirmed {
		logger.Warnf("预订状态不允许取消: id=%d, status=%s", id, booking.Status)
		return nil, errors.New("该预订状态不允许取消")
	}

	now := time.Now()
	if booking.CancelDeadline != nil && now.After(*booking.CancelDeadline) {
		logger.Warnf("已超过免费取消期限: id=%d", id)
		return nil, errors.New("已超过免费取消期限，如需取消请联系客服")
	}

	booking.Status = model.BookingStatusCancelled
	if req.Reason != "" {
		booking.Remarks = booking.Remarks + " | 取消原因: " + req.Reason
	}

	err = s.bookingRepo.Update(booking)
	if err != nil {
		logger.Errorf("取消预订失败: %v", err)
		return nil, errors.New("取消预订失败")
	}

	logger.Infof("预订取消成功: booking_no=%s", booking.BookingNo)
	return booking, nil
}

func (s *bookingService) UpdateBooking(id uint, req *dto.BookingUpdateRequest) (*model.Booking, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		logger.Errorf("更新预订失败，预订不存在: id=%d, err=%v", id, err)
		return nil, errors.New("预订不存在")
	}

	if booking.Status == model.BookingStatusCancelled || booking.Status == model.BookingStatusCompleted {
		logger.Warnf("预订状态不允许修改: id=%d, status=%s", id, booking.Status)
		return nil, errors.New("该预订状态不允许修改")
	}

	if req.RoomID != nil && *req.RoomID != booking.RoomID {
		room, err := s.roomRepo.GetByID(*req.RoomID)
		if err != nil {
			logger.Errorf("更新预订失败，房间不存在: room_id=%d, err=%v", *req.RoomID, err)
			return nil, errors.New("房间不存在")
		}
		if room.Status != model.RoomStatusAvailable {
			logger.Warnf("房间状态不可用: room_id=%d, status=%s", *req.RoomID, room.Status)
			return nil, errors.New("房间状态不可用")
		}

		checkIn := booking.CheckInDate
		checkOut := booking.CheckOutDate
		if !req.CheckInDate.IsZero() {
			checkIn = req.CheckInDate
		}
		if !req.CheckOutDate.IsZero() {
			checkOut = req.CheckOutDate
		}

		available, err := s.bookingRepo.CheckRoomAvailability(*req.RoomID, checkIn, checkOut, &id)
		if err != nil {
			logger.Errorf("检查房间可用性失败: %v", err)
			return nil, errors.New("更新预订失败")
		}
		if !available {
			logger.Warnf("房间在指定日期已被预订: room_id=%d", *req.RoomID)
			return nil, errors.New("房间在指定日期已被预订")
		}

		booking.RoomID = *req.RoomID
	}

	checkIn := booking.CheckInDate
	checkOut := booking.CheckOutDate
	if !req.CheckInDate.IsZero() {
		checkIn = req.CheckInDate
	}
	if !req.CheckOutDate.IsZero() {
		checkOut = req.CheckOutDate
	}

	if checkIn != booking.CheckInDate || checkOut != booking.CheckOutDate {
		available, err := s.bookingRepo.CheckRoomAvailability(booking.RoomID, checkIn, checkOut, &id)
		if err != nil {
			logger.Errorf("检查房间可用性失败: %v", err)
			return nil, errors.New("更新预订失败")
		}
		if !available {
			logger.Warnf("房间在新日期范围内已被预订: room_id=%d", booking.RoomID)
			return nil, errors.New("房间在新日期范围内已被预订")
		}

		days := int(checkOut.Sub(checkIn).Hours() / 24)
		if days <= 0 {
			logger.Warnf("入住天数无效: days=%d", days)
			return nil, errors.New("入住天数必须大于0")
		}

		room, err := s.roomRepo.GetByID(booking.RoomID)
		if err != nil {
			logger.Errorf("获取房间信息失败: %v", err)
			return nil, errors.New("更新预订失败")
		}

		originalPrice := room.Price * float64(days)
		discountRate := 1.0
		totalPrice := originalPrice

		if booking.MemberID != nil {
			var member model.Member
			err = s.db.Preload("Level").First(&member, *booking.MemberID).Error
			if err == nil && member.Status == model.MemberStatusActive && member.Level != nil {
				discountRate = member.Level.DiscountRate
				totalPrice = originalPrice * discountRate
			}
		}

		booking.CheckInDate = checkIn
		booking.CheckOutDate = checkOut
		booking.Days = days
		booking.TotalPrice = totalPrice
		cancelDeadline := checkIn.Add(-24 * time.Hour)
		booking.CancelDeadline = &cancelDeadline
	}

	if req.GuestName != "" {
		booking.GuestName = req.GuestName
	}
	if req.GuestPhone != "" {
		booking.GuestPhone = req.GuestPhone
	}
	if req.GuestIDCard != "" {
		booking.GuestIDCard = req.GuestIDCard
	}
	if req.Remarks != "" {
		booking.Remarks = req.Remarks
	}

	err = s.bookingRepo.Update(booking)
	if err != nil {
		logger.Errorf("更新预订失败: %v", err)
		return nil, errors.New("更新预订失败")
	}

	logger.Infof("预订更新成功: booking_no=%s", booking.BookingNo)
	return booking, nil
}

func (s *bookingService) CalculatePrice(req *dto.BookingPriceCalculationRequest) (*dto.BookingPriceCalculationResponse, error) {
	room, err := s.roomRepo.GetByID(req.RoomID)
	if err != nil {
		logger.Errorf("计算价格失败，房间不存在: room_id=%d, err=%v", req.RoomID, err)
		return nil, errors.New("房间不存在")
	}

	days := int(req.CheckOutDate.Sub(req.CheckInDate).Hours() / 24)
	if days <= 0 {
		logger.Warnf("入住天数无效: days=%d", days)
		return nil, errors.New("入住天数必须大于0")
	}

	originalPrice := room.Price * float64(days)
	discountRate := 1.0
	totalPrice := originalPrice

	if req.MemberID != nil {
		var member model.Member
		err = s.db.Preload("Level").First(&member, *req.MemberID).Error
		if err == nil && member.Status == model.MemberStatusActive && member.Level != nil {
			discountRate = member.Level.DiscountRate
			totalPrice = originalPrice * discountRate
		}
	}

	return &dto.BookingPriceCalculationResponse{
		RoomPrice:     room.Price,
		Days:          days,
		OriginalPrice: originalPrice,
		DiscountRate:  discountRate,
		Discount:      originalPrice - totalPrice,
		TotalPrice:    totalPrice,
	}, nil
}

func (s *bookingService) GetBookingDetail(id uint) (*model.Booking, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取预订详情失败: id=%d, err=%v", id, err)
		return nil, errors.New("预订不存在")
	}
	return booking, nil
}

func (s *bookingService) ListBookings(req *dto.BookingListRequest) ([]model.Booking, int64, error) {
	page := req.GetPage()
	pageSize := req.GetPageSize()

	bookings, total, err := s.bookingRepo.List(
		page, pageSize,
		req.BookingNo,
		req.RoomID,
		req.MemberID,
		req.GuestName,
		req.GuestPhone,
		req.Status,
		req.CheckInDate,
		req.CheckOutDate,
	)
	if err != nil {
		logger.Errorf("获取预订列表失败: %v", err)
		return nil, 0, errors.New("获取预订列表失败")
	}

	return bookings, total, nil
}

func generateBookingNo() (string, error) {
	now := time.Now().Format("20060102150405")
	random, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("BK%s%04d", now, random.Int64()), nil
}
