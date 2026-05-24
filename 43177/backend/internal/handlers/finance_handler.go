package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"repair-platform/internal/models"
	"repair-platform/internal/utils"
	"repair-platform/pkg/config"
	"repair-platform/pkg/logger"
)

type FinanceHandler struct{}

func NewFinanceHandler() *FinanceHandler {
	return &FinanceHandler{}
}

type CreateWithdrawRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	BankAccount string  `json:"bank_account" binding:"required"`
	BankName    string  `json:"bank_name" binding:"required"`
}

func (h *FinanceHandler) CreateWithdrawRequest(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "用户不存在")
		return
	}

	if user.Balance < req.Amount {
		utils.Error(c, http.StatusBadRequest, 400, "余额不足")
		return
	}

	tx := models.DB.Begin()

	withdrawRequest := models.WithdrawRequest{
		RequestNo:   utils.GenerateOrderNo("WD"),
		TechnicianID: userID,
		Amount:      req.Amount,
		BankAccount: req.BankAccount,
		BankName:    req.BankName,
		Status:      "pending",
	}

	if err := tx.Create(&withdrawRequest).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Create withdraw request error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "申请失败")
		return
	}

	tx.Model(&user).Update("balance", gorm.Expr("balance - ?", req.Amount))

	transaction := models.Transaction{
		TransactionNo: utils.GenerateOrderNo("TXN"),
		UserID:        userID,
		Type:          "withdraw",
		Amount:        req.Amount,
		BalanceAfter:  user.Balance - req.Amount,
		Description:   "提现申请: " + withdrawRequest.RequestNo,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "申请失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "申请失败")
		return
	}

	utils.Success(c, gin.H{
		"request_id": withdrawRequest.ID,
		"request_no": withdrawRequest.RequestNo,
		"message":    "提现申请已提交",
	})
}

func (h *FinanceHandler) GetWithdrawRequestList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var requests []models.WithdrawRequest
	query := models.DB.Preload("Technician")

	if role == "technician" {
		query = query.Where("technician_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.WithdrawRequest{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests)

	utils.Success(c, gin.H{
		"list":     requests,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *FinanceHandler) GetWithdrawRequestDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.WithdrawRequest
	if err := models.DB.Preload("Technician").First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	utils.Success(c, request)
}

func (h *FinanceHandler) ApproveWithdraw(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.WithdrawRequest
	if err := models.DB.First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	if request.Status != "pending" {
		utils.Error(c, http.StatusBadRequest, 400, "申请状态不允许审批")
		return
	}

	now := time.Now()
	tx := models.DB.Begin()

	if err := tx.Model(&request).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_at": now,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "审批失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "审批失败")
		return
	}

	utils.Success(c, gin.H{"message": "审批通过"})
}

func (h *FinanceHandler) RejectWithdraw(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	type RejectRequest struct {
		Remark string `json:"remark" binding:"required"`
	}

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var request models.WithdrawRequest
	if err := models.DB.First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	if request.Status != "pending" {
		utils.Error(c, http.StatusBadRequest, 400, "申请状态不允许审批")
		return
	}

	tx := models.DB.Begin()

	if err := tx.Model(&request).Updates(map[string]interface{}{
		"status": "rejected",
		"remark": req.Remark,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	var user models.User
	models.DB.First(&user, request.TechnicianID)
	tx.Model(&user).Update("balance", gorm.Expr("balance + ?", request.Amount))

	transaction := models.Transaction{
		TransactionNo: utils.GenerateOrderNo("TXN"),
		UserID:        request.TechnicianID,
		Type:          "refund",
		Amount:        request.Amount,
		BalanceAfter:  user.Balance + request.Amount,
		Description:   "提现拒绝退款: " + request.RequestNo,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	utils.Success(c, gin.H{"message": "已拒绝"})
}

func (h *FinanceHandler) CompleteWithdraw(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request models.WithdrawRequest
	if err := models.DB.First(&request, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "申请不存在")
		return
	}

	if request.Status != "approved" {
		utils.Error(c, http.StatusBadRequest, 400, "申请状态不允许打款")
		return
	}

	now := time.Now()
	models.DB.Model(&request).Updates(map[string]interface{}{
		"status":        "completed",
		"transferred_at": now,
	})

	utils.Success(c, gin.H{"message": "打款完成"})
}

func (h *FinanceHandler) GetTransactionList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	transType := c.Query("type")

	var transactions []models.Transaction
	query := models.DB.Preload("User").Preload("Order")

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if transType != "" {
		query = query.Where("type = ?", transType)
	}

	var total int64
	query.Model(&models.Transaction{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&transactions)

	utils.Success(c, gin.H{
		"list":     transactions,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *FinanceHandler) GetMonthlyReport(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	var report models.MonthlyReport
	if err := models.DB.Where("month = ?", month).First(&report).Error; err != nil {
		report = h.generateMonthlyReport(month)
	}

	utils.Success(c, report)
}

func (h *FinanceHandler) generateMonthlyReport(month string) models.MonthlyReport {
	startDate := month + "-01"
	endDate := month + "-31"

	var totalOrders int64
	models.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) BETWEEN ? AND ?",
			models.OrderStatusCompleted, startDate, endDate).
		Count(&totalOrders)

	var totalRevenue float64
	models.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) BETWEEN ? AND ?",
			models.OrderStatusCompleted, startDate, endDate).
		Select("COALESCE(SUM(final_price), 0)").Scan(&totalRevenue)

	cfg := config.LoadConfig()
	platformIncome := totalRevenue * cfg.PlatformFee
	technicianPay := totalRevenue - platformIncome

	var totalWithdraw float64
	models.DB.Model(&models.WithdrawRequest{}).
		Where("status = ? AND DATE(created_at) BETWEEN ? AND ?",
			"completed", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalWithdraw)

	var newTechnicians int64
	models.DB.Model(&models.User{}).
		Where("role = ? AND DATE(created_at) BETWEEN ? AND ?",
			models.RoleTech, startDate, endDate).
		Count(&newTechnicians)

	var newCustomers int64
	models.DB.Model(&models.User{}).
		Where("role = ? AND DATE(created_at) BETWEEN ? AND ?",
			models.RoleCustomer, startDate, endDate).
		Count(&newCustomers)

	report := models.MonthlyReport{
		Month:          month,
		TotalOrders:    int(totalOrders),
		TotalRevenue:   totalRevenue,
		PlatformIncome: platformIncome,
		TechnicianPay:  technicianPay,
		TotalWithdraw:  totalWithdraw,
		NewTechnicians: int(newTechnicians),
		NewCustomers:   int(newCustomers),
	}

	models.DB.Create(&report)
	return report
}

func (h *FinanceHandler) GetTechnicianPerformance(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	startDate := month + "-01"
	endDate := month + "-31"

	type TechPerformance struct {
		TechnicianID   uint    `json:"technician_id"`
		Username       string  `json:"username"`
		RealName       string  `json:"real_name"`
		TotalOrders    int64   `json:"total_orders"`
		TotalRevenue   float64 `json:"total_revenue"`
		AvgRating      float64 `json:"avg_rating"`
		CompletionRate float64 `json:"completion_rate"`
	}

	var performances []TechPerformance

	models.DB.Table("technician_profiles tp").
		Select(`tp.user_id as technician_id, 
			u.username as username, 
			u.real_name as real_name,
			COUNT(o.id) as total_orders,
			COALESCE(SUM(o.final_price), 0) as total_revenue,
			tp.rating as avg_rating,
			CASE WHEN tp.total_orders > 0 
				THEN CAST(tp.completed_orders AS FLOAT) / tp.total_orders 
				ELSE 0 END as completion_rate`).
		Joins("LEFT JOIN users u ON u.id = tp.user_id").
		Joins("LEFT JOIN orders o ON o.technician_id = tp.user_id AND o.status = ? AND DATE(o.created_at) BETWEEN ? AND ?",
			models.OrderStatusCompleted, startDate, endDate).
		Group("tp.user_id").
		Order("total_revenue DESC").
		Scan(&performances)

	utils.Success(c, performances)
}

func (h *FinanceHandler) GetBalance(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "用户不存在")
		return
	}

	var pendingWithdraw float64
	models.DB.Model(&models.WithdrawRequest{}).
		Where("technician_id = ? AND status = ?", userID, "pending").
		Select("COALESCE(SUM(amount), 0)").Scan(&pendingWithdraw)

	utils.Success(c, gin.H{
		"balance":         user.Balance,
		"pending_withdraw": pendingWithdraw,
		"available_balance": user.Balance - pendingWithdraw,
	})
}

func (h *FinanceHandler) SettleTechnicianIncome(c *gin.Context) {
	cfg := config.LoadConfig()

	var completedOrders []models.Order
	models.DB.Where("status = ? AND final_price > 0", models.OrderStatusCompleted).
		Where("id NOT IN (SELECT DISTINCT order_id FROM transactions WHERE type = 'income')").
		Find(&completedOrders)

	if len(completedOrders) == 0 {
		utils.Success(c, gin.H{"message": "没有待结算的工单"})
		return
	}

	tx := models.DB.Begin()

	settledCount := 0
	for _, order := range completedOrders {
		if order.TechnicianID == nil {
			continue
		}

		platformFee := order.FinalPrice * cfg.PlatformFee
		techIncome := order.FinalPrice - platformFee

		var techUser models.User
		models.DB.First(&techUser, *order.TechnicianID)

		tx.Model(&techUser).Update("balance", gorm.Expr("balance + ?", techIncome))

		transaction := models.Transaction{
			TransactionNo: utils.GenerateOrderNo("TXN"),
			UserID:        *order.TechnicianID,
			Type:          "income",
			Amount:        techIncome,
			BalanceAfter:  techUser.Balance + techIncome,
			OrderID:       &order.ID,
			Description:   "工单收入: " + order.OrderNo,
		}
		tx.Create(&transaction)

		settledCount++
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "结算失败")
		return
	}

	utils.Success(c, gin.H{
		"settled_count": settledCount,
		"message":       "结算完成",
	})
}
