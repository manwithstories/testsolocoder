package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateOrderNo(prefix string) string {
	return strings.ToUpper(fmt.Sprintf("%s%s%d", prefix, time.Now().Format("20060102150405"), time.Now().UnixNano()%1000))
}

func GenerateUUID() string {
	return uuid.New().String()
}

func StringInSlice(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

func DaysBetween(start, end time.Time) int {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	return int(end.Sub(start).Hours()/24) + 1
}
