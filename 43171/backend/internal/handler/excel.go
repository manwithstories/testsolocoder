package handler

import (
	"drone-rental/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ExcelHandler struct {
	statsService *service.StatsService
}

func NewExcelHandler() *ExcelHandler {
	return &ExcelHandler{
		statsService: service.NewStatsService(),
	}
}

func (h *ExcelHandler) ExportRevenue(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	stats, err := h.statsService.GetRevenueStats(startDate, endDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "日期")
	f.SetCellValue(sheet, "B1", "收入金额")
	f.SetCellValue(sheet, "C1", "订单数量")
	for i, s := range stats {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), s.Date)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), s.Amount)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), s.Count)
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=revenue_stats.xlsx")
	f.Write(c.Writer)
}

func (h *ExcelHandler) ExportRegion(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	stats, err := h.statsService.GetRegionStats(startDate, endDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "区域")
	f.SetCellValue(sheet, "B1", "订单数量")
	f.SetCellValue(sheet, "C1", "收入金额")
	for i, s := range stats {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), s.Region)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), s.Count)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), s.Amount)
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=region_stats.xlsx")
	f.Write(c.Writer)
}

func (h *ExcelHandler) ExportDrone(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	stats, err := h.statsService.GetDroneStats(startDate, endDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "设备名称")
	f.SetCellValue(sheet, "B1", "租赁天数")
	f.SetCellValue(sheet, "C1", "利用率")
	f.SetCellValue(sheet, "D1", "收入金额")
	for i, s := range stats {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), s.DroneName)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), s.TotalDays)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), s.Utilization)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), s.Income)
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=drone_stats.xlsx")
	f.Write(c.Writer)
}
