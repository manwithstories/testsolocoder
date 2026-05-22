package handler

import (
	"strconv"

	"beauty-salon-system/internal/service"
	"beauty-salon-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentService    *service.AppointmentService
	notificationService   *service.NotificationService
	auditService          *service.AuditService
}

func NewAppointmentHandler(
	appointmentService *service.AppointmentService,
	notificationService *service.NotificationService,
	auditService *service.AuditService,
) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService:  appointmentService,
		notificationService: notificationService,
		auditService:        auditService,
	}
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var req service.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.appointmentService.Create(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "create", "appointment", "Create appointment", c.ClientIP())

	h.notificationService.SendAppointmentNotification(result)

	response.Success(c, result)
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.appointmentService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "appointment not found")
		return
	}

	response.Success(c, result)
}

func (h *AppointmentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	filters := make(map[string]interface{})
	if v := c.Query("customer_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filters["customer_id"] = uint(id)
		}
	}
	if v := c.Query("technician_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filters["technician_id"] = uint(id)
		}
	}
	if v := c.Query("service_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filters["service_id"] = uint(id)
		}
	}
	if v := c.Query("status"); v != "" {
		filters["status"] = v
	}
	if v := c.Query("start_date"); v != "" {
		filters["start_date"] = v
	}
	if v := c.Query("end_date"); v != "" {
		filters["end_date"] = v
	}

	result, total, err := h.appointmentService.List(page, pageSize, filters)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  result,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *AppointmentHandler) GetByTechnicianAndDate(c *gin.Context) {
	techID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	date := c.Query("date")

	result, err := h.appointmentService.GetByTechnicianAndDate(uint(techID), date)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AppointmentHandler) GetMyAppointments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, total, err := h.appointmentService.GetByCustomer(userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  result,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
	var req service.CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.appointmentService.Cancel(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	appointment, _ := h.appointmentService.GetByID(req.ID)
	if appointment != nil {
		h.notificationService.SendCancelNotification(appointment)
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "cancel", "appointment", "Cancel appointment", c.ClientIP())

	response.Success(c, nil)
}

func (h *AppointmentHandler) Reschedule(c *gin.Context) {
	var req service.RescheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.appointmentService.Reschedule(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "reschedule", "appointment", "Reschedule appointment", c.ClientIP())

	response.Success(c, result)
}

func (h *AppointmentHandler) GetAvailableSlots(c *gin.Context) {
	techID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	date := c.Query("date")
	duration, _ := strconv.Atoi(c.DefaultQuery("duration", "60"))

	result, err := h.appointmentService.GetAvailableSlots(uint(techID), date, duration)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AppointmentHandler) Complete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.appointmentService.Complete(uint(id)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "complete", "appointment", "Complete appointment", c.ClientIP())

	response.Success(c, nil)
}
