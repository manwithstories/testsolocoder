package handlers

import (
	"bytes"
	"fmt"
	"smart-energy-platform/models"
	"smart-energy-platform/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/tealeg/xlsx/v3"
)

type EnergySummary struct {
	TotalEnergy   float64
	DeviceCount   int
	TopDevices    []DeviceEnergyDetail
	AvgDailyUsage float64
	SavingTips    []string
}

type DeviceEnergyDetail struct {
	DeviceName  string
	DeviceType  string
	TotalEnergy float64
	Percentage  float64
	Power       float64
}

func ExportEnergyReport(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")
	format := c.DefaultQuery("format", "excel")
	period := c.DefaultQuery("period", "month")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var familyIDs []uint
	if familyID > 0 {
		familyIDs = []uint{familyID}
	} else {
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
	}

	if len(familyIDs) == 0 {
		utils.BadRequest(c, "No family found")
		return
	}

	var startTime time.Time
	now := time.Now()
	switch period {
	case "week":
		startTime = now.Add(-7 * 24 * time.Hour)
	case "month":
		startTime = now.Add(-30 * 24 * time.Hour)
	default:
		startTime = now.Add(-7 * 24 * time.Hour)
	}

	summary := generateEnergySummary(familyIDs, startTime, now)

	if format == "pdf" {
		exportPDF(c, summary, period)
	} else {
		exportExcel(c, summary, familyIDs, startTime, now, period)
	}
}

func generateEnergySummary(familyIDs []uint, startTime, endTime time.Time) EnergySummary {
	var totalEnergy float64
	models.DB.Model(&models.EnergyData{}).
		Where("family_id IN ? AND timestamp >= ? AND timestamp <= ?", familyIDs, startTime, endTime).
		Select("COALESCE(SUM(energy_used), 0)").
		Scan(&totalEnergy)

	var deviceCount int64
	models.DB.Model(&models.Device{}).
		Where("family_id IN ?", familyIDs).
		Count(&deviceCount)

	type result struct {
		DeviceID    uint
		TotalEnergy float64
	}
	var results []result
	models.DB.Model(&models.EnergyData{}).
		Select("device_id, COALESCE(SUM(energy_used), 0) as total_energy").
		Where("family_id IN ? AND timestamp >= ? AND timestamp <= ?", familyIDs, startTime, endTime).
		Group("device_id").
		Order("total_energy DESC").
		Limit(10).
		Scan(&results)

	var topDevices []DeviceEnergyDetail
	for _, r := range results {
		var device models.Device
		models.DB.First(&device, r.DeviceID)
		pct := 0.0
		if totalEnergy > 0 {
			pct = r.TotalEnergy / totalEnergy * 100
		}
		topDevices = append(topDevices, DeviceEnergyDetail{
			DeviceName:  device.Name,
			DeviceType:  device.DeviceType,
			TotalEnergy: r.TotalEnergy,
			Percentage:  pct,
			Power:       device.Power,
		})
	}

	days := endTime.Sub(startTime).Hours() / 24
	avgDaily := totalEnergy
	if days > 0 {
		avgDaily = totalEnergy / days
	}

	tips := generateSavingTips(topDevices, totalEnergy)

	return EnergySummary{
		TotalEnergy:   totalEnergy,
		DeviceCount:   int(deviceCount),
		TopDevices:    topDevices,
		AvgDailyUsage: avgDaily,
		SavingTips:    tips,
	}
}

func generateSavingTips(topDevices []DeviceEnergyDetail, totalEnergy float64) []string {
	var tips []string

	if len(topDevices) > 0 && topDevices[0].Percentage > 30 {
		tips = append(tips, fmt.Sprintf("'%s' 占总能耗的 %.1f%%，建议优化其使用时间。", topDevices[0].DeviceName, topDevices[0].Percentage))
	}

	highPowerDevices := 0
	for _, d := range topDevices {
		if d.Power > 1000 {
			highPowerDevices++
		}
	}
	if highPowerDevices > 2 {
		tips = append(tips, "检测到多个大功率设备(>1000W)，建议在非高峰时段使用以降低电费。")
	}

	if totalEnergy > 30 {
		tips = append(tips, "月能耗较高，建议创建节能场景来自动管理设备。")
	}

	if len(tips) == 0 {
		tips = append(tips, "您的能耗处于正常范围，继续保持良好的用电习惯。")
	}

	return tips
}

func exportExcel(c *gin.Context, summary EnergySummary, familyIDs []uint, startTime, endTime time.Time, period string) {
	wb := xlsx.NewFile()

	summarySheet, _ := wb.AddSheet("能耗摘要")
	addRow(summarySheet, "能源消耗报告")
	addRow(summarySheet, fmt.Sprintf("统计周期: %s 至 %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")))
	addRow(summarySheet, "")
	addRow(summarySheet, "总能耗 (kWh)", fmt.Sprintf("%.2f", summary.TotalEnergy))
	addRow(summarySheet, "设备数量", fmt.Sprintf("%d", summary.DeviceCount))
	addRow(summarySheet, "日均能耗 (kWh)", fmt.Sprintf("%.2f", summary.AvgDailyUsage))
	addRow(summarySheet, "")
	addRow(summarySheet, "高能耗设备排名")
	addRow(summarySheet, "排名", "设备名称", "类型", "能耗 (kWh)", "占比 (%)")
	for i, d := range summary.TopDevices {
		addRow(summarySheet, fmt.Sprintf("%d", i+1), d.DeviceName, d.DeviceType,
			fmt.Sprintf("%.2f", d.TotalEnergy), fmt.Sprintf("%.1f", d.Percentage))
	}
	addRow(summarySheet, "")
	addRow(summarySheet, "节能建议")
	for i, tip := range summary.SavingTips {
		addRow(summarySheet, fmt.Sprintf("%d. %s", i+1, tip))
	}

	detailSheet, _ := wb.AddSheet("每日明细")
	addRow(detailSheet, "日期", "设备", "能耗 (kWh)", "峰值功率 (W)")

	type dailyDetail struct {
		Date       string
		DeviceID   uint
		EnergyUsed float64
		PeakPower  float64
	}
	var details []dailyDetail
	models.DB.Model(&models.EnergyData{}).
		Select("date, device_id, COALESCE(SUM(energy_used), 0) as energy_used, COALESCE(MAX(power), 0) as peak_power").
		Where("family_id IN ? AND timestamp >= ? AND timestamp <= ?", familyIDs, startTime, endTime).
		Group("date, device_id").
		Order("date ASC").
		Scan(&details)

	for _, d := range details {
		var device models.Device
		models.DB.First(&device, d.DeviceID)
		addRow(detailSheet, d.Date, device.Name,
			fmt.Sprintf("%.2f", d.EnergyUsed), fmt.Sprintf("%.1f", d.PeakPower))
	}

	var buf bytes.Buffer
	wb.Write(&buf)

	filename := fmt.Sprintf("energy_report_%s_%s.xlsx", period, time.Now().Format("20060102"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func addRow(sheet *xlsx.Sheet, values ...string) {
	row := sheet.AddRow()
	for _, v := range values {
		cell := row.AddCell()
		cell.Value = v
	}
}

func exportPDF(c *gin.Context, summary EnergySummary, period string) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	fontPath := "fonts/NotoSansSC-Regular.ttf"
	err := pdf.AddTTFFont("noto", fontPath)
	if err != nil {
		fontPath = "../fonts/NotoSansSC-Regular.ttf"
		err = pdf.AddTTFFont("noto", fontPath)
	}
	if err != nil {
		fontPath = "../../fonts/NotoSansSC-Regular.ttf"
		err = pdf.AddTTFFont("noto", fontPath)
	}

	fontFamily := "Helvetica"
	if err == nil {
		fontFamily = "noto"
	}

	pdf.SetXY(50, 50)
	pdf.SetFont(fontFamily, "B", 18)
	pdf.Cell(nil, "能源消耗报告")

	pdf.SetXY(50, 70)
	pdf.SetFont(fontFamily, "", 12)
	periodLabel := period
	if period == "week" {
		periodLabel = "本周"
	} else if period == "month" {
		periodLabel = "本月"
	}
	pdf.Cell(nil, fmt.Sprintf("统计周期: %s", periodLabel))
	pdf.SetXY(50, 85)
	pdf.Cell(nil, fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))

	pdf.SetXY(50, 110)
	pdf.SetFont(fontFamily, "B", 14)
	pdf.Cell(nil, "能耗摘要")

	pdf.SetXY(50, 130)
	pdf.SetFont(fontFamily, "", 12)
	pdf.Cell(nil, fmt.Sprintf("总能耗: %.2f kWh", summary.TotalEnergy))
	pdf.SetXY(50, 145)
	pdf.Cell(nil, fmt.Sprintf("设备数量: %d", summary.DeviceCount))
	pdf.SetXY(50, 160)
	pdf.Cell(nil, fmt.Sprintf("日均能耗: %.2f kWh", summary.AvgDailyUsage))

	pdf.SetXY(50, 190)
	pdf.SetFont(fontFamily, "B", 14)
	pdf.Cell(nil, "高能耗设备排行")

	y := 210.0
	pdf.SetXY(50, y)
	pdf.SetFont(fontFamily, "B", 10)
	pdf.Cell(nil, "排名  设备名称                  类型        能耗(kWh)    占比(%)")

	y = 225.0
	for i, d := range summary.TopDevices {
		if y > 750 {
			pdf.AddPage()
			y = 50
		}
		pdf.SetXY(50, y)
		pdf.SetFont(fontFamily, "", 10)
		pdf.Cell(nil, fmt.Sprintf("%-4s %-24s %-12s %-12s %-12s",
			fmt.Sprintf("%d.", i+1),
			truncateString(d.DeviceName, 22),
			d.DeviceType,
			fmt.Sprintf("%.2f", d.TotalEnergy),
			fmt.Sprintf("%.1f%%", d.Percentage)))
		y += 15
	}

	y += 20
	if y > 750 {
		pdf.AddPage()
		y = 50
	}
	pdf.SetXY(50, y)
	pdf.SetFont(fontFamily, "B", 14)
	pdf.Cell(nil, "节能建议")

	y += 20
	for i, tip := range summary.SavingTips {
		if y > 780 {
			pdf.AddPage()
			y = 50
		}
		pdf.SetXY(50, y)
		pdf.SetFont(fontFamily, "", 10)
		pdf.Cell(nil, fmt.Sprintf("%d. %s", i+1, tip))
		y += 15
	}

	var buf bytes.Buffer
	pdf.Write(&buf)

	filename := fmt.Sprintf("energy_report_%s_%s.pdf", period, time.Now().Format("20060102"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/pdf", buf.Bytes())
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
