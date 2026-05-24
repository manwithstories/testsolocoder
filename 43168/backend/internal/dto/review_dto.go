package dto

// CreateReviewRequest 创建售后评价请求
type CreateReviewRequest struct {
	OrderID       uint     `json:"order_id" binding:"required"`
	ProductID     uint     `json:"product_id" binding:"required"`
	ProductRating int      `json:"product_rating" binding:"required,gte=1,lte=5"`
	ServiceRating int      `json:"service_rating" binding:"required,gte=1,lte=5"`
	Content       string   `json:"content"`
	Images        []string `json:"images"`
}

// ReviewListRequest 评价列表查询请求
type ReviewListRequest struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=10" binding:"min=1,max=100"`
	OrderID   uint   `form:"order_id"`
	ProductID uint   `form:"product_id"`
	OwnerID   uint   `form:"owner_id"`
	Role      string `form:"-"`
	UserID    uint   `form:"-"`
}

// ReviewResponse 评价响应 DTO
type ReviewResponse struct {
	ID            uint     `json:"id"`
	OrderID       uint     `json:"order_id"`
	ProductID     uint     `json:"product_id"`
	OwnerID       uint     `json:"owner_id"`
	ProductRating int      `json:"product_rating"`
	ServiceRating int      `json:"service_rating"`
	Content       string   `json:"content"`
	Images        []string `json:"images"`
	CreatedAt     string   `json:"created_at"`
}

// ReviewListResponse 评价列表响应
type ReviewListResponse struct {
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	List     []ReviewResponse `json:"list"`
}
