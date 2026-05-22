package handlers

import (
	"fmt"
	"strconv"
	"time"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBillList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	billType := c.Query("type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := config.DB.Model(&models.Bill{})

	if role == "service_provider" {
		query = query.Where("provider_id = ?", userID)
	}

	if billType != "" {
		query = query.Where("bill_type = ?", billType)
	}

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var bills []models.Bill
	offset := (page - 1) * pageSize
	query.Preload("Order").Offset(offset).Limit(pageSize).Order("id DESC").Find(&bills)

	var incomeTotal, withdrawTotal, penaltyTotal float64
	for _, bill := range bills {
		switch bill.BillType {
		case models.BillTypeIncome:
			incomeTotal += bill.Amount
		case models.BillTypeWithdraw:
			withdrawTotal += bill.Amount
		case models.BillTypePenalty:
			penaltyTotal += bill.Amount
		}
	}

	utils.Success(c, gin.H{
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
		"income_total":  utils.RoundToTwoDecimals(incomeTotal),
		"withdraw_total": utils.RoundToTwoDecimals(withdrawTotal),
		"penalty_total": utils.RoundToTwoDecimals(penaltyTotal),
		"list":          bills,
	})
}

func GetBalance(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if result := config.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var pendingIncome float64
	config.DB.Model(&models.Bill{}).
		Where("provider_id = ? AND bill_type = ? AND status = ?", userID, models.BillTypeIncome, models.BillStatusPending).
		Select("COALESCE(SUM(amount), 0)").Scan(&pendingIncome)

	utils.Success(c, gin.H{
		"balance":        utils.RoundToTwoDecimals(user.Balance),
		"total_income":   utils.RoundToTwoDecimals(user.TotalIncome),
		"pending_income": utils.RoundToTwoDecimals(pendingIncome),
		"order_count":    user.OrderCount,
	})
}

func CreateWithdrawRequest(c *gin.Context) {
	providerID := c.GetUint("user_id")

	var req struct {
		Amount        float64 `json:"amount" binding:"required"`
		BankName      string  `json:"bank_name" binding:"required"`
		BankAccount   string  `json:"bank_account" binding:"required"`
		AccountHolder string  `json:"account_holder" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	cfg := config.Load()
	if req.Amount < cfg.App.MinWithdrawAmount {
		utils.BadRequest(c, fmt.Sprintf("最低提现金额为%.0f元", cfg.App.MinWithdrawAmount))
		return
	}

	var user models.User
	if result := config.DB.First(&user, providerID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if user.Balance < req.Amount {
		utils.BadRequest(c, "余额不足")
		return
	}

	tx := config.DB.Begin()

	withdraw := models.WithdrawRequest{
		ProviderID:    providerID,
		Amount:        req.Amount,
		BankName:      req.BankName,
		BankAccount:   req.BankAccount,
		AccountHolder: req.AccountHolder,
		Status:        models.WithdrawStatusPending,
	}

	if err := tx.Create(&withdraw).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "申请失败")
		return
	}

	if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", req.Amount)).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "申请失败")
		return
	}

	bill := models.Bill{
		ProviderID:    providerID,
		BillType:      models.BillTypeWithdraw,
		Amount:        req.Amount,
		Balance:       -req.Amount,
		Status:        models.BillStatusPending,
		Description:   "提现申请",
		TransactionNo: utils.GenerateTransactionNo(),
	}
	tx.Create(&bill)

	tx.Commit()

	go utils.LogOperation(providerID, "service_provider", "withdraw", "create", &withdraw.ID, "withdraw_request", "创建提现申请", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, withdraw)
}

func GetWithdrawList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := config.DB.Model(&models.WithdrawRequest{})

	if role == "service_provider" {
		query = query.Where("provider_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var withdraws []models.WithdrawRequest
	offset := (page - 1) * pageSize
	query.Preload("Provider").Offset(offset).Limit(pageSize).Order("id DESC").Find(&withdraws)

	utils.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      withdraws,
	})
}

func HandleWithdrawRequest(c *gin.Context) {
	handlerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Approved    bool   `json:"approved"`
		Remark      string `json:"remark"`
		TransferNo  string `json:"transfer_no"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var withdraw models.WithdrawRequest
	if result := config.DB.First(&withdraw, id); result.Error != nil {
		utils.NotFound(c, "提现申请不存在")
		return
	}

	if withdraw.Status != models.WithdrawStatusPending {
		utils.BadRequest(c, "该申请已处理")
		return
	}

	tx := config.DB.Begin()

	now := time.Now()
	if req.Approved {
		withdraw.Status = models.WithdrawStatusCompleted
		withdraw.TransferNo = req.TransferNo
	} else {
		withdraw.Status = models.WithdrawStatusRejected
		tx.Model(&models.User{}).Where("id = ?", withdraw.ProviderID).
			Update("balance", gorm.Expr("balance + ?", withdraw.Amount))
	}
	withdraw.HandlerID = &handlerID
	withdraw.HandleRemark = req.Remark
	withdraw.HandledAt = &now

	if err := tx.Save(&withdraw).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "处理失败")
		return
	}

	var bill models.Bill
	tx.Where("provider_id = ? AND bill_type = ? AND status = ?", withdraw.ProviderID, models.BillTypeWithdraw, models.BillStatusPending).
		Order("id DESC").First(&bill)

	if bill.ID != 0 {
		if req.Approved {
			bill.Status = models.BillStatusCompleted
			bill.SettledAt = &now
		} else {
			bill.Status = models.BillStatusFailed
		}
		tx.Save(&bill)
	}

	tx.Commit()

	if req.Approved {
		utils.SendSystemNotification(withdraw.ProviderID, "提现成功", "您的提现申请已通过，金额："+fmt.Sprintf("%.2f", withdraw.Amount)+"元")
	} else {
		utils.SendSystemNotification(withdraw.ProviderID, "提现申请被拒绝", "您的提现申请被拒绝，原因："+req.Remark)
	}

	go utils.LogOperation(handlerID, "admin", "withdraw", "handle", &withdraw.ID, "withdraw_request", "处理提现申请", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func GetIncomeSummary(c *gin.Context) {
	providerID := c.GetUint("user_id")
	month := c.Query("month")

	var startDate, endDate time.Time
	if month != "" {
		startDate, _ = time.Parse("2006-01", month)
		endDate = startDate.AddDate(0, 1, 0)
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
	}

	var totalIncome float64
	config.DB.Model(&models.Bill{}).
		Where("provider_id = ? AND bill_type = ? AND status = ? AND created_at >= ? AND created_at < ?",
			providerID, models.BillTypeIncome, models.BillStatusCompleted, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)

	var orderCount int64
	config.DB.Model(&models.Order{}).
		Where("provider_id = ? AND status = ? AND completed_at >= ? AND completed_at < ?",
			providerID, models.OrderStatusCompleted, startDate, endDate).
		Count(&orderCount)

	var totalPlatformFee float64
	config.DB.Model(&models.Order{}).
		Where("provider_id = ? AND status = ? AND completed_at >= ? AND completed_at < ?",
			providerID, models.OrderStatusCompleted, startDate, endDate).
		Select("COALESCE(SUM(platform_fee), 0)").Scan(&totalPlatformFee)

	utils.Success(c, gin.H{
		"month":            startDate.Format("2006-01"),
		"total_income":     utils.RoundToTwoDecimals(totalIncome),
		"order_count":      orderCount,
		"total_platform_fee": utils.RoundToTwoDecimals(totalPlatformFee),
	})
}
