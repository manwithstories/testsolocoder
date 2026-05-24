package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(bytes), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidatePhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password must be at most 128 characters")
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else if strings.ContainsRune("!@#$%^&*()_+-=[]{}|;':\",./<>?", char) {
			hasSpecial = true
		}
	}

	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func GenerateQRCode(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data + time.Now().String()))
	return hex.EncodeToString(hash.Sum(nil))[:32]
}

func CalculateWorkHours(startTime, endTime string) float64 {
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return 0
	}
	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return 0
	}

	duration := end.Sub(start)
	hours := duration.Hours()

	if hours > 24 {
		return 0
	}
	return hours
}

func CalculateOvertimeHours(regularHours, actualHours float64) float64 {
	if actualHours > regularHours {
		return actualHours - regularHours
	}
	return 0
}
