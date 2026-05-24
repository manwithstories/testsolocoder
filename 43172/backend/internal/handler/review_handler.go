package handler

import (
	"strconv"

	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/service"
	resp "luxury-trading-platform/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	reviewerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	var req service.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	review, err := h.reviewService.CreateReview(c.Request.Context(), reviewerID.(uint), &req)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Created(c, review)
}

func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	review, err := h.reviewService.GetReview(c.Request.Context(), uint(id))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, review)
}

func (h *ReviewHandler) ListReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var revieweeID, reviewerID *uint
	if revieweeStr := c.Query("reviewee_id"); revieweeStr != "" {
		if id, err := strconv.ParseUint(revieweeStr, 10, 32); err == nil {
			uid := uint(id)
			revieweeID = &uid
		}
	}
	if reviewerStr := c.Query("reviewer_id"); reviewerStr != "" {
		if id, err := strconv.ParseUint(reviewerStr, 10, 32); err == nil {
			uid := uint(id)
			reviewerID = &uid
		}
	}

	var minRating *model.ReviewRating
	if ratingStr := c.Query("min_rating"); ratingStr != "" {
		if rating, err := strconv.Atoi(ratingStr); err == nil && rating >= 1 && rating <= 5 {
			r := model.ReviewRating(rating)
			minRating = &r
		}
	}

	reviews, total, err := h.reviewService.ListReviews(page, pageSize, revieweeID, reviewerID, minRating)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.SuccessWithPage(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) GetUserAverageRating(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	avgRating, err := h.reviewService.GetUserAverageRating(uint(id))
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, gin.H{"user_id": id, "average_rating": avgRating})
}
