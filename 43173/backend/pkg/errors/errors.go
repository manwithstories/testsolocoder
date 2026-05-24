package errors

import "fmt"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("error code: %d, message: %s", e.Code, e.Message)
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrUserNotFound        = NewAppError(1001, "用户不存在")
	ErrUserAlreadyExists   = NewAppError(1002, "用户已存在")
	ErrInvalidPassword     = NewAppError(1003, "密码错误")
	ErrInvalidToken        = NewAppError(1004, "无效的Token")
	ErrTokenExpired        = NewAppError(1005, "Token已过期")
	ErrUnauthorized        = NewAppError(1006, "未授权访问")
	ErrForbidden           = NewAppError(1007, "权限不足")

	ErrWorkNotFound        = NewAppError(2001, "作品不存在")
	ErrWorkAlreadyExists   = NewAppError(2002, "作品已存在")
	ErrInvalidAudioFormat  = NewAppError(2003, "无效的音频格式")
	ErrUploadFailed        = NewAppError(2004, "上传失败")

	ErrAlbumNotFound       = NewAppError(3001, "专辑不存在")
	ErrPlaylistNotFound    = NewAppError(3002, "歌单不存在")

	ErrEventNotFound       = NewAppError(4001, "演出不存在")
	ErrTicketSoldOut       = NewAppError(4002, "票已售罄")
	ErrInvalidSeat         = NewAppError(4003, "座位无效")

	ErrInsufficientBalance = NewAppError(5001, "余额不足")
	ErrWithdrawFailed      = NewAppError(5002, "提现失败")

	ErrDatabaseError       = NewAppError(6001, "数据库错误")
	ErrRedisError          = NewAppError(6002, "Redis错误")
	ErrSystemError         = NewAppError(6003, "系统错误")

	ErrValidationError     = NewAppError(7001, "参数校验失败")
)
