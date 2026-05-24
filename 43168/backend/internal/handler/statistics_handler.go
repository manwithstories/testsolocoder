package handler

import (
	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

// StatisticsHandler 数据统计 HTTP 处理器
type StatisticsHandler struct {
	service *service.StatisticsService
}

// NewStatisticsHandler 创建数据统计处理器
func NewStatisticsHandler(svc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{service: svc}
}

// GetSalesTrend 销售趋势统计
func (h *StatisticsHandler) GetSalesTrend(c *gin.Context) {
	var req dto.StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := h.service.GetSalesTrend(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetCustomerProfile 客户画像统计
func (h *StatisticsHandler) GetCustomerProfile(c *gin.Context) {
	resp, err := h.service.GetCustomerProfile()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, resp)
}

// ExportExcel 导出统计报表
func (h *StatisticsHandler) ExportExcel(c *gin.Context) {
	var req dto.ExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	data, filename, err := h.service.ExportExcel(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
