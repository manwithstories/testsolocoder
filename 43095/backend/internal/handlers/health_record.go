package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HealthRecordHandler struct {
	service *services.HealthRecordService
}

func NewHealthRecordHandler() *HealthRecordHandler {
	return &HealthRecordHandler{
		service: services.NewHealthRecordService(),
	}
}

func RegisterHealthRecordRoutes(api *gin.RouterGroup) {
	handler := NewHealthRecordHandler()

	healthRecords := api.Group("/health-records")
	healthRecords.Use(middleware.Auth())
	{
		healthRecords.GET("/:patient_user_id", handler.Get)
		healthRecords.PUT("/:patient_user_id", middleware.PatientRequired(), handler.Update)
		healthRecords.GET("/:patient_user_id/visit-history", handler.GetVisitHistory)
		healthRecords.GET("/:patient_user_id/export", handler.Export)
	}
}

func (h *HealthRecordHandler) Get(c *gin.Context) {
	patientUserID, err := strconv.ParseUint(c.Param("patient_user_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的患者用户ID")
		return
	}

	user := utils.GetCurrentUser(c)
	healthRecord, err := h.service.GetHealthRecord(uint(patientUserID), user.UserID, user.Role)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, healthRecord)
}

func (h *HealthRecordHandler) Update(c *gin.Context) {
	patientUserID, err := strconv.ParseUint(c.Param("patient_user_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的患者用户ID")
		return
	}

	var req services.UpdateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user := utils.GetCurrentUser(c)
	healthRecord, err := h.service.UpdateHealthRecord(uint(patientUserID), user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, healthRecord)
}

func (h *HealthRecordHandler) GetVisitHistory(c *gin.Context) {
	patientUserID, err := strconv.ParseUint(c.Param("patient_user_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的患者用户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	user := utils.GetCurrentUser(c)
	history, total, err := h.service.GetMedicalVisitHistory(uint(patientUserID), user.UserID, user.Role, page, pageSize)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, history, total, page, pageSize)
}

func (h *HealthRecordHandler) Export(c *gin.Context) {
	patientUserID, err := strconv.ParseUint(c.Param("patient_user_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的患者用户ID")
		return
	}

	user := utils.GetCurrentUser(c)
	html, err := h.service.ExportHealthRecordHTML(uint(patientUserID), user.UserID, user.Role)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	format := c.DefaultQuery("format", "html")
	if format == "html" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=health_record.html")
		c.String(200, html)
	} else {
		utils.Success(c, gin.H{
			"html": html,
		})
	}
}
