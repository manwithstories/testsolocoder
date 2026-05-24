package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateOrderNo() string {
	prefix := "ORD"
	timestamp := time.Now().Format("20060102150405")
	random := rand.Intn(9000) + 1000
	return prefix + timestamp + strconv.Itoa(random)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}

func IsValidISBN(isbn string) bool {
	if len(isbn) != 10 && len(isbn) != 13 {
		return false
	}

	for _, c := range isbn {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func Paginate(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return (page - 1) * pageSize, pageSize
}

func CalculateRating(currentRating float64, ratingCount int, newRating int) (float64, int) {
	total := currentRating * float64(ratingCount)
	total += float64(newRating)
	newCount := ratingCount + 1
	calculatedRating := total / float64(newCount)
	return calculatedRating, newCount
}
