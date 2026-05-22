package handlers

import (
	"medical-platform/internal/config"
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConsultationHandler struct {
	service *services.ConsultationService
}

func NewConsultationHandler() *ConsultationHandler {
	return &ConsultationHandler{
		service: services.NewConsultationService(),
	}
}

func RegisterConsultationRoutes(api *gin.RouterGroup) {
	handler := NewConsultationHandler()

	consultations := api.Group("/consultations")
	consultations.Use(middleware.Auth())
	{
		consultations.POST("", middleware.DoctorRequired(), handler.Create)
		consultations.GET("/:id", handler.Get)
		consultations.PUT("/:id", middleware.DoctorRequired(), handler.Update)
		consultations.POST("/prescriptions", middleware.DoctorRequired(), handler.CreatePrescription)
		consultations.POST("/reports", middleware.DoctorRequired(), handler.UploadReport)
		consultations.GET("/patient/history", middleware.PatientRequired(), handler.GetPatientHistory)
	}
}

func (h *ConsultationHandler) Create(c *gin.Context) {
	var req services.CreateConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user := utils.GetCurrentUser(c)
	consultation, err := h.service.CreateConsultation(user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, consultation)
}

func (h *ConsultationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的问诊记录ID")
		return
	}

	user := utils.GetCurrentUser(c)
	consultation, err := h.service.GetConsultationDetail(uint(id), user.UserID, user.Role)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, consultation)
}

func (h *ConsultationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的问诊记录ID")
		return
	}

	var req services.UpdateConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user := utils.GetCurrentUser(c)
	consultation, err := h.service.UpdateConsultation(uint(id), user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, consultation)
}

func (h *ConsultationHandler) CreatePrescription(c *gin.Context) {
	var req services.CreatePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user := utils.GetCurrentUser(c)
	prescription, err := h.service.CreatePrescription(user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, prescription)
}

func (h *ConsultationHandler) UploadReport(c *gin.Context) {
	consultationID, err := strconv.ParseUint(c.PostForm("consultation_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的问诊记录ID")
		return
	}

	reportType := c.PostForm("report_type")
	if reportType == "" {
		utils.BadRequest(c, "请提供报告类型")
		return
	}

	reportName := c.PostForm("report_name")
	if reportName == "" {
		utils.BadRequest(c, "请提供报告名称")
		return
	}

	findings := c.PostForm("findings")
	conclusion := c.PostForm("conclusion")

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}

	maxSize := config.AppConfig.MaxUploadSize
	if file.Size > maxSize {
		utils.BadRequest(c, "文件大小超过限制")
		return
	}

	req := &services.UploadReportRequest{
		ConsultationID: uint(consultationID),
		ReportType:     reportType,
		ReportName:     reportName,
		Findings:       findings,
		Conclusion:     conclusion,
	}

	user := utils.GetCurrentUser(c)
	report, err := h.service.UploadReport(user.UserID, req, file)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, report)
}

func (h *ConsultationHandler) GetPatientHistory(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	consultations, total, err := h.service.GetPatientConsultationHistory(user.UserID, page, pageSize)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, consultations, total, page, pageSize)
}
