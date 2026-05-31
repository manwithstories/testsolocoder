package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess       = 0
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500
)

const (
	MsgSuccess       = "success"
	MsgBadRequest    = "请求参数错误"
	MsgUnauthorized  = "未授权访问"
	MsgForbidden     = "权限不足"
	MsgNotFound      = "资源不存在"
	MsgInternalError = "服务器内部错误"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, msg string) {
	statusCode := http.StatusInternalServerError
	switch code {
	case CodeBadRequest:
		statusCode = http.StatusBadRequest
	case CodeUnauthorized:
		statusCode = http.StatusUnauthorized
	case CodeForbidden:
		statusCode = http.StatusForbidden
	case CodeNotFound:
		statusCode = http.StatusNotFound
	}

	c.JSON(statusCode, Response{
		Code:    code,
		Message: msg,
	})
}

func FailWithData(c *gin.Context, code int, msg string, data interface{}) {
	statusCode := http.StatusInternalServerError
	switch code {
	case CodeBadRequest:
		statusCode = http.StatusBadRequest
	case CodeUnauthorized:
		statusCode = http.StatusUnauthorized
	case CodeForbidden:
		statusCode = http.StatusForbidden
	case CodeNotFound:
		statusCode = http.StatusNotFound
	}

	c.JSON(statusCode, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
