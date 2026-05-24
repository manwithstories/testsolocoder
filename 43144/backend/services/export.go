package services

import (
	"bytes"
	"embed"
	"fmt"
	"time"

	"pet-adoption-platform/database"
	"pet-adoption-platform/models"

	"github.com/go-pdf/fpdf"
	"github.com/xuri/excelize/v2"
)

//go:embed assets/fonts/NotoSansSC-Regular.ttf
var fontFS embed.FS

func GenerateAdoptionReport(rescueID uint, startDate, endDate string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	f.SetCellValue(sheet, "A1", "宠物领养统计报表")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))
	if startDate != "" || endDate != "" {
		f.SetCellValue(sheet, "A3", fmt.Sprintf("统计周期: %s 至 %s", startDate, endDate))
	}

	f.SetCellValue(sheet, "A5", "编号")
	f.SetCellValue(sheet, "B5", "宠物名称")
	f.SetCellValue(sheet, "C5", "品种")
	f.SetCellValue(sheet, "D5", "性别")
	f.SetCellValue(sheet, "E5", "年龄")
	f.SetCellValue(sheet, "F5", "状态")
	f.SetCellValue(sheet, "G5", "领养人")
	f.SetCellValue(sheet, "H5", "申请时间")
	f.SetCellValue(sheet, "I5", "审批时间")
	f.SetCellValue(sheet, "J5", "签署时间")

	db := database.DB.Model(&models.AdoptionApplication{}).
		Preload("Pet").Preload("Adopter")

	if rescueID > 0 {
		db = db.Where("rescue_id = ?", rescueID)
	}
	if startDate != "" {
		sd, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			db = db.Where("created_at >= ?", sd)
		}
	}
	if endDate != "" {
		ed, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			db = db.Where("created_at <= ?", ed)
		}
	}

	var applications []models.AdoptionApplication
	if err := db.Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, err
	}

	for i, app := range applications {
		row := i + 6
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), app.ID)
		if app.Pet != nil {
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), app.Pet.Name)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), app.Pet.Breed)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), app.Pet.Gender)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), app.Pet.Age)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), app.Pet.Status)
		}
		if app.Adopter != nil {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), app.Adopter.Name)
		}
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), app.CreatedAt.Format("2006-01-02 15:04:05"))
		if app.ReviewedAt != nil {
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), app.ReviewedAt.Format("2006-01-02 15:04:05"))
		}
		if app.SignedAt != nil {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), app.SignedAt.Format("2006-01-02 15:04:05"))
		}
	}

	return f, nil
}

func GenerateHealthReport(petID uint) (*bytes.Buffer, string, error) {
	pet, err := GetPetByID(petID)
	if err != nil {
		return nil, "", err
	}

	records, _, err := ListHealthRecords(&models.HealthRecordListQuery{
		Page:     1,
		PageSize: 1000,
		PetID:    petID,
	})
	if err != nil {
		return nil, "", err
	}

	followUps, err := ListFollowUpRecords(petID)
	if err != nil {
		return nil, "", err
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 20)

	fontData, err := fontFS.ReadFile("assets/fonts/NotoSansSC-Regular.ttf")
	if err != nil {
		return nil, "", fmt.Errorf("failed to read font: %w", err)
	}

	pdf.AddUTF8FontFromBytes("noto", "", fontData)
	pdf.SetFont("noto", "", 16)
	pdf.SetTextColor(24, 144, 255)
	pdf.CellFormat(0, 10, fmt.Sprintf("宠物健康档案 - %s", pet.Name), "", 1, "C", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("noto", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, fmt.Sprintf("档案编号: %s", pet.ArchiveNumber), "", 1, "L", false, 0, "")
	pdf.SetFont("noto", "", 10)
	pdf.SetTextColor(102, 102, 102)
	pdf.CellFormat(0, 6, fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "L", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("noto", "", 14)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 8, "宠物基本信息", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("noto", "", 10)
	pdf.SetTextColor(0, 0, 0)
	infoItems := [][]string{
		{"名称", pet.Name},
		{"品种", pet.Breed},
		{"年龄", fmt.Sprintf("%d", pet.Age)},
		{"性别", string(pet.Gender)},
		{"健康状况", pet.HealthStatus},
	}
	for _, item := range infoItems {
		pdf.SetFillColor(245, 245, 245)
		pdf.CellFormat(30, 8, item[0], "1", 0, "L", true, 0, "")
		pdf.SetFillColor(255, 255, 255)
		pdf.CellFormat(0, 8, item[1], "1", 1, "L", true, 0, "")
	}

	pdf.Ln(4)
	pdf.SetFont("noto", "", 14)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 8, "健康记录", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("noto", "", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(245, 245, 245)
	headers := []string{"日期", "类型", "标题", "疫苗名称", "兽医", "医院", "备注"}
	colWidths := []float64{22, 18, 30, 25, 20, 25, 0}
	for i, h := range headers {
		w := colWidths[i]
		if w == 0 {
			w = 55
		}
		pdf.CellFormat(w, 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("noto", "", 7)
	pdf.SetFillColor(255, 255, 255)
	for _, record := range records {
		vals := []string{
			record.RecordDate.Format("2006-01-02"),
			string(record.RecordType),
			truncate(record.Title, 10),
			truncate(record.VaccineName, 8),
			truncate(record.VetName, 6),
			truncate(record.Hospital, 8),
			truncate(record.Notes, 12),
		}
		for i, v := range vals {
			w := colWidths[i]
			if w == 0 {
				w = 55
			}
			pdf.CellFormat(w, 6, v, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
	}

	pdf.Ln(4)
	pdf.SetFont("noto", "", 14)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 8, "回访记录", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("noto", "", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(245, 245, 245)
	fuHeaders := []string{"回访日期", "健康状况", "生活环境", "备注"}
	fuWidths := []float64{30, 30, 40, 70}
	for i, h := range fuHeaders {
		pdf.CellFormat(fuWidths[i], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("noto", "", 8)
	pdf.SetFillColor(255, 255, 255)
	for _, fu := range followUps {
		var dateStr string
		if fu.FollowUpDate != nil {
			dateStr = fu.FollowUpDate.Format("2006-01-02")
		}
		vals := []string{
			dateStr,
			fu.HealthStatus,
			truncate(fu.LivingCondition, 15),
			truncate(fu.Notes, 25),
		}
		for i, v := range vals {
			pdf.CellFormat(fuWidths[i], 6, v, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("health_report_%s_%s.pdf", pet.ArchiveNumber, time.Now().Format("20060102_150405"))
	return &buf, filename, nil
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n]) + "..."
	}
	return s
}
