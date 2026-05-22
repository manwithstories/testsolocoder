package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"auction-system/internal/dto"
	"auction-system/internal/middleware"
	"auction-system/internal/services"
	"auction-system/pkg/response"
)

type ReviewController struct {
	reviewService *services.ReviewService
}

func NewReviewController() *ReviewController {
	return &ReviewController{
		reviewService: services.NewReviewService(),
	}
}

func (ctrl *ReviewController) Create(c *gin.Context) {
	reviewerID := middleware.GetCurrentUserID(c)

	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	review, err := ctrl.reviewService.CreateReview(reviewerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, review)
}

func (ctrl *ReviewController) GetOrderReviews(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("order_id"))

	reviews, err := ctrl.reviewService.GetOrderReviews(uint(orderID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, reviews)
}

func (ctrl *ReviewController) GetUserReviews(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reviews, total, err := ctrl.reviewService.GetUserReviews(uint(userID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      reviews,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *ReviewController) GetUserRating(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))

	avgRating, total, err := ctrl.reviewService.GetUserAverageRating(uint(userID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"avg_rating": avgRating,
		"total":      total,
	})
}

func (ctrl *ReviewController) GetByID(c *gin.Context) {
	reviewID, _ := strconv.Atoi(c.Param("id"))

	review, err := ctrl.reviewService.GetReviewByID(uint(reviewID))
	if err != nil {
		response.NotFound(c, "评价不存在")
		return
	}

	response.Success(c, review)
}
