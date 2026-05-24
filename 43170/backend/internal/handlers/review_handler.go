package handlers

import (
	"net/http"
	"photo-rental/internal/models"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct{}

type CreateReviewRequest struct {
	OrderID    uint   `json:"orderId" binding:"required"`
	ToUserID   uint   `json:"toUserId" binding:"required"`
	EquipmentID uint  `json:"equipmentId" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Content    string `json:"content"`
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetUint("userId")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, req.OrderID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.Status != "completed" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Can only review completed orders")
		return
	}

	if order.RenterID != userID && order.OwnerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only review orders you participated in")
		return
	}

	if req.ToUserID == userID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot review yourself")
		return
	}

	if order.RenterID != req.ToUserID && order.OwnerID != req.ToUserID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Target user is not part of this order")
		return
	}

	if req.EquipmentID != order.EquipmentID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Equipment does not match this order")
		return
	}

	var existingReview models.Review
	if result := database.DB.Where("order_id = ? AND from_user_id = ?", req.OrderID, userID).First(&existingReview); result.Error == nil {
		utils.ErrorResponse(c, http.StatusConflict, "You have already reviewed this order")
		return
	}

	content := req.Content
	if utils.ContainsSensitiveWords(content) {
		content = utils.FilterSensitiveWords(content)
	}

	review := models.Review{
		OrderID:     req.OrderID,
		FromUserID:  userID,
		ToUserID:    req.ToUserID,
		EquipmentID: req.EquipmentID,
		Rating:      req.Rating,
		Content:     content,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to create review: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create review")
		return
	}

	var avgRating float64
	var reviewCount int64

	tx.Model(&models.Review{}).
		Where("equipment_id = ?", req.EquipmentID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating)

	tx.Model(&models.Review{}).
		Where("equipment_id = ?", req.EquipmentID).
		Count(&reviewCount)

	if err := tx.Model(&models.Equipment{}).
		Where("id = ?", req.EquipmentID).
		Updates(map[string]interface{}{
			"rating":       avgRating,
			"review_count": reviewCount,
		}).Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to update equipment rating: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update rating")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Logger.Error("Failed to commit review: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create review")
		return
	}

	utils.Logger.Info("Review created: %d by user %d for equipment %d", review.ID, userID, req.EquipmentID)
	utils.SuccessResponse(c, review)
}

func (h *ReviewHandler) GetEquipmentReviews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var reviews []models.Review
	var total int64

	db := database.DB.Model(&models.Review{}).Where("equipment_id = ?", id)
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Preload("FromUser").Preload("ToUser").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}

	for i := range reviews {
		if reviews[i].FromUser.Password != "" {
			reviews[i].FromUser.Password = ""
		}
		if reviews[i].ToUser.Password != "" {
			reviews[i].ToUser.Password = ""
		}
	}

	utils.PaginatedSuccessResponse(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var reviews []models.Review
	if err := database.DB.Preload("Equipment").Preload("FromUser").
		Where("to_user_id = ?", id).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}

	for i := range reviews {
		if reviews[i].FromUser.Password != "" {
			reviews[i].FromUser.Password = ""
		}
	}

	utils.SuccessResponse(c, reviews)
}

func (h *ReviewHandler) GetMyReviews(c *gin.Context) {
	userID := c.GetUint("userId")

	var reviews []models.Review
	if err := database.DB.Preload("Equipment").Preload("ToUser").
		Where("from_user_id = ?", userID).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}

	for i := range reviews {
		if reviews[i].ToUser.Password != "" {
			reviews[i].ToUser.Password = ""
		}
	}

	utils.SuccessResponse(c, reviews)
}
