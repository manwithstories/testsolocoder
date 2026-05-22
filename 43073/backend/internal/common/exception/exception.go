package exception

import "fmt"

type BusinessException struct {
	Code    int
	Message string
}

func (e *BusinessException) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewBusinessException(code int, message string) *BusinessException {
	return &BusinessException{
		Code:    code,
		Message: message,
	}
}

func New(message string) *BusinessException {
	return &BusinessException{
		Code:    400,
		Message: message,
	}
}
