package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/config"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type DepositRequest struct {
	Amount        float64 `json:"amount" binding:"required,min=1"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
}

type WithdrawRequest struct {
	Amount        float64 `json:"amount" binding:"required,min=1"`
	Currency      string  `json:"currency"`
	BankName      string  `json:"bankName" binding:"required"`
	BankAccount   string  `json:"bankAccount" binding:"required"`
	AccountHolder string  `json:"accountHolder" binding:"required"`
}

type ProcessWithdrawInput struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
	Notes  string `json:"notes"`
}

func GetWallet(c *gin.Context) {
	userID, _ := c.Get("userId")

	var wallet models.Wallet
	if err := database.DB.Preload("Transactions").Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func GetTransactions(c *gin.Context) {
	userID, _ := c.Get("userId")
	transactionType := c.Query("type")
	status := c.Query("status")

	var transactions []models.Transaction
	query := database.DB.Preload("Booking").Preload("User").Where("user_id = ?", userID)

	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	var transaction models.Transaction
	if err := database.DB.Preload("Booking").Preload("User").Where("id = ?", id).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func Deposit(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wallet models.Wallet
	if err := database.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	transaction := models.Transaction{
		WalletID:      wallet.ID,
		UserID:        userID.(uuid.UUID),
		Type:          models.TransactionTypeDeposit,
		Amount:        req.Amount,
		Currency:      req.Currency,
		BalanceAfter:  wallet.Balance + req.Amount,
		Status:        models.TransactionStatusPending,
		Description:   fmt.Sprintf("Deposit of %.2f %s", req.Amount, req.Currency),
		PaymentMethod: req.PaymentMethod,
		TransactionID: fmt.Sprintf("txn_%s", uuid.New().String()),
	}

	tx := database.DB.Begin()

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	wallet.Balance += req.Amount
	wallet.TotalIncome += req.Amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
		return
	}

	now := time.Now()
	transaction.Status = models.TransactionStatusCompleted
	transaction.CompletedAt = &now
	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete transaction"})
		return
	}

	tx.Commit()

	createNotification(c, userID.(uuid.UUID), models.NotificationTypePaymentReceived, "Deposit Successful", fmt.Sprintf("You have deposited %.2f %s", req.Amount, req.Currency))

	c.JSON(http.StatusCreated, gin.H{"message": "Deposit successful", "transaction": transaction, "wallet": wallet})
}

func Withdraw(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := config.AppConfig
	if req.Amount < cfg.Platform.MinWithdrawAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Minimum withdrawal amount is %.2f", cfg.Platform.MinWithdrawAmount)})
		return
	}

	var wallet models.Wallet
	if err := database.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	if wallet.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	withdrawRequest := models.WithdrawRequest{
		UserID:        userID.(uuid.UUID),
		WalletID:      wallet.ID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		BankName:      req.BankName,
		BankAccount:   req.BankAccount,
		AccountHolder: req.AccountHolder,
		Status:        models.WithdrawStatusPending,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&withdrawRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create withdrawal request"})
		return
	}

	wallet.Balance -= req.Amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Withdrawal request submitted", "withdrawRequest": withdrawRequest})
}

func GetWithdrawRequests(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	var requests []models.WithdrawRequest
	query := database.DB.Preload("User")

	if userRole != models.RoleAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Order("created_at DESC").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch withdrawal requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func ProcessWithdrawRequest(c *gin.Context) {
	id := c.Param("id")
	adminID, _ := c.Get("userId")

	var req ProcessWithdrawInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var withdrawRequest models.WithdrawRequest
	if err := database.DB.Where("id = ?", id).First(&withdrawRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Withdrawal request not found"})
		return
	}

	if withdrawRequest.Status != models.WithdrawStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request already processed"})
		return
	}

	now := time.Now()
	adminUid := adminID.(uuid.UUID)

	updates := map[string]interface{}{
		"status":      req.Status,
		"admin_notes": req.Notes,
		"processed_at": &now,
		"processed_by": &adminUid,
	}

	if err := database.DB.Model(&withdrawRequest).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if req.Status == "rejected" {
		var wallet models.Wallet
		if database.DB.Where("id = ?", withdrawRequest.WalletID).First(&wallet).Error == nil {
			wallet.Balance += withdrawRequest.Amount
			database.DB.Save(&wallet)
		}
		createNotification(c, withdrawRequest.UserID, models.NotificationTypeWithdrawRejected, "Withdrawal Rejected", "Your withdrawal request has been rejected")
	} else if req.Status == "approved" {
		createNotification(c, withdrawRequest.UserID, models.NotificationTypeWithdrawApproved, "Withdrawal Approved", "Your withdrawal request has been approved")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request processed successfully"})
}

func ProcessBookingPayment(c *gin.Context, booking *models.Booking) error {
	var studentWallet models.Wallet
	if err := database.DB.Where("user_id = ?", booking.StudentID).First(&studentWallet).Error; err != nil {
		return fmt.Errorf("student wallet not found")
	}

	if studentWallet.Balance < booking.TotalAmount {
		return fmt.Errorf("insufficient balance")
	}

	cfg := config.AppConfig
	commission := booking.TotalAmount * cfg.Platform.CommissionRate
	teacherAmount := booking.TotalAmount - commission

	tx := database.DB.Begin()

	studentTransaction := models.Transaction{
		WalletID:      studentWallet.ID,
		UserID:        booking.StudentID,
		Type:          models.TransactionTypePayment,
		Amount:        booking.TotalAmount,
		BalanceAfter:  studentWallet.Balance - booking.TotalAmount,
		Status:        models.TransactionStatusCompleted,
		Description:   fmt.Sprintf("Payment for booking with teacher"),
		ReferenceID:   booking.ID.String(),
		ReferenceType: "booking",
		BookingID:     &booking.ID,
		CompletedAt:   &time.Time{},
	}
	studentTransaction.CompletedAt = nil
	now := time.Now()
	studentTransaction.CompletedAt = &now

	if err := tx.Create(&studentTransaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create student transaction")
	}

	studentWallet.Balance -= booking.TotalAmount
	studentWallet.TotalSpent += booking.TotalAmount
	if err := tx.Save(&studentWallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update student wallet")
	}

	var teacherWallet models.Wallet
	if err := tx.Where("user_id = ?", booking.TeacherID).First(&teacherWallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("teacher wallet not found")
	}

	teacherTransaction := models.Transaction{
		WalletID:      teacherWallet.ID,
		UserID:        booking.TeacherID,
		Type:          models.TransactionTypeDeposit,
		Amount:        teacherAmount,
		BalanceAfter:  teacherWallet.Balance + teacherAmount,
		Status:        models.TransactionStatusCompleted,
		Description:   fmt.Sprintf("Earnings from booking (minus %.1f%% commission)", cfg.Platform.CommissionRate*100),
		ReferenceID:   booking.ID.String(),
		ReferenceType: "booking",
		BookingID:     &booking.ID,
		CompletedAt:   &now,
	}
	if err := tx.Create(&teacherTransaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create teacher transaction")
	}

	teacherWallet.Balance += teacherAmount
	teacherWallet.TotalIncome += teacherAmount
	if err := tx.Save(&teacherWallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update teacher wallet")
	}

	platformTransaction := models.Transaction{
		WalletID:      uuid.Nil,
		UserID:        uuid.Nil,
		Type:          models.TransactionTypeCommission,
		Amount:        commission,
		Status:        models.TransactionStatusCompleted,
		Description:   fmt.Sprintf("Platform commission (%.1f%%)", cfg.Platform.CommissionRate*100),
		ReferenceID:   booking.ID.String(),
		ReferenceType: "booking",
		BookingID:     &booking.ID,
		CompletedAt:   &now,
	}
	if err := tx.Create(&platformTransaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create platform transaction")
	}

	tx.Commit()

	createNotification(c, booking.TeacherID, models.NotificationTypePaymentReceived, "New Earnings", fmt.Sprintf("You received %.2f from a completed lesson", teacherAmount))

	return nil
}

func GetEarningsSummary(c *gin.Context) {
	userID, _ := c.Get("userId")

	var wallet models.Wallet
	if err := database.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	var totalEarnings float64
	var monthEarnings float64
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var transactions []models.Transaction
	database.DB.Where("user_id = ? AND type = ? AND status = ?", userID, models.TransactionTypeDeposit, models.TransactionStatusCompleted).
		Find(&transactions)

	for _, t := range transactions {
		totalEarnings += t.Amount
		if t.CreatedAt.After(monthStart) {
			monthEarnings += t.Amount
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":       wallet.Balance,
		"totalEarnings": totalEarnings,
		"monthEarnings": monthEarnings,
		"totalSpent":    wallet.TotalSpent,
		"currency":      wallet.Currency,
	})
}

func GetPaymentConfigs(c *gin.Context) {
	var configs []models.PaymentConfig
	if err := database.DB.Where("is_active = ?", true).Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment configs"})
		return
	}

	for i := range configs {
		configs[i].APISecret = "***"
		configs[i].APIKey = "***"
	}

	c.JSON(http.StatusOK, configs)
}

func UpdatePaymentConfig(c *gin.Context) {
	id := c.Param("id")

	var config models.PaymentConfig
	if err := database.DB.Where("id = ?", id).First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment config not found"})
		return
	}

	var req struct {
		APIKey      string `json:"apiKey"`
		APISecret   string `json:"apiSecret"`
		WebhookURL  string `json:"webhookUrl"`
		IsActive    bool   `json:"isActive"`
		ExtraConfig string `json:"extraConfig"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"api_key":      req.APIKey,
		"api_secret":   req.APISecret,
		"webhook_url":  req.WebhookURL,
		"is_active":    req.IsActive,
		"extra_config": req.ExtraConfig,
	}

	if err := database.DB.Model(&config).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment config updated successfully"})
}
