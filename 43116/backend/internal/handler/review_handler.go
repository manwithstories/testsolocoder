package handler

import (
	"car-rental/internal/middleware"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		reviewService: service.NewReviewService(),
	}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	review, err := h.reviewService.CreateReview(user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) GetReviewByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	review, err := h.reviewService.GetReviewByID(uint(id))
	if err != nil {
		utils.NotFound(c, "评价不存在")
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) GetCarReviews(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Param("carId"), 10, 64)
	page, pageSize, _, _ := utils.ParsePageParams(c)

	reviews, total, err := h.reviewService.GetCarReviews(uint(carID), page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	page, pageSize, _, _ := utils.ParsePageParams(c)

	reviews, total, err := h.reviewService.GetUserReviews(user.UserID, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) GetAllReviews(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	carID, _ := strconv.ParseUint(c.Query("car_id"), 10, 64)
	minRating, _ := strconv.Atoi(c.Query("min_rating"))

	reviews, total, err := h.reviewService.GetAllReviews(page, pageSize, uint(carID), minRating)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.reviewService.UpdateReview(uint(id), user.UserID, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.reviewService.DeleteReview(uint(id), user.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ReviewHandler) ToggleReviewHidden(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		IsHidden bool `json:"is_hidden"`
	}
	c.ShouldBindJSON(&req)

	err := h.reviewService.ToggleReviewHidden(uint(id), req.IsHidden)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ReviewHandler) LikeReview(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.reviewService.LikeReview(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
