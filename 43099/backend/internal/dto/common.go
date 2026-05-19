package dto

import "venue-booking/internal/model"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

type PaginationResponse struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     interface{} `json:"list"`
}

func Success(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func SuccessNoData() Response {
	return Response{
		Code:    200,
		Message: "success",
	}
}

func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

type UserContext struct {
	UserID uint           `json:"user_id"`
	Role   model.UserRole `json:"role"`
}

type DateRangeRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
