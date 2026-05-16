package utils

import (
	"fmt"
	"strings"
	"time"
)

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func FormatHours(hours float64) string {
	return fmt.Sprintf("%.2f", hours)
}

func FormatMoney(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

func ParseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
	}

	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, dateStr, time.Local)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return StartOfDay(t.AddDate(0, 0, -(weekday - 1)))
}

func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(StartOfWeek(t).AddDate(0, 0, 6))
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	firstDay := StartOfMonth(t)
	lastDay := firstDay.AddDate(0, 1, -1)
	return EndOfDay(lastDay)
}

func CalculateHours(start, end time.Time) float64 {
	duration := end.Sub(start)
	return duration.Hours()
}

func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func JoinTags(tags []string) string {
	if len(tags) == 0 {
		return "-"
	}
	return strings.Join(tags, ", ")
}
