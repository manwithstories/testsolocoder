package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"travel-planner/internal/database"
	"travel-planner/internal/logger"
	"travel-planner/internal/middleware"
	"travel-planner/internal/models"
	"travel-planner/internal/utils"
)

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=50"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid register request: %v", err)
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "User with this email or username already exists")
		return
	}

	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := user.HashPassword(req.Password); err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	tx := database.BeginTransaction()
	if tx.Error != nil {
		utils.InternalServerError(c, "Database error")
		return
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to create user: %v", err)
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	var userRole models.Role
	if err := tx.Where("name = ?", "user").First(&userRole).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to find user role: %v", err)
		utils.InternalServerError(c, "Failed to assign role")
		return
	}

	userRoleAssociation := models.UserRole{
		UserID: user.ID,
		RoleID: userRole.ID,
	}
	if err := tx.Create(&userRoleAssociation).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to assign role: %v", err)
		utils.InternalServerError(c, "Failed to assign role")
		return
	}

	if err := tx.Commit().Error; err != nil {
		logger.Errorf("Failed to commit transaction: %v", err)
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	logger.Infof("User registered successfully: %s", user.Email)
	utils.Created(c, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "Validation failed", utils.FormatValidationErrors(err))
		return
	}

	var user models.User
	if err := database.DB.Preload("Roles").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Unauthorized(c, "Invalid email or password")
			return
		}
		logger.Errorf("Database error: %v", err)
		utils.InternalServerError(c, "Database error")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.Unauthorized(c, "Invalid email or password")
		return
	}

	if !user.IsActive {
		utils.Forbidden(c, "Account is disabled")
		return
	}

	token, expiresAt, err := middleware.GenerateToken(user)
	if err != nil {
		logger.Errorf("Failed to generate token: %v", err)
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	logger.Infof("User logged in: %s", user.Email)
	utils.Success(c, AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserInfo{
			ID:        user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	utils.Success(c, UserInfo{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	})
}

func UpdateProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	allowedFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"avatar":     true,
		"timezone":   true,
		"language":   true,
		"currency":   true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(filteredUpdates).Error; err != nil {
		logger.Errorf("Failed to update profile: %v", err)
		utils.InternalServerError(c, "Failed to update profile")
		return
	}

	var updatedUser models.User
	database.DB.First(&updatedUser, userID)

	logger.Infof("User profile updated: %s", userID)
	utils.Success(c, UserInfo{
		ID:        updatedUser.ID.String(),
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Avatar:    updatedUser.Avatar,
	})
}
