package handlers

import (
	"net/http"
	"time"

	"print3d-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	orderService  *service.OrderService
	modelService  *service.ModelService
	printerService *service.PrinterService
}

func NewStatsHandler(
	orderService *service.OrderService,
	modelService *service.ModelService,
	printerService *service.PrinterService,
) *StatsHandler {
	return &StatsHandler{
		orderService:  orderService,
		modelService:  modelService,
		printerService: printerService,
	}
}

func (h *StatsHandler) GetPlatformStats(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM-DD"})
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use YYYY-MM-DD"})
			return
		}
	} else {
		endDate = time.Now()
	}

	orderStats, err := h.orderService.GetOrderStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order stats"})
		return
	}

	modelStats, err := h.modelService.GetStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get model stats"})
		return
	}

	materialStats, err := h.printerService.GetMaterialStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get material stats"})
		return
	}

	hotModels, err := h.modelService.GetHotModels(c.Request.Context(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get hot models"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"order_stats":    orderStats,
		"model_stats":    modelStats,
		"material_stats": materialStats,
		"hot_models":     hotModels,
	})
}

func (h *StatsHandler) GetRevenueStats(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	orderStats, err := h.orderService.GetOrderStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get revenue stats"})
		return
	}

	c.JSON(http.StatusOK, orderStats)
}

func (h *StatsHandler) GetMaterialStats(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	materialStats, err := h.printerService.GetMaterialStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get material stats"})
		return
	}

	c.JSON(http.StatusOK, materialStats)
}

func (h *StatsHandler) ExportStats(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	exportType := c.DefaultQuery("type", "csv")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	orderStats, _ := h.orderService.GetOrderStats(c.Request.Context(), startDate, endDate)
	modelStats, _ := h.modelService.GetStats(c.Request.Context(), startDate, endDate)

	var content string
	var contentType string
	var filename string

	if exportType == "json" {
		contentType = "application/json"
		filename = "stats_report.json"
		content = `{"order_stats":` + interfaceToString(orderStats) + `,"model_stats":` + interfaceToString(modelStats) + `}`
	} else {
		contentType = "text/csv"
		filename = "stats_report.csv"
		content = generateCSVContent(orderStats, modelStats)
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.String(http.StatusOK, content)
}

func generateCSVContent(orderStats, modelStats map[string]interface{}) string {
	content := "Metric,Value\n"
	for k, v := range orderStats {
		content += k + "," + interfaceToString(v) + "\n"
	}
	for k, v := range modelStats {
		content += k + "," + interfaceToString(v) + "\n"
	}
	return content
}

func interfaceToString(v interface{}) string {
	switch val := v.(type) {
	case float64:
		return formatFloat(val)
	case int64:
		return int64ToString(val)
	case string:
		return val
	default:
		return "{}"
	}
}

func formatFloat(f float64) string {
	return formatFloatWithPrecision(f, 2)
}

func formatFloatWithPrecision(f float64, precision int) string {
	format := "%." + intToString(precision) + "f"
	return formatFloatWithFormat(format, f)
}

func formatFloatWithFormat(format string, f float64) string {
	return formatFloatGo(format, f)
}

func formatFloatGo(format string, f float64) string {
	result := ""
	for i := 0; i < len(format); i++ {
		if format[i] == '%' {
			result += format[0:i]
			break
		}
	}
	return sprintfFloat(result, f)
}

func sprintfFloat(format string, f float64) string {
	return formatFloat64(f)
}

func formatFloat64(f float64) string {
	return formatFloatWithStr(f, 2)
}

func formatFloatWithStr(f float64, precision int) string {
	return formatFloat64WithPrecision(f, precision)
}

func formatFloat64WithPrecision(f float64, precision int) string {
	result := ""
	whole := int64(f)
	frac := f - float64(whole)
	result = int64ToString(whole) + "."
	for i := 0; i < precision; i++ {
		frac *= 10
		digit := int64(frac) % 10
		if digit < 0 {
			digit = -digit
		}
		result += intToString(int(digit))
	}
	return result
}

func intToString(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}

func int64ToString(i int64) string {
	if i == 0 {
		return "0"
	}
	result := ""
	negative := i < 0
	if negative {
		i = -i
	}
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
