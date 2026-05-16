package controllers

import (
	"fmt"
	"time"

	"finance-api/middleware"
	"finance-api/models"
	"finance-api/utils"

	"github.com/gin-gonic/gin"
)

type TransactionController struct{}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

type CreateTransactionRequest struct {
	AccountID  uint      `json:"account_id" binding:"required"`
	CategoryID uint      `json:"category_id" binding:"required"`
	Type       string    `json:"type" binding:"required,oneof=income expense"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	Remark     string    `json:"remark" binding:"max=255"`
	Date       time.Time `json:"date"`
}

type UpdateTransactionRequest struct {
	AccountID  uint      `json:"account_id"`
	CategoryID uint      `json:"category_id"`
	Type       string    `json:"type" binding:"omitempty,oneof=income expense"`
	Amount     float64   `json:"amount" binding:"omitempty,gt=0"`
	Remark     string    `json:"remark" binding:"max=255"`
	Date       time.Time `json:"date"`
}

func (ctrl *TransactionController) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	var account models.Account
	if result := utils.DB.Where("id = ? AND user_id = ?", req.AccountID, userID).First(&account); result.Error != nil {
		utils.NotFound(c, "Account not found")
		return
	}

	var category models.Category
	if result := utils.DB.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&category); result.Error != nil {
		utils.NotFound(c, "Category not found")
		return
	}

	if category.Type != req.Type {
		utils.BadRequest(c, "Category type mismatch with transaction type")
		return
	}

	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	transaction := models.Transaction{
		UserID:     userID,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		Type:       req.Type,
		Amount:     req.Amount,
		Remark:     req.Remark,
		Date:       req.Date,
	}

	tx := utils.DB.Begin()

	if result := tx.Create(&transaction); result.Error != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to create transaction: "+result.Error.Error())
		return
	}

	if req.Type == models.TransactionTypeIncome {
		account.Balance += req.Amount
	} else {
		account.Balance -= req.Amount
	}

	if result := tx.Save(&account); result.Error != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update account balance")
		return
	}

	tx.Commit()

	utils.Success(c, transaction)
}

func (ctrl *TransactionController) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	query := utils.DB.Where("user_id = ?", userID)

	if accountID := c.Query("account_id"); accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if transType := c.Query("type"); transType != "" {
		query = query.Where("type = ?", transType)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	var total int64
	query.Model(&models.Transaction{}).Count(&total)

	var transactions []models.Transaction
	offset := (page - 1) * pageSize
	if result := query.Preload("Account").Preload("Category").
		Order("date DESC, id DESC").
		Limit(pageSize).Offset(offset).
		Find(&transactions); result.Error != nil {
		utils.InternalError(c, "Failed to fetch transactions: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"transactions": transactions,
	})
}

func (ctrl *TransactionController) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var transaction models.Transaction
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Account").Preload("Category").
		First(&transaction); result.Error != nil {
		utils.NotFound(c, "Transaction not found")
		return
	}

	utils.Success(c, transaction)
}

func (ctrl *TransactionController) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var transaction models.Transaction
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&transaction); result.Error != nil {
		utils.NotFound(c, "Transaction not found")
		return
	}

	var req UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	tx := utils.DB.Begin()

	var account models.Account
	if result := tx.Where("id = ? AND user_id = ?", transaction.AccountID, userID).First(&account); result.Error != nil {
		tx.Rollback()
		utils.NotFound(c, "Account not found")
		return
	}

	if transaction.Type == models.TransactionTypeIncome {
		account.Balance -= transaction.Amount
	} else {
		account.Balance += transaction.Amount
	}

	newType := transaction.Type
	if req.Type != "" {
		newType = req.Type
	}

	newAmount := transaction.Amount
	if req.Amount > 0 {
		newAmount = req.Amount
	}

	newAccountID := transaction.AccountID
	if req.AccountID > 0 {
		newAccountID = req.AccountID
	}

	if newAccountID != transaction.AccountID {
		if result := tx.Save(&account); result.Error != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to update old account balance")
			return
		}

		if result := tx.Where("id = ? AND user_id = ?", newAccountID, userID).First(&account); result.Error != nil {
			tx.Rollback()
			utils.NotFound(c, "New account not found")
			return
		}
	}

	if newType == models.TransactionTypeIncome {
		account.Balance += newAmount
	} else {
		account.Balance -= newAmount
	}

	if result := tx.Save(&account); result.Error != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update account balance")
		return
	}

	newCategoryID := transaction.CategoryID
	if req.CategoryID > 0 {
		newCategoryID = req.CategoryID
	}

	var category models.Category
	if result := tx.Where("id = ? AND user_id = ?", newCategoryID, userID).First(&category); result.Error != nil {
		tx.Rollback()
		utils.NotFound(c, "Category not found")
		return
	}
	if category.Type != newType {
		tx.Rollback()
		utils.BadRequest(c, "Category type mismatch with transaction type")
		return
	}
	transaction.CategoryID = newCategoryID

	transaction.AccountID = newAccountID
	transaction.Type = newType
	transaction.Amount = newAmount
	if req.Remark != "" {
		transaction.Remark = req.Remark
	}
	if !req.Date.IsZero() {
		transaction.Date = req.Date
	}

	if result := tx.Save(&transaction); result.Error != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update transaction: "+result.Error.Error())
		return
	}

	tx.Commit()

	utils.Success(c, transaction)
}

func (ctrl *TransactionController) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var transaction models.Transaction
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&transaction); result.Error != nil {
		utils.NotFound(c, "Transaction not found")
		return
	}

	tx := utils.DB.Begin()

	var account models.Account
	if result := tx.Where("id = ? AND user_id = ?", transaction.AccountID, userID).First(&account); result.Error != nil {
		tx.Rollback()
		utils.NotFound(c, "Account not found")
		return
	}

	if transaction.Type == models.TransactionTypeIncome {
		account.Balance -= transaction.Amount
	} else {
		account.Balance += transaction.Amount
	}

	if result := tx.Save(&account); result.Error != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update account balance")
		return
	}

	if result := tx.Delete(&transaction); result.Error != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to delete transaction: "+result.Error.Error())
		return
	}

	tx.Commit()

	utils.Success(c, nil)
}
