package dto

// CreateDeliveryRequest 创建配送安装预约请求
type CreateDeliveryRequest struct {
	OrderID      uint   `json:"order_id" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=delivery install both"`
	TimeSlot     string `json:"time_slot" binding:"omitempty,max=128"`
	Address      string `json:"address" binding:"omitempty,max=255"`
	ContactName  string `json:"contact_name" binding:"omitempty,max=64"`
	ContactPhone string `json:"contact_phone" binding:"omitempty,max=20"`
	InstallerID  uint   `json:"installer_id"`
	Remark       string `json:"remark"`
}

// UpdateDeliveryRequest 更新配送安装请求
type UpdateDeliveryRequest struct {
	Type         string `json:"type" binding:"omitempty,oneof=delivery install both"`
	TimeSlot     string `json:"time_slot" binding:"omitempty,max=128"`
	Address      string `json:"address" binding:"omitempty,max=255"`
	ContactName  string `json:"contact_name" binding:"omitempty,max=64"`
	ContactPhone string `json:"contact_phone" binding:"omitempty,max=20"`
	Status       string `json:"status" binding:"omitempty,oneof=pending confirmed completed cancelled"`
	InstallerID  uint   `json:"installer_id"`
	Remark       string `json:"remark"`
}

// DeliveryListRequest 配送安装列表查询请求
type DeliveryListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	OrderID  uint   `form:"order_id"`
	Type     string `form:"type" binding:"omitempty,oneof=delivery install both"`
	Status   string `form:"status" binding:"omitempty,oneof=pending confirmed completed cancelled"`
	Role     string `form:"-"`
	UserID   uint   `form:"-"`
}

// DeliveryResponse 配送安装响应 DTO
type DeliveryResponse struct {
	ID           uint   `json:"id"`
	OrderID      uint   `json:"order_id"`
	OwnerID      uint   `json:"owner_id"`
	Type         string `json:"type"`
	TimeSlot     string `json:"time_slot"`
	Address      string `json:"address"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Status       string `json:"status"`
	InstallerID  uint   `json:"installer_id"`
	Remark       string `json:"remark"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// DeliveryListResponse 配送安装列表响应
type DeliveryListResponse struct {
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	List     []DeliveryResponse `json:"list"`
}
