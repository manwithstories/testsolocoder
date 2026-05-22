package dto

type ReviewCreateRequest struct {
	OrderID     uint   `json:"order_id" binding:"required"`
	OrderItemID uint   `json:"order_item_id" binding:"required"`
	ProductID   uint   `json:"product_id" binding:"required"`
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Content     string `json:"content" binding:"required"`
	Images      []string `json:"images"`
}

type ReviewReplyRequest struct {
	Reply string `json:"reply" binding:"required"`
}

type ReviewInfo struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	ProductID uint   `json:"product_id"`
	Rating    int    `json:"rating"`
	Content   string `json:"content"`
	Images    []string `json:"images"`
	Reply     string `json:"reply,omitempty"`
	ReplyAt   string `json:"reply_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

type CategoryCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Icon     string `json:"icon"`
	ParentID *uint  `json:"parent_id"`
	Sort     int    `json:"sort"`
}

type CategoryUpdateRequest struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Sort   int    `json:"sort"`
	Status string `json:"status"`
}

type CategoryInfo struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon"`
	ParentID *uint          `json:"parent_id"`
	Level    int            `json:"level"`
	Sort     int            `json:"sort"`
	Status   string         `json:"status"`
	Children []CategoryInfo `json:"children,omitempty"`
}

type NotificationInfo struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type DisputeCreateRequest struct {
	OrderID  uint     `json:"order_id" binding:"required"`
	Type     string   `json:"type" binding:"required"`
	Reason   string   `json:"reason" binding:"required"`
	Evidence []string `json:"evidence"`
}

type DisputeResolveRequest struct {
	Result string `json:"result" binding:"required"`
}

type DisputeInfo struct {
	ID         uint     `json:"id"`
	OrderID    uint     `json:"order_id"`
	OrderNo    string   `json:"order_no"`
	UserID     uint     `json:"user_id"`
	Username   string   `json:"username"`
	ShopID     uint     `json:"shop_id"`
	ShopName   string   `json:"shop_name"`
	Type       string   `json:"type"`
	Reason     string   `json:"reason"`
	Evidence   []string `json:"evidence"`
	Status     string   `json:"status"`
	Result     string   `json:"result,omitempty"`
	AdminID    *uint    `json:"admin_id,omitempty"`
	CreatedAt  string   `json:"created_at"`
	ResolvedAt string   `json:"resolved_at,omitempty"`
}

type StatisticsRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type ShopStatistics struct {
	TotalOrders   int64   `json:"total_orders"`
	TotalSales    float64 `json:"total_sales"`
	TotalProducts int64   `json:"total_products"`
	NewOrders     int64   `json:"new_orders"`
	DailySales    []DailySales `json:"daily_sales"`
	TopProducts   []ProductSales `json:"top_products"`
}

type AdminStatistics struct {
	TotalUsers    int64   `json:"total_users"`
	TotalShops    int64   `json:"total_shops"`
	TotalProducts int64   `json:"total_products"`
	TotalOrders   int64   `json:"total_orders"`
	TotalSales    float64 `json:"total_sales"`
	PendingShops  int64   `json:"pending_shops"`
	OpenDisputes  int64   `json:"open_disputes"`
	DailySales    []DailySales `json:"daily_sales"`
}

type DailySales struct {
	Date  string  `json:"date"`
	Amount float64 `json:"amount"`
	Orders int64   `json:"orders"`
}

type ProductSales struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Sales int    `json:"sales"`
	Amount float64 `json:"amount"`
}
