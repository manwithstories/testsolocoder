package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ship-rental-platform/internal/config"
	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FinanceHandler struct{}

func NewFinanceHandler() *FinanceHandler {
	return &FinanceHandler{}
}

func (h *FinanceHandler) CreateTransaction(c *gin.Context) {
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	cfg := config.GetConfig()
	currency := req.Currency
	if currency == "" {
		currency = cfg.Finance.DefaultCurrency
	}

	payerID, _ := uuid.Parse(req.PayerID)
	payeeID, _ := uuid.Parse(req.PayeeID)

	var rentalID *uuid.UUID
	if req.RentalID != "" {
		rid, err := uuid.Parse(req.RentalID)
		if err == nil {
			rentalID = &rid
		}
	}

	var platformFee, dockFee float64
	if req.TransactionType == model.TransactionTypeIncome {
		platformFee = req.Amount * cfg.Finance.PlatformFeeRate
		dockFee = req.Amount * cfg.Finance.DockFeeRate
	}

	transaction := model.Transaction{
		RentalID:        rentalID,
		PayerID:         payerID,
		PayeeID:         payeeID,
		Amount:          req.Amount,
		Currency:        currency,
		TransactionType: req.TransactionType,
		Description:     req.Description,
		Status:          model.TransactionStatusPending,
		PaymentMethod:   req.PaymentMethod,
		PlatformFee:     utils.RoundTo(platformFee, 2),
		DockFee:         utils.RoundTo(dockFee, 2),
		NetAmount:       utils.RoundTo(req.Amount-platformFee-dockFee, 2),
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	utils.Created(c, transaction)
}

func (h *FinanceHandler) GetTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var transactions []model.Transaction
	query := database.DB.Preload("Payer").Preload("Payee").Preload("Rental")

	if role.(string) != string(model.RoleAdmin) {
		uid, _ := uuid.Parse(userID.(string))
		query = query.Where("payer_id = ? OR payee_id = ?", uid, uid)
	}

	if transactionType := c.Query("transaction_type"); transactionType != "" {
		query = query.Where("transaction_type = ?", transactionType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	var total int64
	query.Model(&model.Transaction{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "20"), 20)
	offset := (page - 1) * pageSize

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	utils.Paginated(c, transactions, total, page, pageSize)
}

func (h *FinanceHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var transaction model.Transaction
	if err := database.DB.Preload("Payer").Preload("Payee").Preload("Rental").First(&transaction, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	utils.Success(c, transaction)
}

func (h *FinanceHandler) UpdateTransactionStatus(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var transaction model.Transaction
	if err := database.DB.First(&transaction, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	status := c.PostForm("status")
	validStatuses := map[string]bool{
		"completed": true,
		"failed":    true,
		"refunded":  true,
	}

	if !validStatuses[status] {
		utils.Error(c, http.StatusBadRequest, "Invalid status")
		return
	}

	transaction.Status = model.TransactionStatus(status)

	if status == "completed" {
		now := time.Now()
		transaction.PaidAt = &now
	}

	if err := database.DB.Save(&transaction).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update transaction status")
		return
	}

	utils.Success(c, transaction)
}

func (h *FinanceHandler) GetSettlements(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var settlements []model.Settlement
	query := database.DB.Preload("User")

	if role.(string) != string(model.RoleAdmin) {
		uid, _ := uuid.Parse(userID.(string))
		query = query.Where("user_id = ?", uid)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&model.Settlement{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "20"), 20)
	offset := (page - 1) * pageSize

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&settlements).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch settlements")
		return
	}

	utils.Paginated(c, settlements, total, page, pageSize)
}

func (h *FinanceHandler) CreateSettlement(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := uuid.Parse(userID.(string))

	monthStr := c.PostForm("month")
	yearStr := c.PostForm("year")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	periodStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := utils.EndOfMonth(periodStart)

	var totalIncome, totalExpense float64

	database.DB.Model(&model.Transaction{}).
		Where("payee_id = ? AND transaction_type = ? AND status = ? AND created_at BETWEEN ? AND ?",
			uid, model.TransactionTypeIncome, model.TransactionStatusCompleted, periodStart, periodEnd).
		Select("COALESCE(SUM(net_amount), 0)").
		Scan(&totalIncome)

	database.DB.Model(&model.Transaction{}).
		Where("payer_id = ? AND transaction_type = ? AND status = ? AND created_at BETWEEN ? AND ?",
			uid, model.TransactionTypeExpense, model.TransactionStatusCompleted, periodStart, periodEnd).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	cfg := config.GetConfig()
	settlement := model.Settlement{
		UserID:       uid,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		TotalIncome:  utils.RoundTo(totalIncome, 2),
		TotalExpense: utils.RoundTo(totalExpense, 2),
		Currency:     cfg.Finance.DefaultCurrency,
		Status:       "pending",
	}

	if err := database.DB.Create(&settlement).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create settlement")
		return
	}

	utils.Created(c, settlement)
}

func (h *FinanceHandler) GetFinancialSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	uid, _ := uuid.Parse(userID.(string))

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	var result struct {
		TotalIncome  float64 `json:"total_income"`
		TotalExpense float64 `json:"total_expense"`
		NetBalance   float64 `json:"net_balance"`
		PlatformFee  float64 `json:"platform_fee"`
		DockFee      float64 `json:"dock_fee"`
	}

	baseQuery := database.DB.Model(&model.Transaction{}).
		Where("status = ? AND created_at BETWEEN ? AND ?",
			model.TransactionStatusCompleted, startDate, endDate)

	if role.(string) != string(model.RoleAdmin) {
		baseQuery = baseQuery.Where("payer_id = ? OR payee_id = ?", uid, uid)
	}

	incomeQuery := baseQuery.Where("transaction_type = ?", model.TransactionTypeIncome)
	expenseQuery := baseQuery.Where("transaction_type = ?", model.TransactionTypeExpense)

	incomeQuery.Select("COALESCE(SUM(net_amount), 0)").Scan(&result.TotalIncome)
	expenseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&result.TotalExpense)

	baseQuery.Select("COALESCE(SUM(platform_fee), 0)").Scan(&result.PlatformFee)
	baseQuery.Select("COALESCE(SUM(dock_fee), 0)").Scan(&result.DockFee)

	result.NetBalance = utils.RoundTo(result.TotalIncome-result.TotalExpense, 2)
	result.TotalIncome = utils.RoundTo(result.TotalIncome, 2)
	result.TotalExpense = utils.RoundTo(result.TotalExpense, 2)
	result.PlatformFee = utils.RoundTo(result.PlatformFee, 2)
	result.DockFee = utils.RoundTo(result.DockFee, 2)

	utils.Success(c, result)
}

func (h *FinanceHandler) ExportFinancialReport(c *gin.Context) {
	var req model.FinancialReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var transactions []model.Transaction
	query := database.DB.Preload("Payer").Preload("Payee").
		Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)

	if req.UserID != "" {
		uid, _ := uuid.Parse(req.UserID)
		query = query.Where("payer_id = ? OR payee_id = ?", uid, uid)
	}
	if req.TransactionType != "" {
		query = query.Where("transaction_type = ?", req.TransactionType)
	}

	if err := query.Order("created_at ASC").Find(&transactions).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	if req.Format == "csv" {
		csvData := generateFinancialCSV(transactions)
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=financial_report_%s.csv", time.Now().Format("20060102")))
		c.String(http.StatusOK, csvData)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=financial_report_%s.pdf", time.Now().Format("20060102")))
	c.String(http.StatusOK, "PDF export placeholder - install unipdf for full PDF support")
}

func (h *FinanceHandler) ExportMonthlyReport(c *gin.Context) {
	var req model.MonthlyReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	periodStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := utils.EndOfMonth(periodStart)

	var transactions []model.Transaction
	query := database.DB.Preload("Payer").Preload("Payee").
		Where("created_at BETWEEN ? AND ?", periodStart, periodEnd)

	if req.UserID != "" {
		uid, _ := uuid.Parse(req.UserID)
		query = query.Where("payer_id = ? OR payee_id = ?", uid, uid)
	}

	if err := query.Order("created_at ASC").Find(&transactions).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	if req.Format == "csv" {
		csvData := generateFinancialCSV(transactions)
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=monthly_report_%d_%02d.csv", req.Year, req.Month))
		c.String(http.StatusOK, csvData)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=monthly_report_%d_%02d.pdf", req.Year, req.Month))
	c.String(http.StatusOK, "PDF export placeholder - install unipdf for full PDF support")
}

func generateFinancialCSV(transactions []model.Transaction) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{
		"Date", "Type", "Payer", "Payee", "Amount", "Currency",
		"Platform Fee", "Dock Fee", "Net Amount", "Status", "Description",
	}
	writer.Write(headers)

	for _, t := range transactions {
		row := []string{
			t.CreatedAt.Format("2006-01-02 15:04:05"),
			string(t.TransactionType),
			t.Payer.FullName,
			t.Payee.FullName,
			fmt.Sprintf("%.2f", t.Amount),
			t.Currency,
			fmt.Sprintf("%.2f", t.PlatformFee),
			fmt.Sprintf("%.2f", t.DockFee),
			fmt.Sprintf("%.2f", t.NetAmount),
			string(t.Status),
			t.Description,
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.String()
}
