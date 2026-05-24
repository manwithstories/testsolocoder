package service

import (
	"encoding/json"
	"fmt"
	"time"

	"music-platform/internal/model"
	"music-platform/internal/repository"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/redis"
	"music-platform/pkg/utils"
)

type EventService struct {
	eventRepo *repository.EventRepository
	userRepo  *repository.UserRepository
}

func NewEventService() *EventService {
	return &EventService{
		eventRepo: repository.NewEventRepository(),
		userRepo:  repository.NewUserRepository(),
	}
}

type CreateEventRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	CoverURL     string  `json:"cover_url"`
	Venue        string  `json:"venue" binding:"required"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	StartTime    string  `json:"start_time" binding:"required"`
	EndTime      string  `json:"end_time" binding:"required"`
	DoorTime     string  `json:"door_time"`
	TicketPrice  float64 `json:"ticket_price" binding:"required"`
	TotalTickets int     `json:"total_tickets" binding:"required"`
	MaxPerUser   int     `json:"max_per_user"`
	SeatMap      string  `json:"seat_map"`
}

type UpdateEventRequest struct {
	Title        *string  `json:"title"`
	Description  *string  `json:"description"`
	CoverURL     *string  `json:"cover_url"`
	Venue        *string  `json:"venue"`
	Address      *string  `json:"address"`
	City         *string  `json:"city"`
	Longitude    *float64 `json:"longitude"`
	Latitude     *float64 `json:"latitude"`
	StartTime    *string  `json:"start_time"`
	EndTime      *string  `json:"end_time"`
	DoorTime     *string  `json:"door_time"`
	TicketPrice  *float64 `json:"ticket_price"`
	TotalTickets *int     `json:"total_tickets"`
	MaxPerUser   *int     `json:"max_per_user"`
	SeatMap      *string  `json:"seat_map"`
	Status       *int     `json:"status"`
}

type PurchaseTicketRequest struct {
	EventID  uint     `json:"event_id" binding:"required"`
	Quantity int      `json:"quantity" binding:"required,min=1"`
	Seats    []SeatInfo `json:"seats"`
}

type SeatInfo struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (s *EventService) CreateEvent(userID uint, req *CreateEventRequest) (*model.Event, error) {
	artistInfo, err := s.userRepo.FindArtistInfoByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		return nil, apperrors.NewAppError(4001, "开始时间格式错误")
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		return nil, apperrors.NewAppError(4002, "结束时间格式错误")
	}

	event := &model.Event{
		UserID:       userID,
		ArtistID:     artistInfo.ID,
		Title:        req.Title,
		Description:  req.Description,
		CoverURL:     req.CoverURL,
		Venue:        req.Venue,
		Address:      req.Address,
		City:         req.City,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		StartTime:    startTime,
		EndTime:      endTime,
		TicketPrice:  req.TicketPrice,
		TotalTickets: req.TotalTickets,
		MaxPerUser:   req.MaxPerUser,
		SeatMap:      req.SeatMap,
		Status:       model.EventStatusDraft,
	}

	if req.DoorTime != "" {
		doorTime, _ := time.Parse("2006-01-02 15:04:05", req.DoorTime)
		event.DoorTime = &doorTime
	}

	err = s.eventRepo.Create(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) GetEventByID(id uint) (*model.Event, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.ErrEventNotFound
	}

	_ = s.eventRepo.UpdateViewCount(id)

	return event, nil
}

func (s *EventService) ListEvents(page, pageSize int, keyword string, artistID uint, city string, status int) ([]model.Event, int64, error) {
	return s.eventRepo.List(page, pageSize, keyword, artistID, city, status)
}

func (s *EventService) UpdateEvent(eventID uint, userID uint, req *UpdateEventRequest) error {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return apperrors.ErrEventNotFound
	}

	if event.UserID != userID {
		return apperrors.ErrForbidden
	}

	updates := map[string]interface{}{}

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.CoverURL != nil {
		updates["cover_url"] = *req.CoverURL
	}
	if req.Venue != nil {
		updates["venue"] = *req.Venue
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.StartTime != nil {
		startTime, err := time.Parse("2006-01-02 15:04:05", *req.StartTime)
		if err == nil {
			updates["start_time"] = startTime
		}
	}
	if req.EndTime != nil {
		endTime, err := time.Parse("2006-01-02 15:04:05", *req.EndTime)
		if err == nil {
			updates["end_time"] = endTime
		}
	}
	if req.TicketPrice != nil {
		updates["ticket_price"] = *req.TicketPrice
	}
	if req.TotalTickets != nil {
		updates["total_tickets"] = *req.TotalTickets
	}
	if req.MaxPerUser != nil {
		updates["max_per_user"] = *req.MaxPerUser
	}
	if req.SeatMap != nil {
		updates["seat_map"] = *req.SeatMap
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	return s.eventRepo.UpdateEventInfo(eventID, updates)
}

func (s *EventService) DeleteEvent(eventID uint, userID uint) error {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return apperrors.ErrEventNotFound
	}

	if event.UserID != userID {
		return apperrors.ErrForbidden
	}

	return s.eventRepo.Delete(eventID)
}

func (s *EventService) PublishEvent(eventID uint, userID uint) error {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return apperrors.ErrEventNotFound
	}

	if event.UserID != userID {
		return apperrors.ErrForbidden
	}

	now := time.Now()
	return s.eventRepo.UpdateEventInfo(eventID, map[string]interface{}{
		"status":       model.EventStatusPublished,
		"published_at": &now,
	})
}

func (s *EventService) PurchaseTicket(userID uint, req *PurchaseTicketRequest) (*model.Order, error) {
	event, err := s.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, apperrors.ErrEventNotFound
	}

	if event.Status != model.EventStatusPublished {
		return nil, apperrors.NewAppError(4004, "演出未发布")
	}

	remainingTickets := event.TotalTickets - event.SoldTickets
	if remainingTickets < req.Quantity {
		return nil, apperrors.ErrTicketSoldOut
	}

	if req.Quantity > event.MaxPerUser {
		return nil, apperrors.NewAppError(4005, fmt.Sprintf("每人最多购买%d张票", event.MaxPerUser))
	}

	lockKey := fmt.Sprintf("lock:event:%d", req.EventID)
	lockValue := utils.GenerateUUID()
	locked, err := redis.AcquireLock(lockKey, lockValue, 10*time.Second)
	if err != nil || !locked {
		return nil, apperrors.NewAppError(4006, "系统繁忙，请稍后再试")
	}
	defer redis.ReleaseLock(lockKey, lockValue)

	event, err = s.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, apperrors.ErrEventNotFound
	}

	remainingTickets = event.TotalTickets - event.SoldTickets
	if remainingTickets < req.Quantity {
		return nil, apperrors.ErrTicketSoldOut
	}

	user, _ := s.userRepo.FindByID(userID)
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	orderNo := fmt.Sprintf("T%s%d", time.Now().Format("20060102150405"), userID)

	totalAmount := float64(req.Quantity) * event.TicketPrice

	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		EventID:     req.EventID,
		ArtistID:    event.ArtistID,
		TotalAmount: totalAmount,
		Quantity:    req.Quantity,
		Status:      model.OrderStatusPending,
	}

	err = s.eventRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	var tickets []model.Ticket
	for i := 0; i < req.Quantity; i++ {
		seatLabel := ""
		seatRow := 0
		seatCol := 0

		if i < len(req.Seats) {
			seatRow = req.Seats[i].Row
			seatCol = req.Seats[i].Col
			seatLabel = fmt.Sprintf("%d排%d座", seatRow, seatCol)
		}

		ticket := model.Ticket{
			OrderID:   order.ID,
			UserID:    userID,
			EventID:   req.EventID,
			SeatRow:   seatRow,
			SeatCol:   seatCol,
			SeatLabel: seatLabel,
			QRCode:    utils.GenerateUUID(),
			Status:    model.TicketStatusPending,
		}
		tickets = append(tickets, ticket)
	}

	err = s.eventRepo.BatchCreateTickets(tickets)
	if err != nil {
		return nil, err
	}

	err = s.eventRepo.UpdateSoldTickets(req.EventID, req.Quantity)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *EventService) GetOrderByID(id uint) (*model.Order, error) {
	order, err := s.eventRepo.FindOrderByID(id)
	if err != nil {
		return nil, apperrors.NewAppError(4007, "订单不存在")
	}
	return order, nil
}

func (s *EventService) GetOrdersByUser(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	return s.eventRepo.GetOrdersByUserID(userID, page, pageSize)
}

func (s *EventService) GetOrdersByArtist(artistID uint, page, pageSize int) ([]model.Order, int64, error) {
	return s.eventRepo.GetOrdersByArtistID(artistID, page, pageSize)
}

func (s *EventService) GetTicketsByUser(userID uint, page, pageSize int) ([]model.Ticket, int64, error) {
	return s.eventRepo.GetTicketsByUserID(userID, page, pageSize)
}

func (s *EventService) UseTicket(ticketID uint, userID uint) error {
	ticket, err := s.eventRepo.FindTicketByID(ticketID)
	if err != nil {
		return apperrors.NewAppError(4008, "票不存在")
	}

	if ticket.UserID != userID {
		return apperrors.ErrForbidden
	}

	if ticket.Status != model.TicketStatusPaid {
		return apperrors.NewAppError(4009, "票状态异常")
	}

	now := time.Now()
	return s.eventRepo.UpdateEventInfo(ticketID, map[string]interface{}{
		"status":  model.TicketStatusUsed,
		"used_at": &now,
	})
}

func (s *EventService) GetEventStats(eventID uint) (map[string]interface{}, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, apperrors.ErrEventNotFound
	}

	return map[string]interface{}{
		"event_id":       event.ID,
		"total_tickets":  event.TotalTickets,
		"sold_tickets":   event.SoldTickets,
		"remaining":      event.TotalTickets - event.SoldTickets,
		"revenue":        float64(event.SoldTickets) * event.TicketPrice,
		"view_count":     event.ViewCount,
		"like_count":     event.LikeCount,
	}, nil
}

func (s *EventService) GetSeatAvailability(eventID uint) (map[string]interface{}, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, apperrors.ErrEventNotFound
	}

	var seatMap map[string]interface{}
	if event.SeatMap != "" {
		_ = json.Unmarshal([]byte(event.SeatMap), &seatMap)
	}

	remainingTickets := event.TotalTickets - event.SoldTickets

	return map[string]interface{}{
		"event_id":         event.ID,
		"total_tickets":    event.TotalTickets,
		"sold_tickets":     event.SoldTickets,
		"remaining_tickets": remainingTickets,
		"seat_map":         seatMap,
	}, nil
}
