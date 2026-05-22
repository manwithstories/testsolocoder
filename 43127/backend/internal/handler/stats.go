package handler

import (
	"bytes"
	"fmt"
	"property-management/internal/database"
	"property-management/internal/model"
	"property-management/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

func (h *StatsHandler) GetOverview(c *gin.Context) {
	var totalProperties int64
	database.DB.Model(&model.Property{}).Count(&totalProperties)

	var rentedProperties int64
	database.DB.Model(&model.Property{}).Where("status = ?", 2).Count(&rentedProperties)

	occupancyRate := float64(0)
	if totalProperties > 0 {
		occupancyRate = float64(rentedProperties) / float64(totalProperties) * 100
	}

	var totalIncome float64
	startOfMonth := time.Now().Format("2006-01")
	database.DB.Model(&model.RentRecord{}).
		Where("month = ? AND status = ?", startOfMonth, 1).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	var pendingRepairs int64
	database.DB.Model(&model.RepairOrder{}).Where("status IN ?", []int{1, 2}).Count(&pendingRepairs)

	var activeContracts int64
	database.DB.Model(&model.Contract{}).Where("status = ?", 1).Count(&activeContracts)

	utils.Success(c, gin.H{
		"totalProperties":  totalProperties,
		"rentedProperties": rentedProperties,
		"occupancyRate":     fmt.Sprintf("%.1f", occupancyRate),
		"totalIncome":       totalIncome,
		"pendingRepairs":    pendingRepairs,
		"activeContracts":   activeContracts,
	})
}

func (h *StatsHandler) GetOccupancyTrend(c *gin.Context) {
	months := []string{}
	rates := []float64{}

	for i := 5; i >= 0; i-- {
		date := time.Now().AddDate(0, -i, 0)
		month := date.Format("2006-01")
		months = append(months, month)

		var total int64
		var rented int64
		database.DB.Model(&model.Property{}).Count(&total)
		database.DB.Model(&model.Contract{}).
			Where("status = ? AND start_date <= ? AND end_date >= ?", 1, month+"-31", month+"-01").
			Count(&rented)

		rate := float64(0)
		if total > 0 {
			rate = float64(rented) / float64(total) * 100
		}
		rates = append(rates, rate)
	}

	utils.Success(c, gin.H{
		"months": months,
		"rates":  rates,
	})
}

func (h *StatsHandler) GetIncomeTrend(c *gin.Context) {
	months := []string{}
	incomes := []float64{}

	for i := 5; i >= 0; i-- {
		date := time.Now().AddDate(0, -i, 0)
		month := date.Format("2006-01")
		months = append(months, month)

		var income float64
		database.DB.Model(&model.RentRecord{}).
			Where("month = ? AND status = ?", month, 1).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&income)
		incomes = append(incomes, income)
	}

	utils.Success(c, gin.H{
		"months":  months,
		"incomes": incomes,
	})
}

func (h *StatsHandler) GetRepairStats(c *gin.Context) {
	var stats []model.RepairOrder
	database.DB.Find(&stats)

	byCategory := map[string]int{}
	byStatus := map[string]int{}

	for _, order := range stats {
		category := order.Category
		if category == "" {
			category = "其他"
		}
		byCategory[category]++
		byStatus[fmt.Sprintf("%d", order.Status)]++
	}

	utils.Success(c, gin.H{
		"byCategory": byCategory,
		"byStatus":   byStatus,
		"total":      len(stats),
	})
}

func (h *StatsHandler) ExportRentRecords(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	var records []model.RentRecord
	query := database.DB.Preload("Contract").Preload("Contract.Tenant").Preload("Contract.Property")
	if month != "all" {
		query = query.Where("month = ?", month)
	}
	query.Find(&records)

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "月份")
	f.SetCellValue(sheet, "B1", "房源")
	f.SetCellValue(sheet, "C1", "租户")
	f.SetCellValue(sheet, "D1", "金额")
	f.SetCellValue(sheet, "E1", "状态")
	f.SetCellValue(sheet, "F1", "应缴日期")
	f.SetCellValue(sheet, "G1", "实缴日期")
	f.SetCellValue(sheet, "H1", "滞纳金")

	for i, record := range records {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), record.Month)
		if record.Contract != nil {
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), record.Contract.Property.Title)
			if record.Contract.Tenant != nil {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), record.Contract.Tenant.Name)
			}
		}
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), record.Amount)
		status := "未缴"
		if record.Status == 1 {
			status = "已缴"
		}
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), status)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), record.DueDate.Format("2006-01-02"))
		if record.PaidAt != nil {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), record.PaidAt.Format("2006-01-02"))
		}
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), record.LateFee)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		utils.Error(c, 500, "Failed to export")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=rent_records_%s.xlsx", month))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func (h *StatsHandler) ExportRepairOrders(c *gin.Context) {
	var orders []model.RepairOrder
	database.DB.Preload("Tenant").Preload("Property").Preload("Handler").Find(&orders)

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "标题")
	f.SetCellValue(sheet, "B1", "分类")
	f.SetCellValue(sheet, "C1", "房源")
	f.SetCellValue(sheet, "D1", "租户")
	f.SetCellValue(sheet, "E1", "处理人")
	f.SetCellValue(sheet, "F1", "优先级")
	f.SetCellValue(sheet, "G1", "状态")
	f.SetCellValue(sheet, "H1", "创建时间")
	f.SetCellValue(sheet, "I1", "完成时间")

	for i, order := range orders {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), order.Title)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), order.Category)
		if order.Property != nil {
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), order.Property.Title)
		}
		if order.Tenant != nil {
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), order.Tenant.Name)
		}
		if order.Handler != nil {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), order.Handler.RealName)
		}
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), order.Priority)
		statusText := "待处理"
		switch order.Status {
		case 2:
			statusText = "处理中"
		case 3:
			statusText = "已完成"
		case 4:
			statusText = "已关闭"
		}
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), statusText)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), order.CreatedAt.Format("2006-01-02 15:04:05"))
		if order.CompletedAt != nil {
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), order.CompletedAt.Format("2006-01-02 15:04:05"))
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		utils.Error(c, 500, "Failed to export")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=repair_orders.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
