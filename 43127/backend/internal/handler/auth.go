package handler

import (
	"property-management/internal/database"
	"property-management/internal/model"
	"property-management/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	JWTSecret string
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{JWTSecret: jwtSecret}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"realName" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
	Role     string `json:"role" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	var user model.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Error(c, 401, "Invalid username or password")
		return
	}

	if user.Status != 1 {
		utils.Error(c, 403, "Account is disabled")
		return
	}

	if !database.VerifyPassword(req.Password, user.Password) {
		utils.Error(c, 401, "Invalid username or password")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role, h.JWTSecret)
	if err != nil {
		utils.Error(c, 500, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"realName": user.RealName,
			"role":     user.Role,
			"avatar":   user.Avatar,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	var existing model.User
	if database.DB.Where("username = ?", req.Username).First(&existing).Error == nil {
		utils.Error(c, 400, "Username already exists")
		return
	}

	hashedPassword, _ := hashPassword(req.Password)
	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
		Role:     req.Role,
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.Error(c, 500, "Failed to create user")
		return
	}

	utils.Success(c, nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, 404, "User not found")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"realName": user.RealName,
		"phone":    user.Phone,
		"email":    user.Email,
		"role":     user.Role,
		"avatar":   user.Avatar,
		"status":   user.Status,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, 404, "User not found")
		return
	}

	var req struct {
		RealName string `json:"realName"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	updates := map[string]interface{}{}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	database.DB.Model(&user).Updates(updates)
	utils.Success(c, nil)
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	role := c.Query("role")

	query := database.DB.Model(&model.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	query.Count(&total)

	var users []model.User
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users)

	utils.Success(c, gin.H{
		"list":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AuthHandler) UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&model.User{}).Where("id = ?", id).Update("status", req.Status)
	utils.Success(c, nil)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&model.User{}, id)
	utils.Success(c, nil)
}
