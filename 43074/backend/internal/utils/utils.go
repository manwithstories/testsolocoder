package utils

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUintParam(c *gin.Context, key string) (uint64, error) {
	idStr := c.Param(key)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetIntQuery(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetStringQuery(c *gin.Context, key string, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func CalculateProgress(currentPage, totalPages int) float64 {
	if totalPages <= 0 {
		return 0
	}
	if currentPage >= totalPages {
		return 100
	}
	if currentPage <= 0 {
		return 0
	}
	return float64(currentPage) / float64(totalPages) * 100
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func StartOfYear(year int) time.Time {
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
}

func EndOfYear(year int) time.Time {
	return time.Date(year, 12, 31, 23, 59, 59, 999999999, time.Local)
}

func StartOfMonth(year, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
}

func EndOfMonth(year, month int) time.Time {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)
	return time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 999999999, time.Local)
}

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateTimeStr)
}
