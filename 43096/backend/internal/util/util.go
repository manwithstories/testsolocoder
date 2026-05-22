package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"ticket-system/internal/config"
)

type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID uint64, username, role string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ticket-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GenerateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("ORD%s%06d", now.Format("20060102150405"), GenerateRandomInt(100000, 999999))
}

func GenerateTicketNo() string {
	now := time.Now()
	return fmt.Sprintf("TKT%s%06d", now.Format("20060102150405"), GenerateRandomInt(100000, 999999))
}

func GenerateRefundNo() string {
	now := time.Now()
	return fmt.Sprintf("REF%s%06d", now.Format("20060102150405"), GenerateRandomInt(100000, 999999))
}

func GenerateCouponCode() string {
	return strings.ToUpper(GenerateRandomString(12))
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[idx.Int64()]
	}
	return string(b)
}

func GenerateRandomInt(min, max int) int {
	if min >= max {
		return min
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

func ValidateIDCard(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}

	regex := regexp.MustCompile(`^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
	if !regex.MatchString(idCard) {
		return false
	}

	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i := 0; i < 17; i++ {
		digit := int(idCard[i] - '0')
		sum += digit * weights[i]
	}

	mod := sum % 11
	expectedCode := checkCodes[mod]
	actualCode := strings.ToUpper(string(idCard[17]))

	return expectedCode == actualCode
}

func ValidatePhone(phone string) bool {
	regex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return regex.MatchString(phone)
}

func ValidateEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func ValidateRealName(name string) bool {
	if len(name) < 2 || len(name) > 20 {
		return false
	}
	regex := regexp.MustCompile(`^[\u4e00-\u9fa5a-zA-Z·]+$`)
	return regex.MatchString(name)
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func CalculatePrice(originalPrice float64, discount float64) float64 {
	return originalPrice * discount
}

func RoundFloat(f float64, n int) float64 {
	shift := 1.0
	for i := 0; i < n; i++ {
		shift *= 10
	}
	return float64(int(f*shift+0.5)) / shift
}

func GetAgeFromIDCard(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	birthYear, _ := time.Parse("2006", idCard[6:10])
	now := time.Now()
	age := now.Year() - birthYear.Year()
	birthMonth, _ := time.Parse("01", idCard[10:12])
	birthDay, _ := time.Parse("02", idCard[12:14])
	if now.Month() < birthMonth.Month() || (now.Month() == birthMonth.Month() && now.Day() < birthDay.Day()) {
		age--
	}
	return age
}

func GetGenderFromIDCard(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	genderDigit := int(idCard[16] - '0')
	if genderDigit%2 == 1 {
		return 1
	}
	return 2
}
