package response

import (
	"net/http"
	"ticket-system/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess      = 0
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, message string) {
	logger.Log.WithFields(logrus.Fields{
		"code": code,
		"path": c.Request.URL.Path,
		"ip":   c.ClientIP(),
	}).Warn(message)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Fail(c, CodeBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Fail(c, CodeUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Fail(c, CodeForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Fail(c, CodeNotFound, message)
}

func ServerError(c *gin.Context, message string) {
	logger.Log.WithFields(logrus.Fields{
		"path":  c.Request.URL.Path,
		"ip":    c.ClientIP(),
		"error": message,
	}).Error("Server error")
	c.JSON(http.StatusOK, Response{
		Code:    CodeServerError,
		Message: message,
	})
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
