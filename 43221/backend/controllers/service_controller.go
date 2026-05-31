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

type ServiceController struct {
	serviceService *services.ServiceService
	cfg            *config.Config
}

func NewServiceController(cfg *config.Config) *ServiceController {
	return &ServiceController{
		serviceService: services.NewServiceService(),
		cfg:            cfg,
	}
}

func (ctrl *ServiceController) CreateService(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var req services.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	service, err := ctrl.serviceService.CreateService(professionalID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, service)
}

func (ctrl *ServiceController) GetServiceByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid service ID")
		return
	}

	service, err := ctrl.serviceService.GetServiceByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, 404, "Service not found")
		return
	}

	utils.SuccessResponse(c, service)
}

func (ctrl *ServiceController) UpdateService(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid service ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	service, err := ctrl.serviceService.UpdateService(id, professionalID, updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, service)
}

func (ctrl *ServiceController) DeleteService(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid service ID")
		return
	}

	err = ctrl.serviceService.DeleteService(id, professionalID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *ServiceController) GetProfessionalServices(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	services, total, err := ctrl.serviceService.GetServicesByProfessional(professionalID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    services,
	})
}

func (ctrl *ServiceController) GetAllServices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	serviceType := c.Query("service_type")

	services, total, err := ctrl.serviceService.GetAllServices(page, pageSize, serviceType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, utils.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    services,
	})
}

func (ctrl *ServiceController) CreateSchedule(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var req services.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	schedule, err := ctrl.serviceService.CreateSchedule(professionalID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, schedule)
}

func (ctrl *ServiceController) BatchCreateSchedules(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	var req services.BatchCreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	err = ctrl.serviceService.BatchCreateSchedules(professionalID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

func (ctrl *ServiceController) GetSchedules(c *gin.Context) {
	serviceID, err := uuid.Parse(c.Param("service_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid service ID")
		return
	}

	date := c.Query("date")
	onlyAvailable := c.DefaultQuery("only_available", "false") == "true"

	var schedules interface{}
	var err2 error

	if onlyAvailable {
		schedules, err2 = ctrl.serviceService.GetAvailableSchedules(serviceID, date)
	} else {
		schedules, err2 = ctrl.serviceService.GetSchedulesByServiceIDAndDate(serviceID, date)
	}

	if err2 != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err2.Error())
		return
	}

	utils.SuccessResponse(c, schedules)
}

func (ctrl *ServiceController) DeleteSchedules(c *gin.Context) {
	professionalID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	serviceID, err := uuid.Parse(c.Param("service_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Invalid service ID")
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	err = ctrl.serviceService.DeleteSchedulesByDateRange(professionalID, serviceID, startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, nil)
}

type ServiceControllerInterface interface {
	CreateService(c *gin.Context)
	GetServiceByID(c *gin.Context)
	UpdateService(c *gin.Context)
	DeleteService(c *gin.Context)
	GetProfessionalServices(c *gin.Context)
	GetAllServices(c *gin.Context)
	CreateSchedule(c *gin.Context)
	BatchCreateSchedules(c *gin.Context)
	GetSchedules(c *gin.Context)
	DeleteSchedules(c *gin.Context)
}
