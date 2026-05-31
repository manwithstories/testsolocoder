package handlers

import (
	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type DailyPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

func StatsOverview(c *gin.Context) {
	var orders int64
	database.DB.Model(&models.Order{}).Count(&orders)
	var revenue float64
	database.DB.Model(&models.Settlement{}).Select("COALESCE(SUM(total_amount),0)").Scan(&revenue)
	var customers int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleCustomer).Count(&customers)
	var staff int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleStaff).Count(&staff)
	utils.OK(c, gin.H{
		"orders":    orders,
		"revenue":   revenue,
		"customers": customers,
		"staff":     staff,
	})
}

func StatsRevenueTrend(c *gin.Context) {
	type row struct {
		Date    string  `json:"date"`
		Revenue float64 `json:"revenue"`
	}
	var rows []row
	database.DB.Raw(`SELECT DATE(created_at) as date, COALESCE(SUM(total_amount),0) as revenue
		FROM settlements GROUP BY DATE(created_at) ORDER BY date ASC LIMIT 30`).Scan(&rows)
	utils.OK(c, rows)
}

func StatsOrdersByCategory(c *gin.Context) {
	type row struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	var rows []row
	database.DB.Raw(`SELECT s.category as category, COUNT(b.id) as count
		FROM bookings b JOIN services s ON b.service_id = s.id
		WHERE b.status NOT IN ('canceled','rejected')
		GROUP BY s.category`).Scan(&rows)
	utils.OK(c, rows)
}

func StatsStaffPerf(c *gin.Context) {
	type row struct {
		StaffID    uint    `json:"staff_id"`
		RealName   string  `json:"real_name"`
		OrderCount int64   `json:"order_count"`
		Revenue    float64 `json:"revenue"`
		Rating     float64 `json:"rating"`
	}
	var rows []row
	database.DB.Raw(`SELECT o.staff_id as staff_id, u.real_name as real_name,
		COUNT(o.id) as order_count, SUM(o.total_amount) as revenue, u.rating as rating
		FROM orders o LEFT JOIN users u ON u.id = o.staff_id
		WHERE o.status = 'paid'
		GROUP BY o.staff_id ORDER BY revenue DESC LIMIT 20`).Scan(&rows)
	utils.OK(c, rows)
}

func StatsCompanyDashboard(c *gin.Context) {
	uid, _ := c.Get("uid")
	var totalOrders int64
	database.DB.Model(&models.Order{}).Where("company_id = ?", uid).Count(&totalOrders)
	var totalRevenue float64
	database.DB.Model(&models.Settlement{}).Where("company_id = ?", uid).Select("COALESCE(SUM(company_share),0)").Scan(&totalRevenue)
	var pending int64
	database.DB.Model(&models.Booking{}).Where("status = ?", models.BookingPending).Count(&pending)
	utils.OK(c, gin.H{
		"orders":   totalOrders,
		"revenue":  totalRevenue,
		"pending":  pending,
		"date":     time.Now().Format(time.RFC3339),
	})
}
