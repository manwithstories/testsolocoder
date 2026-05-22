package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type UserHandler struct {
	cfg *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{cfg: cfg}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname"`
	Role     string `json:"role" binding:"required,oneof=student instructor"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Bio      string `json:"bio"`
}

type InstructorApplicationRequest struct {
	RealName      string `json:"real_name" binding:"required"`
	IDCard        string `json:"id_card" binding:"required"`
	Qualification string `json:"qualification" binding:"required"`
	Experience    string `json:"experience"`
	Certificates  string `json:"certificates"`
}

type ReviewInstructorRequest struct {
	Status       string `json:"status" binding:"required,oneof=approved rejected"`
	ReviewRemark string `json:"review_remark"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var existing models.User
	if database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error == nil {
		utils.BadRequest(c, "Username or email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Role:     models.UserRole(req.Role),
		Status:   models.UserStatusActive,
	}

	if req.Role == "instructor" {
		user.InstructorStatus = models.InstructorPending
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, "Failed to create user")
		return
	}

	utils.Created(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Account, req.Account).First(&user).Error; err != nil {
		utils.Unauthorized(c, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "Invalid credentials")
		return
	}

	if user.Status != models.UserStatusActive {
		utils.Forbidden(c, "Account is disabled")
		return
	}

	now := time.Now()
	database.DB.Model(&user).Update("last_login_at", &now)

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), h.cfg.JWT.ExpireHours)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), 168)
	if err != nil {
		utils.InternalError(c, "Failed to generate refresh token")
		return
	}

	utils.Success(c, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":                user.ID,
			"username":          user.Username,
			"email":             user.Email,
			"nickname":          user.Nickname,
			"avatar":            user.Avatar,
			"role":              user.Role,
			"instructor_status": user.InstructorStatus,
		},
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	token, err := utils.GenerateToken(userID.(uuid.UUID), username.(string), role.(string), h.cfg.JWT.ExpireHours)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{"token": token})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, gin.H{
		"id":                user.ID,
		"username":          user.Username,
		"email":             user.Email,
		"nickname":          user.Nickname,
		"avatar":            user.Avatar,
		"role":              user.Role,
		"phone":             user.Phone,
		"bio":               user.Bio,
		"instructor_status": user.InstructorStatus,
		"email_verified":    user.EmailVerified,
		"created_at":        user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update profile")
		return
	}

	utils.Success(c, gin.H{"message": "Profile updated successfully"})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.BadRequest(c, "Old password is incorrect")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	database.DB.Model(&user).Update("password", string(hashedPassword))
	utils.Success(c, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) ApplyInstructor(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "student" {
		utils.BadRequest(c, "Only students can apply for instructor")
		return
	}

	var req InstructorApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var existing models.InstructorApplication
	if database.DB.Where("user_id = ?", userID).First(&existing).Error == nil {
		utils.BadRequest(c, "Application already submitted")
		return
	}

	application := models.InstructorApplication{
		UserID:        userID.(uuid.UUID),
		RealName:      req.RealName,
		IDCard:        req.IDCard,
		Qualification: req.Qualification,
		Experience:    req.Experience,
		Certificates:  req.Certificates,
		Status:        models.InstructorPending,
	}

	if err := database.DB.Create(&application).Error; err != nil {
		utils.InternalError(c, "Failed to submit application")
		return
	}

	utils.Created(c, gin.H{
		"id":     application.ID,
		"status": application.Status,
	})
}

func (h *UserHandler) ListInstructorApplications(c *gin.Context) {
	status := c.Query("status")

	query := database.DB.Model(&models.InstructorApplication{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var applications []models.InstructorApplication
	var total int64

	query.Count(&total)

	page, pageSize := getPagination(c)
	query.Preload("User").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&applications)

	utils.Paginated(c, applications, total, page, pageSize)
}

func (h *UserHandler) ReviewInstructor(c *gin.Context) {
	id := c.Param("id")
	appID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	var req ReviewInstructorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var application models.InstructorApplication
	if err := database.DB.First(&application, appID).Error; err != nil {
		utils.NotFound(c, "Application not found")
		return
	}

	if application.Status != models.InstructorPending {
		utils.BadRequest(c, "Application already reviewed")
		return
	}

	now := time.Now()
	reviewerIDRaw, _ := c.Get("user_id")

	application.Status = models.InstructorStatus(req.Status)
	application.ReviewRemark = req.ReviewRemark
	reviewerID := reviewerIDRaw.(uuid.UUID)
	application.ReviewerID = &reviewerID
	application.ReviewedAt = &now

	tx := database.DB.Begin()
	if err := tx.Save(&application).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to review application")
		return
	}

	if req.Status == "approved" {
		if err := tx.Model(&models.User{}).
			Where("id = ?", application.UserID).
			Updates(map[string]interface{}{
				"role":              "instructor",
				"instructor_status": "approved",
			}).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to update user role")
			return
		}
	}

	tx.Commit()
	utils.Success(c, gin.H{"message": "Application reviewed successfully"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	status := c.Query("status")
	search := c.Query("search")

	query := database.DB.Model(&models.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("username ILIKE ? OR email ILIKE ? OR nickname ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var users []models.User
	var total int64

	query.Count(&total)

	page, pageSize := getPagination(c)
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	query.Order(sort + " " + order).Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	utils.Paginated(c, users, total, page, pageSize)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active disabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("status", req.Status).Error; err != nil {
		utils.InternalError(c, "Failed to update user status")
		return
	}

	utils.Success(c, gin.H{"message": "User status updated"})
}

func (h *UserHandler) GetInstructorStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var courseCount int64
	var totalStudents int64
	var totalRevenue float64

	database.DB.Model(&models.Course{}).Where("instructor_id = ?", userID).Count(&courseCount)
	database.DB.Model(&models.Course{}).Where("instructor_id = ?", userID).Select("COALESCE(SUM(student_count), 0)").Scan(&totalStudents)
	database.DB.Model(&models.Order{}).
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.status = ?", userID, models.OrderPaid).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	utils.Success(c, gin.H{
		"course_count":    courseCount,
		"total_students":  totalStudents,
		"total_revenue":   totalRevenue,
	})
}

func (h *UserHandler) GetAdminStats(c *gin.Context) {
	var userCount int64
	var courseCount int64
	var orderCount int64
	var totalRevenue float64

	database.DB.Model(&models.User{}).Where("role = ?", "student").Count(&userCount)
	database.DB.Model(&models.Course{}).Where("status = ?", models.CoursePublished).Count(&courseCount)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderPaid).Count(&orderCount)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderPaid).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	utils.Success(c, gin.H{
		"user_count":     userCount,
		"course_count":   courseCount,
		"order_count":    orderCount,
		"total_revenue":  totalRevenue,
	})
}

func getPagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10
	if v := c.Query("page"); v != "" {
		if p := parseInt(v, 1); p > 0 {
			page = p
		}
	}
	if v := c.Query("page_size"); v != "" {
		if ps := parseInt(v, 10); ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}
	return page, pageSize
}

func parseInt(s string, defaultVal int) int {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return defaultVal
	}
	return n
}
