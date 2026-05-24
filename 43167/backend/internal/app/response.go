package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func Fail(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{Code: httpStatus, Message: message})
}

func BindFail(c *gin.Context, err error) {
	Fail(c, http.StatusBadRequest, "参数校验失败: "+err.Error())
}

func BizFail(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	var biz *BizError
	if errors.As(err, &biz) {
		status = biz.Status
	}
	Fail(c, status, err.Error())
}

type BizError struct {
	Status  int
	Message string
}

func (e *BizError) Error() string { return e.Message }

func NewBizError(status int, msg string) error {
	return &BizError{Status: status, Message: msg}
}
