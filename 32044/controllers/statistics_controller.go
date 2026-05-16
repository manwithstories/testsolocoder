package controllers

import (
	"time"

	"finance-api/middleware"
	"finance-api/models"
	"finance-api/utils"

	"github.com/gin-gonic/gin"
)

type StatisticsController struct{}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{}
}

type CategoryStat struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type AccountStat struct {
	AccountID   uint    `json:"account_id"`
	AccountName string  `json:"account_name"`
	Currency    string  `json:"currency"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
	Balance     float64 `json:"balance"`
}

type MonthlyStatistics struct {
	Month             string         `json:"month"`
	TotalIncome       float64        `json:"total_income"`
	TotalExpense      float64        `json:"total_expense"`
	NetIncome         float64        `json:"net_income"`
	IncomeByCategory  []CategoryStat `json:"income_by_category"`
	ExpenseByCategory []CategoryStat `json:"expense_by_category"`
	ByAccount         []AccountStat  `json:"by_account"`
}

func (ctrl *StatisticsController) GetMonthly(c *gin.Context) {
	userID := middleware.GetUserID(c)
	month := c.Query("month")

	if month == "" {
		now := time.Now()
		month = now.Format("2006-01")
	}

	startDate, err := time.Parse("2006-01", month)
	if err != nil {
		utils.BadRequest(c, "Invalid month format, use YYYY-MM")
		return
	}
	endDate := startDate.AddDate(0, 1, 0)

	var totalIncome, totalExpense float64

	utils.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeIncome, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	utils.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeExpense, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	incomeByCategory := getCategoryStats(userID, models.TransactionTypeIncome, startDate, endDate, totalIncome)
	expenseByCategory := getCategoryStats(userID, models.TransactionTypeExpense, startDate, endDate, totalExpense)

	byAccount := getAccountStats(userID, startDate, endDate)

	stats := MonthlyStatistics{
		Month:             month,
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		NetIncome:         totalIncome - totalExpense,
		IncomeByCategory:  incomeByCategory,
		ExpenseByCategory: expenseByCategory,
		ByAccount:         byAccount,
	}

	utils.Success(c, stats)
}

func getCategoryStats(userID uint, transType string, startDate, endDate time.Time, total float64) []CategoryStat {
	type CategorySum struct {
		CategoryID uint
		Amount     float64
	}

	var sums []CategorySum
	utils.DB.Model(&models.Transaction{}).
		Select("category_id, COALESCE(SUM(amount), 0) as amount").
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, transType, startDate, endDate).
		Group("category_id").
		Scan(&sums)

	var stats []CategoryStat
	for _, sum := range sums {
		var category models.Category
		utils.DB.Select("name").First(&category, sum.CategoryID)

		percentage := 0.0
		if total > 0 {
			percentage = (sum.Amount / total) * 100
		}

		stats = append(stats, CategoryStat{
			CategoryID:   sum.CategoryID,
			CategoryName: category.Name,
			Amount:       sum.Amount,
			Percentage:   percentage,
		})
	}

	return stats
}

func getAccountStats(userID uint, startDate, endDate time.Time) []AccountStat {
	var accounts []models.Account
	utils.DB.Where("user_id = ?", userID).Find(&accounts)

	var stats []AccountStat
	for _, account := range accounts {
		var income, expense float64

		utils.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND account_id = ? AND type = ? AND date >= ? AND date < ?",
				userID, account.ID, models.TransactionTypeIncome, startDate, endDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&income)

		utils.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND account_id = ? AND type = ? AND date >= ? AND date < ?",
				userID, account.ID, models.TransactionTypeExpense, startDate, endDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&expense)

		stats = append(stats, AccountStat{
			AccountID:   account.ID,
			AccountName: account.Name,
			Currency:    account.Currency,
			Income:      income,
			Expense:     expense,
			Balance:     account.Balance,
		})
	}

	return stats
}
