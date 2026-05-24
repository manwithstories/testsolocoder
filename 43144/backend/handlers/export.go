package handlers

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"pet-adoption-platform/services"
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func ExportAdoptionReport(c *gin.Context) {
	rescueID := c.GetUint("rescue_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	file, err := services.GenerateAdoptionReport(rescueID, startDate, endDate)
	if err != nil {
		utils.InternalError(c, "failed to generate report")
		return
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		utils.InternalError(c, "failed to write report")
		return
	}

	filename := fmt.Sprintf("adoption_report_%s.xlsx", time.Now().Format("20060102_150405"))

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func ExportHealthReport(c *gin.Context) {
	petID, err := strconv.ParseUint(c.Param("pet_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	buf, filename, err := services.GenerateHealthReport(uint(petID))
	if err != nil {
		utils.InternalError(c, "failed to generate health report")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/pdf", buf.Bytes())
}
