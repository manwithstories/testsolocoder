package handlers

import (
	"errand-service/internal/models"
	"errand-service/internal/utils"
	"errand-service/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	db *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

type DepositRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=alipay wechat bank"`
}

type WithdrawRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	AccountType string  `json:"account_type" binding:"required,oneof=alipay wechat bank"`
	AccountNo   string  `json:"account_no" binding:"required"`
	AccountName string  `json:"account_name" binding:"required"`
}

type RefundRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
	Amount  float64 `json:"amount"`
}

type HistoryQuery struct {
	Type     string `form:"type"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (h *PaymentHandler) Deposit(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	tx := h.db.Begin()

	balanceBefore := user.Balance

	transaction := models.Transaction{
		UserID:        userID,
		Type:          models.TxTypeDeposit,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceBefore + req.Amount,
		Status:        models.TxStatusCompleted,
		Description:   "Account deposit",
		PaymentMethod: req.PaymentMethod,
		TransactionNo: utils.GenerateOrderNo(userID),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to create transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to process deposit"})
		return
	}

	if err := tx.Model(&user).Update("balance", gorm.Expr("balance + ?", req.Amount)).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to update balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update balance"})
		return
	}

	now := time.Now()
	tx.Model(&transaction).Update("completed_at", now)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Deposit successful",
		"data": gin.H{
			"transaction": transaction,
			"new_balance": balanceBefore + req.Amount,
		},
	})
}

func (h *PaymentHandler) Withdraw(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	if user.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Insufficient balance"})
		return
	}

	tx := h.db.Begin()

	withdrawRequest := models.WithdrawRequest{
		UserID:      userID,
		Amount:      req.Amount,
		AccountType: req.AccountType,
		AccountNo:   req.AccountNo,
		AccountName: req.AccountName,
		Status:      "pending",
	}

	if err := tx.Create(&withdrawRequest).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to create withdraw request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to submit withdrawal request"})
		return
	}

	if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", req.Amount)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to process withdrawal"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Withdrawal request submitted, waiting for review",
		"data":    withdrawRequest,
	})
}

func (h *PaymentHandler) Refund(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	var order models.Order
	if err := h.db.First(&order, req.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Order not found"})
		return
	}

	if order.PublisherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only the publisher can request refund"})
		return
	}

	if order.Status == models.OrderStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Completed orders cannot be refunded, please contact customer service"})
		return
	}

	refundAmount := req.Amount
	if refundAmount <= 0 || refundAmount > order.Reward {
		refundAmount = order.Reward
	}

	refundRequest := models.RefundRequest{
		OrderID: order.ID,
		UserID:  userID,
		Amount:  refundAmount,
		Reason:  req.Reason,
		Status:  "pending",
	}

	if err := h.db.Create(&refundRequest).Error; err != nil {
		logger.Errorf("Failed to create refund request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to submit refund request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Refund request submitted, waiting for review",
		"data":    refundRequest,
	})
}

func (h *PaymentHandler) History(c *gin.Context) {
	userID := c.GetUint("user_id")

	var query HistoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	db := h.db.Model(&models.Transaction{}).Where("user_id = ?", userID)

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	var total int64
	db.Count(&total)

	var transactions []models.Transaction
	offset := (query.Page - 1) * query.PageSize
	db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":     total,
			"page":      query.Page,
			"page_size": query.PageSize,
			"items":     transactions,
		},
	})
}
