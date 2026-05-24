package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	Nickname string    `json:"nickname"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func GenerateToken(userID uuid.UUID, role, nickname, secret string, accessExpire, refreshExpire int) (*TokenPair, error) {
	now := time.Now()

	accessClaims := Claims{
		UserID:   userID,
		Role:     role,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(accessExpire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "skillshare",
			Subject:   userID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshClaims := Claims{
		UserID:   userID,
		Role:     role,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(refreshExpire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "skillshare",
			Subject:   userID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    accessExpire * 3600,
	}, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

func RefreshToken(refreshToken, secret string, accessExpire, refreshExpire int) (*TokenPair, error) {
	claims, err := ValidateToken(refreshToken, secret)
	if err != nil {
		return nil, fmt.Errorf("无效的刷新令牌: %w", err)
	}

	return GenerateToken(claims.UserID, claims.Role, claims.Nickname, secret, accessExpire, refreshExpire)
}
