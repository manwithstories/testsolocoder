package handlers

import (
	"time"
	"translation-platform/internal/database"
	"translation-platform/internal/models"
	"translation-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListPayments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var payments []models.Payment
	query := database.DB.Preload("Project").Preload("Client").Preload("Translator")

	switch role.(string) {
	case "client":
		query = query.Where("client_id = ?", userID)
	case "translator":
		query = query.Where("translator_id = ?", userID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if projectID := c.Query("project_id"); projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	var total int64
	query.Model(&models.Payment{}).Count(&total)

	page, pageSize := parsePagination(c)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&payments)

	utils.Success(c, utils.PageResult{
		List:     payments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func GetPayment(c *gin.Context) {
	id := c.Param("id")

	var payment models.Payment
	if err := database.DB.Preload("Project").Preload("Client").Preload("Translator").
		First(&payment, id).Error; err != nil {
		utils.NotFound(c, "支付记录不存在")
		return
	}

	utils.Success(c, payment)
}

func ConfirmPayment(c *gin.Context) {
	id := c.Param("id")

	var payment models.Payment
	if err := database.DB.First(&payment, id).Error; err != nil {
		utils.NotFound(c, "支付记录不存在")
		return
	}

	if payment.Status == "paid" {
		utils.BadRequest(c, "已支付，无需重复确认")
		return
	}

	now := time.Now()
	database.DB.Model(&payment).Updates(map[string]interface{}{
		"status":  "paid",
		"paid_at": now,
	})

	utils.Success(c, nil)
}

func CalculateProjectFee(c *gin.Context) {
	var req struct {
		WordCount    int                 `json:"word_count" binding:"required"`
		SourceLang   string              `json:"source_lang" binding:"required"`
		TargetLang   string              `json:"target_lang" binding:"required"`
		Urgency      models.UrgencyLevel `json:"urgency" binding:"required"`
		ExpertiseIDs []uint              `json:"expertise_tag_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	domainNames := getDomainNames(req.ExpertiseIDs)
	result := utils.CalculateFee(req.WordCount, req.SourceLang, req.TargetLang, req.Urgency, domainNames)

	utils.Success(c, result)
}

func GetPaymentStatistics(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var totalRevenue, pendingAmount, paidAmount float64
	var paymentCount int64

	query := database.DB.Model(&models.Payment{})

	switch role.(string) {
	case "client":
		query = query.Where("client_id = ?", userID)
	case "translator":
		query = query.Where("translator_id = ?", userID)
	}

	query.Count(&paymentCount)

	var payments []models.Payment
	query.Find(&payments)

	for _, p := range payments {
		totalRevenue += p.Amount
		if p.Status == "paid" {
			paidAmount += p.Amount
		} else {
			pendingAmount += p.Amount
		}
	}

	utils.Success(c, gin.H{
		"total_count":    paymentCount,
		"total_revenue":  totalRevenue,
		"paid_amount":    paidAmount,
		"pending_amount": pendingAmount,
	})
}
