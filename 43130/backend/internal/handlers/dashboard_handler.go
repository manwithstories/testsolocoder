package handlers

import (
	"time"
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")

	db := database.GetDB()

	weddingID := c.Query("wedding_id")

	var weddingQuery *WeddingQuery
	if weddingID != "" {
		weddingQuery = &WeddingQuery{ID: weddingID}
	}

	var totalWeddings int64
	var totalGuests int64
	var totalVendors int64
	var totalBudget float64
	var totalPaid float64

	weddingDBQuery := db.Model(&models.Wedding{})
	if userRole != "admin" {
		weddingDBQuery = weddingDBQuery.Where("user_id = ?", userID)
	}
	if weddingQuery != nil {
		weddingDBQuery = weddingDBQuery.Where("id = ?", weddingQuery.ID)
	}
	weddingDBQuery.Count(&totalWeddings)

	guestQuery := db.Model(&models.Guest{})
	if weddingQuery != nil {
		guestQuery = guestQuery.Where("wedding_id = ?", weddingQuery.ID)
	} else if userRole != "admin" {
		guestQuery = guestQuery.Joins("JOIN weddings ON weddings.id = guests.wedding_id").
			Where("weddings.user_id = ?", userID)
	}
	guestQuery.Count(&totalGuests)

	vendorQuery := db.Model(&models.Vendor{})
	if userRole != "admin" {
		vendorQuery = vendorQuery.Where("user_id = ?", userID)
	}
	vendorQuery.Count(&totalVendors)

	budgetQuery := db.Model(&models.BudgetItem{})
	if weddingQuery != nil {
		budgetQuery = budgetQuery.Where("wedding_id = ?", weddingQuery.ID)
	} else if userRole != "admin" {
		budgetQuery = budgetQuery.Joins("JOIN weddings ON weddings.id = budget_items.wedding_id").
			Where("weddings.user_id = ?", userID)
	}
	budgetQuery.Select("COALESCE(SUM(actual_cost), 0)").Scan(&totalBudget)
	budgetQuery.Select("COALESCE(SUM(paid_amount), 0)").Scan(&totalPaid)

	var rsvpStats struct {
		Accepted int64
		Declined int64
		Pending  int64
	}
	rsvpBaseQuery := db.Model(&models.Guest{})
	if weddingQuery != nil {
		rsvpBaseQuery = rsvpBaseQuery.Where("wedding_id = ?", weddingQuery.ID)
	} else if userRole != "admin" {
		rsvpBaseQuery = rsvpBaseQuery.Joins("JOIN weddings ON weddings.id = guests.wedding_id").
			Where("weddings.user_id = ?", userID)
	}
	rsvpBaseQuery.Where("rsvp_status = ?", "accepted").Count(&rsvpStats.Accepted)
	rsvpBaseQuery.Where("rsvp_status = ?", "declined").Count(&rsvpStats.Declined)
	rsvpBaseQuery.Where("rsvp_status = ?", "pending").Count(&rsvpStats.Pending)

	var taskStats struct {
		Total     int64
		Completed int64
		Pending   int64
		Overdue   int64
	}
	taskQuery := db.Model(&models.Task{})
	if weddingQuery != nil {
		taskQuery = taskQuery.Where("wedding_id = ?", weddingQuery.ID)
	} else if userRole != "admin" {
		taskQuery = taskQuery.Joins("JOIN weddings ON weddings.id = tasks.wedding_id").
			Where("weddings.user_id = ?", userID)
	}
	taskQuery.Count(&taskStats.Total)
	taskQuery.Where("status = ?", "completed").Count(&taskStats.Completed)
	taskQuery.Where("status = ?", "pending").Count(&taskStats.Pending)
	taskQuery.Where("status != ? AND due_date < ?", "completed", time.Now()).Count(&taskStats.Overdue)

	response.Success(c, gin.H{
		"total_weddings": totalWeddings,
		"total_guests":   totalGuests,
		"total_vendors":  totalVendors,
		"budget": gin.H{
			"total": totalBudget,
			"paid":  totalPaid,
			"rate":  calculateRate(totalPaid, totalBudget),
		},
		"rsvp": rsvpStats,
		"tasks": taskStats,
	})
}

func (h *DashboardHandler) GetBudgetChart(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	weddingID := c.Query("wedding_id")

	db := database.GetDB()

	type CategoryBudget struct {
		Category      string  `json:"category"`
		EstimatedCost float64 `json:"estimated_cost"`
		ActualCost    float64 `json:"actual_cost"`
		PaidAmount    float64 `json:"paid_amount"`
	}

	var results []CategoryBudget

	query := db.Model(&models.BudgetItem{}).
		Select("category, COALESCE(SUM(estimated_cost), 0) as estimated_cost, COALESCE(SUM(actual_cost), 0) as actual_cost, COALESCE(SUM(paid_amount), 0) as paid_amount").
		Group("category")

	if weddingID != "" {
		query = query.Where("wedding_id = ?", weddingID)
	} else if userRole != "admin" {
		query = query.Joins("JOIN weddings ON weddings.id = budget_items.wedding_id").
			Where("weddings.user_id = ?", userID)
	}

	query.Find(&results)

	response.Success(c, results)
}

func (h *DashboardHandler) GetTaskProgress(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	weddingID := c.Query("wedding_id")

	db := database.GetDB()

	type TaskProgress struct {
		Category string `json:"category"`
		Total    int64  `json:"total"`
		Complete int64  `json:"complete"`
	}

	var results []TaskProgress

	query := db.Model(&models.Task{}).
		Select("category, COUNT(*) as total, SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as complete").
		Group("category")

	if weddingID != "" {
		query = query.Where("wedding_id = ?", weddingID)
	} else if userRole != "admin" {
		query = query.Joins("JOIN weddings ON weddings.id = tasks.wedding_id").
			Where("weddings.user_id = ?", userID)
	}

	query.Find(&results)

	response.Success(c, results)
}

func (h *DashboardHandler) GetUpcomingTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	weddingID := c.Query("wedding_id")

	db := database.GetDB()

	var tasks []models.Task

	query := db.Where("status != ?", "completed").
		Where("due_date IS NOT NULL").
		Where("due_date >= ?", time.Now()).
		Order("due_date ASC").
		Limit(10)

	if weddingID != "" {
		query = query.Where("wedding_id = ?", weddingID)
	} else if userRole != "admin" {
		query = query.Joins("JOIN weddings ON weddings.id = tasks.wedding_id").
			Where("weddings.user_id = ?", userID)
	}

	query.Find(&tasks)

	response.Success(c, tasks)
}

func (h *DashboardHandler) GetVendorStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()

	type VendorCategoryStat struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	var results []VendorCategoryStat

	db.Model(&models.Vendor{}).
		Select("category, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("category").
		Find(&results)

	response.Success(c, results)
}

func (h *DashboardHandler) ExportReport(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	weddingID := c.Query("wedding_id")

	db := database.GetDB()

	f := excelize.NewFile()

	var weddingQuery string
	if weddingID != "" {
		weddingQuery = weddingID
	}

	sheet1 := "预算报告"
	f.NewSheet(sheet1)

	var budgetItems []models.BudgetItem
	budgetQuery := db.Where("1 = 1")
	if weddingQuery != "" {
		budgetQuery = budgetQuery.Where("wedding_id = ?", weddingQuery)
	} else if userRole != "admin" {
		budgetQuery = budgetQuery.Joins("JOIN weddings ON weddings.id = budget_items.wedding_id").
			Where("weddings.user_id = ?", userID)
	}
	budgetQuery.Find(&budgetItems)

	f.SetCellValue(sheet1, "A1", "预算报告")
	f.SetCellValue(sheet1, "A3", "类别")
	f.SetCellValue(sheet1, "B3", "描述")
	f.SetCellValue(sheet1, "C3", "预估费用")
	f.SetCellValue(sheet1, "D3", "实际费用")
	f.SetCellValue(sheet1, "E3", "已付金额")
	f.SetCellValue(sheet1, "F3", "状态")

	var totalEstimated, totalActual, totalPaid float64
	for i, item := range budgetItems {
		row := i + 4
		f.SetCellValue(sheet1, "A"+toStr(row), item.Category)
		f.SetCellValue(sheet1, "B"+toStr(row), item.Description)
		f.SetCellValue(sheet1, "C"+toStr(row), item.EstimatedCost)
		f.SetCellValue(sheet1, "D"+toStr(row), item.ActualCost)
		f.SetCellValue(sheet1, "E"+toStr(row), item.PaidAmount)
		f.SetCellValue(sheet1, "F"+toStr(row), item.Status)
		totalEstimated += item.EstimatedCost
		totalActual += item.ActualCost
		totalPaid += item.PaidAmount
	}

	lastRow := len(budgetItems) + 4
	f.SetCellValue(sheet1, "A"+toStr(lastRow), "总计")
	f.SetCellValue(sheet1, "C"+toStr(lastRow), totalEstimated)
	f.SetCellValue(sheet1, "D"+toStr(lastRow), totalActual)
	f.SetCellValue(sheet1, "E"+toStr(lastRow), totalPaid)

	sheet2 := "任务清单"
	f.NewSheet(sheet2)

	var tasks []models.Task
	taskQuery := db.Where("1 = 1")
	if weddingQuery != "" {
		taskQuery = taskQuery.Where("wedding_id = ?", weddingQuery)
	} else if userRole != "admin" {
		taskQuery = taskQuery.Joins("JOIN weddings ON weddings.id = tasks.wedding_id").
			Where("weddings.user_id = ?", userID)
	}
	taskQuery.Find(&tasks)

	f.SetCellValue(sheet2, "A1", "任务清单")
	f.SetCellValue(sheet2, "A3", "标题")
	f.SetCellValue(sheet2, "B3", "分类")
	f.SetCellValue(sheet2, "C3", "负责人")
	f.SetCellValue(sheet2, "D3", "截止日期")
	f.SetCellValue(sheet2, "E3", "优先级")
	f.SetCellValue(sheet2, "F3", "状态")

	for i, task := range tasks {
		row := i + 4
		f.SetCellValue(sheet2, "A"+toStr(row), task.Title)
		f.SetCellValue(sheet2, "B"+toStr(row), task.Category)
		f.SetCellValue(sheet2, "C"+toStr(row), task.Assignee)
		if task.DueDate != nil {
			f.SetCellValue(sheet2, "D"+toStr(row), task.DueDate.Format("2006-01-02"))
		}
		f.SetCellValue(sheet2, "E"+toStr(row), task.Priority)
		f.SetCellValue(sheet2, "F"+toStr(row), task.Status)
	}

	sheet3 := "嘉宾名单"
	f.NewSheet(sheet3)

	var guests []models.Guest
	guestQuery := db.Where("1 = 1")
	if weddingQuery != "" {
		guestQuery = guestQuery.Where("wedding_id = ?", weddingQuery)
	} else if userRole != "admin" {
		guestQuery = guestQuery.Joins("JOIN weddings ON weddings.id = guests.wedding_id").
			Where("weddings.user_id = ?", userID)
	}
	guestQuery.Find(&guests)

	f.SetCellValue(sheet3, "A1", "嘉宾名单")
	f.SetCellValue(sheet3, "A3", "姓名")
	f.SetCellValue(sheet3, "B3", "邮箱")
	f.SetCellValue(sheet3, "C3", "电话")
	f.SetCellValue(sheet3, "D3", "分组")
	f.SetCellValue(sheet3, "E3", "RSVP状态")

	for i, guest := range guests {
		row := i + 4
		f.SetCellValue(sheet3, "A"+toStr(row), guest.FullName)
		f.SetCellValue(sheet3, "B"+toStr(row), guest.Email)
		f.SetCellValue(sheet3, "C"+toStr(row), guest.Phone)
		f.SetCellValue(sheet3, "D"+toStr(row), guest.Group)
		f.SetCellValue(sheet3, "E"+toStr(row), guest.RSVPStatus)
	}

	f.DeleteSheet("Sheet1")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=wedding_report.xlsx")
	f.Write(c.Writer)
}

type WeddingQuery struct {
	ID string
}

func calculateRate(paid, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (paid / total) * 100
}

func toStr(n int) string {
	return intToStr(n)
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
