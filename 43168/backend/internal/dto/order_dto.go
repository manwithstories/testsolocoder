package dto

// CreateInquiryRequest 创建询价请求
type CreateInquiryRequest struct {
	DesignerID   uint             `json:"designer_id"`
	Address      string           `json:"address" binding:"omitempty,max=255"`
	ContactName  string           `json:"contact_name" binding:"omitempty,max=64"`
	ContactPhone string           `json:"contact_phone" binding:"omitempty,max=20"`
	Remark       string           `json:"remark"`
	Items        []InquiryItemDTO `json:"items" binding:"required,min=1"`
}

// InquiryItemDTO 询价项 DTO
type InquiryItemDTO struct {
	ProductID     uint              `json:"product_id" binding:"required"`
	Quantity      int               `json:"quantity" binding:"required,gte=1"`
	CustomOptions []CustomOptionDTO `json:"custom_options"`
}

// CustomOptionDTO 定制选项 DTO
type CustomOptionDTO struct {
	OptionType  string  `json:"option_type" binding:"required"`
	OptionValue string  `json:"option_value" binding:"required"`
	PriceAdjust float64 `json:"price_adjust"`
}

// QuoteRequest 厂商报价请求
type QuoteRequest struct {
	Discount   float64          `json:"discount" binding:"gte=0"`
	ItemPrices []QuoteItemPrice `json:"item_prices" binding:"required,min=1"`
	Remark     string           `json:"remark"`
}

// QuoteItemPrice 订单项报价
type QuoteItemPrice struct {
	OrderItemID uint    `json:"order_item_id" binding:"required"`
	BasePrice   float64 `json:"base_price" binding:"gte=0"`
	Subtotal    float64 `json:"subtotal" binding:"gte=0"`
}

// OrderListRequest 订单列表查询请求
type OrderListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status" binding:"omitempty,oneof=inquiry quoted confirmed paid producing shipped completed cancelled"`
	Role     string `form:"-"`
	UserID   uint   `form:"-"`
}

// UpdateStatusRequest 更新订单状态请求
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=inquiry quoted confirmed paid producing shipped completed cancelled"`
	Remark string `json:"remark"`
}

// OrderItemDTO 订单项响应 DTO
type OrderItemDTO struct {
	ID            uint              `json:"id"`
	OrderID       uint              `json:"order_id"`
	ProductID     uint              `json:"product_id"`
	ProductName   string            `json:"product_name"`
	BasePrice     float64           `json:"base_price"`
	CustomOptions []CustomOptionDTO `json:"custom_options"`
	Quantity      int               `json:"quantity"`
	Subtotal      float64           `json:"subtotal"`
}

// OrderHistoryDTO 订单历史响应 DTO
type OrderHistoryDTO struct {
	ID           uint   `json:"id"`
	OrderID      uint   `json:"order_id"`
	Status       string `json:"status"`
	OperatorID   uint   `json:"operator_id"`
	OperatorRole string `json:"operator_role"`
	Remark       string `json:"remark"`
	CreatedAt    string `json:"created_at"`
}

// OrderResponse 订单响应 DTO（列表）
type OrderResponse struct {
	ID             uint    `json:"id"`
	OrderNo        string  `json:"order_no"`
	OwnerID        uint    `json:"owner_id"`
	DesignerID     uint    `json:"designer_id"`
	ManufacturerID uint    `json:"manufacturer_id"`
	TotalAmount    float64 `json:"total_amount"`
	Discount       float64 `json:"discount"`
	FinalAmount    float64 `json:"final_amount"`
	Status         string  `json:"status"`
	Address        string  `json:"address"`
	ContactName    string  `json:"contact_name"`
	ContactPhone   string  `json:"contact_phone"`
	Remark         string  `json:"remark"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// OrderDetailResponse 订单详情响应 DTO
type OrderDetailResponse struct {
	OrderResponse
	Items     []OrderItemDTO    `json:"items"`
	Histories []OrderHistoryDTO `json:"histories"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	List     []OrderResponse `json:"list"`
}
