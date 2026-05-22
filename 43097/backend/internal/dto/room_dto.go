package dto

import "hotel-system/internal/model"

type RoomTypeCreateRequest struct {
	Name        string   `json:"name" binding:"required,max=50"`
	Description string   `json:"description" binding:"max=500"`
	BasePrice   float64  `json:"basePrice" binding:"required,min=0"`
	BedCount    int      `json:"bedCount" binding:"min=1"`
	MaxGuests   int      `json:"maxGuests" binding:"min=1"`
	Facilities  []string `json:"facilities"`
}

type RoomTypeUpdateRequest struct {
	Name        string   `json:"name" binding:"omitempty,max=50"`
	Description string   `json:"description" binding:"omitempty,max=500"`
	BasePrice   float64  `json:"basePrice" binding:"omitempty,min=0"`
	BedCount    int      `json:"bedCount" binding:"omitempty,min=1"`
	MaxGuests   int      `json:"maxGuests" binding:"omitempty,min=1"`
	Facilities  []string `json:"facilities"`
}

type RoomTypeResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	BasePrice   float64  `json:"basePrice"`
	BedCount    int      `json:"bedCount"`
	MaxGuests   int      `json:"maxGuests"`
	Facilities  []string `json:"facilities"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type RoomCreateRequest struct {
	RoomNo     string              `json:"roomNo" binding:"required,max=20"`
	Floor      int                 `json:"floor" binding:"required,min=1"`
	RoomTypeID uint                `json:"roomTypeId" binding:"required,min=1"`
	Price      float64             `json:"price" binding:"min=0"`
	Facilities []string            `json:"facilities"`
	Status     model.RoomStatus    `json:"status" binding:"omitempty,oneof=available occupied reserved maintenance"`
}

type RoomUpdateRequest struct {
	RoomNo     string              `json:"roomNo" binding:"omitempty,max=20"`
	Floor      int                 `json:"floor" binding:"omitempty,min=1"`
	RoomTypeID uint                `json:"roomTypeId" binding:"omitempty,min=1"`
	Price      float64             `json:"price" binding:"omitempty,min=0"`
	Facilities []string            `json:"facilities"`
	Status     model.RoomStatus    `json:"status" binding:"omitempty,oneof=available occupied reserved maintenance"`
}

type RoomResponse struct {
	ID         uint              `json:"id"`
	RoomNo     string            `json:"roomNo"`
	Floor      int               `json:"floor"`
	RoomTypeID uint              `json:"roomTypeId"`
	RoomType   *RoomTypeResponse `json:"roomType,omitempty"`
	Status     model.RoomStatus  `json:"status"`
	Price      float64           `json:"price"`
	Facilities []string          `json:"facilities"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  string            `json:"updatedAt"`
}

type RoomListRequest struct {
	PaginationRequest
	RoomNo     string           `form:"roomNo"`
	Floor      int              `form:"floor"`
	RoomTypeID uint             `form:"roomTypeId"`
	Status     model.RoomStatus `form:"status"`
}

type RoomBatchImportRequest struct {
	Rooms []RoomCreateRequest `json:"rooms" binding:"required,min=1,dive"`
}

type RoomBatchImportResponse struct {
	SuccessCount int      `json:"successCount"`
	FailCount    int      `json:"failCount"`
	FailReasons  []string `json:"failReasons,omitempty"`
}

type RoomStatusUpdateRequest struct {
	Status model.RoomStatus `json:"status" binding:"required,oneof=available occupied reserved maintenance"`
}

type RoomDashboardResponse struct {
	Total       int64 `json:"total"`
	Available   int64 `json:"available"`
	Occupied    int64 `json:"occupied"`
	Reserved    int64 `json:"reserved"`
	Maintenance int64 `json:"maintenance"`
}
