package controllers

import (
	"net/http"
	"time"

	"consultation-platform/services"
	"consultation-platform/utils"

	"github.com/gin-gonic/gin"
)

type StatisticsController struct {
	statisticsService *services.StatisticsService
}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{
		statisticsService: services.NewStatisticsService(),
	}
}

func (ctrl *StatisticsController) GetProfessionalStats(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	stats, err := ctrl.statisticsService.GetProfessionalStats(userID, startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}

func (ctrl *StatisticsController) GetAdminStats(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	stats, err := ctrl.statisticsService.GetAdminStats(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}

func (ctrl *StatisticsController) ExportAppointments(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	userRole, _ := utils.GetUserRoleFromContext(c)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	data, err := ctrl.statisticsService.ExportAppointments(userID, userRole, startDate, endDate, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=appointments.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (ctrl *StatisticsController) ExportRevenue(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	userRole, _ := utils.GetUserRoleFromContext(c)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := ctrl.statisticsService.ExportRevenue(userID, userRole, startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=revenue.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

type StatisticsControllerInterface interface {
	GetProfessionalStats(c *gin.Context)
	GetAdminStats(c *gin.Context)
	ExportAppointments(c *gin.Context)
	ExportRevenue(c *gin.Context)
}
