package handlers

import (
	"net/http"
	"time"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"max=20"`
	Role     string `json:"role" binding:"required,oneof=beekeeper buyer inspector"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string    `json:"token"`
	Expires  time.Time `json:"expires"`
	UserInfo UserInfo  `json:"user_info"`
}

type UserInfo struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Avatar     string `json:"avatar"`
	Reputation float64 `json:"reputation"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		utils.Fail(c, http.StatusConflict, "username or email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to process password")
		return
	}

	user := models.User{
		Username:   req.Username,
		Password:   string(hashedPassword),
		Email:      req.Email,
		Phone:      req.Phone,
		Role:       req.Role,
		Reputation: 5.0,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, 24)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.Success(c, LoginResponse{
		Token:   token,
		Expires: time.Now().Add(24 * time.Hour),
		UserInfo: UserInfo{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Role:       user.Role,
			Avatar:     user.Avatar,
			Reputation: user.Reputation,
		},
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "user not found")
		return
	}

	utils.Success(c, UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Role:       user.Role,
		Avatar:     user.Avatar,
		Reputation: user.Reputation,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Phone   *string `json:"phone" binding:"omitempty,max=20"`
		Address *string `json:"address"`
		Avatar  *string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "user not found")
		return
	}

	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Address != nil {
		user.Address = *req.Address
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	database.DB.Save(&user)

	utils.Success(c, UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Role:       user.Role,
		Avatar:     user.Avatar,
		Reputation: user.Reputation,
	})
}
