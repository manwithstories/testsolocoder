package handlers

import (
	"net/http"
	"smart-energy-platform/config"
	"smart-energy-platform/middleware"
	"smart-energy-platform/models"
	"smart-energy-platform/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	Token string `json:"token" binding:"required"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=50"`
}

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
	middleware.InitJWT(cfg)
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var existing models.User
	if models.DB.Where("email = ?", req.Email).First(&existing).Error == nil {
		utils.Error(c, http.StatusConflict, 409, "Email already registered")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "Failed to process password")
		return
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Status:   1,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, "Failed to create user")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, 401, "Invalid email or password")
		return
	}

	if user.Status != 1 {
		utils.Error(c, http.StatusForbidden, 403, "Account is disabled")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, 401, "Invalid email or password")
		return
	}

	token, expireTime, err := middleware.GenerateToken(cfg, user.ID, user.Email, user.Username)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token":      token,
		"expireAt":   expireTime,
		"userId":     user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"avatar":     user.Avatar,
	})
}

func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	claims, err := middleware.ParseToken(req.Token)
	if err != nil {
		utils.Unauthorized(c, "Invalid token")
		return
	}

	var user models.User
	if err := models.DB.First(&user, claims.UserID).Error; err != nil {
		utils.Unauthorized(c, "User not found")
		return
	}

	token, expireTime, err := middleware.GenerateToken(cfg, user.ID, user.Email, user.Username)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token":    token,
		"expireAt": expireTime,
	})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("userId")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"username":  user.Username,
		"phone":     user.Phone,
		"avatar":    user.Avatar,
		"status":    user.Status,
		"createdAt": user.CreatedAt,
	})
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userId")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update profile")
		return
	}

	var user models.User
	models.DB.First(&user, userID)

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
	})
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("userId")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "Old password is incorrect")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "Failed to process password")
		return
	}

	models.DB.Model(&user).Update("password", string(hashedPassword))

	utils.Success(c, nil)
}

func getCurrentTime() time.Time {
	return time.Now()
}
