package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"repair-platform/pkg/config"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, role string) (string, error) {
	cfg := config.LoadConfig()

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "repair-platform",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	cfg := config.LoadConfig()

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func GenerateOrderNo(prefix string) string {
	return prefix + time.Now().Format("20060102150405") + RandomString(6)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1)
	}
	return string(b)
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0

	dLat := (lat2 - lat1) * 3.141592653589793 / 180.0
	dLon := (lon2 - lon1) * 3.141592653589793 / 180.0

	a := dLat*dLat/2 + dLon*dLon/2*cos(lat1*3.141592653589793/180.0)*cos(lat2*3.141592653589793/180.0)/2
	c := 2 * asin(sqrt(a))

	return earthRadius * c
}

func cos(x float64) float64 {
	term := 1.0
	sum := 1.0
	for i := 1; i <= 10; i++ {
		term *= -x * x / float64((2*i-1)*(2*i))
		sum += term
	}
	return sum
}

func sin(x float64) float64 {
	term := x
	sum := x
	for i := 1; i <= 10; i++ {
		term *= -x * x / float64((2*i)*(2*i+1))
		sum += term
	}
	return sum
}

func asin(x float64) float64 {
	if x > 1.0 {
		x = 1.0
	}
	if x < -1.0 {
		x = -1.0
	}
	return x + x*x*x/6 + 3*x*x*x*x*x/40 + 5*x*x*x*x*x*x*x/112
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 20; i++ {
		z = (z + x/z) / 2
	}
	return z
}
