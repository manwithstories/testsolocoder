package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"sports-league/models"
)

type Claims struct {
	UserID   uint           `json:"user_id"`
	Email    string         `json:"email"`
	Role     models.Role    `json:"role"`
	TeamID   *uint          `json:"team_id,omitempty"`
	FullName string         `json:"full_name"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	hours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if hours <= 0 {
		hours = 72
	}
	claims := Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		TeamID:   user.TeamID,
		FullName: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(hours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sports-league",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
