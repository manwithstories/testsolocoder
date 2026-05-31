package handler

import (
	"net/http"
	"time"

	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct{}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	reviewerID, _ := uuid.Parse(userID.(string))

	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	rentalID, _ := uuid.Parse(req.RentalID)

	var rental model.Rental
	if err := database.DB.First(&rental, "id = ?", rentalID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Rental not found")
		return
	}

	if rental.TenantID != reviewerID {
		utils.Error(c, http.StatusForbidden, "You can only review your own rentals")
		return
	}

	if rental.Status != model.RentalStatusCompleted {
		utils.Error(c, http.StatusBadRequest, "Can only review completed rentals")
		return
	}

	var existingReview model.Review
	if err := database.DB.Where("rental_id = ? AND target_type = ?", rentalID, req.TargetType).First(&existingReview).Error; err == nil {
		utils.Error(c, http.StatusConflict, "You have already reviewed this rental")
		return
	}

	targetID, _ := uuid.Parse(req.TargetID)

	review := model.Review{
		RentalID:      rentalID,
		ReviewerID:    reviewerID,
		TargetType:    req.TargetType,
		TargetID:      targetID,
		Rating:        req.Rating,
		Content:       req.Content,
		IsRecommended: req.IsRecommended,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create review")
		return
	}

	h.updateAverageRating(req.TargetType, targetID)

	utils.Created(c, review)
}

func (h *ReviewHandler) GetReviews(c *gin.Context) {
	var req model.SearchReviewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	targetID, _ := uuid.Parse(req.TargetID)

	var reviews []model.Review
	query := database.DB.Preload("Reviewer").
		Where("target_type = ? AND target_id = ? AND is_deleted = ?", req.TargetType, targetID, false)

	if req.MinRating > 0 {
		query = query.Where("rating >= ?", req.MinRating)
	}
	if req.MaxRating > 0 {
		query = query.Where("rating <= ?", req.MaxRating)
	}

	var total int64
	query.Model(&model.Review{}).Count(&total)

	sortBy := req.SortBy
	validSorts := map[string]bool{"created_at": true, "rating": true, "helpful_count": true}
	if !validSorts[sortBy] {
		sortBy = "created_at"
	}

	sortOrder := req.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	orderClause := sortBy + " " + sortOrder

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order(orderClause).Offset(offset).Limit(req.PageSize).Find(&reviews).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}

	utils.Paginated(c, reviews, total, req.Page, req.PageSize)
}

func (h *ReviewHandler) GetReview(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid review ID")
		return
	}

	var review model.Review
	if err := database.DB.Preload("Reviewer").Preload("Responder").First(&review, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Review not found")
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) RespondToReview(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid review ID")
		return
	}

	userID, _ := c.Get("user_id")
	responderID, _ := uuid.Parse(userID.(string))

	var review model.Review
	if err := database.DB.First(&review, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Review not found")
		return
	}

	var req model.RespondToReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	review.Response = req.Response
	review.ResponseBy = &responderID
	review.ResponseAt = &now

	if err := database.DB.Save(&review).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to respond to review")
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) MarkHelpful(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid review ID")
		return
	}

	result := database.DB.Model(&model.Review{}).
		Where("id = ?", id).
		UpdateColumn("helpful_count", database.DB.Raw("helpful_count + 1"))

	if result.Error != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update helpful count")
		return
	}

	utils.Success(c, gin.H{"message": "Marked as helpful"})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid review ID")
		return
	}

	userID, _ := c.Get("user_id")
	reviewerID, _ := uuid.Parse(userID.(string))

	var review model.Review
	if err := database.DB.First(&review, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Review not found")
		return
	}

	if review.ReviewerID != reviewerID {
		role, _ := c.Get("role")
		if role.(string) != string(model.RoleAdmin) {
			utils.Error(c, http.StatusForbidden, "You can only delete your own reviews")
			return
		}
	}

	now := time.Now()
	review.IsDeleted = true
	review.DeletedAt = &now

	if err := database.DB.Save(&review).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete review")
		return
	}

	h.updateAverageRating(review.TargetType, review.TargetID)

	utils.Success(c, gin.H{"message": "Review deleted successfully"})
}

func (h *ReviewHandler) GetMyReviews(c *gin.Context) {
	userID, _ := c.Get("user_id")
	reviewerID, _ := uuid.Parse(userID.(string))

	var reviews []model.Review
	if err := database.DB.Preload("Reviewer").
		Where("reviewer_id = ? AND is_deleted = ?", reviewerID, false).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}

	utils.Success(c, reviews)
}

func (h *ReviewHandler) updateAverageRating(targetType model.ReviewType, targetID uuid.UUID) {
	var result struct {
		AvgRating  float64
		ReviewCount int
	}

	database.DB.Model(&model.Review{}).
		Where("target_type = ? AND target_id = ? AND is_deleted = ?", targetType, targetID, false).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as review_count").
		Scan(&result)

	switch targetType {
	case model.ReviewTypeShip:
		database.DB.Model(&model.Ship{}).
			Where("id = ?", targetID).
			Updates(map[string]interface{}{
				"average_rating": result.AvgRating,
				"review_count":   result.ReviewCount,
			})
	case model.ReviewTypeDock:
		database.DB.Model(&model.Dock{}).
			Where("id = ?", targetID).
			Updates(map[string]interface{}{
				"average_rating": result.AvgRating,
				"review_count":   result.ReviewCount,
			})
	}
}
