package dto

type PaymentListRequest struct {
	PaginationRequest
	Status string `form:"status"`
}

type ConfirmPaymentRequest struct {
	TransactionNo string  `json:"transaction_no" binding:"required,max=64"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=wechat alipay cash"`
	Amount        float64 `json:"amount" binding:"required,min=0"`
}

type PaymentExportRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
