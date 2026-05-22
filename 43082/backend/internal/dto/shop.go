package dto

type ShopApplyRequest struct {
	Name            string `json:"name" binding:"required,max=100"`
	Description     string `json:"description"`
	Logo            string `json:"logo"`
	ContactName     string `json:"contact_name" binding:"required"`
	ContactPhone    string `json:"contact_phone" binding:"required"`
	Address         string `json:"address"`
	IDCardFront     string `json:"id_card_front"`
	IDCardBack      string `json:"id_card_back"`
	BusinessLicense string `json:"business_license"`
}

type ShopUpdateRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Logo            string `json:"logo"`
	ContactName     string `json:"contact_name"`
	ContactPhone    string `json:"contact_phone"`
	Address         string `json:"address"`
}

type ShopReviewRequest struct {
	Status       string `json:"status" binding:"required,oneof=approved rejected"`
	RejectReason string `json:"reject_reason"`
}

type ShopInfo struct {
	ID           uint    `json:"id"`
	UserID       uint    `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Logo         string  `json:"logo"`
	Status       string  `json:"status"`
	Rating       float64 `json:"rating"`
	CreatedAt    string  `json:"created_at"`
	ContactName  string  `json:"contact_name"`
	ContactPhone string  `json:"contact_phone"`
	Address      string  `json:"address"`
	ProductCount int64   `json:"product_count"`
	SoldCount    int64   `json:"sold_count"`
}

type ShopDetail struct {
	ShopInfo
	RejectReason string `json:"reject_reason,omitempty"`
	IDCardFront  string `json:"id_card_front"`
	IDCardBack   string `json:"id_card_back"`
	BusinessLicense string `json:"business_license"`
}
