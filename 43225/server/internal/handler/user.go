package handler

import (
	"net/http"
	"time"

	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/middleware"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser model.User
	if err := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		utils.Error(c, http.StatusConflict, "Email or username already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		FullName: req.FullName,
		Phone:    req.Phone,
		Company:  req.Company,
		Timezone: "UTC",
		IsActive: true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.Created(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"full_name": user.FullName,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user model.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !user.IsActive {
		utils.Error(c, http.StatusForbidden, "Account is deactivated")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	token, err := middleware.GenerateToken(&user, 24)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"full_name":  user.FullName,
			"avatar_url": user.AvatarURL,
			"timezone":   user.Timezone,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user model.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, gin.H{
		"id":               user.ID,
		"username":         user.Username,
		"email":            user.Email,
		"role":             user.Role,
		"full_name":        user.FullName,
		"phone":            user.Phone,
		"avatar_url":       user.AvatarURL,
		"company":          user.Company,
		"address":          user.Address,
		"city":             user.City,
		"country":          user.Country,
		"timezone":         user.Timezone,
		"is_active":        user.IsActive,
		"is_email_verified": user.IsEmailVerified,
		"last_login_at":    user.LastLoginAt,
		"created_at":       user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user model.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Company != "" {
		user.Company = req.Company
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.Country != "" {
		user.Country = req.Country
	}
	if req.Timezone != "" {
		user.Timezone = req.Timezone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.Success(c, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":         user.ID,
			"full_name":  user.FullName,
			"phone":      user.Phone,
			"company":    user.Company,
			"address":    user.Address,
			"city":       user.City,
			"country":    user.Country,
			"timezone":   user.Timezone,
			"avatar_url": user.AvatarURL,
		},
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user model.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		utils.Error(c, http.StatusBadRequest, "Current password is incorrect")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	utils.Success(c, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []model.User
	query := database.DB.Select("id, username, email, role, full_name, avatar_url, is_active, created_at")

	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	query.Model(&model.User{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "10"), 10)
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.Paginated(c, users, total, page, pageSize)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user model.User
	if err := database.DB.Select("id, username, email, role, full_name, avatar_url, company, is_active, created_at").First(&user, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) ToggleUserStatus(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userID, _ := uuid.Parse(id)

	var user model.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	user.IsActive = !user.IsActive
	if err := database.DB.Save(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update user status")
		return
	}

	utils.Success(c, gin.H{
		"message":   "User status updated successfully",
		"is_active": user.IsActive,
	})
}
