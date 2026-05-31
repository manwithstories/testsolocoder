package controllers

import (
	"health-platform/services"
	"health-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppointmentController struct {
	appointmentService *services.AppointmentService
}

func NewAppointmentController() *AppointmentController {
	return &AppointmentController{
		appointmentService: services.NewAppointmentService(),
	}
}

func (ctrl *AppointmentController) CreateAppointment(c *gin.Context) {
	var req services.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	appointment, err := ctrl.appointmentService.CreateAppointment(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func (ctrl *AppointmentController) GetAppointment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	appointment, err := ctrl.appointmentService.GetAppointment(uint(id))
	if err != nil {
		utils.Error(c, 404, "预约不存在")
		return
	}

	utils.Success(c, appointment)
}

func (ctrl *AppointmentController) GetEmployeeAppointments(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	appointments, total, err := ctrl.appointmentService.GetEmployeeAppointments(uint(employeeID), page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, appointments)
}

func (ctrl *AppointmentController) GetCompanyAppointments(c *gin.Context) {
	companyID := c.GetUint("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	appointments, total, err := ctrl.appointmentService.GetCompanyAppointments(companyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, appointments)
}

func (ctrl *AppointmentController) GetAgencyAppointments(c *gin.Context) {
	agencyID := c.GetUint("agency_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	appointments, total, err := ctrl.appointmentService.GetAgencyAppointments(agencyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, appointments)
}

func (ctrl *AppointmentController) RescheduleAppointment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	var req services.RescheduleAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.appointmentService.RescheduleAppointment(uint(id), &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AppointmentController) CancelAppointment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	var req services.CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.appointmentService.CancelAppointment(uint(id), &req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AppointmentController) CompleteAppointment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	if err := ctrl.appointmentService.CompleteAppointment(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AppointmentController) GetEmployeeAppointmentStatus(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	
	status, err := ctrl.appointmentService.GetEmployeeAppointmentStatus(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, status)
}

func (ctrl *AppointmentController) CheckQuota(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	packageID, _ := strconv.ParseUint(c.Query("package_id"), 10, 64)
	
	canBook, err := ctrl.appointmentService.CheckQuota(uint(employeeID), uint(packageID))
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, gin.H{"can_book": canBook})
}
