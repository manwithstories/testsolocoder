package controllers

import (
	"net/http"
	"strconv"

	"consultation-platform/config"
	"consultation-platform/services"
	"consultation-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentController struct {
	appointmentService *services.AppointmentService
	notificationService *services.NotificationService
	cfg                *config.Config
}

func NewAppointmentController(cfg *config.Config) *AppointmentController {
	return &AppointmentController{
		appointmentService:  services.NewAppointmentService(cfg),
		notificationService: services.NewNotificationService(cfg),
		cfg:                 cfg,
	}
}

func (ctrl *AppointmentController) CreateAppointment(c *gin.Context) {
	clientID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var req services.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	appointment, payment, err := ctrl.appointmentService.CreateAppointment(clientID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	ctrl.notificationService.NotifyAppointmentSuccess(clientID, appointment)
	ctrl.notificationService.NotifyAppointmentSuccess(appointment.ProfessionalID, appointment)

	utils.SuccessResponse(c, gin.H{
		"appointment": appointment,
		"payment":     payment,
	})
}

func (ctrl *AppointmentController) GetAppointmentByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid appointment ID")
		return
	}

	appointment, err := ctrl.appointmentService.GetAppointmentByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, 404, "Appointment not found")
		return
	}

	utils.SuccessResponse(c, appointment)
}

func (ctrl *AppointmentController) GetClientAppointments(c *gin.Context) {
	clientID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	appointments, total, err := ctrl.appointmentService.GetClientAppointments(clientID, page, pageSize, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    appointments,
	})
}

func (ctrl *AppointmentController) GetProfessionalAppointments(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	appointments, total, err := ctrl.appointmentService.GetProfessionalAppointments(professionalID, page, pageSize, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    appointments,
	})
}

func (ctrl *AppointmentController) ConfirmAppointment(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var req services.ConfirmAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	appointment, err := ctrl.appointmentService.ConfirmAppointment(professionalID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, appointment)
}

func (ctrl *AppointmentController) CancelAppointment(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	userRole, _ := utils.GetUserRoleFromContext(c)

	var req services.CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	appointment, err := ctrl.appointmentService.CancelAppointment(userID, userRole, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	ctrl.notificationService.NotifyAppointmentCancel(userID, appointment, req.Reason)
	ctrl.notificationService.NotifyAppointmentCancel(appointment.ProfessionalID, appointment, req.Reason)

	utils.SuccessResponse(c, appointment)
}

func (ctrl *AppointmentController) CompleteAppointment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid appointment ID")
		return
	}

	appointment, err := ctrl.appointmentService.CompleteAppointment(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, appointment)
}

func (ctrl *AppointmentController) ProcessPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid appointment ID")
		return
	}

	var req struct {
		TransactionID string `json:"transaction_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	err = ctrl.appointmentService.ProcessPayment(id, req.TransactionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *AppointmentController) RefundPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid appointment ID")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	err = ctrl.appointmentService.RefundPayment(id, req.Reason)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

type AppointmentControllerInterface interface {
	CreateAppointment(c *gin.Context)
	GetAppointmentByID(c *gin.Context)
	GetClientAppointments(c *gin.Context)
	GetProfessionalAppointments(c *gin.Context)
	ConfirmAppointment(c *gin.Context)
	CancelAppointment(c *gin.Context)
	CompleteAppointment(c *gin.Context)
	ProcessPayment(c *gin.Context)
	RefundPayment(c *gin.Context)
}
