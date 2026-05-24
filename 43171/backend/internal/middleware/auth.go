package middleware

import (
	"net/http"
	"time"

	"drone-rental/internal/config"
	"drone-rental/internal/model"
	"drone-rental/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint      `json:"user_id"`
	Username string    `json:"username"`
	Role     model.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Cfg.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "drone-rental",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}
		if tokenStr == "" {
			response.ErrAuth(c, "未登录或Token无效")
			c.Abort()
			return
		}
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}
		claims, err := ParseToken(tokenStr)
		if err != nil {
			response.ErrAuth(c, "Token已过期或无效")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleAuth(roles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.ErrAuth(c, "未登录")
			c.Abort()
			return
		}
		userRole := role.(model.Role)
		allowed := false
		for _, r := range roles {
			if userRole == r {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, response.Response{
				Code:    403,
				Message: "权限不足",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	id, _ := c.Get("user_id")
	return id.(uint)
}

func GetRole(c *gin.Context) model.Role {
	r, _ := c.Get("role")
	return r.(model.Role)
}

func GetUsername(c *gin.Context) string {
	u, _ := c.Get("username")
	return u.(string)
}
