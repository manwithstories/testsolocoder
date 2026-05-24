package handlers

import (
	"net/http"
	"photo-rental/internal/models"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettlementHandler struct{}

type CreateSettlementRequest struct {
	OrderID       uint    `json:"orderId" binding:"required"`
	DamageFee     float64 `json:"damageFee"`
	RefundDeposit float64 `json:"refundDeposit"`
	Remark        string  `json:"remark"`
}

func NewSettlementHandler() *SettlementHandler {
	return &SettlementHandler{}
}

func (h *SettlementHandler) CreateSettlement(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	var req CreateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, req.OrderID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.OwnerID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only the equipment owner can create settlement")
		return
	}

	if order.Status != "completed" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Can only create settlement for completed orders")
		return
	}

	var existingSettlement models.Settlement
	if result := database.DB.Where("order_id = ?", req.OrderID).First(&existingSettlement); result.Error == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Settlement already exists for this order")
		return
	}

	refundDeposit := order.Deposit
	if req.RefundDeposit > 0 {
		refundDeposit = req.RefundDeposit
	}

	if req.DamageFee > 0 {
		refundDeposit -= req.DamageFee
	}

	if refundDeposit < 0 {
		refundDeposit = 0
	}

	finalAmount := order.TotalRent + (order.Deposit - refundDeposit)

	settlement := models.Settlement{
		OrderID:       req.OrderID,
		TotalRent:     order.TotalRent,
		Deposit:       order.Deposit,
		RefundDeposit: refundDeposit,
		DamageFee:     req.DamageFee,
		FinalAmount:   finalAmount,
		Status:        "completed",
		Remark:        req.Remark,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&settlement).Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to create settlement: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create settlement")
		return
	}

	if err := tx.Model(&models.User{}).Where("id = ?", order.RenterID).
		Update("deposit_balance", gorm.Expr("deposit_balance + ?", refundDeposit)).Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to update renter deposit balance: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to refund deposit")
		return
	}

	if err := tx.Model(&models.User{}).Where("id = ?", order.OwnerID).
		Update("deposit_balance", gorm.Expr("deposit_balance + ?", order.TotalRent)).Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to update owner balance: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process payment")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to commit settlement: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create settlement")
		return
	}

	utils.Logger.Info("Settlement created: %d for order %d", settlement.ID, order.ID)
	utils.SuccessResponse(c, settlement)
}

func (h *SettlementHandler) GetSettlement(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid settlement ID")
		return
	}

	var settlement models.Settlement
	if err := database.DB.Preload("Order").First(&settlement, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Settlement not found")
		return
	}

	if userRole != "admin" && settlement.Order.RenterID != userID && settlement.Order.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only view your own settlements")
		return
	}

	utils.SuccessResponse(c, settlement)
}

func (h *SettlementHandler) GetOrderSettlement(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var settlement models.Settlement
	if err := database.DB.Preload("Order").Where("order_id = ?", orderID).First(&settlement).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Settlement not found")
		return
	}

	if userRole != "admin" && settlement.Order.RenterID != userID && settlement.Order.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only view your own settlements")
		return
	}

	utils.SuccessResponse(c, settlement)
}

func (h *SettlementHandler) GetMySettlements(c *gin.Context) {
	userID := c.GetUint("userId")

	var settlements []models.Settlement
	if err := database.DB.Preload("Order").
		Joins("JOIN orders ON orders.id = settlements.order_id").
		Where("orders.renter_id = ? OR orders.owner_id = ?", userID, userID).
		Order("settlements.created_at DESC").
		Find(&settlements).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch settlements")
		return
	}

	utils.SuccessResponse(c, settlements)
}
