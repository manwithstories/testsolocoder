package handlers

import (
	"errand-service/internal/models"
	"errand-service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	db *gorm.DB
}

func NewReviewHandler(db *gorm.DB) *ReviewHandler {
	return &ReviewHandler{db: db}
}

type CreateReviewRequest struct {
	OrderID    uint   `json:"order_id" binding:"required"`
	ReviewType string `json:"review_type" binding:"required,oneof=courier publisher"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Content    string `json:"content"`
	Tags       string `json:"tags"`
}

type ListReviewsQuery struct {
	UserID   uint   `form:"user_id"`
	Type     string `form:"type"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (h *ReviewHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	var order models.Order
	if err := h.db.Preload("Task").First(&order, req.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Order not found"})
		return
	}

	if order.Status != models.OrderStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Can only review completed orders"})
		return
	}

	var existing models.Review
	if err := h.db.Where("order_id = ? AND reviewer_id = ?", req.OrderID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "You have already reviewed this order"})
		return
	}

	var revieweeID uint
	if req.ReviewType == "courier" {
		if order.PublisherID != userID {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only the publisher can review the courier"})
			return
		}
		revieweeID = order.CourierID
	} else {
		if order.CourierID != userID {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only the courier can review the publisher"})
			return
		}
		revieweeID = order.PublisherID
	}

	tx := h.db.Begin()

	review := models.Review{
		OrderID:    req.OrderID,
		ReviewerID: userID,
		RevieweeID: revieweeID,
		ReviewType: models.ReviewType(req.ReviewType),
		Rating:     req.Rating,
		Content:    req.Content,
		Tags:       req.Tags,
	}

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to create review: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create review"})
		return
	}

	var avgRating struct {
		AvgRating float64 `gorm:"column:avg_rating"`
	}
	tx.Model(&models.Review{}).
		Where("reviewee_id = ?", revieweeID).
		Select("COALESCE(AVG(rating), 5) as avg_rating").
		Scan(&avgRating)

	newRating := avgRating.AvgRating
	if newRating == 0 {
		newRating = 5
	}

	if err := tx.Model(&models.User{}).Where("id = ?", revieweeID).Update("rating", newRating).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to update user rating: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update rating"})
		return
	}

	if req.ReviewType == "courier" {
		if err := tx.Model(&models.CourierProfile{}).Where("user_id = ?", revieweeID).Update("rating", newRating).Error; err != nil {
			tx.Rollback()
			logger.Errorf("Failed to update courier rating: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update courier rating"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Review submitted successfully",
		"data":    review,
	})
}

func (h *ReviewHandler) List(c *gin.Context) {
	var query ListReviewsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	db := h.db.Model(&models.Review{}).Preload("Reviewer").Preload("Order")

	if query.UserID > 0 {
		db = db.Where("reviewee_id = ?", query.UserID)
	}
	if query.Type != "" {
		db = db.Where("review_type = ?", query.Type)
	}

	var total int64
	db.Count(&total)

	var reviews []models.Review
	offset := (query.Page - 1) * query.PageSize
	db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&reviews)

	var avgRating float64
	var reviewCount int64
	if query.UserID > 0 {
		h.db.Model(&models.Review{}).Where("reviewee_id = ?", query.UserID).
			Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
		h.db.Model(&models.Review{}).Where("reviewee_id = ?", query.UserID).Count(&reviewCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":        total,
			"page":         query.Page,
			"page_size":    query.PageSize,
			"items":        reviews,
			"avg_rating":   avgRating,
			"review_count": reviewCount,
		},
	})
}
