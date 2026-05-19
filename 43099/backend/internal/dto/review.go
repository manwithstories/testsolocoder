package dto

type ReviewCreateRequest struct {
	OrderID uint   `json:"order_id" binding:"required,min=1"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content" binding:"max=1000"`
}

type ReviewListRequest struct {
	PaginationRequest
	Status string `form:"status"`
	ItemID uint   `form:"item_id"`
	Type   string `form:"type"`
}

type ReviewActionRequest struct {
	Note string `json:"note" binding:"max=500"`
}
