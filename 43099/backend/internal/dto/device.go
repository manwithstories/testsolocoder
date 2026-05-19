package dto

type DeviceCategoryCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type DeviceCreateRequest struct {
	CategoryID      uint    `json:"category_id" binding:"required,min=1"`
	Name            string  `json:"name" binding:"required,max=100"`
	Description     string  `json:"description"`
	Specification   string  `json:"specification"`
	StockQuantity   int     `json:"stock_quantity" binding:"min=0"`
	RentalPrice     float64 `json:"rental_price" binding:"min=0"`
	DepositAmount   float64 `json:"deposit_amount" binding:"min=0"`
}

type DeviceUpdateRequest struct {
	CategoryID      uint    `json:"category_id" binding:"min=1"`
	Name            string  `json:"name" binding:"max=100"`
	Description     string  `json:"description"`
	Specification   string  `json:"specification"`
	StockQuantity   int     `json:"stock_quantity" binding:"min=0"`
	RentalPrice     float64 `json:"rental_price" binding:"min=0"`
	DepositAmount   float64 `json:"deposit_amount" binding:"min=0"`
	Status          string  `json:"status" binding:"omitempty,oneof=online offline"`
}

type DeviceStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=online offline"`
}

type DeviceAvailabilityRequest struct {
	Date string `form:"date" binding:"required"`
}

type DeviceAvailabilityResponse struct {
	Date              string `json:"date"`
	DeviceID          uint   `json:"device_id"`
	TotalStock        int    `json:"total_stock"`
	AvailableQuantity int    `json:"available_quantity"`
	BookedQuantity    int    `json:"booked_quantity"`
}

type DeviceListRequest struct {
	PaginationRequest
	CategoryID uint   `form:"category_id"`
	Status     string `form:"status"`
}

type DeviceBatchImportItem struct {
	CategoryName    string  `json:"category_name"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Specification   string  `json:"specification"`
	StockQuantity   int     `json:"stock_quantity"`
	RentalPrice     float64 `json:"rental_price"`
	DepositAmount   float64 `json:"deposit_amount"`
}

type DeviceBatchImportResponse struct {
	SuccessCount int      `json:"success_count"`
	FailCount    int      `json:"fail_count"`
	Errors       []string `json:"errors"`
}
