package handler

import (
	"net/http"
	"strconv"
	"venue-booking/internal/dto"
	"venue-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
	logService    *service.OperationLogService
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		reviewService: service.NewReviewService(),
		logService:    service.NewOperationLogService(),
	}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var req dto.ReviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	userID, _ := c.Get("userID")
	review, err := h.reviewService.Create(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, userID.(uint), "create_review", "review", map[string]interface{}{
		"review_id": review.ID,
		"order_id":  review.OrderID,
		"rating":    review.Rating,
	})

	c.JSON(http.StatusOK, dto.Success(review))
}

func (h *ReviewHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid review ID"))
		return
	}

	review, err := h.reviewService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "Review not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(review))
}

func (h *ReviewHandler) List(c *gin.Context) {
	var req dto.ReviewListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	reviews, total, err := h.reviewService.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get reviews"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     reviews,
	}))
}

func (h *ReviewHandler) Approve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid review ID"))
		return
	}

	var req dto.ReviewActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	userID, _ := c.Get("userID")
	err = h.reviewService.Approve(uint(id), userID.(uint), req.Note)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, userID.(uint), "approve_review", "review", map[string]interface{}{
		"review_id": id,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *ReviewHandler) Reject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid review ID"))
		return
	}

	var req dto.ReviewActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	userID, _ := c.Get("userID")
	err = h.reviewService.Reject(uint(id), userID.(uint), req.Note)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	h.logService.Log(c, userID.(uint), "reject_review", "review", map[string]interface{}{
		"review_id": id,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}
