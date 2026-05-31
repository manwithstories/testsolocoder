package handlers

import (
	"net/http"
	"strconv"
	"time"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/models"
	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if !utils.IsValidUsername(req.Username) {
		utils.Error(c, http.StatusBadRequest, "用户名只能包含字母、数字和下划线，长度3-50")
		return
	}

	if !utils.IsValidEmail(req.Email) {
		utils.Error(c, http.StatusBadRequest, "邮箱格式不正确")
		return
	}

	var existingUser models.User
	if database.DB.Where("username = ?", req.Username).First(&existingUser).Error == nil {
		utils.Error(c, http.StatusConflict, "用户名已存在")
		return
	}

	if database.DB.Where("email = ?", req.Email).First(&existingUser).Error == nil {
		utils.Error(c, http.StatusConflict, "邮箱已被注册")
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Role:     models.RoleUser,
		Status:   models.UserStatusActive,
	}

	if err := user.SetPassword(req.Password); err != nil {
		utils.Error(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "注册失败")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), &h.cfg.JWT)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var user models.User
	query := database.DB.Where("username = ? OR email = ?", req.Account, req.Account)
	if err := query.First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if user.Status != models.UserStatusActive {
		utils.Error(c, http.StatusForbidden, "账户已被禁用")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), &h.cfg.JWT)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.Preload("Certification").First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	var user models.User
	database.DB.First(&user, userID)
	utils.Success(c, user)
}

type UserHandler struct {
	cfg *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{cfg: cfg}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))
	role := c.Query("role")
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var users []models.User
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	utils.PaginatedResponse(c, users, total, page, pageSize)
}

func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var user models.User
	if err := database.DB.Preload("Certification").First(&user, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Status models.UserStatus `json:"status" binding:"required,oneof=active disabled pending"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).
		Update("status", req.Status).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "状态更新成功", nil)
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Role models.UserRole `json:"role" binding:"required,oneof=admin roaster user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).
		Update("role", req.Role).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "角色更新成功", nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currentUserID := c.GetUint("user_id")

	if uint(id) == currentUserID {
		utils.Error(c, http.StatusBadRequest, "不能删除自己")
		return
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *UserHandler) GetUserActivity(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 || days > 365 {
		days = 30
	}

	type ActivityItem struct {
		Date       string `json:"date"`
		UserCount  int64  `json:"user_count"`
		LoginCount int64  `json:"login_count"`
	}

	var results []ActivityItem
	since := time.Now().AddDate(0, 0, -days)

	database.DB.Table("operation_logs").
		Select("DATE(created_at) as date, COUNT(DISTINCT user_id) as user_count, COUNT(*) as login_count").
		Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results)

	utils.Success(c, results)
}
