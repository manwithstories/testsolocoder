package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateBookingNo() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	random := rand.Intn(10000)
	return fmt.Sprintf("BK%s%04d", now.Format("20060102150405"), random)
}

func GenerateOrderNo() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	random := rand.Intn(10000)
	return fmt.Sprintf("OR%s%04d", now.Format("20060102150405"), random)
}

func GeneratePromoCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ValidatePhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func ValidateIDCard(idCard string) bool {
	pattern := `^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	return matched
}

func CalculateDays(start, end time.Time) int {
	start = start.Truncate(time.Hour)
	end = end.Truncate(time.Hour)
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)
	if duration.Hours() > float64(days*24) {
		days++
	}
	if days <= 0 {
		days = 1
	}
	return days
}

func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func ParsePageParams(c *gin.Context) (int, int, string, string) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	if !strings.EqualFold(sortOrder, "asc") && !strings.EqualFold(sortOrder, "desc") {
		sortOrder = "desc"
	}
	return page, pageSize, sortBy, sortOrder
}

func GenerateFileName(originalName string) string {
	ext := ""
	if idx := strings.LastIndex(originalName, "."); idx >= 0 {
		ext = originalName[idx:]
	}
	timestamp := time.Now().UnixNano()
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d%s", timestamp, originalName)))
	hash := hex.EncodeToString(h.Sum(nil))[:16]
	return fmt.Sprintf("%d_%s%s", timestamp, hash, ext)
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}