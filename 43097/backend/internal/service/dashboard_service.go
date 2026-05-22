package service

import (
	"errors"
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type DashboardService interface {
	GetRoomStatusBoard() (map[int][]dto.RoomStatusItem, error)
	GetDashboardStats() (*dto.DashboardStatsResponse, error)
	GetFloorRooms(floor int) (*dto.FloorRoomsResponse, error)
}

type dashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) DashboardService {
	return &dashboardService{db: db}
}

func (s *dashboardService) GetRoomStatusBoard() (map[int][]dto.RoomStatusItem, error) {
	var rooms []model.Room
	err := s.db.Preload("RoomType").Order("floor ASC, room_no ASC").Find(&rooms).Error
	if err != nil {
		logger.Errorf("获取房间列表失败: %v", err)
		return nil, errors.New("获取房间状态看板失败")
	}

	roomIDs := make([]uint, len(rooms))
	for i, room := range rooms {
		roomIDs[i] = room.ID
	}

	activeCheckIns := make(map[uint]*model.CheckIn)
	var checkIns []model.CheckIn
	err = s.db.Where("room_id IN ? AND status = ?", roomIDs, model.CheckInStatusActive).Find(&checkIns).Error
	if err != nil {
		logger.Errorf("获取活跃入住信息失败: %v", err)
	}
	for i := range checkIns {
		activeCheckIns[checkIns[i].RoomID] = &checkIns[i]
	}

	activeBookings := make(map[uint]*model.Booking)
	var bookings []model.Booking
	today := time.Now()
	err = s.db.Where(
		"room_id IN ? AND status IN ? AND check_in_date <= ? AND check_out_date > ?",
		roomIDs,
		[]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed},
		today,
		today,
	).Find(&bookings).Error
	if err != nil {
		logger.Errorf("获取活跃预订信息失败: %v", err)
	}
	for i := range bookings {
		activeBookings[bookings[i].RoomID] = &bookings[i]
	}

	result := make(map[int][]dto.RoomStatusItem)
	for _, room := range rooms {
		item := dto.RoomStatusItem{
			ID:         room.ID,
			RoomNo:     room.RoomNo,
			Floor:      room.Floor,
			RoomTypeID: room.RoomTypeID,
			Status:     room.Status,
			Price:      room.Price,
		}

		if room.RoomType != nil {
			item.RoomType = room.RoomType.Name
		}

		if checkIn, ok := activeCheckIns[room.ID]; ok {
			item.Guest = &dto.GuestInfo{
				Name:         checkIn.GuestName,
				Phone:        checkIn.GuestPhone,
				CheckInTime:  checkIn.CheckInTime.Format("2006-01-02 15:04:05"),
				CheckOutTime: checkIn.ExpectedCheckOut.Format("2006-01-02 15:04:05"),
			}
		}

		if booking, ok := activeBookings[room.ID]; ok {
			item.Booking = &dto.GuestInfo{
				Name:         booking.GuestName,
				Phone:        booking.GuestPhone,
				CheckInTime:  booking.CheckInDate.Format("2006-01-02"),
				CheckOutTime: booking.CheckOutDate.Format("2006-01-02"),
			}
		}

		result[room.Floor] = append(result[room.Floor], item)
	}

	return result, nil
}

func (s *dashboardService) GetDashboardStats() (*dto.DashboardStatsResponse, error) {
	stats := &dto.DashboardStatsResponse{}

	var totalRooms int64
	err := s.db.Model(&model.Room{}).Count(&totalRooms).Error
	if err != nil {
		logger.Errorf("获取房间总数失败: %v", err)
		return nil, errors.New("获取看板统计失败")
	}
	stats.RoomStatus.Total = totalRooms

	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	err = s.db.Model(&model.Room{}).Select("status, count(*) as count").Group("status").Scan(&statusCounts).Error
	if err != nil {
		logger.Errorf("获取房间状态统计失败: %v", err)
		return nil, errors.New("获取看板统计失败")
	}
	for _, sc := range statusCounts {
		switch sc.Status {
		case string(model.RoomStatusAvailable):
			stats.RoomStatus.Available = sc.Count
		case string(model.RoomStatusOccupied):
			stats.RoomStatus.Occupied = sc.Count
		case string(model.RoomStatusReserved):
			stats.RoomStatus.Reserved = sc.Count
		case string(model.RoomStatusMaintenance):
			stats.RoomStatus.Maintenance = sc.Count
		}
	}

	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	err = s.db.Model(&model.CheckIn{}).Where("check_in_time >= ? AND check_in_time < ?", todayStart, todayEnd).Count(&stats.TodayCheckIns).Error
	if err != nil {
		logger.Errorf("获取今日入住数失败: %v", err)
		return nil, errors.New("获取看板统计失败")
	}

	err = s.db.Model(&model.CheckIn{}).Where("actual_check_out >= ? AND actual_check_out < ?", todayStart, todayEnd).Count(&stats.TodayCheckOuts).Error
	if err != nil {
		logger.Errorf("获取今日退房数失败: %v", err)
		return nil, errors.New("获取看板统计失败")
	}

	err = s.db.Model(&model.CheckIn{}).Where("status = ?", model.CheckInStatusActive).Count(&stats.CurrentGuests).Error
	if err != nil {
		logger.Errorf("获取当前在住人数失败: %v", err)
		return nil, errors.New("获取看板统计失败")
	}

	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	nextMonth := monthStart.AddDate(0, 1, 0)

	type RevenueResult struct {
		Total float64
	}
	var revenueResult RevenueResult
	err = s.db.Model(&model.Payment{}).
		Select("COALESCE(SUM(CASE WHEN payment_type = ? THEN amount ELSE 0 END) - COALESCE(SUM(CASE WHEN payment_type = ? THEN amount ELSE 0 END), 0), 0) as total",
			model.PaymentTypePrepaid, model.PaymentTypeRefund).
		Where("status = ? AND created_at >= ? AND created_at < ?", model.PaymentStatusCompleted, monthStart, nextMonth).
		Scan(&revenueResult).Error
	if err != nil {
		logger.Errorf("获取本月营收失败: %v", err)
		return nil, errors.New("获取看板统计失败")
	}
	stats.MonthlyRevenue = revenueResult.Total

	var occupiedRoomNights float64
	err = s.db.Model(&model.CheckIn{}).
		Select("COALESCE(SUM(EXTRACT(EPOCH FROM (COALESCE(actual_check_out, NOW()) - check_in_time)) / 86400), 0)").
		Where("check_in_time >= ? AND check_in_time < ?", monthStart, nextMonth).
		Scan(&occupiedRoomNights).Error
	if err != nil {
		logger.Errorf("计算入住率失败: %v", err)
		return nil, errors.New("获取看板统计失败")
	}

	daysInMonth := float64(nextMonth.Sub(monthStart).Hours() / 24)
	totalRoomNights := float64(totalRooms) * daysInMonth
	if totalRoomNights > 0 {
		stats.MonthlyOccupancy = (occupiedRoomNights / totalRoomNights) * 100
	}

	return stats, nil
}

func (s *dashboardService) GetFloorRooms(floor int) (*dto.FloorRoomsResponse, error) {
	var rooms []model.Room
	err := s.db.Preload("RoomType").Where("floor = ?", floor).Order("room_no ASC").Find(&rooms).Error
	if err != nil {
		logger.Errorf("获取楼层房间列表失败: floor=%d, err=%v", floor, err)
		return nil, errors.New("获取楼层房间列表失败")
	}

	roomIDs := make([]uint, len(rooms))
	for i, room := range rooms {
		roomIDs[i] = room.ID
	}

	activeCheckIns := make(map[uint]*model.CheckIn)
	var checkIns []model.CheckIn
	err = s.db.Where("room_id IN ? AND status = ?", roomIDs, model.CheckInStatusActive).Find(&checkIns).Error
	if err != nil {
		logger.Errorf("获取活跃入住信息失败: %v", err)
	}
	for i := range checkIns {
		activeCheckIns[checkIns[i].RoomID] = &checkIns[i]
	}

	activeBookings := make(map[uint]*model.Booking)
	var bookings []model.Booking
	today := time.Now()
	err = s.db.Where(
		"room_id IN ? AND status IN ? AND check_in_date <= ? AND check_out_date > ?",
		roomIDs,
		[]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed},
		today,
		today,
	).Find(&bookings).Error
	if err != nil {
		logger.Errorf("获取活跃预订信息失败: %v", err)
	}
	for i := range bookings {
		activeBookings[bookings[i].RoomID] = &bookings[i]
	}

	roomItems := make([]dto.RoomStatusItem, 0, len(rooms))
	for _, room := range rooms {
		item := dto.RoomStatusItem{
			ID:         room.ID,
			RoomNo:     room.RoomNo,
			Floor:      room.Floor,
			RoomTypeID: room.RoomTypeID,
			Status:     room.Status,
			Price:      room.Price,
		}

		if room.RoomType != nil {
			item.RoomType = room.RoomType.Name
		}

		if checkIn, ok := activeCheckIns[room.ID]; ok {
			item.Guest = &dto.GuestInfo{
				Name:         checkIn.GuestName,
				Phone:        checkIn.GuestPhone,
				CheckInTime:  checkIn.CheckInTime.Format("2006-01-02 15:04:05"),
				CheckOutTime: checkIn.ExpectedCheckOut.Format("2006-01-02 15:04:05"),
			}
		}

		if booking, ok := activeBookings[room.ID]; ok {
			item.Booking = &dto.GuestInfo{
				Name:         booking.GuestName,
				Phone:        booking.GuestPhone,
				CheckInTime:  booking.CheckInDate.Format("2006-01-02"),
				CheckOutTime: booking.CheckOutDate.Format("2006-01-02"),
			}
		}

		roomItems = append(roomItems, item)
	}

	return &dto.FloorRoomsResponse{
		Floor: floor,
		Rooms: roomItems,
	}, nil
}
