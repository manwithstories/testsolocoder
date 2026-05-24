package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"repair-platform/internal/models"
	"repair-platform/internal/utils"
	"repair-platform/pkg/logger"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
	RealName string `json:"real_name"`
	Role     string `json:"role" binding:"required,oneof=customer technician"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误: "+err.Error())
		return
	}

	var existingUser models.User
	if err := models.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.Error(c, http.StatusConflict, 409, "用户名已存在")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Password hash error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "注册失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		Role:     models.UserRole(req.Role),
		Status:   models.UserStatusActive,
	}

	if req.Role == "technician" {
		user.Status = models.UserStatusPending
	}

	if err := models.DB.Create(&user).Error; err != nil {
		logger.Errorf("Create user error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "注册失败")
		return
	}

	if req.Role == "technician" {
		techProfile := models.TechnicianProfile{
			UserID:       user.ID,
			VerifyStatus: "pending",
			MaxActiveOrders: 5,
			ServiceRadius: 50.0,
		}
		if err := models.DB.Create(&techProfile).Error; err != nil {
			logger.Errorf("Create technician profile error: %v", err)
		}
	}

	logger.Infof("User registered: %s, role: %s", user.Username, user.Role)
	utils.Success(c, gin.H{
		"user_id": user.ID,
		"message": "注册成功",
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string     `json:"token"`
	UserID   uint       `json:"user_id"`
	Username string     `json:"username"`
	Role     string     `json:"role"`
	Status   string     `json:"status"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, 401, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, 401, "用户名或密码错误")
		return
	}

	if user.Status == models.UserStatusDisabled {
		utils.Error(c, http.StatusForbidden, 403, "账号已被禁用")
		return
	}

	if user.Role == models.RoleTech && user.Status == models.UserStatusPending {
		utils.Error(c, http.StatusForbidden, 403, "技师资质审核中，请耐心等待")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		logger.Errorf("Generate token error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "登录失败")
		return
	}

	utils.Success(c, LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		Status:   string(user.Status),
	})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "用户不存在")
		return
	}

	response := gin.H{
		"id":       user.ID,
		"username": user.Username,
		"phone":    user.Phone,
		"email":    user.Email,
		"real_name": user.RealName,
		"avatar":   user.Avatar,
		"role":     user.Role,
		"status":   user.Status,
		"balance":  user.Balance,
		"address":  user.Address,
	}

	if user.Role == models.RoleTech {
		var profile models.TechnicianProfile
		if err := models.DB.Where("user_id = ?", user.ID).First(&profile).Error; err == nil {
			response["technician_profile"] = profile
		}
	}

	utils.Success(c, response)
}

type UpdateProfileRequest struct {
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	RealName string  `json:"real_name"`
	Avatar   string  `json:"avatar"`
	Address  string  `json:"address"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Longitude != 0 {
		updates["longitude"] = req.Longitude
	}
	if req.Latitude != 0 {
		updates["latitude"] = req.Latitude
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		logger.Errorf("Update profile error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "更新失败")
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

type SubmitCertificateRequest struct {
	CertificateImage string `json:"certificate_image" binding:"required"`
	CertificateNo    string `json:"certificate_no" binding:"required"`
	Specialty        string `json:"specialty"`
	ExperienceYears  int    `json:"experience_years"`
}

func (h *UserHandler) SubmitCertificate(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "technician" {
		utils.Error(c, http.StatusForbidden, 403, "只有技师可以提交资质证书")
		return
	}

	var req SubmitCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var profile models.TechnicianProfile
	if err := models.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "技师资料不存在")
		return
	}

	updates := map[string]interface{}{
		"certificate_image": req.CertificateImage,
		"certificate_no":    req.CertificateNo,
		"specialty":         req.Specialty,
		"experience_years":  req.ExperienceYears,
		"verify_status":     "pending",
	}

	if err := models.DB.Model(&profile).Updates(updates).Error; err != nil {
		logger.Errorf("Update certificate error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "提交失败")
		return
	}

	utils.Success(c, gin.H{"message": "资质证书已提交，等待审核"})
}

func (h *UserHandler) GetTechnicianList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	specialty := c.Query("specialty")

	var profiles []models.TechnicianProfile
	query := models.DB.Preload("User").Where("is_verified = ?", true)

	if specialty != "" {
		query = query.Where("specialty LIKE ?", "%"+specialty+"%")
	}

	var total int64
	query.Model(&models.TechnicianProfile{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("rating DESC").Offset(offset).Limit(pageSize).Find(&profiles)

	type TechListItem struct {
		ID              uint    `json:"id"`
		UserID          uint    `json:"user_id"`
		Username        string  `json:"username"`
		RealName        string  `json:"real_name"`
		Avatar          string  `json:"avatar"`
		Rating          float64 `json:"rating"`
		Specialty       string  `json:"specialty"`
		ExperienceYears int     `json:"experience_years"`
		CompletedOrders int     `json:"completed_orders"`
	}

	var result []TechListItem
	for _, p := range profiles {
		result = append(result, TechListItem{
			ID:              p.ID,
			UserID:          p.UserID,
			Username:        p.User.Username,
			RealName:        p.User.RealName,
			Avatar:          p.User.Avatar,
			Rating:          p.Rating,
			Specialty:       p.Specialty,
			ExperienceYears: p.ExperienceYears,
			CompletedOrders: p.CompletedOrders,
		})
	}

	utils.Success(c, gin.H{
		"list":     result,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *UserHandler) GetTechnicianDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var profile models.TechnicianProfile
	if err := models.DB.Preload("User").First(&profile, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "技师不存在")
		return
	}

	var orders []models.Order
	models.DB.Where("technician_id = ? AND status = ?", profile.UserID, models.OrderStatusCompleted).
		Order("created_at DESC").Limit(10).Find(&orders)

	var reviews []models.Review
	models.DB.Preload("Customer").Where("technician_id = ?", profile.UserID).
		Order("created_at DESC").Limit(10).Find(&reviews)

	utils.Success(c, gin.H{
		"profile": profile,
		"recent_orders": orders,
		"recent_reviews": reviews,
	})
}

type ReviewTechRequest struct {
	OrderID    uint   `json:"order_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Content    string `json:"content"`
	Images     string `json:"images"`
}

func (h *UserHandler) ReviewTechnician(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req ReviewTechRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, req.OrderID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.CustomerID != userID {
		utils.Error(c, http.StatusForbidden, 403, "只能评价自己的工单")
		return
	}

	if order.Status != models.OrderStatusCompleted {
		utils.Error(c, http.StatusBadRequest, 400, "只能评价已完成的工单")
		return
	}

	var existingReview models.Review
	if err := models.DB.Where("order_id = ?", req.OrderID).First(&existingReview).Error; err == nil {
		utils.Error(c, http.StatusBadRequest, 400, "已评价过该工单")
		return
	}

	review := models.Review{
		OrderID:      req.OrderID,
		CustomerID:   userID,
		TechnicianID: *order.TechnicianID,
		Rating:       req.Rating,
		Content:      req.Content,
		Images:       req.Images,
	}

	if err := models.DB.Create(&review).Error; err != nil {
		logger.Errorf("Create review error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "评价失败")
		return
	}

	var profile models.TechnicianProfile
	models.DB.Where("user_id = ?", *order.TechnicianID).First(&profile)

	var reviews []models.Review
	models.DB.Where("technician_id = ?", *order.TechnicianID).Find(&reviews)

	totalRating := 0
	for _, r := range reviews {
		totalRating += r.Rating
	}
	newRating := float64(totalRating) / float64(len(reviews))
	newRating = math.Round(newRating*10) / 10

	models.DB.Model(&profile).Update("rating", newRating)

	logger.Infof("Technician %d rating updated to %.1f", *order.TechnicianID, newRating)

	utils.Success(c, gin.H{
		"review_id": review.ID,
		"new_rating": newRating,
		"message":   "评价成功",
	})
}

func (h *UserHandler) GetReviewList(c *gin.Context) {
	technicianID, _ := strconv.Atoi(c.Query("technician_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := models.DB.Preload("Customer").Preload("Order")
	if technicianID > 0 {
		query = query.Where("technician_id = ?", technicianID)
	}

	var total int64
	query.Model(&models.Review{}).Count(&total)

	var reviews []models.Review
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews)

	utils.Success(c, gin.H{
		"list":     reviews,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

type ReplyReviewRequest struct {
	Reply string `json:"reply" binding:"required"`
}

func (h *UserHandler) ReplyReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	reviewID, _ := strconv.Atoi(c.Param("id"))

	var req ReplyReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var review models.Review
	if err := models.DB.First(&review, reviewID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "评价不存在")
		return
	}

	if review.TechnicianID != userID {
		utils.Error(c, http.StatusForbidden, 403, "只能回复自己收到的评价")
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"reply":      req.Reply,
		"replied_at": now,
	}

	if err := models.DB.Model(&review).Updates(updates).Error; err != nil {
		logger.Errorf("Reply review error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "回复失败")
		return
	}

	utils.Success(c, gin.H{"message": "回复成功"})
}
