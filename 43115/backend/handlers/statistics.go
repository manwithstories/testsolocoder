package handlers

import (
	"fmt"
	"time"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetDashboardStats(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var totalOrders int64
	config.DB.Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate+" 23:59:59").
		Count(&totalOrders)

	var completedOrders int64
	config.DB.Model(&models.Order{}).
		Where("status = ? AND created_at >= ? AND created_at <= ?",
			models.OrderStatusCompleted, startDate, endDate+" 23:59:59").
		Count(&completedOrders)

	var totalRevenue float64
	config.DB.Model(&models.Order{}).
		Where("status = ? AND completed_at >= ? AND completed_at <= ?",
			models.OrderStatusCompleted, startDate, endDate+" 23:59:59").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	var totalPlatformFee float64
	config.DB.Model(&models.Order{}).
		Where("status = ? AND completed_at >= ? AND completed_at <= ?",
			models.OrderStatusCompleted, startDate, endDate+" 23:59:59").
		Select("COALESCE(SUM(platform_fee), 0)").Scan(&totalPlatformFee)

	var totalCustomers int64
	config.DB.Model(&models.User{}).
		Where("role = ? AND created_at >= ? AND created_at <= ?",
			models.RoleCustomer, startDate, endDate+" 23:59:59").
		Count(&totalCustomers)

	var totalProviders int64
	config.DB.Model(&models.User{}).
		Where("role = ? AND provider_status = ?",
			models.RoleServiceProvider, models.ProviderStatusApproved).
		Count(&totalProviders)

	var orderTrend []struct {
		Date       string `json:"date"`
		OrderCount int64  `json:"order_count"`
		Revenue    float64 `json:"revenue"`
	}
	config.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as order_count,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN total_amount ELSE 0 END), 0) as revenue
		FROM orders
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, startDate, endDate+" 23:59:59").Scan(&orderTrend)

	var serviceTypeDistribution []struct {
		CategoryName string `json:"category_name"`
		OrderCount   int64  `json:"order_count"`
		Percentage   float64 `json:"percentage"`
	}
	config.DB.Raw(`
		SELECT 
			sc.name as category_name,
			COUNT(o.id) as order_count,
			ROUND(COUNT(o.id) * 100.0 / NULLIF((SELECT COUNT(*) FROM orders WHERE created_at >= ? AND created_at <= ?), 0), 2) as percentage
		FROM orders o
		JOIN service_items si ON o.service_item_id = si.id
		JOIN service_categories sc ON si.category_id = sc.id
		WHERE o.created_at >= ? AND o.created_at <= ?
		GROUP BY sc.name
		ORDER BY order_count DESC
	`, startDate, endDate+" 23:59:59", startDate, endDate+" 23:59:59").Scan(&serviceTypeDistribution)

	var topProviders []struct {
		ID          uint    `json:"id"`
		Nickname    string  `json:"nickname"`
		Rating      float64 `json:"rating"`
		OrderCount  int64   `json:"order_count"`
		TotalIncome float64 `json:"total_income"`
	}
	config.DB.Raw(`
		SELECT 
			u.id,
			u.nickname,
			u.rating,
			u.order_count,
			COALESCE(SUM(o.provider_income), 0) as total_income
		FROM users u
		LEFT JOIN orders o ON u.id = o.provider_id AND o.status = 'completed' AND o.completed_at >= ? AND o.completed_at <= ?
		WHERE u.role = 'service_provider' AND u.provider_status = 'approved'
		GROUP BY u.id, u.nickname, u.rating, u.order_count
		ORDER BY u.rating DESC, u.order_count DESC
		LIMIT 10
	`, startDate, endDate+" 23:59:59").Scan(&topProviders)

	utils.Success(c, gin.H{
		"total_orders":              totalOrders,
		"completed_orders":          completedOrders,
		"total_revenue":             utils.RoundToTwoDecimals(totalRevenue),
		"total_platform_fee":        utils.RoundToTwoDecimals(totalPlatformFee),
		"new_customers":             totalCustomers,
		"active_providers":          totalProviders,
		"order_trend":               orderTrend,
		"service_type_distribution": serviceTypeDistribution,
		"top_providers":             topProviders,
	})
}

func GetOrderStatistics(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	status := c.Query("status")
	serviceType := c.Query("service_type")

	query := config.DB.Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate+" 23:59:59")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if serviceType != "" {
		var serviceItemIDs []uint
		config.DB.Model(&models.ServiceItem{}).Where("category_id = ?", serviceType).Pluck("id", &serviceItemIDs)
		query = query.Where("service_item_id IN ?", serviceItemIDs)
	}

	var totalOrders int64
	query.Count(&totalOrders)

	var totalAmount float64
	config.DB.Model(&models.Order{}).
		Where("status = ? AND completed_at >= ? AND completed_at <= ?",
			models.OrderStatusCompleted, startDate, endDate+" 23:59:59").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalAmount)

	var statusDistribution []struct {
		Status     string `json:"status"`
		OrderCount int64  `json:"order_count"`
	}
	config.DB.Raw(`
		SELECT status, COUNT(*) as order_count
		FROM orders
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY status
	`, startDate, endDate+" 23:59:59").Scan(&statusDistribution)

	var hourlyDistribution []struct {
		Hour       int   `json:"hour"`
		OrderCount int64 `json:"order_count"`
	}
	config.DB.Raw(`
		SELECT EXTRACT(HOUR FROM created_at)::int as hour, COUNT(*) as order_count
		FROM orders
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY EXTRACT(HOUR FROM created_at)
		ORDER BY hour
	`, startDate, endDate+" 23:59:59").Scan(&hourlyDistribution)

	utils.Success(c, gin.H{
		"total_orders":         totalOrders,
		"total_amount":         utils.RoundToTwoDecimals(totalAmount),
		"status_distribution":  statusDistribution,
		"hourly_distribution":  hourlyDistribution,
	})
}

func GetUserStatistics(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var totalCustomers int64
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleCustomer).Count(&totalCustomers)

	var totalProviders int64
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleServiceProvider).Count(&totalProviders)

	var newCustomers int64
	config.DB.Model(&models.User{}).
		Where("role = ? AND created_at >= ? AND created_at <= ?",
			models.RoleCustomer, startDate, endDate+" 23:59:59").
		Count(&newCustomers)

	var newProviders int64
	config.DB.Model(&models.User{}).
		Where("role = ? AND created_at >= ? AND created_at <= ?",
			models.RoleServiceProvider, startDate, endDate+" 23:59:59").
		Count(&newProviders)

	var activeCustomers int64
	config.DB.Raw(`
		SELECT COUNT(DISTINCT customer_id)
		FROM orders
		WHERE created_at >= ? AND created_at <= ?
	`, startDate, endDate+" 23:59:59").Scan(&activeCustomers)

	var activeProviders int64
	config.DB.Raw(`
		SELECT COUNT(DISTINCT provider_id)
		FROM orders
		WHERE status = 'confirmed' OR status = 'in_service'
	`, startDate, endDate+" 23:59:59").Scan(&activeProviders)

	var customerGrowth []struct {
		Date      string `json:"date"`
		NewUsers  int64  `json:"new_users"`
	}
	config.DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as new_users
		FROM users
		WHERE role = 'customer' AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, startDate, endDate+" 23:59:59").Scan(&customerGrowth)

	var ratingDistribution []struct {
		Rating      float64 `json:"rating"`
		ProviderCount int64 `json:"provider_count"`
	}
	config.DB.Raw(`
		SELECT 
			ROUND(rating, 1) as rating,
			COUNT(*) as provider_count
		FROM users
		WHERE role = 'service_provider' AND provider_status = 'approved'
		GROUP BY ROUND(rating, 1)
		ORDER BY rating DESC
	`).Scan(&ratingDistribution)

	utils.Success(c, gin.H{
		"total_customers":      totalCustomers,
		"total_providers":      totalProviders,
		"new_customers":        newCustomers,
		"new_providers":        newProviders,
		"active_customers":     activeCustomers,
		"active_providers":     activeProviders,
		"customer_growth":      customerGrowth,
		"rating_distribution":  ratingDistribution,
	})
}

func ExportStatistics(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	exportType := c.DefaultQuery("type", "orders")

	f := excelize.NewFile()

	switch exportType {
	case "orders":
		exportOrderStatistics(f, startDate, endDate)
	case "users":
		exportUserStatistics(f, startDate, endDate)
	case "revenue":
		exportRevenueStatistics(f, startDate, endDate)
	default:
		exportOrderStatistics(f, startDate, endDate)
	}

	fileName := fmt.Sprintf("statistics_%s_%s_%s.xlsx", exportType, startDate, endDate)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fileName)

	if err := f.Write(c.Writer); err != nil {
		utils.InternalError(c, "导出失败")
		return
	}
}

func exportOrderStatistics(f *excelize.File, startDate, endDate string) {
	sheet := "订单统计"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"订单号", "客户", "服务人员", "服务类型", "金额", "状态", "预约时间", "创建时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	var orders []struct {
		OrderNo         string
		CustomerName    string
		ProviderName    string
		ServiceName     string
		TotalAmount     float64
		Status          string
		AppointmentTime time.Time
		CreatedAt       time.Time
	}

	config.DB.Raw(`
		SELECT o.order_no, c.nickname as customer_name, p.nickname as provider_name,
			si.name as service_name, o.total_amount, o.status, o.appointment_time, o.created_at
		FROM orders o
		LEFT JOIN users c ON o.customer_id = c.id
		LEFT JOIN users p ON o.provider_id = p.id
		LEFT JOIN service_items si ON o.service_item_id = si.id
		WHERE o.created_at >= ? AND o.created_at <= ?
		ORDER BY o.id DESC
	`, startDate, endDate+" 23:59:59").Scan(&orders)

	for i, order := range orders {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), order.OrderNo)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), order.CustomerName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), order.ProviderName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), order.ServiceName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), order.TotalAmount)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), order.Status)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), order.AppointmentTime.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), order.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func exportUserStatistics(f *excelize.File, startDate, endDate string) {
	sheet := "用户统计"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"用户ID", "昵称", "手机号", "角色", "状态", "注册时间", "订单数", "评分"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	var users []struct {
		ID        uint
		Nickname  string
		Phone     string
		Role      string
		IsActive  bool
		CreatedAt time.Time
		OrderCount int
		Rating    float64
	}

	config.DB.Raw(`
		SELECT id, nickname, phone, role, is_active, created_at, order_count, rating
		FROM users
		WHERE created_at >= ? AND created_at <= ?
		ORDER BY id DESC
	`, startDate, endDate+" 23:59:59").Scan(&users)

	for i, user := range users {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), user.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), user.Nickname)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), user.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), user.Role)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), user.IsActive)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), user.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), user.OrderCount)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), user.Rating)
	}
}

func exportRevenueStatistics(f *excelize.File, startDate, endDate string) {
	sheet := "收入统计"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"日期", "订单数", "总收入", "平台佣金", "服务人员收入"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	var revenues []struct {
		Date          string
		OrderCount    int64
		TotalAmount   float64
		PlatformFee   float64
		ProviderIncome float64
	}

	config.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as order_count,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN total_amount ELSE 0 END), 0) as total_amount,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN platform_fee ELSE 0 END), 0) as platform_fee,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN provider_income ELSE 0 END), 0) as provider_income
		FROM orders
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, startDate, endDate+" 23:59:59").Scan(&revenues)

	for i, revenue := range revenues {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), revenue.Date)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), revenue.OrderCount)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), revenue.TotalAmount)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), revenue.PlatformFee)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), revenue.ProviderIncome)
	}
}
