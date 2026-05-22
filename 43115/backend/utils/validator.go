package utils

import (
	"regexp"
	"time"
)

func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func IsValidIDCard(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}
	pattern := `^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	return matched
}

func IsValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	hasLetter := false
	hasNumber := false
	for _, ch := range password {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			hasLetter = true
		}
		if ch >= '0' && ch <= '9' {
			hasNumber = true
		}
	}
	return hasLetter && hasNumber
}

func IsFutureTime(t time.Time) bool {
	return t.After(time.Now())
}

func IsWithinRange(value, min, max float64) bool {
	return value >= min && value <= max
}

func ValidateRating(rating int) bool {
	return rating >= 1 && rating <= 5
}

func ValidatePrice(price float64) bool {
	return price >= 0
}

func ValidateDuration(duration int) bool {
	return duration >= 30 && duration <= 480
}
