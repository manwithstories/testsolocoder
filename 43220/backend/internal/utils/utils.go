package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, username, role, secret string, expireHour int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pet-board",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func GenerateOrderNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d",
		prefix,
		now.Format("20060102150405"),
		rand.Intn(1000000),
	)
}

func HashAmount(amount float64, salt string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%.2f-%s", amount, salt)))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyAmountHash(amount float64, salt, hash string) bool {
	return HashAmount(amount, salt) == hash
}

func DaysBetween(start, end time.Time) int {
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	return int(endDate.Sub(startDate).Hours()/24) + 1
}

func IsDateRangeOverlap(start1, end1, start2, end2 time.Time) bool {
	return !start1.After(end2) && !start2.After(end1)
}

func ParseID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func StringPtr(s string) *string {
	return &s
}

func FloatPtr(f float64) *float64 {
	return &f
}

func IntPtr(i int) *int {
	return &i
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func AtoiOrZero(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
