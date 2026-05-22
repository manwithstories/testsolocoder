package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserInfo UserInfo `json:"user_info"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	Balance  float64 `json:"balance"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type CreateAuctionItemRequest struct {
	Title        string    `json:"title" binding:"required,max=200"`
	Description  string    `json:"description"`
	CategoryID   uint      `json:"category_id" binding:"required"`
	StartPrice   float64   `json:"start_price" binding:"required,min=0"`
	ReservePrice float64   `json:"reserve_price" binding:"required,min=0"`
	Location     string    `json:"location"`
	Condition    string    `json:"condition"`
}

type UpdateAuctionItemRequest struct {
	Title        string  `json:"title" binding:"max=200"`
	Description  string  `json:"description"`
	CategoryID   uint    `json:"category_id"`
	StartPrice   float64 `json:"start_price" binding:"min=0"`
	ReservePrice float64 `json:"reserve_price" binding:"min=0"`
	Location     string  `json:"location"`
	Condition    string  `json:"condition"`
}

type AuctionItemQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	Keyword    string `form:"keyword"`
	CategoryID uint   `form:"category_id"`
	Status     *int   `form:"status"`
	SortBy     string `form:"sort_by,default=created_at"`
	SortOrder  string `form:"sort_order,default=desc"`
}

type CreateSessionRequest struct {
	Name         string  `json:"name" binding:"required,max=200"`
	Description  string  `json:"description"`
	StartTime    string  `json:"start_time" binding:"required"`
	EndTime      string  `json:"end_time" binding:"required"`
	MinIncrement float64 `json:"min_increment" binding:"required,min=0"`
	ExtendTime   int     `json:"extend_time" binding:"min=0"`
}

type UpdateSessionRequest struct {
	Name         string  `json:"name" binding:"max=200"`
	Description  string  `json:"description"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	MinIncrement float64 `json:"min_increment" binding:"min=0"`
	ExtendTime   int     `json:"extend_time" binding:"min=0"`
}

type AddItemToSessionRequest struct {
	AuctionItemIDs []uint `json:"auction_item_ids" binding:"required"`
}

type BidRequest struct {
	Amount     float64 `json:"amount" binding:"required,min=0"`
	MaxAutoBid float64 `json:"max_auto_bid"`
}

type SetAutoBidRequest struct {
	AuctionItemID uint    `json:"auction_item_id" binding:"required"`
	MaxPrice      float64 `json:"max_price" binding:"required,min=0"`
}

type CreateOrderRequest struct {
	AuctionItemID uint   `json:"auction_item_id" binding:"required"`
	ShippingInfo  string `json:"shipping_info"`
}

type PayOrderRequest struct {
	Method string `json:"method" binding:"required"`
}

type CreateReviewRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
}

type StatisticsQuery struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type StatisticsResponse struct {
	TotalAuctions      int64   `json:"total_auctions"`
	TotalBids          int64   `json:"total_bids"`
	TotalOrders        int64   `json:"total_orders"`
	TotalAmount        float64 `json:"total_amount"`
	SuccessRate        float64 `json:"success_rate"`
	ActiveUsers        int64   `json:"active_users"`
	NewUsers           int64   `json:"new_users"`
	AverageBidAmount   float64 `json:"average_bid_amount"`
}

type NotificationQuery struct {
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=20"`
	IsRead   *int `form:"is_read"`
}
