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

type ErrorCode struct {
	Code    int
	Message string
}

var (
	Success              = ErrorCode{0, "success"}
	BadRequest           = ErrorCode{40001, "请求参数错误"}
	Unauthorized         = ErrorCode{40101, "未授权"}
	Forbidden            = ErrorCode{40301, "权限不足"}
	NotFound             = ErrorCode{40401, "资源不存在"}
	InternalServerError  = ErrorCode{50001, "服务器内部错误"}
	UserAlreadyExists    = ErrorCode{40002, "用户已存在"}
	UserNotFound         = ErrorCode{40003, "用户不存在"}
	InvalidCredentials   = ErrorCode{40004, "用户名或密码错误"}
	QuestionNotFound     = ErrorCode{40005, "问题不存在"}
	AnswerNotFound       = ErrorCode{40006, "回答不存在"}
	CommentNotFound      = ErrorCode{40007, "评论不存在"}
	InsufficientPoints   = ErrorCode{40008, "积分不足"}
	RewardNotFound       = ErrorCode{40009, "奖品不存在"}
	ReportNotFound       = ErrorCode{40010, "举报不存在"}
	CategoryNotFound     = ErrorCode{40011, "分类不存在"}
	AlreadyFavorited     = ErrorCode{40012, "已收藏"}
	AlreadyFollowed      = ErrorCode{40013, "已关注"}
	ContentAuditPending  = ErrorCode{40014, "内容审核中"}
	ContentAuditRejected = ErrorCode{40015, "内容审核未通过"}
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    Success.Code,
		Message: Success.Message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, errCode ErrorCode) {
	c.JSON(http.StatusOK, Response{
		Code:    errCode.Code,
		Message: errCode.Message,
	})
}

func ErrorResponseWithMessage(c *gin.Context, errCode ErrorCode, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    errCode.Code,
		Message: message,
	})
}
