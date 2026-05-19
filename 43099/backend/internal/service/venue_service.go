package service

import (
	"errors"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/internal/repository"
)

type VenueService struct {
	venueRepo *repository.VenueRepository
	orderRepo *repository.OrderRepository
}

func NewVenueService() *VenueService {
	return &VenueService{
		venueRepo: repository.NewVenueRepository(),
		orderRepo: repository.NewOrderRepository(),
	}
}

func (s *VenueService) Create(req *dto.VenueCreateRequest, userID uint) (*model.Venue, error) {
	venue := &model.Venue{
		Name:        req.Name,
		Location:    req.Location,
		Capacity:    req.Capacity,
		Facilities:  req.Facilities,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Status:      model.VenueStatusOnline,
		CreatedBy:   userID,
	}

	err := s.venueRepo.Create(venue)
	return venue, err
}

func (s *VenueService) GetByID(id uint) (*model.Venue, error) {
	return s.venueRepo.GetByID(id)
}

func (s *VenueService) List(req *dto.VenueListRequest) ([]model.Venue, int64, error) {
	return s.venueRepo.List(req)
}

func (s *VenueService) Update(id uint, req *dto.VenueUpdateRequest) (*model.Venue, error) {
	venue, err := s.venueRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		venue.Name = req.Name
	}
	if req.Location != "" {
		venue.Location = req.Location
	}
	if req.Capacity > 0 {
		venue.Capacity = req.Capacity
	}
	if req.Facilities != "" {
		venue.Facilities = req.Facilities
	}
	if req.Description != "" {
		venue.Description = req.Description
	}
	if req.CoverImage != "" {
		venue.CoverImage = req.CoverImage
	}

	err = s.venueRepo.Update(venue)
	return venue, err
}

func (s *VenueService) Delete(id uint) error {
	return s.venueRepo.Delete(id)
}

func (s *VenueService) UpdateStatus(id uint, status string) error {
	return s.venueRepo.UpdateStatus(id, model.VenueStatus(status))
}

func (s *VenueService) SetPrice(venueID uint, req *dto.VenuePriceRequest) error {
	return s.venueRepo.SetPrice(venueID, req.DayOfWeek, req.TimeSlots)
}

func (s *VenueService) GetAvailability(venueID uint, date string) (*dto.VenueAvailabilityResponse, error) {
	_, err := s.venueRepo.GetByID(venueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	t, err := parseDate(date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	dayOfWeek := int(t.Weekday())

	priceSlots, err := s.venueRepo.GetPrice(venueID, dayOfWeek)
	if err != nil {
		priceSlots = []dto.TimeSlotPrice{}
	}

	bookedOrders, err := s.orderRepo.GetVenueBookingsForDate(venueID, date)
	if err != nil {
		return nil, err
	}

	bookedSlots := make([]dto.BookedTimeSlot, 0, len(bookedOrders))
	for _, order := range bookedOrders {
		bookedSlots = append(bookedSlots, dto.BookedTimeSlot{
			Start:   order.StartTime.Format("15:04"),
			End:     order.EndTime.Format("15:04"),
			OrderID: order.ID,
		})
	}

	availableSlots := make([]dto.AvailableTimeSlot, 0)
	for _, slot := range priceSlots {
		available := true
		for _, booked := range bookedSlots {
			if isTimeOverlap(slot.Start, slot.End, booked.Start, booked.End) {
				available = false
				break
			}
		}
		if available {
			availableSlots = append(availableSlots, dto.AvailableTimeSlot{
				Start: slot.Start,
				End:   slot.End,
				Price: slot.Price,
			})
		}
	}

	return &dto.VenueAvailabilityResponse{
		Date:      date,
		VenueID:   venueID,
		Available: availableSlots,
		Booked:    bookedSlots,
	}, nil
}
