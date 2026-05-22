package dto

type ProductCreateRequest struct {
	CategoryID  uint                 `json:"category_id" binding:"required"`
	Name        string               `json:"name" binding:"required,max=200"`
	Description string              `json:"description"`
	MainImage   string               `json:"main_image"`
	Price       float64              `json:"price" binding:"required,min=0"`
	Stock       int                  `json:"stock" binding:"min=0"`
	Weight      float64              `json:"weight"`
	Images      []string             `json:"images"`
	Specs       []ProductSpecCreate `json:"specs"`
	Skus        []SKUCreate        `json:"skus"`
	IsHot       bool                 `json:"is_hot"`
	IsRecommend bool                 `json:"is_recommend"`
}

type ProductSpecCreate struct {
	Name   string   `json:"name" binding:"required"`
	Values []string `json:"values" binding:"required"`
}

type SKUCreate struct {
	Specs  map[string]string `json:"specs" binding:"required"`
	Price  float64         `json:"price" binding:"required,min=0"`
	Stock  int             `json:"stock" binding:"min=0"`
	SKUCode string        `json:"sku_code"`
}

type ProductUpdateRequest struct {
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MainImage   string  `json:"main_image"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Status      string  `json:"status"`
	IsHot       *bool   `json:"is_hot"`
	IsRecommend *bool   `json:"is_recommend"`
}

type ProductInfo struct {
	ID          uint    `json:"id"`
	ShopID      uint    `json:"shop_id"`
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name"`
	MainImage   string  `json:"main_image"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Sales       int     `json:"sales"`
	Status      string  `json:"status"`
	IsHot       bool    `json:"is_hot"`
	IsRecommend bool    `json:"is_recommend"`
	CreatedAt   string  `json:"created_at"`
	ShopName    string  `json:"shop_name"`
	CategoryName string `json:"category_name"`
}

type ProductDetail struct {
	ProductInfo
	Description string                 `json:"description"`
	Images      []string             `json:"images"`
	Specs       []ProductSpecInfo  `json:"specs"`
	Skus        []SKUInfo          `json:"skus"`
	Shop        ShopInfo             `json:"shop"`
}

type ProductSpecInfo struct {
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type SKUInfo struct {
	ID      uint              `json:"id"`
	Specs   map[string]string `json:"specs"`
	Price   float64           `json:"price"`
	Stock   int               `json:"stock"`
	SKUCode string            `json:"sku_code"`
}

type ProductQuery struct {
	Keyword    string  `form:"keyword"`
	CategoryID *uint   `form:"category_id"`
	ShopID     *uint   `form:"shop_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Status     string  `form:"status"`
	SortBy     string  `form:"sort_by"`
	SortOrder  string  `form:"sort_order"`
	Page       int     `form:"page,default=1"`
	PageSize   int     `form:"page_size,default=10"`
	IsHot      *bool   `form:"is_hot"`
}
