package handlers

import (
	"strconv"
	"time"
)

func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, parseError("invalid date format")
}

type parseErrorType string

func parseError(msg string) error {
	return parseErrorType(msg)
}

func (e parseErrorType) Error() string {
	return string(e)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
