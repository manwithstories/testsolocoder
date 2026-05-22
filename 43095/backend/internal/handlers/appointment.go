package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	service *services.AppointmentService
}

func NewAppointmentHandler() *AppointmentHandler {
	return &AppointmentHandler{
		service: services.NewAppointmentService(),
	}
}

func RegisterAppointmentRoutes(api *gin.RouterGroup) {
	handler := NewAppointmentHandler()

	appointments := api.Group("/appointments")
	appointments.Use(middleware.Auth())
	{
		appointments.POST("", middleware.PatientRequired(), handler.Create)
		appointments.GET("", handler.List)
		appointments.GET("/:id", handler.Get)
		appointments.PUT("/:id/cancel", middleware.PatientRequired(), handler.Cancel)
		appointments.PUT("/:id/reschedule", middleware.PatientRequired(), handler.Reschedule)
		appointments.PUT("/:id/confirm", middleware.DoctorRequired(), handler.Confirm)
		appointments.PUT("/:id/complete", middleware.DoctorRequired(), handler.Complete)
		appointments.GET("/check-availability", handler.CheckAvailability)
	}
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var req services.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user := utils.GetCurrentUser(c)
	appointment, err := h.service.CreateAppointment(user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func (h *AppointmentHandler) List(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	appointments, total, err := h.service.GetAppointmentList(user.UserID, user.Role, page, pageSize, status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, appointments, total, page, pageSize)
}

func (h *AppointmentHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的预约ID")
		return
	}

	user := utils.GetCurrentUser(c)
	appointment, err := h.service.GetAppointmentDetail(uint(id), user.UserID, user.Role)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的预约ID")
		return
	}

	var req services.CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user := utils.GetCurrentUser(c)
	if err := h.service.CancelAppointment(uint(id), user.UserID, user.Role, req.CancelReason); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *AppointmentHandler) Reschedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的预约ID")
		return
	}

	var req services.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user := utils.GetCurrentUser(c)
	appointment, err := h.service.RescheduleAppointment(uint(id), user.UserID, user.Role, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func (h *AppointmentHandler) Confirm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的预约ID")
		return
	}

	user := utils.GetCurrentUser(c)
	if err := h.service.ConfirmAppointment(uint(id), user.UserID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *AppointmentHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的预约ID")
		return
	}

	user := utils.GetCurrentUser(c)
	if err := h.service.CompleteAppointment(uint(id), user.UserID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *AppointmentHandler) CheckAvailability(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Query("doctor_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的医生ID")
		return
	}

	dateStr := c.Query("appointment_date")
	if dateStr == "" {
		utils.BadRequest(c, "请提供预约日期")
		return
	}

	appointmentDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.BadRequest(c, "日期格式错误，请使用YYYY-MM-DD格式")
		return
	}

	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	if startTime == "" || endTime == "" {
		utils.BadRequest(c, "请提供开始时间和结束时间")
		return
	}

	availability, err := h.service.CheckTimeSlotAvailability(uint(doctorID), appointmentDate, startTime, endTime)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, availability)
}
