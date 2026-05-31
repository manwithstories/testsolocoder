package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ship-rental-platform/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func CalculateDays(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}
	return int(math.Ceil(end.Sub(start).Hours() / 24))
}

func CalculateHours(start, end time.Time) float64 {
	if end.Before(start) {
		return 0
	}
	return end.Sub(start).Hours()
}

func RoundTo(value float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*shift) / shift
}

func ConvertCurrency(amount float64, from, to string, rates map[string]float64) float64 {
	if from == to {
		return amount
	}

	fromRate, exists := rates[from]
	if !exists {
		return amount
	}

	usdAmount := amount / fromRate

	toRate, exists := rates[to]
	if !exists {
		return usdAmount
	}

	return RoundTo(usdAmount*toRate, 2)
}

func SaveUploadedFile(c *gin.Context, fieldName string, uploadPath string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileID := uuid.New().String()
	newFilename := fmt.Sprintf("%s%s", fileID, ext)
	filePath := filepath.Join(uploadPath, newFilename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}

func IsAllowedFileType(filename string, allowedTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ContainsString(allowedTypes, ext)
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, -1)
}

func IsTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

func GetConfig() *config.Config {
	return config.AppConfig
}
