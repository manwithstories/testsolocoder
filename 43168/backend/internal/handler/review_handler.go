package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

// ReviewHandler 售后评价 HTTP 处理器
type ReviewHandler struct {
	service *service.ReviewService
}

// NewReviewHandler 创建售后评价处理器
func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: svc}
}

// toReviewResponse 将模型转为响应 DTO
func toReviewResponse(r *model.Review) dto.ReviewResponse {
	images, _ := r.ParseImages()
	if images == nil {
		images = []string{}
	}
	return dto.ReviewResponse{
		ID:            r.ID,
		OrderID:       r.OrderID,
		ProductID:     r.ProductID,
		OwnerID:       r.OwnerID,
		ProductRating: r.ProductRating,
		ServiceRating: r.ServiceRating,
		Content:       r.Content,
		Images:        images,
		CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// parseReviewID 解析评价 ID 路径参数
func parseReviewID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// CreateReview 业主创建评价
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	review, err := h.service.CreateReview(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toReviewResponse(review))
}

// GetReview 获取评价详情
func (h *ReviewHandler) GetReview(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseReviewID(c)
	if err != nil {
		response.BadRequest(c, "评价 ID 格式错误")
		return
	}

	review, err := h.service.GetByID(id, userID, role)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toReviewResponse(review))
}

// ListReviews 分页查询评价列表
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.ReviewListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}
	req.UserID = userID
	req.Role = role

	list, total, err := h.service.List(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	resp := make([]dto.ReviewResponse, 0, len(list))
	for _, r := range list {
		resp = append(resp, toReviewResponse(r))
	}
	response.Success(c, dto.ReviewListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     resp,
	})
}
