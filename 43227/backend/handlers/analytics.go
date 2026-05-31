package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct{}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{}
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var totalBeehives int64
	database.DB.Model(&models.Beehive{}).Where("user_id = ?", userID).Count(&totalBeehives)

	var totalHarvest float64
	database.DB.Model(&models.Harvest{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(quantity), 0)").Scan(&totalHarvest)

	var totalProducts int64
	database.DB.Model(&models.Product{}).Where("user_id = ?", userID).Count(&totalProducts)

	var totalOrders int64
	database.DB.Model(&models.Order{}).Where("seller_id = ?", userID).Count(&totalOrders)

	var totalRevenue float64
	database.DB.Model(&models.Order{}).
		Where("seller_id = ? AND status IN ?", userID, []string{"completed", "delivered", "shipped"}).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	utils.Success(c, gin.H{
		"total_beehives": totalBeehives,
		"total_harvest":  totalHarvest,
		"total_products": totalProducts,
		"total_orders":   totalOrders,
		"total_revenue":  totalRevenue,
	})
}

func (h *AnalyticsHandler) GetProductionStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := database.DB.Model(&models.Harvest{}).Where("user_id = ?", userID)

	if startDate != "" {
		query = query.Where("harvest_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("harvest_date <= ?", endDate)
	}

	var totalQuantity float64
	query.Select("COALESCE(SUM(quantity), 0)").Scan(&totalQuantity)

	type HoneyTypeStats struct {
		HoneyType string  `json:"honey_type"`
		TotalQty  float64 `json:"total_qty"`
		Count     int64   `json:"count"`
	}

	var honeyTypeStats []HoneyTypeStats
	query.Select("honey_type, SUM(quantity) as total_qty, COUNT(*) as count").
		Group("honey_type").Find(&honeyTypeStats)

	type MonthlyStats struct {
		Month     string  `json:"month"`
		TotalQty  float64 `json:"total_qty"`
		Count     int64   `json:"count"`
	}

	var monthlyStats []MonthlyStats
	database.DB.Model(&models.Harvest{}).
		Where("user_id = ?", userID).
		Select("to_char(harvest_date, 'YYYY-MM') as month, COALESCE(SUM(quantity), 0) as total_qty, COUNT(*) as count").
		Group("month").
		Order("month").
		Scan(&monthlyStats)

	utils.Success(c, gin.H{
		"total_quantity":   totalQuantity,
		"honey_type_stats": honeyTypeStats,
		"monthly_stats":    monthlyStats,
	})
}

func (h *AnalyticsHandler) GetDiseaseStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var totalRecords int64
	database.DB.Model(&models.HealthRecord{}).
		Joins("JOIN beehives ON health_records.beehive_id = beehives.id").
		Where("beehives.user_id = ?", userID).
		Count(&totalRecords)

	var diseaseRecords int64
	database.DB.Model(&models.HealthRecord{}).
		Joins("JOIN beehives ON health_records.beehive_id = beehives.id").
		Where("beehives.user_id = ? AND health_records.has_disease = ?", userID, true).
		Count(&diseaseRecords)

	diseaseRate := float64(0)
	if totalRecords > 0 {
		diseaseRate = float64(diseaseRecords) / float64(totalRecords) * 100
	}

	type DiseaseTypeStats struct {
		DiseaseType string `json:"disease_type"`
		Count       int64  `json:"count"`
	}

	var diseaseTypeStats []DiseaseTypeStats
	database.DB.Model(&models.HealthRecord{}).
		Joins("JOIN beehives ON health_records.beehive_id = beehives.id").
		Where("beehives.user_id = ? AND health_records.has_disease = ?", userID, true).
		Select("health_records.disease_type, COUNT(*) as count").
		Group("health_records.disease_type").
		Find(&diseaseTypeStats)

	type SeasonalDiseaseStats struct {
		Season         string `json:"season"`
		TotalRecords   int64  `json:"total_records"`
		DiseaseRecords int64  `json:"disease_records"`
	}

	var seasonalDiseaseStats []SeasonalDiseaseStats
	database.DB.Table("health_records").
		Select(`health_records.season, 
			COUNT(*) as total_records,
			COUNT(CASE WHEN health_records.has_disease = true THEN 1 END) as disease_records`).
		Joins("JOIN beehives ON health_records.beehive_id = beehives.id").
		Where("beehives.user_id = ?", userID).
		Group("health_records.season").
		Scan(&seasonalDiseaseStats)

	utils.Success(c, gin.H{
		"total_records":          totalRecords,
		"disease_records":        diseaseRecords,
		"disease_rate":           fmt.Sprintf("%.2f%%", diseaseRate),
		"disease_type_stats":     diseaseTypeStats,
		"seasonal_disease_stats": seasonalDiseaseStats,
	})
}

func (h *AnalyticsHandler) GetSalesStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := database.DB.Model(&models.Order{}).Where("seller_id = ?", userID)

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	var totalOrders int64
	query.Count(&totalOrders)

	var totalRevenue float64
	query.Where("status IN ?", []string{"completed", "delivered", "shipped"}).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	type OrderStatusStats struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
		Amount float64 `json:"amount"`
	}

	var orderStatusStats []OrderStatusStats
	database.DB.Model(&models.Order{}).
		Where("seller_id = ?", userID).
		Select("status, COUNT(*) as count, COALESCE(SUM(total_amount), 0) as amount").
		Group("status").
		Find(&orderStatusStats)

	type MonthlySalesStats struct {
		Month        string  `json:"month"`
		OrdersCount  int64   `json:"orders_count"`
		TotalRevenue float64 `json:"total_revenue"`
	}

	var monthlySalesStats []MonthlySalesStats
	database.DB.Model(&models.Order{}).
		Where("seller_id = ?", userID).
		Select("to_char(created_at, 'YYYY-MM') as month, COUNT(*) as orders_count, COALESCE(SUM(total_amount), 0) as total_revenue").
		Group("month").
		Order("month").
		Scan(&monthlySalesStats)

	utils.Success(c, gin.H{
		"total_orders":          totalOrders,
		"total_revenue":         totalRevenue,
		"order_status_stats":    orderStatusStats,
		"monthly_sales_stats":   monthlySalesStats,
	})
}

func (h *AnalyticsHandler) ExportReport(c *gin.Context) {
	userID, _ := c.Get("user_id")
	reportType := c.Query("type")

	filename := fmt.Sprintf("%s_report_%s.csv", reportType, time.Now().Format("20060102150405"))

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	switch reportType {
	case "harvest":
		writer.Write([]string{"日期", "蜂箱", "蜂蜜类型", "产量(kg)", "批次号", "质量等级"})

		var harvests []models.Harvest
		database.DB.Where("user_id = ?", userID).
			Preload("Beehive").
			Order("harvest_date DESC").
			Find(&harvests)

		for _, h := range harvests {
			writer.Write([]string{
				h.HarvestDate.Format("2006-01-02"),
				h.Beehive.Name,
				h.HoneyType,
				strconv.FormatFloat(h.Quantity, 'f', 2, 64),
				h.BatchCode,
				h.Quality,
			})
		}

	case "inventory":
		writer.Write([]string{"批次号", "蜂蜜类型", "库存数量(kg)", "等级", "状态", "保质期"})

		var inventories []models.Inventory
		database.DB.Where("user_id = ?", userID).Find(&inventories)

		for _, inv := range inventories {
			writer.Write([]string{
				inv.BatchCode,
				inv.HoneyType,
				strconv.FormatFloat(inv.Quantity, 'f', 2, 64),
				inv.Grade,
				inv.Status,
				inv.ExpiryDate.Format("2006-01-02"),
			})
		}

	case "orders":
		writer.Write([]string{"订单号", "买家", "商品", "数量", "单价", "总金额", "状态", "创建时间"})

		var orders []models.Order
		database.DB.Where("seller_id = ?", userID).
			Preload("Buyer").Preload("Product").
			Order("created_at DESC").
			Find(&orders)

		for _, order := range orders {
			writer.Write([]string{
				order.OrderNo,
				order.Buyer.Username,
				order.Product.Title,
				strconv.FormatFloat(order.Quantity, 'f', 2, 64),
				strconv.FormatFloat(order.UnitPrice, 'f', 2, 64),
				strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
				order.Status,
				order.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report type"})
		return
	}
}
