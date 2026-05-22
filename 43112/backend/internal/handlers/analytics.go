package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type AnalyticsHandler struct {
	cfg *config.Config
}

func NewAnalyticsHandler(cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{cfg: cfg}
}

func (h *AnalyticsHandler) GetInstructorDashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var courseCount int64
	database.DB.Model(&models.Course{}).Where("instructor_id = ?", userID).Count(&courseCount)

	var totalStudents int64
	database.DB.Model(&models.Course{}).Where("instructor_id = ?", userID).
		Select("COALESCE(SUM(student_count), 0)").Scan(&totalStudents)

	var totalRevenue float64
	database.DB.Model(&models.Order{}).
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.status = ?", userID, models.OrderPaid).
		Select("COALESCE(SUM(orders.amount), 0)").Scan(&totalRevenue)

	var orderCount int64
	database.DB.Model(&models.Order{}).
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.status = ?", userID, models.OrderPaid).
		Count(&orderCount)

	var courses []models.Course
	database.DB.Where("instructor_id = ?", userID).
		Order("created_at DESC").Limit(5).Find(&courses)

	var recentOrders []models.Order
	database.DB.Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ?", userID).
		Order("orders.created_at DESC").Limit(10).
		Preload("User").Preload("Course").Find(&recentOrders)

	utils.Success(c, gin.H{
		"course_count":    courseCount,
		"total_students":  totalStudents,
		"total_revenue":   totalRevenue,
		"order_count":     orderCount,
		"recent_courses":  courses,
		"recent_orders":   recentOrders,
	})
}

func (h *AnalyticsHandler) GetInstructorRevenueReport(c *gin.Context) {
	userID, _ := c.Get("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := database.DB.Model(&models.Order{}).
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.status = ?", userID, models.OrderPaid)

	if startDate != "" {
		query = query.Where("orders.paid_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("orders.paid_at <= ?", endDate)
	}

	type RevenueByCourse struct {
		CourseID    uuid.UUID `json:"course_id"`
		CourseTitle string    `json:"course_title"`
		OrderCount  int64     `json:"order_count"`
		TotalAmount float64   `json:"total_amount"`
	}

	var results []RevenueByCourse
	query.Select("courses.id as course_id, courses.title as course_title, COUNT(orders.id) as order_count, COALESCE(SUM(orders.amount), 0) as total_amount").
		Group("courses.id, courses.title").
		Order("total_amount DESC").
		Scan(&results)

	var totalRevenue float64
	var totalOrders int64
	for _, r := range results {
		totalRevenue += r.TotalAmount
		totalOrders += r.OrderCount
	}

	utils.Success(c, gin.H{
		"by_course":    results,
		"total_revenue": totalRevenue,
		"total_orders":  totalOrders,
	})
}

func (h *AnalyticsHandler) ExportInstructorReport(c *gin.Context) {
	userID, _ := c.Get("user_id")

	type ReportRow struct {
		OrderNo     string    `json:"order_no"`
		CourseTitle string    `json:"course_title"`
		StudentName string    `json:"student_name"`
		Amount      float64   `json:"amount"`
		PaidAt      time.Time `json:"paid_at"`
	}

	var rows []ReportRow
	database.DB.Model(&models.Order{}).
		Select("orders.order_no, courses.title as course_title, users.username as student_name, orders.amount, orders.paid_at").
		Joins("JOIN courses ON courses.id = orders.course_id").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("courses.instructor_id = ? AND orders.status = ?", userID, models.OrderPaid).
		Order("orders.paid_at DESC").
		Scan(&rows)

	columns := []utils.ExcelColumn{
		{Header: "订单号", Key: "order_no", Width: 25},
		{Header: "课程名称", Key: "course_title", Width: 30},
		{Header: "学员", Key: "student_name", Width: 15},
		{Header: "金额", Key: "amount", Width: 12},
		{Header: "支付时间", Key: "paid_at", Width: 20},
	}

	data := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		data[i] = map[string]interface{}{
			"order_no":     row.OrderNo,
			"course_title": row.CourseTitle,
			"student_name": row.StudentName,
			"amount":       fmt.Sprintf("%.2f", row.Amount),
			"paid_at":      row.PaidAt.Format("2006-01-02 15:04:05"),
		}
	}

	buf, err := utils.ExportExcel("instructor_report.xlsx", columns, data)
	if err != nil {
		utils.InternalError(c, "Failed to export Excel")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=instructor_report.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}

func (h *AnalyticsHandler) GetAdminDashboard(c *gin.Context) {
	var userCount int64
	database.DB.Model(&models.User{}).Where("role = ?", "student").Count(&userCount)

	var instructorCount int64
	database.DB.Model(&models.User{}).Where("role = ?", "instructor").Count(&instructorCount)

	var courseCount int64
	database.DB.Model(&models.Course{}).Where("status = ?", models.CoursePublished).Count(&courseCount)

	var totalRevenue float64
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderPaid).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	var orderCount int64
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderPaid).Count(&orderCount)

	var pendingApplications int64
	database.DB.Model(&models.InstructorApplication{}).Where("status = ?", models.InstructorPending).Count(&pendingApplications)

	var recentOrders []models.Order
	database.DB.Order("created_at DESC").Limit(10).
		Preload("User").Preload("Course").Find(&recentOrders)

	var recentUsers []models.User
	database.DB.Where("role = ?", "student").Order("created_at DESC").Limit(10).Find(&recentUsers)

	utils.Success(c, gin.H{
		"user_count":           userCount,
		"instructor_count":     instructorCount,
		"course_count":         courseCount,
		"total_revenue":        totalRevenue,
		"order_count":          orderCount,
		"pending_applications": pendingApplications,
		"recent_orders":        recentOrders,
		"recent_users":         recentUsers,
	})
}

func (h *AnalyticsHandler) GetAdminRevenueReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := database.DB.Model(&models.Order{}).Where("status = ?", models.OrderPaid)
	if startDate != "" {
		query = query.Where("paid_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("paid_at <= ?", endDate)
	}

	type RevenueByDate struct {
		Date        string  `json:"date"`
		TotalAmount float64 `json:"total_amount"`
		OrderCount  int64   `json:"order_count"`
	}

	var results []RevenueByDate
	query.Select("DATE(paid_at) as date, COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as order_count").
		Group("DATE(paid_at)").
		Order("date DESC").
		Scan(&results)

	type RevenueByCourse struct {
		CourseID    uuid.UUID `json:"course_id"`
		CourseTitle string    `json:"course_title"`
		Instructor  string    `json:"instructor"`
		OrderCount  int64     `json:"order_count"`
		TotalAmount float64   `json:"total_amount"`
	}

	var byCourse []RevenueByCourse
	query.Select("courses.id as course_id, courses.title as course_title, users.username as instructor, COUNT(orders.id) as order_count, COALESCE(SUM(orders.amount), 0) as total_amount").
		Joins("JOIN courses ON courses.id = orders.course_id").
		Joins("JOIN users ON users.id = courses.instructor_id").
		Group("courses.id, courses.title, users.username").
		Order("total_amount DESC").
		Scan(&byCourse)

	utils.Success(c, gin.H{
		"by_date":   results,
		"by_course": byCourse,
	})
}

func (h *AnalyticsHandler) ExportAdminReport(c *gin.Context) {
	type ReportRow struct {
		OrderNo      string    `json:"order_no"`
		CourseTitle  string    `json:"course_title"`
		Instructor   string    `json:"instructor"`
		StudentName  string    `json:"student_name"`
		Amount       float64   `json:"amount"`
		PaidAt       time.Time `json:"paid_at"`
		Status       string    `json:"status"`
	}

	var rows []ReportRow
	database.DB.Model(&models.Order{}).
		Select("orders.order_no, courses.title as course_title, instructors.username as instructor, "+
			"users.username as student_name, orders.amount, orders.paid_at, orders.status").
		Joins("JOIN courses ON courses.id = orders.course_id").
		Joins("JOIN users ON users.id = orders.user_id").
		Joins("JOIN users instructors ON instructors.id = courses.instructor_id").
		Order("orders.created_at DESC").
		Scan(&rows)

	columns := []utils.ExcelColumn{
		{Header: "订单号", Key: "order_no", Width: 25},
		{Header: "课程名称", Key: "course_title", Width: 30},
		{Header: "讲师", Key: "instructor", Width: 15},
		{Header: "学员", Key: "student_name", Width: 15},
		{Header: "金额", Key: "amount", Width: 12},
		{Header: "支付时间", Key: "paid_at", Width: 20},
		{Header: "状态", Key: "status", Width: 12},
	}

	data := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		data[i] = map[string]interface{}{
			"order_no":     row.OrderNo,
			"course_title": row.CourseTitle,
			"instructor":   row.Instructor,
			"student_name": row.StudentName,
			"amount":       fmt.Sprintf("%.2f", row.Amount),
			"paid_at":      row.PaidAt.Format("2006-01-02 15:04:05"),
			"status":       row.Status,
		}
	}

	buf, err := utils.ExportExcel("admin_report.xlsx", columns, data)
	if err != nil {
		utils.InternalError(c, "Failed to export Excel")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=admin_report.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}

func (h *AnalyticsHandler) ExportUserReport(c *gin.Context) {
	type UserRow struct {
		ID               uuid.UUID `json:"id"`
		Username         string    `json:"username"`
		Email            string    `json:"email"`
		Nickname         string    `json:"nickname"`
		Role             string    `json:"role"`
		InstructorStatus string    `json:"instructor_status"`
		CreatedAt        time.Time `json:"created_at"`
	}

	var rows []UserRow
	database.DB.Model(&models.User{}).
		Select("id, username, email, nickname, role, instructor_status, created_at").
		Order("created_at DESC").
		Scan(&rows)

	columns := []utils.ExcelColumn{
		{Header: "ID", Key: "id", Width: 40},
		{Header: "用户名", Key: "username", Width: 15},
		{Header: "邮箱", Key: "email", Width: 25},
		{Header: "昵称", Key: "nickname", Width: 15},
		{Header: "角色", Key: "role", Width: 10},
		{Header: "讲师状态", Key: "instructor_status", Width: 12},
		{Header: "注册时间", Key: "created_at", Width: 20},
	}

	data := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		data[i] = map[string]interface{}{
			"id":                row.ID.String(),
			"username":          row.Username,
			"email":             row.Email,
			"nickname":          row.Nickname,
			"role":              row.Role,
			"instructor_status": row.InstructorStatus,
			"created_at":        row.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	buf, err := utils.ExportExcel("users_report.xlsx", columns, data)
	if err != nil {
		utils.InternalError(c, "Failed to export Excel")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=users_report.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}
