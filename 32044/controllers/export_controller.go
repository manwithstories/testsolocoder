package controllers

import (
	"encoding/csv"
	"fmt"
	"time"

	"finance-api/middleware"
	"finance-api/models"
	"finance-api/utils"

	"github.com/gin-gonic/gin"
)

type ExportController struct{}

func NewExportController() *ExportController {
	return &ExportController{}
}

func (ctrl *ExportController) ExportTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.BadRequest(c, "start_date and end_date are required (format: YYYY-MM-DD)")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD")
		return
	}

	endDate = endDate.Add(24 * time.Hour)

	var transactions []models.Transaction
	if result := utils.DB.Where("user_id = ? AND date >= ? AND date < ?",
		userID, startDate, endDate).
		Preload("Account").Preload("Category").
		Order("date ASC").
		Find(&transactions); result.Error != nil {
		utils.InternalError(c, "Failed to fetch transactions: "+result.Error.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	filename := fmt.Sprintf("transactions_%s_%s.csv", startDateStr, endDateStr)
	c.Header("Content-Disposition", "attachment; filename="+filename)

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	headers := []string{"ID", "日期", "类型", "金额", "账户", "分类", "备注"}
	if err := writer.Write(headers); err != nil {
		utils.InternalError(c, "Failed to write CSV headers")
		return
	}

	for _, t := range transactions {
		typeStr := "支出"
		if t.Type == models.TransactionTypeIncome {
			typeStr = "收入"
		}

		row := []string{
			fmt.Sprintf("%d", t.ID),
			t.Date.Format("2006-01-02 15:04:05"),
			typeStr,
			fmt.Sprintf("%.2f", t.Amount),
			t.Account.Name,
			t.Category.Name,
			t.Remark,
		}
		if err := writer.Write(row); err != nil {
			utils.InternalError(c, "Failed to write CSV row")
			return
		}
	}
}
