package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

func InternalServerError(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}

var (
	ErrUserNotFound        = NotFound("用户不存在")
	ErrInvalidCredentials  = Unauthorized("用户名或密码错误")
	ErrUserExists          = BadRequest("用户已存在")
	ErrShopExists          = BadRequest("店铺名称已存在")
	ErrShopNotFound        = NotFound("店铺不存在")
	ErrProductNotFound     = NotFound("商品不存在")
	ErrOrderNotFound       = NotFound("订单不存在")
	ErrInsufficientStock   = BadRequest("库存不足")
	ErrInvalidOrderStatus  = BadRequest("订单状态不正确")
	ErrInvalidAmount       = BadRequest("金额不正确")
	ErrPermissionDenied    = Forbidden("权限不足")
	ErrUnauthorized        = Unauthorized("未登录")
	ErrInvalidToken        = Unauthorized("无效的token")
	ErrCategoryNotFound    = NotFound("分类不存在")
	ErrCartEmpty           = BadRequest("购物车为空")
)
