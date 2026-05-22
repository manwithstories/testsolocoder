package handlers

import (
	"multishop/internal/database"
	"multishop/internal/dto"
	"multishop/internal/middleware"
	"multishop/internal/models"
	"multishop/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type StatisticsHandler struct{}

func NewStatisticsHandler() *StatisticsHandler {
	return &StatisticsHandler{}
}

func (h *StatisticsHandler) GetShopStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, 404, "店铺不存在")
		return
	}

	var totalProducts int64
	database.DB.Model(&models.Product{}).Where("shop_id = ?", shop.ID).Count(&totalProducts)

	var totalOrders int64
	database.DB.Model(&models.Order{}).Where("shop_id = ? AND status != ?", shop.ID, models.OrderStatusCancelled).
		Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59").
		Count(&totalOrders)

	var totalSales float64
	database.DB.Model(&models.Order{}).Where("shop_id = ?", shop.ID).
		Where("status IN ?", []string{models.OrderStatusCompleted, models.OrderStatusRefunded}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalSales)

	var newOrders int64
	database.DB.Model(&models.Order{}).Where("shop_id = ? AND status = ?", shop.ID, models.OrderStatusPendingShip).
		Count(&newOrders)

	type DailySalesResult struct {
		Date   string  `gorm:"column:date"`
		Amount float64 `gorm:"column:amount"`
		Orders int64   `gorm:"column:orders"`
	}
	var dailySalesResults []DailySalesResult
	database.DB.Model(&models.Order{}).
		Select("DATE(created_at) as date, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as orders").
		Where("shop_id = ? AND status IN ?", shop.ID, []string{models.OrderStatusCompleted, models.OrderStatusRefunded}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59").
		Group("DATE(created_at)").Order("date ASC").
		Scan(&dailySalesResults)

	dailySales := make([]dto.DailySales, 0, len(dailySalesResults))
	for _, ds := range dailySalesResults {
		dailySales = append(dailySales, dto.DailySales{
			Date:   ds.Date,
			Amount: ds.Amount,
			Orders: ds.Orders,
		})
	}

	type ProductSalesResult struct {
		ID    uint    `gorm:"column:id"`
		Name  string  `gorm:"column:name"`
		Sales int     `gorm:"column:sales"`
		Amount float64 `gorm:"column:amount"`
	}
	var topProductsResults []ProductSalesResult
	database.DB.Table("order_items oi").
		Select("oi.product_id as id, p.name, SUM(oi.quantity) as sales, SUM(oi.subtotal) as amount").
		Joins("LEFT JOIN products p ON oi.product_id = p.id").
		Joins("LEFT JOIN orders o ON oi.order_id = o.id").
		Where("o.shop_id = ? AND o.status IN ?", shop.ID, []string{models.OrderStatusCompleted, models.OrderStatusRefunded}).
		Where("o.created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59").
		Group("oi.product_id, p.name").Order("sales DESC").Limit(10).
		Scan(&topProductsResults)

	topProducts := make([]dto.ProductSales, 0, len(topProductsResults))
	for _, ps := range topProductsResults {
		topProducts = append(topProducts, dto.ProductSales{
			ID:     ps.ID,
			Name:   ps.Name,
			Sales:  ps.Sales,
			Amount: ps.Amount,
		})
	}

	utils.Success(c, dto.ShopStatistics{
		TotalOrders:   totalOrders,
		TotalSales:    totalSales,
		TotalProducts: totalProducts,
		NewOrders:     newOrders,
		DailySales:    dailySales,
		TopProducts:   topProducts,
	})
}

func (h *StatisticsHandler) ExportShopOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, 404, "店铺不存在")
		return
	}

	var orders []models.Order
	query := database.DB.Preload("Items").Preload("User").Where("shop_id = ?", shop.ID)
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59")
	}
	query.Order("created_at DESC").Find(&orders)

	f := excelize.NewFile()
	sheetName := "订单数据"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"订单号", "下单时间", "买家", "商品", "金额", "状态", "收货人", "电话", "地址", "备注"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	for i, order := range orders {
		row := i + 2
		rowStr := strconv.Itoa(row)
		items := ""
		for j, item := range order.Items {
			if j > 0 {
				items += "; "
			}
			items += item.ProductName + " x" + strconv.Itoa(item.Quantity)
		}

		statusText := map[string]string{
			models.OrderStatusPendingPayment: "待付款",
			models.OrderStatusPendingShip:    "待发货",
			models.OrderStatusShipped:        "已发货",
			models.OrderStatusCompleted:      "已完成",
			models.OrderStatusRefunded:       "已退款",
			models.OrderStatusCancelled:      "已取消",
		}[order.Status]

		f.SetCellValue(sheetName, "A"+rowStr, order.OrderNo)
		f.SetCellValue(sheetName, "B"+rowStr, order.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, "C"+rowStr, order.User.Username)
		f.SetCellValue(sheetName, "D"+rowStr, items)
		f.SetCellValue(sheetName, "E"+rowStr, order.TotalAmount)
		f.SetCellValue(sheetName, "F"+rowStr, statusText)
		f.SetCellValue(sheetName, "G"+rowStr, order.ReceiverName)
		f.SetCellValue(sheetName, "H"+rowStr, order.ReceiverPhone)
		f.SetCellValue(sheetName, "I"+rowStr, order.ReceiverAddress)
		f.SetCellValue(sheetName, "J"+rowStr, order.Remark)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=shop_orders.xlsx")
	f.Write(c.Writer)
}
