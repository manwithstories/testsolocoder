package dto

import "hotel-system/internal/model"

type RoomStatusSummary struct {
	Total       int64 `json:"total"`
	Available   int64 `json:"available"`
	Occupied    int64 `json:"occupied"`
	Reserved    int64 `json:"reserved"`
	Maintenance int64 `json:"maintenance"`
}

type DashboardStatsResponse struct {
	RoomStatus      RoomStatusSummary `json:"roomStatus"`
	TodayCheckIns   int64             `json:"todayCheckIns"`
	TodayCheckOuts  int64             `json:"todayCheckOuts"`
	CurrentGuests   int64             `json:"currentGuests"`
	MonthlyRevenue  float64           `json:"monthlyRevenue"`
	MonthlyOccupancy float64          `json:"monthlyOccupancy"`
}

type GuestInfo struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	CheckInTime  string `json:"checkInTime,omitempty"`
	CheckOutTime string `json:"checkOutTime,omitempty"`
}

type RoomStatusItem struct {
	ID         uint              `json:"id"`
	RoomNo     string            `json:"roomNo"`
	Floor      int               `json:"floor"`
	RoomTypeID uint              `json:"roomTypeId"`
	RoomType   string            `json:"roomType"`
	Status     model.RoomStatus  `json:"status"`
	Price      float64           `json:"price"`
	Guest      *GuestInfo        `json:"guest,omitempty"`
	Booking    *GuestInfo        `json:"booking,omitempty"`
}

type FloorRoomsResponse struct {
	Floor int                `json:"floor"`
	Rooms []RoomStatusItem   `json:"rooms"`
}
