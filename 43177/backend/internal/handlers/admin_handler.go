package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"repair-platform/internal/models"
	"repair-platform/internal/utils"
	"repair-platform/pkg/logger"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	var totalUsers int64
	models.DB.Model(&models.User{}).Count(&totalUsers)

	var totalCustomers int64
	models.DB.Model(&models.User{}).Where("role = ?", models.RoleCustomer).Count(&totalCustomers)

	var totalTechnicians int64
	models.DB.Model(&models.User{}).Where("role = ?", models.RoleTech).Count(&totalTechnicians)

	var verifiedTechnicians int64
	models.DB.Model(&models.TechnicianProfile{}).Where("is_verified = ?", true).Count(&verifiedTechnicians)

	var pendingTechnicians int64
	models.DB.Model(&models.TechnicianProfile{}).Where("verify_status = ?", "pending").Count(&pendingTechnicians)

	var totalOrders int64
	models.DB.Model(&models.Order{}).Count(&totalOrders)

	var pendingOrders int64
	models.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPending).Count(&pendingOrders)

	var completedOrders int64
	models.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusCompleted).Count(&completedOrders)

	var totalRevenue float64
	models.DB.Model(&models.Order{}).
		Where("status = ?", models.OrderStatusCompleted).
		Select("COALESCE(SUM(final_price), 0)").Scan(&totalRevenue)

	var pendingWithdraw int64
	models.DB.Model(&models.WithdrawRequest{}).Where("status = ?", "pending").Count(&pendingWithdraw)

	var lowStockParts int64
	models.DB.Model(&models.Part{}).Where("stock <= min_stock AND status = ?", true).Count(&lowStockParts)

	var totalParts int64
	models.DB.Model(&models.Part{}).Where("status = ?", true).Count(&totalParts)

	utils.Success(c, gin.H{
		"total_users":          totalUsers,
		"total_customers":      totalCustomers,
		"total_technicians":    totalTechnicians,
		"verified_technicians": verifiedTechnicians,
		"pending_technicians":  pendingTechnicians,
		"total_orders":         totalOrders,
		"pending_orders":       pendingOrders,
		"completed_orders":     completedOrders,
		"total_revenue":        totalRevenue,
		"pending_withdraw":     pendingWithdraw,
		"low_stock_parts":      lowStockParts,
		"total_parts":          totalParts,
	})
}

func (h *AdminHandler) GetTechnicianVerifyList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var profiles []models.TechnicianProfile
	query := models.DB.Preload("User")

	if status != "" {
		query = query.Where("verify_status = ?", status)
	} else {
		query = query.Where("verify_status = ?", "pending")
	}

	var total int64
	query.Model(&models.TechnicianProfile{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&profiles)

	utils.Success(c, gin.H{
		"list":     profiles,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

type VerifyTechnicianRequest struct {
	IsVerified bool   `json:"is_verified"`
	VerifyRemark string `json:"verify_remark"`
}

func (h *AdminHandler) VerifyTechnician(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req VerifyTechnicianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var profile models.TechnicianProfile
	if err := models.DB.First(&profile, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "技师资料不存在")
		return
	}

	verifyStatus := "rejected"
	if req.IsVerified {
		verifyStatus = "approved"
	}

	tx := models.DB.Begin()

	if err := tx.Model(&profile).Updates(map[string]interface{}{
		"is_verified":  req.IsVerified,
		"verify_status": verifyStatus,
		"verify_remark": req.VerifyRemark,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "审核失败")
		return
	}

	if req.IsVerified {
		var user models.User
		models.DB.First(&user, profile.UserID)
		if user.Status == models.UserStatusPending {
			tx.Model(&user).Update("status", models.UserStatusActive)
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "审核失败")
		return
	}

	logger.Infof("Technician %d verified: %v", id, req.IsVerified)

	utils.Success(c, gin.H{"message": "审核完成"})
}

func (h *AdminHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")
	status := c.Query("status")

	var users []models.User
	query := models.DB

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.User{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)

	type UserListItem struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		RealName  string `json:"real_name"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		Status    string `json:"status"`
		Balance   float64 `json:"balance"`
		CreatedAt string `json:"created_at"`
	}

	var result []UserListItem
	for _, u := range users {
		result = append(result, UserListItem{
			ID:        u.ID,
			Username:  u.Username,
			RealName:  u.RealName,
			Phone:     u.Phone,
			Email:     u.Email,
			Role:      string(u.Role),
			Status:    string(u.Status),
			Balance:   u.Balance,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.Success(c, gin.H{
		"list":     result,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "用户不存在")
		return
	}

	response := gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"real_name":  user.RealName,
		"phone":      user.Phone,
		"email":      user.Email,
		"avatar":     user.Avatar,
		"role":       user.Role,
		"status":     user.Status,
		"balance":    user.Balance,
		"address":    user.Address,
		"created_at": user.CreatedAt,
	}

	if user.Role == models.RoleTech {
		var profile models.TechnicianProfile
		models.DB.Where("user_id = ?", user.ID).First(&profile)
		response["technician_profile"] = profile
	}

	var recentOrders []models.Order
	models.DB.Where("customer_id = ? OR technician_id = ?", user.ID, user.ID).
		Order("created_at DESC").Limit(10).Find(&recentOrders)
	response["recent_orders"] = recentOrders

	utils.Success(c, response)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	type UpdateStatusRequest struct {
		Status string `json:"status" binding:"required,oneof=active disabled"`
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", id).
		Update("status", req.Status).Error; err != nil {
		logger.Errorf("Update user status error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "更新失败")
		return
	}

	logger.Infof("User %d status updated to %s", id, req.Status)

	utils.Success(c, gin.H{"message": "更新成功"})
}

func (h *AdminHandler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var orders []models.Order
	query := models.DB.Preload("Customer").Preload("Technician").Preload("ServiceItem")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	utils.Success(c, gin.H{
		"list":     orders,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *AdminHandler) GetRefundList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var orders []models.Order
	query := models.DB.Preload("Customer").Preload("Technician").
		Where("status = ?", models.OrderStatusRefunding)

	var total int64
	query.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	utils.Success(c, gin.H{
		"list":     orders,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *AdminHandler) HandleLowRatingReview(c *gin.Context) {
	reviewID, _ := strconv.Atoi(c.Param("id"))

	var review models.Review
	if err := models.DB.First(&review, reviewID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "评价不存在")
		return
	}

	if review.Rating > 2 {
		utils.Error(c, http.StatusBadRequest, 400, "只有差评需要平台介入")
		return
	}

	type InterveneRequest struct {
		Note string `json:"note" binding:"required"`
	}

	var req InterveneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	models.DB.Model(&review).Updates(map[string]interface{}{
		"is_intervened":  true,
		"intervene_note": req.Note,
	})

	logger.Infof("Review %d intervened, note: %s", reviewID, req.Note)

	utils.Success(c, gin.H{"message": "已介入处理"})
}

func (h *AdminHandler) GetLowRatingReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var reviews []models.Review
	query := models.DB.Preload("Customer").Preload("Technician").Preload("Order").
		Where("rating <= ?", 2)

	var total int64
	query.Model(&models.Review{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("is_intervened ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews)

	utils.Success(c, gin.H{
		"list":     reviews,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *AdminHandler) GetOrderLogs(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Query("order_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	var logs []models.OrderLog
	query := models.DB

	if orderID > 0 {
		query = query.Where("order_id = ?", orderID)
	}

	var total int64
	query.Model(&models.OrderLog{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)

	utils.Success(c, gin.H{
		"list":     logs,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}
