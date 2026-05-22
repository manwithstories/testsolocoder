package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"meeting-room/internal/models"
	"meeting-room/internal/repositories"
	"meeting-room/internal/utils"

	"gorm.io/gorm"
)

type BookingService struct {
	bookingRepo      *repositories.BookingRepository
	roomRepo         *repositories.RoomRepository
	userRepo         *repositories.UserRepository
	notificationRepo *repositories.NotificationRepository
}

func NewBookingService() *BookingService {
	return &BookingService{
		bookingRepo:      repositories.NewBookingRepository(),
		roomRepo:         repositories.NewRoomRepository(),
		userRepo:         repositories.NewUserRepository(),
		notificationRepo: repositories.NewNotificationRepository(),
	}
}

type CreateBookingRequest struct {
	RoomID         uint              `json:"room_id" binding:"required"`
	Title          string            `json:"title" binding:"required"`
	Description    string            `json:"description"`
	StartTime      time.Time         `json:"start_time" binding:"required"`
	EndTime        time.Time         `json:"end_time" binding:"required"`
	RecurrenceType models.RecurrenceType `json:"recurrence_type"`
	RecurrenceEnd  *time.Time        `json:"recurrence_end"`
	Attendees      []string          `json:"attendees"`
}

type RescheduleRequest struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type CancelRequest struct {
	Reason string `json:"reason"`
}

func (s *BookingService) CreateBooking(req *CreateBookingRequest, userID uint) (*models.Booking, error) {
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("预订时间不能早于当前时间")
	}

	duration := req.EndTime.Sub(req.StartTime)
	if duration < 15*time.Minute {
		return nil, errors.New("预订时长不能少于15分钟")
	}

	room, err := s.roomRepo.FindByID(req.RoomID)
	if err != nil {
		return nil, errors.New("会议室不存在")
	}

	if room.Status != models.RoomStatusActive {
		return nil, errors.New("会议室当前不可用")
	}

	if !s.isWithinAvailableHours(req.StartTime, req.EndTime, room) {
		return nil, fmt.Errorf("预订时间超出会议室可用时段 (%s - %s)", room.AvailableStart, room.AvailableEnd)
	}

	if req.RecurrenceType != "" && req.RecurrenceType != models.RecurrenceNone {
		if req.RecurrenceEnd == nil {
			return nil, errors.New("周期预订必须指定结束日期")
		}
		return s.createRecurringBookings(req, userID, room)
	}

	return s.createSingleBooking(req, userID, room)
}

func (s *BookingService) isWithinAvailableHours(startTime, endTime time.Time, room *models.Room) bool {
	availableStart, _ := time.Parse("15:04", room.AvailableStart)
	availableEnd, _ := time.Parse("15:04", room.AvailableEnd)

	bookingStart, _ := time.Parse("15:04", startTime.Format("15:04"))
	bookingEnd, _ := time.Parse("15:04", endTime.Format("15:04"))

	return !bookingStart.Before(availableStart) && !bookingEnd.After(availableEnd)
}

func (s *BookingService) createSingleBooking(req *CreateBookingRequest, userID uint, room *models.Room) (*models.Booking, error) {
	conflicts, err := s.bookingRepo.CheckConflict(req.RoomID, req.StartTime, req.EndTime, 0)
	if err != nil {
		return nil, err
	}
	if len(conflicts) > 0 {
		return nil, errors.New("该时间段会议室已被预订")
	}

	booking := &models.Booking{
		RoomID:         req.RoomID,
		UserID:         userID,
		Title:          req.Title,
		Description:    req.Description,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		RecurrenceType: req.RecurrenceType,
		Status:         models.BookingStatusConfirmed,
		TotalPrice:     calculatePrice(req.StartTime, req.EndTime, room.PricePerHour),
	}

	if len(req.Attendees) > 0 {
		attendeesJSON, _ := json.Marshal(req.Attendees)
		booking.Attendees = string(attendeesJSON)
	}

	err = s.bookingRepo.Create(booking)
	if err != nil {
		return nil, err
	}

	s.cacheRoomStatus(room.ID, req.StartTime, req.EndTime, booking.ID)

	s.enqueueNotification(booking, room, userID, models.NotificationTypeBookingConfirmation)

	return booking, nil
}

func (s *BookingService) createRecurringBookings(req *CreateBookingRequest, userID uint, room *models.Room) (*models.Booking, error) {
	var dates []time.Time
	current := req.StartTime

	for !current.After(*req.RecurrenceEnd) {
		dates = append(dates, current)
		switch req.RecurrenceType {
		case models.RecurrenceDaily:
			current = current.AddDate(0, 0, 1)
		case models.RecurrenceWeekly:
			current = current.AddDate(0, 0, 7)
		case models.RecurrenceBiweekly:
			current = current.AddDate(0, 0, 14)
		case models.RecurrenceMonthly:
			current = current.AddDate(0, 1, 0)
		}
	}

	duration := req.EndTime.Sub(req.StartTime)
	var firstBooking *models.Booking

	err := utils.DB.Transaction(func(tx *gorm.DB) error {
		for _, date := range dates {
			start := date
			end := start.Add(duration)

			conflicts, err := s.bookingRepo.CheckConflict(req.RoomID, start, end, 0)
			if err != nil {
				return err
			}
			if len(conflicts) > 0 {
				return fmt.Errorf("日期 %s 时段会议室已被预订", start.Format("2006-01-02 15:04"))
			}

			booking := &models.Booking{
				RoomID:         req.RoomID,
				UserID:         userID,
				Title:          req.Title,
				Description:    req.Description,
				StartTime:      start,
				EndTime:        end,
				RecurrenceType: req.RecurrenceType,
				Status:         models.BookingStatusConfirmed,
				TotalPrice:     calculatePrice(start, end, room.PricePerHour),
			}

			if firstBooking != nil {
				parentID := firstBooking.ID
				booking.ParentBookingID = &parentID
			}

			if len(req.Attendees) > 0 {
				attendeesJSON, _ := json.Marshal(req.Attendees)
				booking.Attendees = string(attendeesJSON)
			}

			if err := s.bookingRepo.CreateInTx(tx, booking); err != nil {
				return err
			}

			if firstBooking == nil {
				firstBooking = booking
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if firstBooking != nil {
		s.enqueueNotification(firstBooking, room, userID, models.NotificationTypeBookingConfirmation)
	}

	return firstBooking, nil
}

func (s *BookingService) GetBooking(id uint) (*models.Booking, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("预订不存在")
	}
	return booking, nil
}

func (s *BookingService) ListBookings(page, pageSize int, userID uint, roomID uint, status int) ([]models.Booking, int64, error) {
	return s.bookingRepo.List(page, pageSize, userID, roomID, status)
}

func (s *BookingService) CancelBooking(id uint, userID uint, reason string) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预订不存在")
	}

	if booking.Status == models.BookingStatusCancelled {
		return errors.New("预订已被取消")
	}

	if booking.UserID != userID {
		return errors.New("无权取消此预订")
	}

	if booking.StartTime.Before(time.Now()) {
		return errors.New("会议已开始，无法取消")
	}

	now := time.Now()
	booking.Status = models.BookingStatusCancelled
	booking.CancelledBy = userID
	booking.CancelledAt = &now
	booking.CancelReason = reason

	err = s.bookingRepo.Update(booking)
	if err != nil {
		return err
	}

	s.removeRoomStatusCache(booking.RoomID, booking.ID)

	admins, _, _ := s.userRepo.List(1, 100, "admin")
	for _, admin := range admins {
		notification := &models.Notification{
			UserID:    admin.ID,
			BookingID: &booking.ID,
			Type:      models.NotificationTypeBookingCancellation,
			Subject:   fmt.Sprintf("预订已取消: %s", booking.Title),
			Content:   fmt.Sprintf("用户取消了会议室预订: %s, 时间: %s - %s", booking.Title, booking.StartTime.Format("2006-01-02 15:04"), booking.EndTime.Format("15:04")),
		}
		s.notificationRepo.Create(notification)
	}

	return nil
}

func (s *BookingService) RescheduleBooking(id uint, userID uint, req *RescheduleRequest) (*models.Booking, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("预订不存在")
	}

	if booking.Status == models.BookingStatusCancelled {
		return nil, errors.New("预订已被取消，无法改期")
	}

	if booking.UserID != userID {
		return nil, errors.New("无权修改此预订")
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("新预订时间不能早于当前时间")
	}

	conflicts, err := s.bookingRepo.CheckConflict(booking.RoomID, req.StartTime, req.EndTime, booking.ID)
	if err != nil {
		return nil, err
	}
	if len(conflicts) > 0 {
		return nil, errors.New("该时间段会议室已被预订")
	}

	s.removeRoomStatusCache(booking.RoomID, booking.ID)

	oldStart := booking.StartTime
	oldEnd := booking.EndTime
	booking.StartTime = req.StartTime
	booking.EndTime = req.EndTime

	room, _ := s.roomRepo.FindByID(booking.RoomID)
	booking.TotalPrice = calculatePrice(req.StartTime, req.EndTime, room.PricePerHour)

	err = s.bookingRepo.Update(booking)
	if err != nil {
		return nil, err
	}

	s.cacheRoomStatus(booking.RoomID, req.StartTime, req.EndTime, booking.ID)

	notification := &models.Notification{
		UserID:    userID,
		BookingID: &booking.ID,
		Type:      models.NotificationTypeBookingModification,
		Subject:   fmt.Sprintf("预订已改期: %s", booking.Title),
		Content:   fmt.Sprintf("您的会议室预订已从 %s - %s 改期为 %s - %s", oldStart.Format("2006-01-02 15:04"), oldEnd.Format("15:04"), req.StartTime.Format("2006-01-02 15:04"), req.EndTime.Format("15:04")),
	}
	s.notificationRepo.Create(notification)

	return booking, nil
}

func (s *BookingService) GetBookingsByDateRange(startTime, endTime time.Time, roomID uint, floor string) ([]models.Booking, error) {
	return s.bookingRepo.GetByDateRange(startTime, endTime, roomID, floor)
}

func (s *BookingService) GetBookingStats(startTime, endTime time.Time, department string) ([]models.Booking, error) {
	return s.bookingRepo.GetStats(startTime, endTime, department)
}

func (s *BookingService) ApproveBooking(id uint) error {
	return s.bookingRepo.UpdateStatus(id, models.BookingStatusConfirmed)
}

func (s *BookingService) CompleteBooking(id uint) error {
	return s.bookingRepo.UpdateStatus(id, models.BookingStatusCompleted)
}

func (s *BookingService) cacheRoomStatus(roomID uint, startTime, endTime time.Time, bookingID uint) {
	cacheKey := fmt.Sprintf("room_status:%d", roomID)
	statusKey := fmt.Sprintf("%d:%d:%d", bookingID, startTime.Unix(), endTime.Unix())
	utils.RedisHSet(cacheKey, statusKey, "1")
	utils.RedisExpire(cacheKey, 24*time.Hour)
}

func (s *BookingService) removeRoomStatusCache(roomID uint, bookingID uint) {
	cacheKey := fmt.Sprintf("room_status:%d", roomID)
	data, err := utils.RedisHGetAll(cacheKey)
	if err != nil {
		return
	}
	for key := range data {
		if len(key) > 0 && key[:len(fmt.Sprintf("%d:", bookingID))] == fmt.Sprintf("%d:", bookingID) {
			utils.RedisClient.HDel(utils.Ctx, cacheKey, key)
		}
	}
}

func (s *BookingService) enqueueNotification(booking *models.Booking, room *models.Room, userID uint, notifType models.NotificationType) {
	notification := &models.Notification{
		UserID:    userID,
		BookingID: &booking.ID,
		Type:      notifType,
		Subject:   fmt.Sprintf("预订确认: %s", booking.Title),
		Content:   fmt.Sprintf("您的会议室预订已确认！\n会议室: %s\n时间: %s - %s", room.Name, booking.StartTime.Format("2006-01-02 15:04"), booking.EndTime.Format("15:04")),
	}
	s.notificationRepo.Create(notification)

	utils.RedisLPush("notification_queue", notification.ID)
}

func calculatePrice(startTime, endTime time.Time, pricePerHour float64) float64 {
	duration := endTime.Sub(startTime).Hours()
	return float64(int(duration*pricePerHour*100)) / 100
}
