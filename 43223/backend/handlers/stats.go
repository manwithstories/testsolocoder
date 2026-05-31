package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/models"
	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/tealeg/xlsx"
)

type StatsHandler struct {
	cfg *config.Config
}

func NewStatsHandler(cfg *config.Config) *StatsHandler {
	return &StatsHandler{cfg: cfg}
}

type SalesStats struct {
	TotalOrders     int64   `json:"total_orders"`
	TotalAmount     float64 `json:"total_amount"`
	TotalProducts   int64   `json:"total_products"`
	TotalUsers      int64   `json:"total_users"`
	TotalRoasters   int64   `json:"total_roasters"`
	TodayOrders     int64   `json:"today_orders"`
	TodayAmount     float64 `json:"today_amount"`
}

func (h *StatsHandler) GetSalesStats(c *gin.Context) {
	var stats SalesStats

	database.DB.Model(&models.Order{}).
		Where("status IN ?", []models.OrderStatus{
			models.OrderStatusPaid,
			models.OrderStatusProcessing,
			models.OrderStatusShipped,
			models.OrderStatusDelivered,
		}).
		Count(&stats.TotalOrders)

	database.DB.Model(&models.Order{}).
		Where("status IN ?", []models.OrderStatus{
			models.OrderStatusPaid,
			models.OrderStatusProcessing,
			models.OrderStatusShipped,
			models.OrderStatusDelivered,
		}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalAmount)

	database.DB.Model(&models.Product{}).Count(&stats.TotalProducts)
	database.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	database.DB.Model(&models.User{}).Where("role = ? AND is_certified = ?", models.RoleRoaster, true).
		Count(&stats.TotalRoasters)

	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) = ? AND status IN ?", today, []models.OrderStatus{
			models.OrderStatusPaid, models.OrderStatusProcessing,
			models.OrderStatusShipped, models.OrderStatusDelivered,
		}).
		Count(&stats.TodayOrders)

	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) = ? AND status IN ?", today, []models.OrderStatus{
			models.OrderStatusPaid, models.OrderStatusProcessing,
			models.OrderStatusShipped, models.OrderStatusDelivered,
		}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TodayAmount)

	utils.Success(c, stats)
}

func (h *StatsHandler) GetSalesTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 || days > 365 {
		days = 30
	}

	since := time.Now().AddDate(0, 0, -days)

	type TrendItem struct {
		Date       string  `json:"date"`
		OrderCount int64   `json:"order_count"`
		Amount     float64 `json:"amount"`
	}

	var results []TrendItem
	database.DB.Table("orders").
		Select("DATE(created_at) as date, COUNT(*) as order_count, COALESCE(SUM(total_amount), 0) as amount").
		Where("created_at >= ? AND status IN ?", since, []models.OrderStatus{
			models.OrderStatusPaid, models.OrderStatusProcessing,
			models.OrderStatusShipped, models.OrderStatusDelivered,
		}).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results)

	utils.Success(c, results)
}

func (h *StatsHandler) GetOriginDistribution(c *gin.Context) {
	type OriginItem struct {
		Origin string `json:"origin"`
		Count  int64  `json:"count"`
	}

	var results []OriginItem
	database.DB.Model(&models.Product{}).
		Select("origin, COUNT(*) as count").
		Where("status = ?", models.ProductStatusOnSale).
		Group("origin").
		Order("count DESC").
		Scan(&results)

	utils.Success(c, results)
}

func (h *StatsHandler) GetUserActivityStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 || days > 365 {
		days = 30
	}

	type ActivityItem struct {
		Date       string `json:"date"`
		UserCount  int64  `json:"user_count"`
		LoginCount int64  `json:"login_count"`
	}

	var results []ActivityItem
	since := time.Now().AddDate(0, 0, -days)

	database.DB.Table("operation_logs").
		Select("DATE(created_at) as date, COUNT(DISTINCT user_id) as user_count, COUNT(*) as login_count").
		Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results)

	utils.Success(c, results)
}

func (h *StatsHandler) GetTopProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	type TopProduct struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		TotalSold   int64   `json:"total_sold"`
		TotalAmount float64 `json:"total_amount"`
	}

	var results []TopProduct
	database.DB.Table("order_items oi").
		Select("oi.product_id, oi.product_name, SUM(oi.quantity) as total_sold, SUM(oi.subtotal) as total_amount").
		Joins("JOIN orders o ON oi.order_id = o.id").
		Where("o.status IN ?", []models.OrderStatus{
			models.OrderStatusPaid, models.OrderStatusProcessing,
			models.OrderStatusShipped, models.OrderStatusDelivered,
		}).
		Group("oi.product_id, oi.product_name").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&results)

	utils.Success(c, results)
}

func (h *StatsHandler) ExportExcel(c *gin.Context) {
	reportType := c.Query("type")

	switch reportType {
	case "sales":
		h.exportSalesExcel(c)
	case "orders":
		h.exportOrdersExcel(c)
	case "products":
		h.exportProductsExcel(c)
	case "users":
		h.exportUsersExcel(c)
	default:
		utils.Error(c, http.StatusBadRequest, "不支持的报表类型")
	}
}

func (h *StatsHandler) exportSalesExcel(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := database.DB.Model(&models.Order{}).
		Where("status IN ?", []models.OrderStatus{
			models.OrderStatusPaid, models.OrderStatusProcessing,
			models.OrderStatusShipped, models.OrderStatusDelivered,
		})

	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	var orders []models.Order
	query.Preload("Items").Find(&orders)

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("销售报表")

	headers := []string{"订单号", "用户ID", "总金额", "状态", "支付状态", "创建时间"}
	headerRow := sheet.AddRow()
	for _, h := range headers {
		cell := headerRow.AddCell()
		cell.Value = h
	}

	for _, order := range orders {
		row := sheet.AddRow()
		row.AddCell().Value = order.OrderNo
		row.AddCell().Value = strconv.FormatUint(uint64(order.UserID), 10)
		row.AddCell().Value = fmt.Sprintf("%.2f", order.TotalAmount)
		row.AddCell().Value = string(order.Status)
		row.AddCell().Value = string(order.PaymentStatus)
		row.AddCell().Value = order.CreatedAt.Format("2006-01-02 15:04:05")
	}

	filename := fmt.Sprintf("sales_report_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	var buf []byte
	file.WriteToBuffer(&buf)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}

func (h *StatsHandler) exportOrdersExcel(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("Items.Product").Find(&orders)

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("订单列表")

	headers := []string{"订单号", "用户", "商品", "数量", "单价", "小计", "总金额", "状态", "创建时间"}
	headerRow := sheet.AddRow()
	for _, h := range headers {
		cell := headerRow.AddCell()
		cell.Value = h
	}

	for _, order := range orders {
		for _, item := range order.Items {
			row := sheet.AddRow()
			row.AddCell().Value = order.OrderNo
			row.AddCell().Value = strconv.FormatUint(uint64(order.UserID), 10)
			row.AddCell().Value = item.ProductName
			row.AddCell().Value = strconv.Itoa(item.Quantity)
			row.AddCell().Value = fmt.Sprintf("%.2f", item.Price)
			row.AddCell().Value = fmt.Sprintf("%.2f", item.Subtotal)
			row.AddCell().Value = fmt.Sprintf("%.2f", order.TotalAmount)
			row.AddCell().Value = string(order.Status)
			row.AddCell().Value = order.CreatedAt.Format("2006-01-02 15:04:05")
		}
	}

	filename := fmt.Sprintf("orders_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	var buf []byte
	file.WriteToBuffer(&buf)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}

func (h *StatsHandler) exportProductsExcel(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("商品列表")

	headers := []string{"ID", "名称", "产地", "处理法", "烘焙度", "杯测评分", "价格", "重量", "库存", "状态"}
	headerRow := sheet.AddRow()
	for _, h := range headers {
		cell := headerRow.AddCell()
		cell.Value = h
	}

	for _, p := range products {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.FormatUint(uint64(p.ID), 10)
		row.AddCell().Value = p.Name
		row.AddCell().Value = p.Origin
		row.AddCell().Value = string(p.ProcessMethod)
		row.AddCell().Value = string(p.RoastLevel)
		row.AddCell().Value = fmt.Sprintf("%.2f", p.CuppingScore)
		row.AddCell().Value = fmt.Sprintf("%.2f", p.Price)
		row.AddCell().Value = strconv.Itoa(p.Weight)
		row.AddCell().Value = strconv.Itoa(p.Stock)
		row.AddCell().Value = string(p.Status)
	}

	filename := fmt.Sprintf("products_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	var buf []byte
	file.WriteToBuffer(&buf)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}

func (h *StatsHandler) exportUsersExcel(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("用户列表")

	headers := []string{"ID", "用户名", "邮箱", "手机", "昵称", "角色", "状态", "是否认证", "注册时间"}
	headerRow := sheet.AddRow()
	for _, h := range headers {
		cell := headerRow.AddCell()
		cell.Value = h
	}

	for _, u := range users {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.FormatUint(uint64(u.ID), 10)
		row.AddCell().Value = u.Username
		row.AddCell().Value = u.Email
		row.AddCell().Value = u.Phone
		row.AddCell().Value = u.Nickname
		row.AddCell().Value = string(u.Role)
		row.AddCell().Value = string(u.Status)
		row.AddCell().Value = strconv.FormatBool(u.IsCertified)
		row.AddCell().Value = u.CreatedAt.Format("2006-01-02 15:04:05")
	}

	filename := fmt.Sprintf("users_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	var buf []byte
	file.WriteToBuffer(&buf)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}

func (h *StatsHandler) ExportPDF(c *gin.Context) {
	reportType := c.Query("type")

	c.Header("Content-Type", "application/pdf")
	filename := fmt.Sprintf("%s_report_%s.pdf", reportType, time.Now().Format("20060102150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.NewRect(595.28, 841.89)})

	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 20)
	pdf.Cell(nil, "Coffee Platform Report")
	pdf.Br(30)

	pdf.SetFont("Helvetica", "", 12)
	pdf.Cell(nil, fmt.Sprintf("Report Type: %s", reportType))
	pdf.Br(20)
	pdf.Cell(nil, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Br(30)

	var stats SalesStats
	database.DB.Model(&models.Order{}).
		Where("status IN ?", []models.OrderStatus{
			models.OrderStatusPaid, models.OrderStatusProcessing,
			models.OrderStatusShipped, models.OrderStatusDelivered,
		}).
		Count(&stats.TotalOrders)
	database.DB.Model(&models.Order{}).
		Where("status IN ?", []models.OrderStatus{
			models.OrderStatusPaid, models.OrderStatusProcessing,
			models.OrderStatusShipped, models.OrderStatusDelivered,
		}).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalAmount)

	pdf.Cell(nil, fmt.Sprintf("Total Orders: %d", stats.TotalOrders))
	pdf.Br(15)
	pdf.Cell(nil, fmt.Sprintf("Total Amount: %.2f", stats.TotalAmount))
	pdf.Br(15)

	pdf.WritePdf(c.Writer)
}
