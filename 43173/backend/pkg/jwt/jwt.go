package jwt

import (
	"errors"
	"time"

	"music-platform/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, role string) (string, error) {
	cfg := config.AppConfig.JWT

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.AppConfig.JWT

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetUserID(c *gin.Context) uint {
	claims, exists := c.Get("claims")
	if !exists {
		return 0
	}
	return claims.(*Claims).UserID
}

func GetUsername(c *gin.Context) string {
	claims, exists := c.Get("claims")
	if !exists {
		return ""
	}
	return claims.(*Claims).Username
}

func GetUserRole(c *gin.Context) string {
	claims, exists := c.Get("claims")
	if !exists {
		return ""
	}
	return claims.(*Claims).Role
}
