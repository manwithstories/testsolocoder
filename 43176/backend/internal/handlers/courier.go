package handlers

import (
	"errand-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CourierHandler struct {
	db *gorm.DB
}

func NewCourierHandler(db *gorm.DB) *CourierHandler {
	return &CourierHandler{db: db}
}

type ApplyRequest struct {
	IDCardNo    string `json:"id_card_no" binding:"required"`
	IDCardName  string `json:"id_card_name" binding:"required"`
	IDCardFront string `json:"id_card_front" binding:"required"`
	IDCardBack  string `json:"id_card_back" binding:"required"`
	Experience  string `json:"experience"`
	Vehicle     string `json:"vehicle"`
}

type MyTasksQuery struct {
	Status   string `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (h *CourierHandler) Apply(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	if user.Role == "courier" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "You are already a courier"})
		return
	}

	var req ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	tx := h.db.Begin()

	user.Role = "courier"
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update user role"})
		return
	}

	courierProfile := models.CourierProfile{
		UserID: userID,
		Status: "pending",
		Level:  1,
		Rating: 5.0,
	}
	if err := tx.Create(&courierProfile).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create courier profile"})
		return
	}

	verification := models.Verification{
		UserID:      userID,
		IDCardNo:    req.IDCardNo,
		IDCardName:  req.IDCardName,
		IDCardFront: req.IDCardFront,
		IDCardBack:  req.IDCardBack,
		Status:      "pending",
	}
	if err := tx.Create(&verification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create verification"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Courier application submitted, waiting for approval",
		"data": gin.H{
			"courier_profile": courierProfile,
			"verification":    verification,
		},
	})
}

func (h *CourierHandler) GetMyTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	if userRole != "courier" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only couriers can access this endpoint"})
		return
	}

	var query MyTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	db := h.db.Model(&models.Order{}).
		Where("courier_id = ?", userID).
		Preload("Task").Preload("Task.Images").Preload("Publisher")

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	db.Count(&total)

	var orders []models.Order
	offset := (query.Page - 1) * query.PageSize
	db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":     total,
			"page":      query.Page,
			"page_size": query.PageSize,
			"items":     orders,
		},
	})
}
