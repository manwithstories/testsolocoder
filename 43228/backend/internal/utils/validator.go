package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		msg := parseValidationError(err)
		Fail(c, CodeBadRequest, msg)
		return false
	}
	return true
}

func ValidateQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		msg := parseValidationError(err)
		Fail(c, CodeBadRequest, msg)
		return false
	}
	return true
}

func ValidateUri(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		msg := parseValidationError(err)
		Fail(c, CodeBadRequest, msg)
		return false
	}
	return true
}

func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}

func parseValidationError(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var msgs []string
		for _, e := range validationErrs {
			msgs = append(msgs, formatFieldError(e))
		}
		return strings.Join(msgs, "; ")
	}
	return fmt.Sprintf("参数校验失败: %v", err)
}

func formatFieldError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("字段 %s 是必填的", field)
	case "email":
		return fmt.Sprintf("字段 %s 格式不正确", field)
	case "min":
		return fmt.Sprintf("字段 %s 长度不能小于 %s", field, param)
	case "max":
		return fmt.Sprintf("字段 %s 长度不能大于 %s", field, param)
	case "len":
		return fmt.Sprintf("字段 %s 长度必须为 %s", field, param)
	case "gt":
		return fmt.Sprintf("字段 %s 必须大于 %s", field, param)
	case "gte":
		return fmt.Sprintf("字段 %s 必须大于或等于 %s", field, param)
	case "lt":
		return fmt.Sprintf("字段 %s 必须小于 %s", field, param)
	case "lte":
		return fmt.Sprintf("字段 %s 必须小于或等于 %s", field, param)
	case "oneof":
		return fmt.Sprintf("字段 %s 必须是 [%s] 之一", field, param)
	default:
		return fmt.Sprintf("字段 %s 校验失败", field)
	}
}
