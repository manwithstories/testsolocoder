package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type RegisterRequest struct {
	Username string          `json:"username" binding:"required,min=3,max=50"`
	Email    string          `json:"email" binding:"required,email"`
	Phone    string          `json:"phone"`
	Password string          `json:"password" binding:"required,min=8"`
	RealName string          `json:"real_name" binding:"required"`
	Role     models.UserRole `json:"role" binding:"required"`
	Company  string          `json:"company"`
	Address  string          `json:"address"`
	IDCard   string          `json:"id_card"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token  string      `json:"token"`
	User   interface{} `json:"user"`
	Expire time.Time   `json:"expire"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	if !utils.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid email format",
		})
		return
	}

	if req.Phone != "" && !utils.ValidatePhone(req.Phone) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid phone format",
		})
		return
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	validRoles := map[models.UserRole]bool{
		models.RoleEmployer:  true,
		models.RoleAgent:     true,
		models.RoleTemporary: true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid role. Must be employer, agent, or temporary",
		})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.Response{
			Code:    409,
			Message: "Username or email already exists",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to process password",
		})
		return
	}

	user := models.User{
		Username:    req.Username,
		Email:       req.Email,
		Phone:       req.Phone,
		Password:    hashedPassword,
		RealName:    req.RealName,
		Role:        req.Role,
		Company:     req.Company,
		Address:     req.Address,
		IDCard:      req.IDCard,
		Status:      models.UserStatusActive,
		CreditScore: 100,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Registration successful",
		Data:    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ? OR email = ? OR phone = ?", req.Login, req.Login, req.Login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "Invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Internal server error",
		})
		return
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    401,
			Message: "Invalid credentials",
		})
		return
	}

	if user.Status != models.UserStatusActive {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "Account is not active",
		})
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to generate token",
		})
		return
	}

	now := time.Now()
	database.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": c.ClientIP(),
	})

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Login successful",
		Data: LoginResponse{
			Token: token,
			User:  user,
			Expire: now.Add(24 * time.Hour),
		},
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data:    user,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uuid.UUID)

	var req struct {
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Company  string `json:"company"`
		Address  string `json:"address"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "No valid fields to update",
		})
		return
	}

	database.DB.Model(&models.User{}).Where("id = ?", uid).Updates(updates)

	var updatedUser models.User
	database.DB.First(&updatedUser, uid)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Profile updated successfully",
		Data:    updatedUser,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uuid.UUID)

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "User not found",
		})
		return
	}

	if !utils.CheckPassword(user.Password, req.OldPassword) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Old password is incorrect",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to process password",
		})
		return
	}

	database.DB.Model(&user).Update("password", hashedPassword)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Password changed successfully",
	})
}
