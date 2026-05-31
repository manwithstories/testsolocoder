package controllers

import (
	"health-platform/services"
	"health-platform/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type StatisticsController struct {
	statisticsService *services.StatisticsService
}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{
		statisticsService: services.NewStatisticsService(),
	}
}

func (ctrl *StatisticsController) GetCompanyStatistics(c *gin.Context) {
	companyID := c.GetUint("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	query := services.StatisticsQuery{
		CompanyID: companyID,
		Year:      year,
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" {
		query.StartDate = startDate
	}
	if endDate != "" {
		query.EndDate = endDate
	}

	result, err := ctrl.statisticsService.GetCompanyStatistics(query)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

func (ctrl *StatisticsController) GetDepartmentStatistics(c *gin.Context) {
	companyID := c.GetUint("company_id")
	departmentID, _ := strconv.ParseUint(c.Param("department_id"), 10, 64)
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	result, err := ctrl.statisticsService.GetDepartmentStatistics(companyID, uint(departmentID), year)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

func (ctrl *StatisticsController) GetAgeDistribution(c *gin.Context) {
	companyID := c.GetUint("company_id")

	result, err := ctrl.statisticsService.GetAgeDistribution(companyID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

func (ctrl *StatisticsController) GetGenderDistribution(c *gin.Context) {
	companyID := c.GetUint("company_id")

	result, err := ctrl.statisticsService.GetGenderDistribution(companyID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

func (ctrl *StatisticsController) GetAbnormalDistribution(c *gin.Context) {
	companyID := c.GetUint("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	result, err := ctrl.statisticsService.GetAbnormalDistribution(companyID, year)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

func (ctrl *StatisticsController) GetAgencyRating(c *gin.Context) {
	result, err := ctrl.statisticsService.GetAgencyRatingStatistics()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

func (ctrl *StatisticsController) GetPackageRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	packages, err := ctrl.statisticsService.GetPackageRanking(limit)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, packages)
}

func (ctrl *StatisticsController) ExportStatistics(c *gin.Context) {
	companyID := c.GetUint("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "统计报表"
	f.SetCellValue(sheetName, "A1", "企业健康管理统计报表")
	f.SetCellValue(sheetName, "A2", "年份")
	f.SetCellValue(sheetName, "B2", year)

	f.SetCellValue(sheetName, "A4", "员工总数")
	f.SetCellValue(sheetName, "A5", "体检完成率")
	f.SetCellValue(sheetName, "A6", "异常指标人数")
	f.SetCellValue(sheetName, "A7", "总费用")

	query := services.StatisticsQuery{
		CompanyID: companyID,
		Year:      year,
	}

	stats, err := ctrl.statisticsService.GetCompanyStatistics(query)
	if err == nil {
		if total, ok := stats["total_employees"].(int64); ok {
			f.SetCellValue(sheetName, "B4", total)
		}
		if rate, ok := stats["completion_rate"].(float64); ok {
			f.SetCellValue(sheetName, "B5", rate)
		}
		if abnormal, ok := stats["abnormal_count"].(int); ok {
			f.SetCellValue(sheetName, "B6", abnormal)
		}
		if total, ok := stats["total_billing"].(float64); ok {
			f.SetCellValue(sheetName, "B7", total)
		}
	}

	fileName := "statistics_report.xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	
	if err := f.Write(c.Writer); err != nil {
		utils.Error(c, 500, "导出失败")
		return
	}
}

func (ctrl *StatisticsController) ClearStatisticsCache(c *gin.Context) {
	if err := ctrl.statisticsService.ClearCache(); err != nil {
		utils.Error(c, 500, "清理缓存失败")
		return
	}

	utils.Success(c, nil)
}
