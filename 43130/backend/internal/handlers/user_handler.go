package handlers

import (
	"net/http"
	"time"
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"
	"wedding-planner/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var existingUser models.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		response.BadRequest(c, "Username or email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "Failed to process password")
		return
	}

	role := req.Role
	if role == "" {
		role = "couple"
	}
	if !utils.StringInSlice(role, []string{"couple", "planner", "vendor", "admin"}) {
		role = "couple"
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     role,
		Status:   "active",
	}

	if err := db.Create(&user).Error; err != nil {
		response.InternalError(c, "Failed to create user")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.InternalError(c, "Failed to generate token")
		return
	}

	c.JSON(http.StatusCreated, response.Response{
		Code:    0,
		Message: "Registration successful",
		Data: gin.H{
			"token": token,
			"user": gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"role":      user.Role,
			},
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var user models.User
	if err := db.Where("username = ? OR email = ?", req.Account, req.Account).First(&user).Error; err != nil {
		response.Unauthorized(c, "Invalid account or password")
		return
	}

	if user.Status != "active" {
		response.Unauthorized(c, "Account is not active")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "Invalid account or password")
		return
	}

	now := time.Now()
	db.Model(&user).Update("last_login_at", now)

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.InternalError(c, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
			"avatar":    user.Avatar,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"full_name": user.FullName,
		"phone":     user.Phone,
		"avatar":    user.Avatar,
		"role":      user.Role,
		"status":    user.Status,
		"created_at": user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	result := db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"full_name": req.FullName,
		"phone":     req.Phone,
		"avatar":    req.Avatar,
	})

	if result.Error != nil {
		response.InternalError(c, "Failed to update profile")
		return
	}

	response.Success(c, gin.H{"message": "Profile updated successfully"})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		response.BadRequest(c, "Old password is incorrect")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "Failed to process password")
		return
	}

	db.Model(&user).Update("password", string(hashedPassword))

	response.Success(c, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	db := database.GetDB()

	var users []models.User
	var total int64

	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")

	search := c.Query("search")
	role := c.Query("role")

	query := db.Model(&models.User{})

	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users)

	type UserResponse struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		FullName  string `json:"full_name"`
		Role      string `json:"role"`
		Status    string `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FullName:  u.FullName,
			Role:      u.Role,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		})
	}

	response.Paginated(c, userResponses, total, page, pageSize)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	type StatusRequest struct {
		Status string `json:"status" binding:"required,oneof=active disabled"`
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid status")
		return
	}

	db := database.GetDB()

	result := db.Model(&models.User{}).Where("id = ?", userID).Update("status", req.Status)
	if result.Error != nil {
		response.InternalError(c, "Failed to update user status")
		return
	}

	response.Success(c, gin.H{"message": "User status updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()

	result := db.Delete(&models.User{}, userID)
	if result.Error != nil {
		response.InternalError(c, "Failed to delete user")
		return
	}

	response.Success(c, gin.H{"message": "User deleted successfully"})
}

var _ = dbTransactionCheck

func dbTransactionCheck(db *gorm.DB) {
	_ = db.Transaction
}
