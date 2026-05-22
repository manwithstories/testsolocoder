package dto

type CartAddRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	SKU_ID    *uint   `json:"sku_id"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
}

type CartUpdateRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartItemInfo struct {
	ID         uint      `json:"id"`
	ProductID  uint      `json:"product_id"`
	SKU_ID     *uint     `json:"sku_id"`
	Quantity   int       `json:"quantity"`
	Product    ProductInfo `json:"product"`
	SKU        *SKUInfo  `json:"sku,omitempty"`
}

type OrderCreateRequest struct {
	CartIDs       []uint `json:"cart_ids"`
	ReceiverName  string `json:"receiver_name" binding:"required"`
	ReceiverPhone string `json:"receiver_phone" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	Remark        string `json:"remark"`
}

type OrderCreateResponse struct {
	OrderNos []string `json:"order_nos"`
	TotalAmount float64 `json:"total_amount"`
}

type OrderInfo struct {
	ID              uint          `json:"id"`
	OrderNo         string        `json:"order_no"`
	ShopID          uint          `json:"shop_id"`
	ShopName        string        `json:"shop_name"`
	TotalAmount     float64       `json:"total_amount"`
	ShippingFee     float64       `json:"shipping_fee"`
	Status          string        `json:"status"`
	StatusText      string        `json:"status_text"`
	ReceiverName    string        `json:"receiver_name"`
	ReceiverPhone   string        `json:"receiver_phone"`
	ReceiverAddress string        `json:"receiver_address"`
	Remark          string        `json:"remark"`
	CreatedAt       string        `json:"created_at"`
	PaidAt          string        `json:"paid_at,omitempty"`
	ShippedAt       string        `json:"shipped_at,omitempty"`
	CompletedAt     string        `json:"completed_at,omitempty"`
	TrackingNo      string        `json:"tracking_no,omitempty"`
	TrackingCompany string        `json:"tracking_company,omitempty"`
	Items           []OrderItemInfo `json:"items,omitempty"`
}

type OrderDetail struct {
	OrderInfo
	Items []OrderItemInfo `json:"items"`
}

type OrderItemInfo struct {
	ID            uint              `json:"id"`
	ProductID     uint              `json:"product_id"`
	SKU_ID        *uint             `json:"sku_id"`
	ProductName   string            `json:"product_name"`
	ProductImage  string            `json:"product_image"`
	Specs         map[string]string `json:"specs,omitempty"`
	Price         float64           `json:"price"`
	Quantity      int               `json:"quantity"`
	Subtotal      float64           `json:"subtotal"`
	Reviewed      bool              `json:"reviewed"`
}

type OrderShipRequest struct {
	TrackingNo      string `json:"tracking_no" binding:"required"`
	TrackingCompany string `json:"tracking_company" binding:"required"`
}

type OrderQuery struct {
	Status   string `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type PaymentRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Method  string `json:"method" binding:"required"`
}

type RefundRequest struct {
	OrderID     uint   `json:"order_id" binding:"required"`
	OrderItemID *uint  `json:"order_item_id"`
	Reason      string `json:"reason" binding:"required"`
	Type        string `json:"type" binding:"required"`
}

type RefundReviewRequest struct {
	Status       string `json:"status" binding:"required,oneof=approved rejected"`
	RejectReason string `json:"reject_reason"`
}
