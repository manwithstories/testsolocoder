package handlers

import (
	"errand-service/internal/models"
	"errand-service/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewOrderHandler(db *gorm.DB, redisClient *redis.Client) *OrderHandler {
	return &OrderHandler{db: db, redis: redisClient}
}

type ListOrdersQuery struct {
	Status   string `form:"status"`
	Role     string `form:"role"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (h *OrderHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var query ListOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	db := h.db.Model(&models.Order{}).
		Preload("Task").Preload("Publisher").Preload("Courier").Preload("Tracks").Preload("ProofImages")

	if userRole == "courier" {
		db = db.Where("courier_id = ?", userID)
	} else {
		db = db.Where("publisher_id = ?", userID)
	}

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

func (h *OrderHandler) Get(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := h.db.Preload("Task").Preload("Task.Images").Preload("Publisher").
		Preload("Courier").Preload("Tracks").Preload("ProofImages").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Order not found"})
		return
	}

	if order.PublisherID != userID && order.CourierID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": order})
}

type TrackRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Address   string  `json:"address"`
	Message   string  `json:"message"`
	EventType string  `json:"event_type"`
}

func (h *OrderHandler) Track(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid order ID"})
		return
	}

	var req TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	var order models.Order
	if err := h.db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Order not found"})
		return
	}

	if order.CourierID != userID && order.PublisherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Access denied"})
		return
	}

	if req.EventType == "start" && userRole == "courier" {
		now := time.Now()
		order.Status = models.OrderStatusInProgress
		order.StartTime = &now
		h.db.Save(&order)

		var task models.Task
		h.db.First(&task, order.TaskID)
		task.Status = models.TaskStatusInProgress
		h.db.Save(&task)
	}

	track := models.OrderTrack{
		OrderID:   order.ID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Address:   req.Address,
		Message:   req.Message,
		EventType: req.EventType,
	}

	if err := h.db.Create(&track).Error; err != nil {
		logger.Errorf("Failed to create track: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save tracking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Tracking saved successfully",
		"data":    track,
	})
}
