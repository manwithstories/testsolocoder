package handlers

import (
	"errors"
	"net/http"
	"photo-rental/internal/models"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	redispkg "photo-rental/pkg/redis"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct{}

type CreateOrderRequest struct {
	EquipmentID     uint   `json:"equipmentId" binding:"required"`
	StartDate       string `json:"startDate" binding:"required"`
	EndDate         string `json:"endDate" binding:"required"`
	DeliveryMethod  string `json:"deliveryMethod" binding:"required"`
	DeliveryAddress string `json:"deliveryAddress"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed rented completed cancelled"`
	Reason string `json:"reason,omitempty"`
}

type RejectOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("userId")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end date format")
		return
	}

	if startDate.Before(time.Now().Truncate(24*time.Hour)) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Start date cannot be in the past")
		return
	}

	if endDate.Before(startDate) {
		utils.ErrorResponse(c, http.StatusBadRequest, "End date must be after start date")
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, req.EquipmentID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Equipment not found")
		return
	}

	if equipment.OwnerID == userID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot rent your own equipment")
		return
	}

	if equipment.Status != "available" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Equipment is not available for rent")
		return
	}

	conflicts, err := redispkg.CheckEquipmentConflict(req.EquipmentID, startDate, endDate)
	if err != nil {
		utils.Logger.Error("Failed to check equipment conflict: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check availability")
		return
	}

	if len(conflicts) > 0 {
		utils.ErrorResponse(c, http.StatusConflict, "Equipment is not available for the selected dates")
		return
	}

	days := int(endDate.Sub(startDate).Hours()/24) + 1
	totalRent := equipment.DailyRent * float64(days)

	order := models.Order{
		EquipmentID:     req.EquipmentID,
		RenterID:        userID,
		OwnerID:         equipment.OwnerID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		TotalRent:       totalRent,
		Deposit:         equipment.Deposit,
		Status:          "pending",
		DeliveryMethod:  req.DeliveryMethod,
		DeliveryAddress: req.DeliveryAddress,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to create order: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	var dates []string
	current := startDate
	for !current.After(endDate) {
		dates = append(dates, current.Format("2006-01-02"))
		current = current.AddDate(0, 0, 1)
	}

	if err := redispkg.SetEquipmentAvailability(req.EquipmentID, dates, order.ID); err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to set equipment availability: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reserve dates")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		redispkg.RemoveEquipmentAvailability(req.EquipmentID, startDate, endDate)
		utils.Logger.Error("Failed to commit order: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	utils.Logger.Info("Order created: %d by user %d", order.ID, userID)
	utils.SuccessResponse(c, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.Preload("Equipment.Images").Preload("Renter").Preload("Owner").First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if userRole != "admin" && order.RenterID != userID && order.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only view your own orders")
		return
	}

	order.Renter.Password = ""
	order.Renter.IDCard = ""
	order.Owner.Password = ""
	order.Owner.IDCard = ""

	utils.SuccessResponse(c, order)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	status := c.Query("status")

	var orders []models.Order
	db := database.DB.Preload("Equipment.Images").Preload("Renter").Preload("Owner")

	switch userRole {
	case "owner":
		db = db.Where("owner_id = ?", userID)
	case "renter":
		db = db.Where("renter_id = ?", userID)
	case "admin":
	default:
		db = db.Where("renter_id = ? OR owner_id = ?", userID, userID)
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Order("created_at DESC").Find(&orders).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	for i := range orders {
		orders[i].Renter.Password = ""
		orders[i].Renter.IDCard = ""
		orders[i].Owner.Password = ""
		orders[i].Owner.IDCard = ""
	}

	utils.SuccessResponse(c, orders)
}

func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.OwnerID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only the equipment owner can confirm orders")
		return
	}

	if order.Status != "pending" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only pending orders can be confirmed")
		return
	}

	tx := database.DB.Begin()

	if err := tx.Model(&order).Update("status", "confirmed").Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to confirm order")
		return
	}

	if err := tx.Model(&models.Equipment{}).Where("id = ?", order.EquipmentID).Update("status", "rented").Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update equipment status")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to confirm order")
		return
	}

	utils.Logger.Info("Order confirmed: %d by owner %d", id, userID)
	utils.SuccessResponse(c, nil)
}

func (h *OrderHandler) RejectOrder(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req RejectOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Rejection reason is required")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.OwnerID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only the equipment owner can reject orders")
		return
	}

	if order.Status != "pending" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only pending orders can be rejected")
		return
	}

	startDate, _ := time.Parse("2006-01-02", order.StartDate)
	endDate, _ := time.Parse("2006-01-02", order.EndDate)

	tx := database.DB.Begin()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":        "cancelled",
		"reject_reason": req.Reason,
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reject order")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reject order")
		return
	}

	if err := redispkg.RemoveEquipmentAvailability(order.EquipmentID, startDate, endDate); err != nil {
		utils.Logger.Error("Failed to remove equipment availability: %v", err)
	}

	utils.Logger.Info("Order rejected: %d by owner %d", id, userID)
	utils.SuccessResponse(c, nil)
}

func (h *OrderHandler) StartRental(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.RenterID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only the renter can start rental")
		return
	}

	if order.Status != "confirmed" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only confirmed orders can be started")
		return
	}

	if err := database.DB.Model(&order).Update("status", "rented").Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to start rental")
		return
	}

	utils.Logger.Info("Rental started: %d by renter %d", id, userID)
	utils.SuccessResponse(c, nil)
}

func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.OwnerID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only the equipment owner can complete orders")
		return
	}

	if order.Status != "rented" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only rented orders can be completed")
		return
	}

	startDate, _ := time.Parse("2006-01-02", order.StartDate)
	endDate, _ := time.Parse("2006-01-02", order.EndDate)

	tx := database.DB.Begin()

	if err := tx.Model(&order).Update("status", "completed").Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to complete order")
		return
	}

	if err := tx.Model(&models.Equipment{}).Where("id = ?", order.EquipmentID).Update("status", "available").Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update equipment status")
		return
	}

	settlement := models.Settlement{
		OrderID:       order.ID,
		TotalRent:     order.TotalRent,
		Deposit:       order.Deposit,
		RefundDeposit: order.Deposit,
		FinalAmount:   order.TotalRent,
		Status:        "completed",
	}

	if err := tx.Create(&settlement).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create settlement")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to complete order")
		return
	}

	if err := redispkg.RemoveEquipmentAvailability(order.EquipmentID, startDate, endDate); err != nil {
		utils.Logger.Error("Failed to remove equipment availability: %v", err)
	}

	utils.Logger.Info("Order completed: %d", id)
	utils.SuccessResponse(c, nil)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.RenterID != userID && order.OwnerID != userID && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only cancel your own orders")
		return
	}

	if order.Status != "pending" && order.Status != "confirmed" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only pending or confirmed orders can be cancelled")
		return
	}

	startDate, _ := time.Parse("2006-01-02", order.StartDate)
	endDate, _ := time.Parse("2006-01-02", order.EndDate)

	tx := database.DB.Begin()

	if err := tx.Model(&order).Update("status", "cancelled").Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	if order.Status == "confirmed" {
		if err := tx.Model(&models.Equipment{}).Where("id = ?", order.EquipmentID).Update("status", "available").Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update equipment status")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	if err := redispkg.RemoveEquipmentAvailability(order.EquipmentID, startDate, endDate); err != nil {
		utils.Logger.Error("Failed to remove equipment availability: %v", err)
	}

	utils.Logger.Info("Order cancelled: %d by user %d", id, userID)
	utils.SuccessResponse(c, nil)
}

func (h *OrderHandler) GetOrderStatusHistory(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if userRole != "admin" && order.RenterID != userID && order.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only view your own orders")
		return
	}

	statusFlow := []map[string]interface{}{
		{"status": "pending", "timestamp": order.CreatedAt},
	}

	if order.UpdatedAt != order.CreatedAt {
		statusFlow = append(statusFlow, map[string]interface{}{
			"status":    order.Status,
			"timestamp": order.UpdatedAt,
		})
	}

	utils.SuccessResponse(c, statusFlow)
}

func validateOrderTransition(currentStatus, newStatus string) error {
	validTransitions := map[string][]string{
		"pending":   {"confirmed", "cancelled"},
		"confirmed": {"rented", "cancelled"},
		"rented":    {"completed"},
		"completed": {},
		"cancelled": {},
	}

	validNext, ok := validTransitions[currentStatus]
	if !ok {
		return errors.New("invalid current status")
	}

	for _, status := range validNext {
		if status == newStatus {
			return nil
		}
	}

	return errors.New("invalid status transition")
}
