package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func ErrParam(c *gin.Context, msg string) {
	Fail(c, 400, msg)
}

func ErrAuth(c *gin.Context, msg string) {
	Fail(c, 401, msg)
}

func ErrForbidden(c *gin.Context, msg string) {
	Fail(c, 403, msg)
}

func ErrNotFound(c *gin.Context, msg string) {
	Fail(c, 404, msg)
}

func ErrServer(c *gin.Context, msg string) {
	Fail(c, 500, msg)
}
