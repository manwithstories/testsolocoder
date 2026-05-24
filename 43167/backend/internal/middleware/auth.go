package middleware

import (
	"net/http"
	"strings"
	"time"

	"watchplatform/internal/app"
	"watchplatform/internal/config"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Claims struct {
	UserID uint       `json:"uid"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(u *model.User) (string, error) {
	claims := Claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Cfg.JWTExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.Cfg.AppName,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(config.Cfg.JWTSecret))
}

func JWTAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			app.Fail(c, http.StatusUnauthorized, "缺少 Authorization")
			c.Abort()
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		claims := &Claims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JWTSecret), nil
		})
		if err != nil {
			app.Fail(c, http.StatusUnauthorized, "token 无效")
			c.Abort()
			return
		}
		var user model.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			app.Fail(c, http.StatusUnauthorized, "用户不存在")
			c.Abort()
			return
		}
		c.Set("user", &user)
		c.Next()
	}
}

func RequireRoles(roles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := c.Get("user")
		if !ok {
			app.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		user := u.(*model.User)
		for _, r := range roles {
			if user.Role == r {
				c.Next()
				return
			}
		}
		app.Fail(c, http.StatusForbidden, "权限不足")
		c.Abort()
	}
}

func CurrentUser(c *gin.Context) *model.User {
	u, _ := c.Get("user")
	if u == nil {
		return nil
	}
	return u.(*model.User)
}
