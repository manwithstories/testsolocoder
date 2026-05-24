package validator

import (
	"regexp"
	"time"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^1[3-9]\d{9}$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9@$!%*?&]{6,20}$`)
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ValidationErrors []*ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	return e[0].Message
}

func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "邮箱不能为空"}
	}
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "邮箱格式不正确"}
	}
	return nil
}

func ValidatePhone(phone string) error {
	if phone == "" {
		return &ValidationError{Field: "phone", Message: "手机号不能为空"}
	}
	if !phoneRegex.MatchString(phone) {
		return &ValidationError{Field: "phone", Message: "手机号格式不正确"}
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return &ValidationError{Field: "password", Message: "密码不能为空"}
	}
	if !passwordRegex.MatchString(password) {
		return &ValidationError{Field: "password", Message: "密码必须是6-20位字母、数字或特殊字符"}
	}
	return nil
}

func ValidateNickname(nickname string) error {
	if nickname == "" {
		return &ValidationError{Field: "nickname", Message: "昵称不能为空"}
	}
	if len(nickname) > 50 {
		return &ValidationError{Field: "nickname", Message: "昵称不能超过50个字符"}
	}
	return nil
}

func ValidatePrice(price float64) error {
	if price < 0 {
		return &ValidationError{Field: "price", Message: "价格不能为负数"}
	}
	return nil
}

func ValidateTimeSlot(startTime, endTime time.Time) error {
	if startTime.After(endTime) {
		return &ValidationError{Field: "time", Message: "开始时间不能晚于结束时间"}
	}
	if startTime.Before(time.Now()) {
		return &ValidationError{Field: "time", Message: "预约时间不能早于当前时间"}
	}
	return nil
}

func ValidateBookingConflict(existingBookings []struct{ Start, End time.Time }, start, end time.Time) error {
	for _, booking := range existingBookings {
		if start.Before(booking.End) && end.After(booking.Start) {
			return &ValidationError{Field: "time", Message: "预约时间与已有预约冲突"}
		}
	}
	return nil
}

func ValidateRating(rating int) error {
	if rating < 1 || rating > 5 {
		return &ValidationError{Field: "rating", Message: "评分必须在1-5之间"}
	}
	return nil
}

func ValidateReviewContent(content string) error {
	if len(content) > 1000 {
		return &ValidationError{Field: "content", Message: "评价内容不能超过1000个字符"}
	}
	return nil
}
