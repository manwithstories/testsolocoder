package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService  *services.ReviewService
	patientService *services.PatientService
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		reviewService:  services.NewReviewService(),
		patientService: services.NewPatientService(),
	}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	currentUser := utils.GetCurrentUser(c)
	if currentUser == nil {
		utils.Unauthorized(c, "请先登录")
		return
	}

	if string(currentUser.Role) != "patient" {
		utils.Forbidden(c, "只有患者可以创建评价")
		return
	}

	patient, err := h.patientService.GetPatientByUserID(currentUser.UserID)
	if err != nil {
		utils.Forbidden(c, "获取患者信息失败")
		return
	}

	var req services.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	review, err := h.reviewService.CreateReview(patient.ID, req)
	if err != nil {
		if err.Error() == "该预约已存在评价" ||
			err.Error() == "预约不存在" ||
			err.Error() == "只能对已完成的预约进行评价" ||
			err.Error() == "只能评价自己的预约" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) GetReviewList(c *gin.Context) {
	var query services.ReviewListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	reviews, total, err := h.reviewService.GetReviewList(query)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.SuccessWithPagination(c, reviews, total, query.Page, query.PageSize)
}

func (h *ReviewHandler) GetReviewDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的评价ID")
		return
	}

	review, err := h.reviewService.GetReviewByID(uint(id))
	if err != nil {
		if err.Error() == "评价不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的评价ID")
		return
	}

	if err := h.reviewService.DeleteReview(uint(id)); err != nil {
		if err.Error() == "评价不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, nil)
}

func RegisterReviewRoutes(api *gin.RouterGroup) {
	handler := NewReviewHandler()

	reviews := api.Group("/reviews")
	{
		reviews.GET("", handler.GetReviewList)
		reviews.GET("/:id", handler.GetReviewDetail)

		patient := reviews.Group("")
		patient.Use(middleware.Auth())
		patient.Use(middleware.PatientRequired())
		{
			patient.POST("", handler.CreateReview)
		}

		admin := reviews.Group("")
		admin.Use(middleware.Auth())
		admin.Use(middleware.AdminRequired())
		{
			admin.DELETE("/:id", handler.DeleteReview)
		}
	}
}
