package utils

import (
	"errors"
	"time"

	"print3d-platform/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID     `json:"user_id"`
	Email    string        `json:"email"`
	Username string        `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func GenerateTokenPair(userID uuid.UUID, email, username string, role models.UserRole, secretKey string, accessTokenHours, refreshTokenHours int) (*TokenPair, error) {
	accessTokenExp := time.Now().Add(time.Duration(accessTokenHours) * time.Hour)
	accessClaims := Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "print3d-platform",
			Subject:   "access_token",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenExp := time.Now().Add(time.Duration(refreshTokenHours) * time.Hour)
	refreshClaims := Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "print3d-platform",
			Subject:   "refresh_token",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessTokenExp.Unix(),
	}, nil
}

func ParseToken(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshTokenPair(refreshToken, secretKey string, accessTokenHours, refreshTokenHours int) (*TokenPair, error) {
	claims, err := ParseToken(refreshToken, secretKey)
	if err != nil {
		return nil, err
	}

	if claims.Subject != "refresh_token" {
		return nil, errors.New("invalid token type")
	}

	return GenerateTokenPair(claims.UserID, claims.Email, claims.Username, claims.Role, secretKey, accessTokenHours, refreshTokenHours)
}
