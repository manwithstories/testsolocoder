package handlers

import (
	"time"
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BudgetHandler struct{}

func NewBudgetHandler() *BudgetHandler {
	return &BudgetHandler{}
}

type BudgetItemRequest struct {
	Category      string  `json:"category" binding:"required,max=100"`
	Description   string  `json:"description"`
	EstimatedCost float64 `json:"estimated_cost"`
	ActualCost    float64 `json:"actual_cost"`
	PaidAmount    float64 `json:"paid_amount"`
	Status        string  `json:"status"`
	DueDate       string  `json:"due_date"`
	Notes         string  `json:"notes"`
	VendorID      *uint   `json:"vendor_id"`
}

type PaymentRequest struct {
	BudgetItemID uint    `json:"budget_item_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Method       string  `json:"method"`
	Status       string  `json:"status"`
	PaidAt       string  `json:"paid_at"`
	Reference    string  `json:"reference"`
	Notes        string  `json:"notes"`
}

func (h *BudgetHandler) CreateBudgetItem(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	var req BudgetItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	item := models.BudgetItem{
		WeddingID:     weddingID,
		Category:      req.Category,
		Description:   req.Description,
		EstimatedCost: req.EstimatedCost,
		ActualCost:    req.ActualCost,
		PaidAmount:    req.PaidAmount,
		Status:        req.Status,
		Notes:         req.Notes,
		VendorID:      req.VendorID,
	}

	if item.Status == "" {
		item.Status = "pending"
	}

	if req.DueDate != "" {
		dueDate, err := parseDate(req.DueDate)
		if err == nil {
			item.DueDate = &dueDate
		}
	}

	if err := db.Create(&item).Error; err != nil {
		response.InternalError(c, "Failed to create budget item")
		return
	}

	response.Created(c, item)
}

func (h *BudgetHandler) GetBudgetItems(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var items []models.BudgetItem
	db.Where("wedding_id = ?", weddingID).Order("created_at DESC").Find(&items)

	response.Success(c, items)
}

func (h *BudgetHandler) GetBudgetSummary(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var wedding models.Wedding
	db.First(&wedding, weddingID)

	var result struct {
		TotalBudget     float64 `json:"total_budget"`
		TotalEstimated  float64 `json:"total_estimated"`
		TotalActual     float64 `json:"total_actual"`
		TotalPaid       float64 `json:"total_paid"`
		RemainingBudget float64 `json:"remaining_budget"`
		OverBudget      bool    `json:"over_budget"`
		ItemsCount      int64   `json:"items_count"`
	}

	result.TotalBudget = wedding.Budget

	db.Model(&models.BudgetItem{}).Where("wedding_id = ?", weddingID).
		Select("COALESCE(SUM(estimated_cost), 0)").Scan(&result.TotalEstimated)
	db.Model(&models.BudgetItem{}).Where("wedding_id = ?", weddingID).
		Select("COALESCE(SUM(actual_cost), 0)").Scan(&result.TotalActual)
	db.Model(&models.BudgetItem{}).Where("wedding_id = ?", weddingID).
		Select("COALESCE(SUM(paid_amount), 0)").Scan(&result.TotalPaid)

	result.RemainingBudget = wedding.Budget - result.TotalActual
	result.OverBudget = result.TotalActual > wedding.Budget
	db.Model(&models.BudgetItem{}).Where("wedding_id = ?", weddingID).Count(&result.ItemsCount)

	response.Success(c, result)
}

func (h *BudgetHandler) UpdateBudgetItem(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	var req BudgetItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var item models.BudgetItem
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&item).Error; err != nil {
		response.NotFound(c, "Budget item not found")
		return
	}

	updates := map[string]interface{}{
		"category":       req.Category,
		"description":    req.Description,
		"estimated_cost": req.EstimatedCost,
		"actual_cost":    req.ActualCost,
		"paid_amount":    req.PaidAmount,
		"status":         req.Status,
		"notes":          req.Notes,
		"vendor_id":      req.VendorID,
	}

	if req.DueDate != "" {
		dueDate, err := parseDate(req.DueDate)
		if err == nil {
			updates["due_date"] = dueDate
		}
	}

	db.Model(&item).Updates(updates)

	response.Success(c, item)
}

func (h *BudgetHandler) DeleteBudgetItem(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var item models.BudgetItem
	if err := db.Where("id = ? AND wedding_id = ?", id, weddingID).First(&item).Error; err != nil {
		response.NotFound(c, "Budget item not found")
		return
	}

	db.Delete(&item)

	response.Success(c, gin.H{"message": "Budget item deleted successfully"})
}

func (h *BudgetHandler) RecordPayment(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var budgetItem models.BudgetItem
	if err := db.Where("id = ? AND wedding_id = ?", req.BudgetItemID, weddingID).First(&budgetItem).Error; err != nil {
		response.NotFound(c, "Budget item not found")
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		payment := models.Payment{
			BudgetItemID: req.BudgetItemID,
			Amount:       req.Amount,
			Method:       req.Method,
			Status:       req.Status,
			Reference:    req.Reference,
			Notes:        req.Notes,
		}

		if payment.Status == "" {
			payment.Status = "pending"
		}

		if req.PaidAt != "" {
			paidAt, err := parseDate(req.PaidAt)
			if err == nil {
				payment.PaidAt = &paidAt
				if payment.Status == "pending" {
					payment.Status = "completed"
				}
			}
		}

		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		if payment.Status == "completed" {
			newPaidAmount := budgetItem.PaidAmount + req.Amount
			if err := tx.Model(&budgetItem).Update("paid_amount", newPaidAmount).Error; err != nil {
				return err
			}

			if newPaidAmount >= budgetItem.ActualCost {
				tx.Model(&budgetItem).Update("status", "paid")
			} else if newPaidAmount > 0 {
				tx.Model(&budgetItem).Update("status", "partial")
			}
		}

		return nil
	})

	if err != nil {
		response.InternalError(c, "Failed to record payment")
		return
	}

	response.Created(c, gin.H{"message": "Payment recorded successfully"})
}

func (h *BudgetHandler) GetPayments(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var payments []models.Payment
	db.Joins("JOIN budget_items ON budget_items.id = payments.budget_item_id").
		Where("budget_items.wedding_id = ?", weddingID).
		Order("payments.created_at DESC").
		Find(&payments)

	response.Success(c, payments)
}

func (h *BudgetHandler) GetBudgetCategories(c *gin.Context) {
	categories := []string{
		"场地租赁", "餐饮服务", "摄影摄像", "花艺布置", "司仪主持",
		"化妆造型", "婚纱礼服", "婚车租赁", "婚礼蛋糕", "场地布置",
		"音乐娱乐", "请柬印刷", "回礼礼品", "蜜月旅行", "其他",
	}

	response.Success(c, categories)
}

func (h *BudgetHandler) CheckBudgetAlerts(c *gin.Context) {
	weddingID := c.GetUint("wedding_id")

	db := database.GetDB()

	var wedding models.Wedding
	db.First(&wedding, weddingID)

	var alerts []gin.H

	var totalActual float64
	db.Model(&models.BudgetItem{}).Where("wedding_id = ?", weddingID).
		Select("COALESCE(SUM(actual_cost), 0)").Scan(&totalActual)

	if totalActual > wedding.Budget && wedding.Budget > 0 {
		overPercentage := ((totalActual - wedding.Budget) / wedding.Budget) * 100
		alerts = append(alerts, gin.H{
			"type":    "budget_overrun",
			"message": "预算超支 " + formatFloat(overPercentage) + "%",
			"amount":  totalActual - wedding.Budget,
		})
	}

	var overdueItems []models.BudgetItem
	db.Where("wedding_id = ? AND due_date < ? AND status NOT IN ('paid', 'cancelled')", weddingID, time.Now()).
		Find(&overdueItems)

	for _, item := range overdueItems {
		alerts = append(alerts, gin.H{
			"type":      "payment_overdue",
			"message":   "付款逾期: " + item.Category + " - " + item.Description,
			"item_id":   item.ID,
			"due_date":  item.DueDate,
			"amount":    item.ActualCost - item.PaidAmount,
		})
	}

	response.Success(c, alerts)
}
