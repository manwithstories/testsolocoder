package handlers

import (
	"net/http"
	"photo-rental/internal/middleware"
	"photo-rental/internal/models"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Role     string `json:"role" binding:"required,oneof=renter owner admin"`
	RealName string `json:"realName"`
	IDCard   string `json:"idCard"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	User     interface{} `json:"user"`
	ExpiresAt int64     `json:"expiresAt"`
}

type UpdateProfileRequest struct {
	RealName string `json:"realName"`
	IDCard   string `json:"idCard"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser models.User
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already registered")
		return
	}

	if result := database.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Username already taken")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		RealName: req.RealName,
		IDCard:   req.IDCard,
		Phone:    req.Phone,
		Verified: false,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		utils.Logger.Error("Failed to create user: %v", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.Logger.Info("New user registered: %s (%s)", user.Username, user.Email)
	utils.SuccessResponse(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := utils.VerifyPassword(req.Password, user.Password); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := middleware.JWTManager.GenerateToken(user.ID, user.Username, user.Role, user.Verified)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	utils.Logger.Info("User logged in: %s", user.Email)
	utils.SuccessResponse(c, LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userId")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	user.Password = ""
	utils.SuccessResponse(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userId")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.IDCard != "" {
		updates["id_card"] = req.IDCard
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "No fields to update")
		return
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		utils.Logger.Error("Failed to update profile: %v", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	var user models.User
	database.DB.First(&user, userID)
	user.Password = ""

	utils.Logger.Info("User profile updated: %d", userID)
	utils.SuccessResponse(c, user)
}

func (h *UserHandler) VerifyUser(c *gin.Context) {
	userRole := c.GetString("role")

	if userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only admin can verify users")
		return
	}

	var req struct {
		UserID   uint `json:"userId" binding:"required"`
		Verified bool `json:"verified"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", req.UserID).Update("verified", req.Verified)
	if result.Error != nil {
		utils.Logger.Error("Failed to verify user: %v", result.Error)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to verify user")
		return
	}

	utils.Logger.Info("User %d verification status updated to: %v", req.UserID, req.Verified)
	utils.SuccessResponse(c, nil)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	userRole := c.GetString("role")
	if userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only admin can view all users")
		return
	}

	var users []models.User
	database.DB.Find(&users)

	for i := range users {
		users[i].Password = ""
		users[i].IDCard = ""
	}

	utils.SuccessResponse(c, users)
}
