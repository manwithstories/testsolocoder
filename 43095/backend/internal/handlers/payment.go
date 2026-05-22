package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: services.NewPaymentService(),
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req services.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	payment, err := h.paymentService.CreatePayment(req)
	if err != nil {
		if err.Error() == "该预约已存在支付记录" || err.Error() == "预约不存在" || err.Error() == "支付金额必须大于0" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, payment)
}

func (h *PaymentHandler) GetPaymentDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的支付记录ID")
		return
	}

	payment, err := h.paymentService.GetPaymentByID(uint(id))
	if err != nil {
		if err.Error() == "支付记录不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	currentUser := utils.GetCurrentUser(c)
	if currentUser != nil && string(currentUser.Role) == "patient" {
		if payment.Appointment.Patient.UserID != currentUser.UserID {
			utils.Forbidden(c, "无权查看该支付记录")
			return
		}
	}

	utils.Success(c, payment)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的支付记录ID")
		return
	}

	var req services.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	payment, err := h.paymentService.UpdatePaymentStatus(uint(id), req)
	if err != nil {
		if err.Error() == "支付记录不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, payment)
}

func (h *PaymentHandler) GetPaymentList(c *gin.Context) {
	var query services.PaymentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	currentUser := utils.GetCurrentUser(c)
	if currentUser != nil && string(currentUser.Role) == "patient" {
		patientService := services.NewPatientService()
		patient, err := patientService.GetPatientByUserID(currentUser.UserID)
		if err != nil {
			utils.Forbidden(c, "获取患者信息失败")
			return
		}
		query.PatientID = patient.ID
	}

	payments, total, err := h.paymentService.GetPaymentList(query)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.SuccessWithPagination(c, payments, total, query.Page, query.PageSize)
}

func (h *PaymentHandler) GetFeeReport(c *gin.Context) {
	var query services.FeeReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	report, err := h.paymentService.GenerateFeeReport(query)
	if err != nil {
		if err.Error() == "无效的分组方式" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, report)
}

func (h *PaymentHandler) ExportFeeReportCSV(c *gin.Context) {
	var query services.FeeReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	filePath, err := h.paymentService.ExportFeeReportCSV(query)
	if err != nil {
		if err.Error() == "无效的分组方式" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.File(filePath)
}

func RegisterPaymentRoutes(api *gin.RouterGroup) {
	handler := NewPaymentHandler()

	payments := api.Group("/payments")
	payments.Use(middleware.Auth())
	{
		payments.GET("", handler.GetPaymentList)
		payments.GET("/:id", handler.GetPaymentDetail)

		admin := payments.Group("")
		admin.Use(middleware.AdminRequired())
		{
			admin.POST("", handler.CreatePayment)
			admin.PUT("/:id/status", handler.UpdatePaymentStatus)
		}
	}

	reports := api.Group("/reports")
	reports.Use(middleware.Auth())
	reports.Use(middleware.AdminRequired())
	{
		reports.GET("/fees", handler.GetFeeReport)
		reports.GET("/fees/export", handler.ExportFeeReportCSV)
	}
}
