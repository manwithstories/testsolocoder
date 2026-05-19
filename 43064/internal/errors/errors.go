package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	CodeSuccess ErrorCode = 0

	CodeBadRequest          ErrorCode = 40000
	CodeInvalidParams       ErrorCode = 40001
	CodeMissingRequiredField  ErrorCode = 40002
	CodeInvalidFormat       ErrorCode = 40003

	CodeUnauthorized        ErrorCode = 40100
	CodeInvalidToken        ErrorCode = 40101

	CodeForbidden           ErrorCode = 40300
	CodePermissionDenied   ErrorCode = 40301

	CodeNotFound            ErrorCode = 40400
	CodeChannelNotFound     ErrorCode = 40401
	CodeTemplateNotFound  ErrorCode = 40402
	CodeRecipientNotFound ErrorCode = 40403

	CodeConflict            ErrorCode = 40900
	CodeDuplicateChannel    ErrorCode = 40901

	CodeTooManyRequests     ErrorCode = 42900
	CodeRateLimitExceeded  ErrorCode = 42901

	CodeInternalServerError ErrorCode = 50000
	CodeDatabaseError     ErrorCode = 50001
	CodeChannelError      ErrorCode = 50002
	CodeTemplateError   ErrorCode = 50003
	CodeSendError       ErrorCode = 50004
	CodeQueueError      ErrorCode = 50005
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Stack   string    `json:"stack,omitempty"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) HTTPStatus() int {
	switch e.Code / 100 {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 409:
		return http.StatusConflict
	case 429:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func New(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func WithDetails(e *AppError, details string) *AppError {
	e.Details = details
	return e
}

func WithStack(e *AppError, stack string) *AppError {
	e.Stack = stack
	return e
}

func BadRequest(message string, err error) *AppError {
	return New(CodeBadRequest, message, err)
}

func InvalidParams(message string, err error) *AppError {
	return New(CodeInvalidParams, message, err)
}

func MissingRequiredField(field string) *AppError {
	return New(CodeMissingRequiredField, fmt.Sprintf("missing required field: %s", field), nil)
}

func InvalidFormat(field string, err error) *AppError {
	return New(CodeInvalidFormat, fmt.Sprintf("invalid format for %s", field), err)
}

func Unauthorized(message string, err error) *AppError {
	return New(CodeUnauthorized, message, err)
}

func InvalidToken(err error) *AppError {
	return New(CodeInvalidToken, "invalid or expired token", err)
}

func Forbidden(message string, err error) *AppError {
	return New(CodeForbidden, message, err)
}

func PermissionDenied() *AppError {
	return New(CodePermissionDenied, "permission denied", nil)
}

func NotFound(resource string, err error) *AppError {
	return New(CodeNotFound, fmt.Sprintf("%s not found", resource), err)
}

func ChannelNotFound(id uint) *AppError {
	return New(CodeChannelNotFound, fmt.Sprintf("channel with id %d not found", id), nil)
}

func TemplateNotFound(id uint) *AppError {
	return New(CodeTemplateNotFound, fmt.Sprintf("template with id %d not found", id), nil)
}

func RecipientNotFound(id uint) *AppError {
	return New(CodeRecipientNotFound, fmt.Sprintf("recipient with id %d not found", id), nil)
}

func Conflict(message string, err error) *AppError {
	return New(CodeConflict, message, err)
}

func DuplicateChannel(name string) *AppError {
	return New(CodeDuplicateChannel, fmt.Sprintf("channel with name %s already exists", name), nil)
}

func RateLimitExceeded(limit string) *AppError {
	return New(CodeRateLimitExceeded, fmt.Sprintf("rate limit exceeded: %s", limit), nil)
}

func InternalServerError(message string, err error) *AppError {
	return New(CodeInternalServerError, message, err)
}

func DatabaseError(err error) *AppError {
	return New(CodeDatabaseError, "database operation failed", err)
}

func ChannelError(message string, err error) *AppError {
	return New(CodeChannelError, message, err)
}

func TemplateError(message string, err error) *AppError {
	return New(CodeTemplateError, message, err)
}

func SendError(message string, err error) *AppError {
	return New(CodeSendError, message, err)
}

func QueueError(message string, err error) *AppError {
	return New(CodeQueueError, message, err)
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	appErr, ok := err.(*AppError)
	if !ok {
		return false
	}
	return appErr.Code == CodeSendError || appErr.Code == CodeChannelError
}
