package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func GenerateOrderNo(userID uint) string {
	now := time.Now()
	h := md5.New()
	h.Write([]byte(now.String()))
	h.Write([]byte(string(userID)))
	return hex.EncodeToString(h.Sum(nil))[:20]
}

func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func CalculateServiceFee(reward float64) float64 {
	feeRate := 0.10
	if reward >= 100 {
		feeRate = 0.08
	}
	return math.Round(reward*feeRate*100) / 100
}
