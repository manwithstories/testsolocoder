package handlers

import (
	"errand-service/internal/middleware"
	"errand-service/internal/models"
	"errand-service/internal/utils"
	"errand-service/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Nickname string `json:"nickname" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=publisher courier"`
	Code     string `json:"code"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	if !utils.ValidatePhone(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid phone number format"})
		return
	}

	var existing models.User
	if err := h.db.Where("phone = ?", req.Phone).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Phone number already registered"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to process request"})
		return
	}

	user := models.User{
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     models.UserRole(req.Role),
		Status:   models.UserStatusActive,
		Rating:   5.0,
	}

	tx := h.db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create user"})
		return
	}

	if req.Role == "courier" {
		courierProfile := models.CourierProfile{
			UserID: user.ID,
			Status: "pending",
			Level:  1,
			Rating: 5.0,
		}
		if err := tx.Create(&courierProfile).Error; err != nil {
			tx.Rollback()
			logger.Errorf("Failed to create courier profile: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create courier profile"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Registration successful",
		"data":    UserInfo{ID: user.ID, Phone: user.Phone, Nickname: user.Nickname, Role: string(user.Role)},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters"})
		return
	}

	var user models.User
	if err := h.db.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid phone or password"})
		return
	}

	if user.Status == models.UserStatusFrozen {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Account is frozen, please contact customer service"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid phone or password"})
		return
	}

	jwtSecret := c.GetString("jwt_secret")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
	}

	accessToken, err := middleware.GenerateToken(user.ID, string(user.Role), jwtSecret, 24)
	if err != nil {
		logger.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to generate token"})
		return
	}

	refreshToken, _ := middleware.GenerateToken(user.ID, string(user.Role), jwtSecret, 168)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			User: UserInfo{
				ID:       user.ID,
				Phone:    user.Phone,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
				Role:     string(user.Role),
				Status:   string(user.Status),
			},
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Refresh token is required"})
		return
	}

	jwtSecret := c.GetString("jwt_secret")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
	}

	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid or expired refresh token"})
		return
	}

	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "User not found"})
		return
	}

	accessToken, _ := middleware.GenerateToken(user.ID, string(user.Role), jwtSecret, 24)
	newRefreshToken, _ := middleware.GenerateToken(user.ID, string(user.Role), jwtSecret, 168)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		},
	})
}
