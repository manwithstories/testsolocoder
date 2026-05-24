package utils

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	"temp-staff-platform/config"
	"temp-staff-platform/models"
)

func ExportSchedulesToExcel(schedules []models.Schedule, filename string) (string, error) {
	f := excelize.NewFile()
	sheet := "Schedules"

	headers := []string{"Date", "Time", "Location", "Position", "Employee", "Status"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#4472C4"}},
	})
	f.SetCellStyle(sheet, "A1", "F1", style)

	for i, s := range schedules {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), s.ShiftDate.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s - %s", s.StartTime, s.EndTime))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), s.Location)
		if s.JobPosting != nil {
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.JobPosting.Position)
		}
		if s.Temporary != nil {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), s.Temporary.RealName)
		}
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), s.Status)
	}

	for col := 'A'; col <= 'F'; col++ {
		f.SetColWidth(sheet, string(col), string(col), 20)
	}

	dir := config.AppConfig.UploadDir + "/exports"
	os.MkdirAll(dir, 0755)
	filepath := dir + "/" + filename

	if err := f.SaveAs(filepath); err != nil {
		return "", fmt.Errorf("failed to save excel: %v", err)
	}

	return filepath, nil
}

func ExportSalaryToPDF(salary *models.SalaryRecord, details []models.SalaryDetail, filename string) (string, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(nil, "Salary Slip")
	pdf.Br(20)

	pdf.SetFont("Helvetica", "", 10)
	if salary.Temporary != nil {
		pdf.Cell(nil, fmt.Sprintf("Employee: %s", salary.Temporary.RealName))
		pdf.Br(10)
	}
	pdf.Cell(nil, fmt.Sprintf("Period: %s to %s", salary.PeriodStart.Format("2006-01-02"), salary.PeriodEnd.Format("2006-01-02")))
	pdf.Br(15)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(nil, "Details:")
	pdf.Br(10)

	for _, d := range details {
		pdf.SetFont("Helvetica", "", 9)
		pdf.Cell(nil, fmt.Sprintf("%s | %s | %.2f hrs | ¥%.2f", d.Date.Format("2006-01-02"), d.Type, d.WorkHours, d.Amount))
		pdf.Br(8)
	}

	pdf.Br(10)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(nil, fmt.Sprintf("Total: ¥%.2f", salary.TotalSalary))

	dir := config.AppConfig.UploadDir + "/exports"
	os.MkdirAll(dir, 0755)
	filepath := dir + "/" + filename

	if err := pdf.WritePdf(filepath); err != nil {
		return "", fmt.Errorf("failed to save pdf: %v", err)
	}

	return filepath, nil
}

func ServeFile(c *gin.Context, filepath string, filename string) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}
