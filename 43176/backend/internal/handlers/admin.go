package handlers

import (
	"errand-service/internal/models"
	"errand-service/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

type ListUsersQuery struct {
	Role     string `form:"role"`
	Status   string `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var query ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	db := h.db.Model(&models.User{})

	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	db.Count(&total)

	var users []models.User
	offset := (query.Page - 1) * query.PageSize
	db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":     total,
			"page":      query.Page,
			"page_size": query.PageSize,
			"items":     users,
		},
	})
}

func (h *AdminHandler) FreezeUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid user ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Reason is required"})
		return
	}

	result := h.db.Model(&models.User{}).Where("id = ?", id).
		Update("status", models.UserStatusFrozen)
	if result.Error != nil {
		logger.Errorf("Failed to freeze user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to freeze user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "User frozen successfully"})
}

func (h *AdminHandler) UnfreezeUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid user ID"})
		return
	}

	result := h.db.Model(&models.User{}).Where("id = ?", id).
		Update("status", models.UserStatusActive)
	if result.Error != nil {
		logger.Errorf("Failed to unfreeze user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to unfreeze user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "User unfrozen successfully"})
}

func (h *AdminHandler) ApproveCourier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid courier ID"})
		return
	}

	tx := h.db.Begin()

	if err := tx.Model(&models.CourierProfile{}).Where("user_id = ?", id).
		Update("status", "approved").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to approve courier"})
		return
	}

	if err := tx.Model(&models.User{}).Where("id = ?", id).
		Update("status", models.UserStatusVerified).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update user status"})
		return
	}

	var verification models.Verification
	tx.Where("user_id = ? AND status = ?", id, "pending").First(&verification)
	if verification.ID > 0 {
		now := time.Now()
		verification.Status = "approved"
		verification.ReviewedAt = &now
		tx.Save(&verification)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Courier approved successfully"})
}

func (h *AdminHandler) RejectCourier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid courier ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Rejection reason is required"})
		return
	}

	tx := h.db.Begin()

	if err := tx.Model(&models.CourierProfile{}).Where("user_id = ?", id).
		Updates(map[string]interface{}{"status": "rejected"}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to reject courier"})
		return
	}

	var verification models.Verification
	tx.Where("user_id = ? AND status = ?", id, "pending").First(&verification)
	if verification.ID > 0 {
		now := time.Now()
		verification.Status = "rejected"
		verification.Reason = req.Reason
		verification.ReviewedAt = &now
		tx.Save(&verification)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Courier rejected successfully"})
}
