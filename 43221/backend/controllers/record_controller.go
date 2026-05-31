package controllers

import (
	"net/http"
	"strconv"

	"consultation-platform/services"
	"consultation-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecordController struct {
	recordService *services.RecordService
}

func NewRecordController(sensitiveWords []string) *RecordController {
	return &RecordController{
		recordService: services.NewRecordService(sensitiveWords),
	}
}

func (ctrl *RecordController) CreateConsultRecord(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var req services.CreateConsultRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	record, err := ctrl.recordService.CreateConsultRecord(professionalID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, record)
}

func (ctrl *RecordController) GetConsultRecordByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid record ID")
		return
	}

	record, err := ctrl.recordService.GetConsultRecordByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, 404, "Record not found")
		return
	}

	utils.SuccessResponse(c, record)
}

func (ctrl *RecordController) GetClientConsultRecords(c *gin.Context) {
	clientID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	records, total, err := ctrl.recordService.GetClientConsultRecords(clientID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    records,
	})
}

func (ctrl *RecordController) GetProfessionalConsultRecords(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	records, total, err := ctrl.recordService.GetProfessionalConsultRecords(professionalID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    records,
	})
}

func (ctrl *RecordController) CreateReview(c *gin.Context) {
	clientID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var req services.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	review, err := ctrl.recordService.CreateReview(clientID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, review)
}

func (ctrl *RecordController) GetProfessionalReviews(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	reviews, total, err := ctrl.recordService.GetProfessionalReviews(professionalID, page, pageSize, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    reviews,
	})
}

func (ctrl *RecordController) GetServiceReviews(c *gin.Context) {
	serviceID, err := uuid.Parse(c.Param("service_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid service ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reviews, total, err := ctrl.recordService.GetServiceReviews(serviceID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    reviews,
	})
}

func (ctrl *RecordController) GetPendingReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reviews, total, err := ctrl.recordService.GetPendingReviews(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    reviews,
	})
}

func (ctrl *RecordController) UpdateReviewStatus(c *gin.Context) {
	var req services.UpdateReviewStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	err := ctrl.recordService.UpdateReviewStatus(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *RecordController) GetProfessionalReviewStats(c *gin.Context) {
	professionalID, err := uuid.Parse(c.Param("professional_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid professional ID")
		return
	}

	avgRating, count, err := ctrl.recordService.GetProfessionalReviewStats(professionalID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"average_rating": avgRating,
		"total_reviews":  count,
	})
}

type RecordControllerInterface interface {
	CreateConsultRecord(c *gin.Context)
	GetConsultRecordByID(c *gin.Context)
	GetClientConsultRecords(c *gin.Context)
	GetProfessionalConsultRecords(c *gin.Context)
	CreateReview(c *gin.Context)
	GetProfessionalReviews(c *gin.Context)
	GetServiceReviews(c *gin.Context)
	GetPendingReviews(c *gin.Context)
	UpdateReviewStatus(c *gin.Context)
	GetProfessionalReviewStats(c *gin.Context)
}
