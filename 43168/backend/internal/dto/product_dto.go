package dto

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name        string               `json:"name" binding:"required,max=128"`
	Category    string               `json:"category" binding:"required,max=64"`
	Description string               `json:"description"`
	BasePrice   float64              `json:"base_price" binding:"required,gte=0"`
	Stock       int                  `json:"stock" binding:"gte=0"`
	Status      int                  `json:"status" binding:"required,oneof=1 2 3"`
	IsHot       bool                 `json:"is_hot"`
	Options     []ProductOptionDTO   `json:"options"`
	Images      []ProductImageDTO    `json:"images"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name        string             `json:"name" binding:"omitempty,max=128"`
	Category    string             `json:"category" binding:"omitempty,max=64"`
	Description string             `json:"description"`
	BasePrice   *float64           `json:"base_price" binding:"omitempty,gte=0"`
	Stock       *int               `json:"stock" binding:"omitempty,gte=0"`
	Status      *int               `json:"status" binding:"omitempty,oneof=1 2 3"`
	IsHot       *bool              `json:"is_hot"`
	Options     []ProductOptionDTO `json:"options"`
	Images      []ProductImageDTO  `json:"images"`
}

// ProductOptionDTO 产品选项 DTO
type ProductOptionDTO struct {
	ID          uint    `json:"id"`
	OptionType  string  `json:"option_type" binding:"required,oneof=size material color"`
	OptionValue string  `json:"option_value" binding:"required,max=128"`
	PriceAdjust float64 `json:"price_adjust"`
	Sort        int     `json:"sort"`
}

// ProductImageDTO 产品图片 DTO
type ProductImageDTO struct {
	ID       uint   `json:"id"`
	ImageURL string `json:"image_url" binding:"required,max=255"`
	Sort     int    `json:"sort"`
}

// ProductListRequest 产品列表查询请求
type ProductListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Category string `form:"category"`
	Status   *int   `form:"status"`
	IsHot    *bool  `form:"is_hot"`
}

// ProductResponse 产品响应
type ProductResponse struct {
	ID             uint                `json:"id"`
	ManufacturerID uint                `json:"manufacturer_id"`
	Name           string              `json:"name"`
	Category       string              `json:"category"`
	Description    string              `json:"description"`
	BasePrice      float64             `json:"base_price"`
	Stock          int                 `json:"stock"`
	Status         int                 `json:"status"`
	IsHot          bool                `json:"is_hot"`
	Options        []ProductOptionDTO  `json:"options"`
	Images         []ProductImageDTO   `json:"images"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
}

// ProductListResponse 分页产品列表响应
type ProductListResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	List     []ProductResponse `json:"list"`
}

// HotProductResponse 热门产品简版响应（用于列表展示）
type HotProductResponse struct {
	ID             uint    `json:"id"`
	ManufacturerID uint    `json:"manufacturer_id"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	BasePrice      float64 `json:"base_price"`
	Cover          string  `json:"cover"`
}
