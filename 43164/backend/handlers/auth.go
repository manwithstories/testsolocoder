package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/middleware"
	"tutoring-platform/models"
)

type RegisterRequest struct {
	Email     string         `json:"email" binding:"required,email"`
	Password  string         `json:"password" binding:"required,min=6"`
	FirstName string         `json:"firstName" binding:"required"`
	LastName  string         `json:"lastName" binding:"required"`
	Phone     string         `json:"phone"`
	Role      models.UserRole `json:"role" binding:"required,oneof=teacher student"`
	Timezone  string         `json:"timezone"`
	Language  string         `json:"language"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token  string      `json:"token"`
	User   models.User `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatarUrl"`
	Timezone  string `json:"timezone"`
	Language  string `json:"language"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if database.DB.Where("email = ?", req.Email).First(&existingUser).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      req.Role,
		Timezone:  req.Timezone,
		Language:  req.Language,
		Status:    models.UserStatusPending,
	}

	if user.Timezone == "" {
		user.Timezone = "UTC"
	}
	if user.Language == "" {
		user.Language = "en"
	}

	tx := database.DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if user.Role == models.RoleTeacher {
		teacherProfile := models.TeacherProfile{
			UserID:         user.ID,
			ApprovalStatus: "pending",
		}
		if err := tx.Create(&teacherProfile).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teacher profile"})
			return
		}
	} else if user.Role == models.RoleStudent {
		studentProfile := models.StudentProfile{
			UserID:           user.ID,
			AssessmentStatus: "not_started",
		}
		if err := tx.Create(&studentProfile).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student profile"})
			return
		}
	}

	wallet := models.Wallet{
		UserID:   user.ID,
		Balance:  0,
		Currency: "USD",
	}
	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	tx.Commit()

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Preload("TeacherProfile").Preload("StudentProfile").Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !user.ComparePassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if user.Status == models.UserStatusBanned {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is banned"})
		return
	}

	if user.Status == models.UserStatusInactive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is inactive"})
		return
	}

	now := time.Now()
	database.DB.Model(&user).Update("last_login_at", now)

	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user.Password = ""

	cfg := middleware.ParseToken
	_ = cfg

	c.JSON(http.StatusOK, AuthResponse{
		Token:  token,
		User:   user,
		ExpiresAt: now.Add(24 * time.Hour),
	})
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var user models.User
	if err := database.DB.Preload("TeacherProfile.Subjects.Subject").Preload("TeacherProfile.Availabilities").Preload("StudentProfile.LearningGoals.Subject").Preload("Wallet").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Timezone != "" {
		updates["timezone"] = req.Timezone
	}
	if req.Language != "" {
		updates["language"] = req.Language
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.DB.Preload("TeacherProfile.Subjects.Subject").Preload("StudentProfile.LearningGoals.Subject").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
