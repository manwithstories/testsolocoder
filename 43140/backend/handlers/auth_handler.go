package handlers

import (
	"strconv"

	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

type RegisterRequest struct {
	Email       string       `json:"email" binding:"required,email"`
	Password    string       `json:"password" binding:"required,min=6"`
	Name        string       `json:"name" binding:"required"`
	Phone       string       `json:"phone"`
	Role        models.Role  `json:"role" binding:"required,oneof=company jobseeker"`
	CompanyName string       `json:"company_name"`
	Industry    string       `json:"industry"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var existingUser models.User
	if h.DB.Where("email = ?", req.Email).First(&existingUser).Error == nil {
		utils.BadRequest(c, "Email already registered")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	tx := h.DB.Begin()

	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     req.Role,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to create user")
		return
	}

	if req.Role == models.RoleCompany {
		company := models.Company{
			UserID:      user.ID,
			CompanyName: req.CompanyName,
			Industry:    req.Industry,
		}
		if err := tx.Create(&company).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to create company profile")
			return
		}

		defaultDept := models.Department{
			CompanyID: company.ID,
			Name:      "Default",
		}
		tx.Create(&defaultDept)
	} else if req.Role == models.RoleJobSeeker {
		jobSeeker := models.JobSeeker{
			UserID: user.ID,
		}
		if err := tx.Create(&jobSeeker).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to create job seeker profile")
			return
		}
	}

	tx.Commit()

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := h.DB.Preload("Company").Preload("JobSeeker").
		Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.BadRequest(c, "Invalid email or password")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.BadRequest(c, "Invalid email or password")
		return
	}

	if user.Status != "active" {
		utils.Forbidden(c, "Account is not active")
		return
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := h.DB.Preload("Company").Preload("JobSeeker").
		First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if err := h.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update profile")
		return
	}

	h.DB.Preload("Company").Preload("JobSeeker").First(&user, userID)
	utils.Success(c, user)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	if !utils.CheckPassword(req.CurrentPassword, user.Password) {
		utils.BadRequest(c, "Current password is incorrect")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	h.DB.Model(&user).Update("password", hashedPassword)
	utils.Success(c, nil)
}

func (h *AuthHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.DB.Preload("Company").Preload("JobSeeker").
		First(&user, id).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user)
}

func (h *AuthHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive suspended"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		utils.InternalError(c, "Failed to update user status")
		return
	}

	utils.Success(c, nil)
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	var users []models.User
	query := h.DB.Preload("Company").Preload("JobSeeker")

	role := c.Query("role")
	if role != "" {
		query = query.Where("role = ?", role)
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.User{}).Count(&total)

	page := getPage(c)
	pageSize := getPageSize(c)

	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	utils.SuccessWithPagination(c, users, page, pageSize, total)
}

func getPage(c *gin.Context) int {
	page := 1
	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	return page
}

func getPageSize(c *gin.Context) int {
	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if psNum, err := strconv.Atoi(ps); err == nil && psNum > 0 && psNum <= 100 {
			pageSize = psNum
		}
	}
	return pageSize
}
