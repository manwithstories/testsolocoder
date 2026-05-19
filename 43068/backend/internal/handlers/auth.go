package handlers

import (
	"net/http"

	"freelancer-management/internal/database"
	"freelancer-management/internal/logger"
	"freelancer-management/internal/middleware"
	"freelancer-management/internal/models"
	"freelancer-management/internal/utils"
	jwtpkg "freelancer-management/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{db: database.GetDB()}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !utils.IsValidEmail(req.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid email format")
		return
	}

	valid, msg := utils.CheckPasswordStrength(req.Password)
	if !valid {
		utils.ErrorResponse(c, http.StatusBadRequest, msg)
		return
	}

	var existingUser models.User
	result := h.db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already registered")
		return
	}

	user := models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	if err := user.SetPassword(req.Password); err != nil {
		logger.LogError("Failed to hash password: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		logger.LogError("Failed to create user: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	logger.LogOperation(user.ID, "register", "User registered successfully")

	tokens, err := jwtpkg.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		logger.LogError("Failed to generate tokens: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"user": UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		"tokens": tokens,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	result := h.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	logger.LogOperation(user.ID, "login", "User logged in")

	tokens, err := jwtpkg.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		logger.LogError("Failed to generate tokens: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"user": UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		"tokens": tokens,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokens, err := jwtpkg.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	utils.SuccessResponse(c, tokens)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}
