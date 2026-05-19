package service

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"
)

var (
	orderCounter uint64
	orderMu      sync.Mutex
)

func parseDate(dateStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", dateStr, time.Local)
}

func parseDateTime(dateTimeStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04", dateTimeStr, time.Local)
}

func isTimeOverlap(start1, end1, start2, end2 string) bool {
	s1, _ := time.Parse("15:04", start1)
	e1, _ := time.Parse("15:04", end1)
	s2, _ := time.Parse("15:04", start2)
	e2, _ := time.Parse("15:04", end2)

	return s1.Before(e2) && s2.Before(e1)
}

func generateOrderNo(userID uint) string {
	orderMu.Lock()
	defer orderMu.Unlock()

	now := time.Now()
	counter := atomic.AddUint64(&orderCounter, 1)

	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)

	return "ORD" + now.Format("20060102150405") + randomStr + formatUint64(counter, 4)
}

func formatUint64(n uint64, width int) string {
	const digits = "0123456789"
	buf := make([]byte, width)
	for i := width - 1; i >= 0; i-- {
		buf[i] = digits[n%10]
		n /= 10
	}
	return string(buf)
}
