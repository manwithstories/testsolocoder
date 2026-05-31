package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"health-platform/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWT Claims struct {
	UserID   uint              `json:"user_id"`
	Username string            `json:"username"`
	Role     string            `json:"role"`
	CompanyID *uint           `json:"company_id"`
	AgencyID  *uint           `json:"agency_id"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID uint, username, role string, companyID, agencyID *uint) (string, error) {
	claims := JWT{
		UserID:   userID,
		Username: username,
		Role:     role,
		CompanyID: companyID,
		AgencyID:  agencyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GlobalConfig.JWT.ExpireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "health-platform",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

func ParseToken(tokenString string) (*JWT, error) {
	claims := &JWT{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GenerateOrderNo(prefix string) string {
	now := time.Now()
	randBytes := make([]byte, 4)
	rand.Read(randBytes)
	return fmt.Sprintf("%s%s%06d%s",
		prefix,
		now.Format("20060102150405"),
		now.Nanosecond()/1000%1000000,
		hex.EncodeToString(randBytes))
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandomString(length int) string {
	randBytes := make([]byte, length/2)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}

func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:4] + "********" + idCard[len(idCard)-4:]
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func RoundFloat(val float64, precision int) float64 {
	divisor := 1.0
	for i := 0; i < precision; i++ {
		divisor *= 10
	}
	return float64(int(val*divisor)) / divisor
}

func CalcBMI(height, weight float64) float64 {
	if height <= 0 {
		return 0
	}
	heightM := height / 100.0
	return RoundFloat(weight/(heightM*heightM), 2)
}

func AgeFromBirthday(birthday *time.Time) int {
	if birthday == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}
