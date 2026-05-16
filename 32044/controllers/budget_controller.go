package controllers

import (
	"time"

	"finance-api/middleware"
	"finance-api/models"
	"finance-api/utils"

	"github.com/gin-gonic/gin"
)

type BudgetController struct{}

func NewBudgetController() *BudgetController {
	return &BudgetController{}
}

type CreateBudgetRequest struct {
	CategoryID uint    `json:"category_id" binding:"required"`
	Month      string  `json:"month" binding:"required,len=7"`
	Limit      float64 `json:"limit" binding:"required,gt=0"`
}

type UpdateBudgetRequest struct {
	Limit float64 `json:"limit" binding:"required,gt=0"`
}

type BudgetWithUsage struct {
	ID         uint    `json:"id"`
	CategoryID uint    `json:"category_id"`
	Category   string  `json:"category"`
	Month      string  `json:"month"`
	Limit      float64 `json:"limit"`
	Spent      float64 `json:"spent"`
	Remaining  float64 `json:"remaining"`
	IsOver     bool    `json:"is_over"`
	Percentage float64 `json:"percentage"`
}

func (ctrl *BudgetController) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	if _, err := time.Parse("2006-01", req.Month); err != nil {
		utils.BadRequest(c, "Invalid month format, use YYYY-MM")
		return
	}

	var category models.Category
	if result := utils.DB.Where("id = ? AND user_id = ? AND type = ?", req.CategoryID, userID, models.CategoryTypeExpense).First(&category); result.Error != nil {
		utils.NotFound(c, "Expense category not found")
		return
	}

	var existingBudget models.Budget
	if result := utils.DB.Where("user_id = ? AND category_id = ? AND month = ?", userID, req.CategoryID, req.Month).First(&existingBudget); result.Error == nil {
		utils.BadRequest(c, "Budget for this category and month already exists")
		return
	}

	budget := models.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Month:      req.Month,
		Limit:      req.Limit,
	}

	if result := utils.DB.Create(&budget); result.Error != nil {
		utils.InternalError(c, "Failed to create budget: "+result.Error.Error())
		return
	}

	utils.Success(c, budget)
}

func (ctrl *BudgetController) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	month := c.Query("month")

	if month == "" {
		now := time.Now()
		month = now.Format("2006-01")
	}

	var budgets []models.Budget
	if result := utils.DB.Where("user_id = ? AND month = ?", userID, month).
		Preload("Category").
		Find(&budgets); result.Error != nil {
		utils.InternalError(c, "Failed to fetch budgets: "+result.Error.Error())
		return
	}

	startDate, _ := time.Parse("2006-01", month)
	endDate := startDate.AddDate(0, 1, 0)

	var result []BudgetWithUsage
	for _, budget := range budgets {
		var spent float64
		utils.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND category_id = ? AND type = ? AND date >= ? AND date < ?",
				userID, budget.CategoryID, models.TransactionTypeExpense, startDate, endDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&spent)

		remaining := budget.Limit - spent
		percentage := 0.0
		if budget.Limit > 0 {
			percentage = (spent / budget.Limit) * 100
		}

		result = append(result, BudgetWithUsage{
			ID:         budget.ID,
			CategoryID: budget.CategoryID,
			Category:   budget.Category.Name,
			Month:      budget.Month,
			Limit:      budget.Limit,
			Spent:      spent,
			Remaining:  remaining,
			IsOver:     remaining < 0,
			Percentage: percentage,
		})
	}

	utils.Success(c, result)
}

func (ctrl *BudgetController) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var budget models.Budget
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Category").
		First(&budget); result.Error != nil {
		utils.NotFound(c, "Budget not found")
		return
	}

	startDate, _ := time.Parse("2006-01", budget.Month)
	endDate := startDate.AddDate(0, 1, 0)

	var spent float64
	utils.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, budget.CategoryID, models.TransactionTypeExpense, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&spent)

	remaining := budget.Limit - spent
	percentage := 0.0
	if budget.Limit > 0 {
		percentage = (spent / budget.Limit) * 100
	}

	result := BudgetWithUsage{
		ID:         budget.ID,
		CategoryID: budget.CategoryID,
		Category:   budget.Category.Name,
		Month:      budget.Month,
		Limit:      budget.Limit,
		Spent:      spent,
		Remaining:  remaining,
		IsOver:     remaining < 0,
		Percentage: percentage,
	}

	utils.Success(c, result)
}

func (ctrl *BudgetController) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var budget models.Budget
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&budget); result.Error != nil {
		utils.NotFound(c, "Budget not found")
		return
	}

	var req UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	budget.Limit = req.Limit

	if result := utils.DB.Save(&budget); result.Error != nil {
		utils.InternalError(c, "Failed to update budget: "+result.Error.Error())
		return
	}

	utils.Success(c, budget)
}

func (ctrl *BudgetController) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var budget models.Budget
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&budget); result.Error != nil {
		utils.NotFound(c, "Budget not found")
		return
	}

	if result := utils.DB.Delete(&budget); result.Error != nil {
		utils.InternalError(c, "Failed to delete budget: "+result.Error.Error())
		return
	}

	utils.Success(c, nil)
}
