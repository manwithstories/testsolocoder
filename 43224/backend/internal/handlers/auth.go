package handlers

import (
	"net/http"
	"translation-platform/internal/database"
	"translation-platform/internal/middleware"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Password string   `json:"password" binding:"required,min=6,max=128"`
	Email    string   `json:"email" binding:"required,email"`
	Phone    string   `json:"phone"`
	RealName string   `json:"real_name"`
	Role     string   `json:"role" binding:"required,oneof=client translator pm"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	var existing models.User
	if database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error == nil {
		utils.BadRequest(c, "用户名或邮箱已存在")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码加密失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		RealName: req.RealName,
		Role:     models.UserRole(req.Role),
		Status:   "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, "创建用户失败")
		return
	}

	utils.Success(c, gin.H{"id": user.ID, "username": user.Username})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	if user.Status != "active" {
		utils.Forbidden(c, "账户已被禁用")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		utils.InternalError(c, "生成令牌失败")
		return
	}

	utils.Success(c, LoginResponse{
		Token: token,
		User: gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"role":      user.Role,
			"real_name": user.RealName,
			"avatar":    user.Avatar,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.Preload("LanguagePairs").Preload("ExpertiseTags").
		First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var req struct {
		RealName        string   `json:"real_name"`
		Phone           string   `json:"phone"`
		Email           string   `json:"email" binding:"omitempty,email"`
		Avatar          string   `json:"avatar"`
		LanguagePairIDs []uint   `json:"language_pair_ids"`
		ExpertiseTagIDs []uint   `json:"expertise_tag_ids"`
		DailyCapacity   int      `json:"daily_capacity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	updates := map[string]interface{}{
		"real_name":      req.RealName,
		"phone":          req.Phone,
		"email":          req.Email,
		"avatar":         req.Avatar,
		"daily_capacity": req.DailyCapacity,
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalError(c, "更新用户资料失败")
		return
	}

	if len(req.LanguagePairIDs) > 0 {
		var pairs []models.LanguagePair
		database.DB.Find(&pairs, req.LanguagePairIDs)
		database.DB.Model(&user).Association("LanguagePairs").Replace(pairs)
	}

	if len(req.ExpertiseTagIDs) > 0 {
		var tags []models.ExpertiseTag
		database.DB.Find(&tags, req.ExpertiseTagIDs)
		database.DB.Model(&user).Association("ExpertiseTags").Replace(tags)
	}

	utils.Success(c, user)
}

func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.BadRequest(c, "原密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码加密失败")
		return
	}

	database.DB.Model(&user).Update("password", string(hashedPassword))
	utils.Success(c, nil)
}

func ListUsers(c *gin.Context) {
	var users []models.User
	query := database.DB.Preload("LanguagePairs").Preload("ExpertiseTags")

	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Model(&models.User{}).Count(&total)

	page, pageSize := parsePagination(c)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	utils.Success(c, utils.PageResult{
		List:     users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active disabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	database.DB.Model(&user).Update("status", req.Status)
	utils.Success(c, nil)
}

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if n, err := parseUint(p); err == nil {
			page = int(n)
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if n, err := parseUint(ps); err == nil {
			pageSize = int(n)
		}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func parseUint(s string) (uint, error) {
	var n uint
	_, err := http.ParseUint(s)
	return n, err
}
