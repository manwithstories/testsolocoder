package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type StatisticsController struct {
	statisticsService *services.StatisticsService
}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{
		statisticsService: services.NewStatisticsService(),
	}
}

func (ctrl *StatisticsController) GetOverviewStats(c *gin.Context) {
	var startDate, endDate *time.Time

	if v := c.Query("startDate"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			startDate = &t
		}
	}
	if v := c.Query("endDate"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			endDate = &t
		}
	}

	stats, err := ctrl.statisticsService.GetOverviewStats(startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (ctrl *StatisticsController) GetStatusDistribution(c *gin.Context) {
	distributions, err := ctrl.statisticsService.GetStatusDistribution()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, distributions)
}

func (ctrl *StatisticsController) GetCompanyTypeDistribution(c *gin.Context) {
	distributions, err := ctrl.statisticsService.GetCompanyTypeDistribution()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, distributions)
}

func (ctrl *StatisticsController) GetAgentPerformance(c *gin.Context) {
	var startDate, endDate *time.Time

	if v := c.Query("startDate"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			startDate = &t
		}
	}
	if v := c.Query("endDate"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			endDate = &t
		}
	}

	performances, err := ctrl.statisticsService.GetAgentPerformance(startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, performances)
}

func (ctrl *StatisticsController) GetApplicationTimeSeries(c *gin.Context) {
	startDate, err := time.Parse("2006-01-02", c.DefaultQuery("startDate", time.Now().AddDate(0, -1, 0).Format("2006-01-02")))
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	endDate, err := time.Parse("2006-01-02", c.DefaultQuery("endDate", time.Now().Format("2006-01-02")))
	if err != nil {
		endDate = time.Now()
	}

	interval := c.DefaultQuery("interval", "day")

	data, err := ctrl.statisticsService.GetApplicationTimeSeries(startDate, endDate, interval)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, data)
}

func (ctrl *StatisticsController) GetRevenueStats(c *gin.Context) {
	startDate, err := time.Parse("2006-01-02", c.DefaultQuery("startDate", time.Now().AddDate(0, -1, 0).Format("2006-01-02")))
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	endDate, err := time.Parse("2006-01-02", c.DefaultQuery("endDate", time.Now().Format("2006-01-02")))
	if err != nil {
		endDate = time.Now()
	}

	interval := c.DefaultQuery("interval", "day")

	data, err := ctrl.statisticsService.GetApplicationTimeSeries(startDate, endDate, interval)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, data)
}
