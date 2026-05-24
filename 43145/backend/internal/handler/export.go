package handler

import (
	"net/http"
	"strconv"
	"survey-platform/internal/service"
	"survey-platform/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	exportService *service.ExportService
}

func NewExportHandler(exportService *service.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

func (h *ExportHandler) ExportExcel(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	data, err := h.exportService.ExportExcel(uint(surveyID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	filename := "survey_data_" + time.Now().Format("20060102_150405") + ".xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *ExportHandler) ExportPDF(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	report, err := h.exportService.ExportPDF(uint(surveyID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, report)
}

func (h *ExportHandler) ExportChartImages(c *gin.Context) {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid survey id")
		return
	}

	charts, err := h.exportService.ExportChartImages(uint(surveyID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, charts)
}
