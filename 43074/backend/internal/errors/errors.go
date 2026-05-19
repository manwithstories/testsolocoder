package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "资源不存在",
	}
	ErrBadRequest = &AppError{
		Code:    http.StatusBadRequest,
		Message: "请求参数错误",
	}
	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "未授权访问",
	}
	ErrForbidden = &AppError{
		Code:    http.StatusForbidden,
		Message: "禁止访问",
	}
	ErrInternalServer = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
	}
	ErrConflict = &AppError{
		Code:    http.StatusConflict,
		Message: "资源冲突",
	}
	ErrValidation = &AppError{
		Code:    http.StatusBadRequest,
		Message: "数据验证失败",
	}
	ErrISBNInvalid = &AppError{
		Code:    http.StatusBadRequest,
		Message: "ISBN格式无效",
	}
	ErrISBNNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "未找到该ISBN对应的书籍信息",
	}
	ErrFileTooLarge = &AppError{
		Code:    http.StatusBadRequest,
		Message: "文件大小超出限制",
	}
	ErrFileTypeNotAllowed = &AppError{
		Code:    http.StatusBadRequest,
		Message: "文件类型不允许",
	}
	ErrBookBorrowed = &AppError{
		Code:    http.StatusConflict,
		Message: "图书已被借出",
	}
	ErrBookNotBorrowed = &AppError{
		Code:    http.StatusBadRequest,
		Message: "图书未被借出",
	}
	ErrInvalidPageNumber = &AppError{
		Code:    http.StatusBadRequest,
		Message: "页码无效",
	}
)

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewWithDetails(code int, message string, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.Code, appErr)
		return
	}
	c.JSON(http.StatusInternalServerError, ErrInternalServer)
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, &AppError{
		Code:    code,
		Message: message,
	})
}
