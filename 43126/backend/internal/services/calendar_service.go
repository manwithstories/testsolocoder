package services

import (
	"encoding/json"
	"fmt"
	"time"

	"meeting-room/internal/models"
	"meeting-room/internal/repositories"
	"meeting-room/internal/utils"
)

type CalendarService struct {
	bookingRepo *repositories.BookingRepository
	roomRepo    *repositories.RoomRepository
}

func NewCalendarService() *CalendarService {
	return &CalendarService{
		bookingRepo: repositories.NewBookingRepository(),
		roomRepo:    repositories.NewRoomRepository(),
	}
}

type CalendarBooking struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	StartTime string `json:"start"`
	EndTime   string `json:"end"`
	RoomID    uint   `json:"room_id"`
	RoomName  string `json:"room_name"`
	Floor     string `json:"floor"`
	Status    int    `json:"status"`
	Color     string `json:"color"`
}

func (s *CalendarService) GetWeekCalendar(date time.Time, roomID uint, floor string) ([]CalendarBooking, error) {
	startOfWeek := getStartOfWeek(date)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	cacheKey := fmt.Sprintf("calendar:week:%s:%d:%s", startOfWeek.Format("2006-01-02"), roomID, floor)
	cached, err := utils.RedisGet(cacheKey)
	if err == nil && cached != "" {
		var bookings []CalendarBooking
		json.Unmarshal([]byte(cached), &bookings)
		return bookings, nil
	}

	bookings, err := s.bookingRepo.GetByDateRange(startOfWeek, endOfWeek, roomID, floor)
	if err != nil {
		return nil, err
	}

	result := s.convertToCalendarBookings(bookings)

	cachedData, _ := json.Marshal(result)
	utils.RedisSet(cacheKey, string(cachedData), 5*time.Minute)

	return result, nil
}

func (s *CalendarService) GetMonthCalendar(year int, month int, roomID uint, floor string) ([]CalendarBooking, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	cacheKey := fmt.Sprintf("calendar:month:%d-%02d:%d:%s", year, month, roomID, floor)
	cached, err := utils.RedisGet(cacheKey)
	if err == nil && cached != "" {
		var bookings []CalendarBooking
		json.Unmarshal([]byte(cached), &bookings)
		return bookings, nil
	}

	bookings, err := s.bookingRepo.GetByDateRange(startOfMonth, endOfMonth, roomID, floor)
	if err != nil {
		return nil, err
	}

	result := s.convertToCalendarBookings(bookings)

	cachedData, _ := json.Marshal(result)
	utils.RedisSet(cacheKey, string(cachedData), 5*time.Minute)

	return result, nil
}

func (s *CalendarService) GetRoomAvailability(roomID uint, date time.Time) map[string]bool {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return map[string]bool{}
	}

	bookings, err := s.bookingRepo.GetByDateRange(startOfDay, endOfDay, roomID, "")
	if err != nil {
		return map[string]bool{}
	}

	availableStart, _ := time.Parse("15:04", room.AvailableStart)
	availableEnd, _ := time.Parse("15:04", room.AvailableEnd)

	startHour := availableStart.Hour()
	endHour := availableEnd.Hour()
	if availableEnd.Minute() > 0 {
		endHour++
	}

	availability := make(map[string]bool)
	for hour := startHour; hour < endHour; hour++ {
		slot := fmt.Sprintf("%02d:00", hour)
		availability[slot] = true
	}

	for _, booking := range bookings {
		if booking.Status == models.BookingStatusCancelled {
			continue
		}
		bookingStartHour := booking.StartTime.Hour()
		bookingEndHour := booking.EndTime.Hour()
		if booking.EndTime.Minute() > 0 {
			bookingEndHour++
		}
		for hour := bookingStartHour; hour < bookingEndHour; hour++ {
			slot := fmt.Sprintf("%02d:00", hour)
			availability[slot] = false
		}
	}

	return availability
}

func (s *CalendarService) getStatusColor(status models.BookingStatus) string {
	switch status {
	case models.BookingStatusConfirmed:
		return "#4CAF50"
	case models.BookingStatusPending:
		return "#FFC107"
	case models.BookingStatusCancelled:
		return "#9E9E9E"
	case models.BookingStatusCompleted:
		return "#2196F3"
	default:
		return "#9E9E9E"
	}
}

func (s *CalendarService) convertToCalendarBookings(bookings []models.Booking) []CalendarBooking {
	result := make([]CalendarBooking, 0, len(bookings))
	for _, booking := range bookings {
		roomName := ""
		floor := ""
		if booking.Room != nil {
			roomName = booking.Room.Name
			floor = booking.Room.Floor
		}

		result = append(result, CalendarBooking{
			ID:        booking.ID,
			Title:     booking.Title,
			StartTime: booking.StartTime.Format(time.RFC3339),
			EndTime:   booking.EndTime.Format(time.RFC3339),
			RoomID:    booking.RoomID,
			RoomName:  roomName,
			Floor:     floor,
			Status:    int(booking.Status),
			Color:     s.getStatusColor(booking.Status),
		})
	}
	return result
}

func getStartOfWeek(date time.Time) time.Time {
	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}
	return time.Date(date.Year(), date.Month(), date.Day()-dayOfWeek+1, 0, 0, 0, 0, time.Local)
}
