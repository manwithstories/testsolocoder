package handler

import (
	"hotel-system/internal/dto"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GetOccupancyRate(c *gin.Context) {
	var req dto.OccupancyRateRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("入住率统计参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	data, err := h.reportService.GetOccupancyRateReport(req.StartDate, req.EndDate)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, data)
}

func (h *ReportHandler) GetRevenueReport(c *gin.Context) {
	var req dto.RevenueRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("营收统计参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	data, err := h.reportService.GetRevenueReport(req.StartDate, req.EndDate)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, data)
}

func (h *ReportHandler) ExportReport(c *gin.Context) {
	var req dto.ReportExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("报表导出参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	data, err := h.reportService.ExportReportToExcel(req.StartDate, req.EndDate)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	filename := "报表_" + time.Now().Format("20060102150405") + ".xlsx"

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
