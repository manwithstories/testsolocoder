package handler

import (
	"beauty-salon-system/internal/service"
	"beauty-salon-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (h *ReportHandler) GetRevenueReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := h.reportService.GetRevenueReport(startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ReportHandler) GetTechnicianPerformance(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := h.reportService.GetTechnicianPerformance(startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ReportHandler) GetServiceRanking(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := h.reportService.GetServiceRanking(startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ReportHandler) GetFullReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := h.reportService.GetFullReport(startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ReportHandler) ExportExcel(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.reportService.ExportExcel(startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=report.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *ReportHandler) ExportPDF(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.reportService.ExportPDF(startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=report.pdf")
	c.Data(200, "application/pdf", data)
}
