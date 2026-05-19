package dto

type VenueCreateRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Location    string `json:"location" binding:"max=255"`
	Capacity    int    `json:"capacity" binding:"min=0"`
	Facilities  string `json:"facilities"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

type VenueUpdateRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Location    string `json:"location" binding:"max=255"`
	Capacity    int    `json:"capacity" binding:"min=0"`
	Facilities  string `json:"facilities"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

type VenueStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=online offline"`
}

type TimeSlotPrice struct {
	Start string  `json:"start" binding:"required"`
	End   string  `json:"end" binding:"required"`
	Price float64 `json:"price" binding:"min=0"`
}

type VenuePriceRequest struct {
	DayOfWeek int             `json:"day_of_week" binding:"required,min=0,max=6"`
	TimeSlots []TimeSlotPrice `json:"time_slots" binding:"required"`
}

type VenueAvailabilityRequest struct {
	Date string `form:"date" binding:"required"`
}

type VenueAvailabilityResponse struct {
	Date      string                  `json:"date"`
	VenueID   uint                    `json:"venue_id"`
	Available []AvailableTimeSlot `json:"available"`
	Booked    []BookedTimeSlot    `json:"booked"`
}

type AvailableTimeSlot struct {
	Start string  `json:"start"`
	End   string  `json:"end"`
	Price float64 `json:"price"`
}

type BookedTimeSlot struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	OrderID uint   `json:"order_id"`
}

type VenueListRequest struct {
	PaginationRequest
	Status string `form:"status"`
}
