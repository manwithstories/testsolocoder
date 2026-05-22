package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func ValidatePhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func TimeOverlap(start1, end1, start2, end2 string) bool {
	format := "15:04"
	s1, _ := time.Parse(format, start1)
	e1, _ := time.Parse(format, end1)
	s2, _ := time.Parse(format, start2)
	e2, _ := time.Parse(format, end2)
	return s1.Before(e2) && s2.Before(e1)
}

func AddMinutesToTime(timeStr string, minutes int) string {
	format := "15:04"
	t, _ := time.Parse(format, timeStr)
	t = t.Add(time.Duration(minutes) * time.Minute)
	return t.Format(format)
}

func TimeToMinutes(timeStr string) int {
	format := "15:04"
	t, _ := time.Parse(format, timeStr)
	return t.Hour()*60 + t.Minute()
}

func MinutesToTime(mins int) string {
	hours := mins / 60
	minutes := mins % 60
	return time.Date(0, 1, 1, hours, minutes, 0, 0, time.UTC).Format("15:04")
}
