package service

import (
	"errors"
	"time"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/internal/repository"
)

type OrderService struct {
	orderRepo  *repository.OrderRepository
	venueRepo  *repository.VenueRepository
	deviceRepo *repository.DeviceRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:  repository.NewOrderRepository(),
		venueRepo:  repository.NewVenueRepository(),
		deviceRepo: repository.NewDeviceRepository(),
	}
}

func (s *OrderService) Create(req *dto.CreateOrderRequest, userID uint) (*model.Order, error) {
	startTime, err := parseDateTime(req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start time format, use YYYY-MM-DD HH:MM")
	}

	endTime, err := parseDateTime(req.EndTime)
	if err != nil {
		return nil, errors.New("invalid end time format, use YYYY-MM-DD HH:MM")
	}

	if endTime.Before(startTime) {
		return nil, errors.New("end time must be after start time")
	}

	if startTime.Before(time.Now()) {
		return nil, errors.New("cannot book past time")
	}

	totalHours := endTime.Sub(startTime).Hours()

	var itemName string
	var totalAmount, depositAmount float64
	var quantity int = 1

	if req.Type == string(model.OrderTypeVenue) {
		venue, err := s.venueRepo.GetByID(req.ItemID)
		if err != nil {
			return nil, errors.New("venue not found")
		}
		if venue.Status != model.VenueStatusOnline {
			return nil, errors.New("venue is not available")
		}

		conflict, err := s.orderRepo.CheckVenueConflict(req.ItemID, req.StartTime, req.EndTime)
		if err != nil {
			return nil, err
		}
		if conflict {
			return nil, errors.New("this time slot is already booked")
		}

		itemName = venue.Name
		dayOfWeek := int(startTime.Weekday())
		priceSlots, err := s.venueRepo.GetPrice(req.ItemID, dayOfWeek)
		if err == nil && len(priceSlots) > 0 {
			totalAmount = priceSlots[0].Price * totalHours
		}
		depositAmount = 0

	} else if req.Type == string(model.OrderTypeDevice) {
		device, err := s.deviceRepo.GetByID(req.ItemID)
		if err != nil {
			return nil, errors.New("device not found")
		}
		if device.Status != model.DeviceStatusOnline {
			return nil, errors.New("device is not available")
		}

		availability, err := s.deviceRepo.GetAvailability(req.ItemID, startTime.Format("2006-01-02"))
		if err != nil {
			return nil, err
		}
		if req.Quantity > availability.AvailableQuantity {
			return nil, errors.New("insufficient available quantity")
		}

		quantity = req.Quantity
		itemName = device.Name
		totalAmount = device.RentalPrice * totalHours * float64(quantity)
		depositAmount = device.DepositAmount * float64(quantity)
	} else {
		return nil, errors.New("invalid order type")
	}

	orderNo := generateOrderNo(userID)

	order := &model.Order{
		OrderNo:      orderNo,
		UserID:       userID,
		Type:         model.OrderType(req.Type),
		ItemID:       req.ItemID,
		ItemName:     itemName,
		StartTime:    startTime,
		EndTime:      endTime,
		TotalHours:   totalHours,
		Quantity:     quantity,
		TotalAmount:  totalAmount,
		DepositAmount: depositAmount,
		Status:       model.OrderStatusPending,
		Purpose:      req.Purpose,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	if req.Type == string(model.OrderTypeDevice) {
		s.deviceRepo.UpdateAvailableQuantity(req.ItemID, -quantity)
	}

	return order, nil
}

func (s *OrderService) GetByID(id uint) (*model.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) List(req *dto.OrderListRequest, userID *uint) ([]model.Order, int64, error) {
	return s.orderRepo.List(req, userID)
}

func (s *OrderService) Cancel(id uint, userID uint, reason string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("you can only cancel your own orders")
	}

	if order.Status == model.OrderStatusCancelled || order.Status == model.OrderStatusCompleted {
		return errors.New("order cannot be cancelled")
	}

	cancelDeadline := order.StartTime.Add(-24 * time.Hour)
	if time.Now().After(cancelDeadline) {
		return errors.New("cancellation is only allowed 24 hours before start time")
	}

	now := time.Now()
	order.Status = model.OrderStatusCancelled
	order.CancelReason = reason
	order.CancelledAt = &now

	err = s.orderRepo.Update(order)
	if err != nil {
		return err
	}

	if order.Type == model.OrderTypeDevice {
		s.deviceRepo.UpdateAvailableQuantity(order.ItemID, order.Quantity)
	}

	return nil
}

func (s *OrderService) Confirm(id uint, adminID uint, note string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != model.OrderStatusPending {
		return errors.New("only pending orders can be confirmed")
	}

	now := time.Now()
	order.Status = model.OrderStatusConfirmed
	order.ReviewedBy = &adminID
	order.ReviewNote = note
	order.ReviewedAt = &now

	return s.orderRepo.Update(order)
}

func (s *OrderService) Complete(id uint, adminID uint) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != model.OrderStatusPaid && order.Status != model.OrderStatusConfirmed {
		return errors.New("only confirmed or paid orders can be completed")
	}

	order.Status = model.OrderStatusCompleted
	return s.orderRepo.Update(order)
}

func (s *OrderService) GetCalendar(req *dto.CalendarRequest) ([]dto.CalendarEvent, error) {
	orders, err := s.orderRepo.GetCalendarBookings(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	events := make([]dto.CalendarEvent, 0, len(orders))
	for _, order := range orders {
		if req.Type != "" && string(order.Type) != req.Type {
			continue
		}
		if req.ItemID > 0 && order.ItemID != req.ItemID {
			continue
		}

		color := getStatusColor(order.Status)

		events = append(events, dto.CalendarEvent{
			ID:       order.ID,
			Title:    order.ItemName + " (" + string(order.Status) + ")",
			Start:    order.StartTime.Format(time.RFC3339),
			End:      order.EndTime.Format(time.RFC3339),
			Status:   string(order.Status),
			Type:     string(order.Type),
			ItemID:   order.ItemID,
			ItemName: order.ItemName,
			Color:    color,
		})
	}

	return events, nil
}

func getStatusColor(status model.OrderStatus) string {
	switch status {
	case model.OrderStatusPending:
		return "#faad14"
	case model.OrderStatusConfirmed:
		return "#1890ff"
	case model.OrderStatusPaid:
		return "#52c41a"
	case model.OrderStatusCompleted:
		return "#8c8c8c"
	case model.OrderStatusCancelled:
		return "#ff4d4f"
	default:
		return "#1890ff"
	}
}
