package dto

// CreateTicketRequest 创建售后工单请求
type CreateTicketRequest struct {
	OrderID uint     `json:"order_id" binding:"required"`
	Type    string   `json:"type" binding:"required,oneof=product_quality delivery service other"`
	Title   string   `json:"title" binding:"required,max=255"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

// UpdateTicketRequest 更新售后工单请求
type UpdateTicketRequest struct {
	Type    string   `json:"type" binding:"omitempty,oneof=product_quality delivery service other"`
	Title   string   `json:"title" binding:"omitempty,max=255"`
	Content string   `json:"content"`
	Status  string   `json:"status" binding:"omitempty,oneof=open processing resolved closed"`
	Images  []string `json:"images"`
}

// TicketListRequest 工单列表查询请求
type TicketListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	OrderID  uint   `form:"order_id"`
	Type     string `form:"type" binding:"omitempty,oneof=product_quality delivery service other"`
	Status   string `form:"status" binding:"omitempty,oneof=open processing resolved closed"`
	Role     string `form:"-"`
	UserID   uint   `form:"-"`
}

// TicketResponse 工单响应 DTO
type TicketResponse struct {
	ID        uint     `json:"id"`
	OrderID   uint     `json:"order_id"`
	OwnerID   uint     `json:"owner_id"`
	Type      string   `json:"type"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Status    string   `json:"status"`
	Images    []string `json:"images"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// TicketListResponse 工单列表响应
type TicketListResponse struct {
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	List     []TicketResponse `json:"list"`
}
