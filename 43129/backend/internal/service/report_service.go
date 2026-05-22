package service

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"beauty-salon-system/internal/repository"

	"github.com/xuri/excelize/v2"
	"github.com/jung-kurt/gofpdf/v2"
)

type ReportService struct {
	paymentRepo      *repository.PaymentRepository
	appointmentRepo  *repository.AppointmentRepository
	serviceRepo      *repository.ServiceRepository
	technicianRepo   *repository.TechnicianRepository
}

func NewReportService(
	paymentRepo *repository.PaymentRepository,
	appointmentRepo *repository.AppointmentRepository,
	serviceRepo *repository.ServiceRepository,
	technicianRepo *repository.TechnicianRepository,
) *ReportService {
	return &ReportService{
		paymentRepo:     paymentRepo,
		appointmentRepo: appointmentRepo,
		serviceRepo:     serviceRepo,
		technicianRepo:  technicianRepo,
	}
}

type RevenueReport struct {
	TotalRevenue float64       `json:"total_revenue"`
	DateRange    string        `json:"date_range"`
	DailyData    []DailyRevenue `json:"daily_data"`
}

type DailyRevenue struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Count   int64   `json:"count"`
}

type TechnicianPerformance struct {
	TechnicianID   uint    `json:"technician_id"`
	TechnicianName string  `json:"technician_name"`
	ServiceCount   int64   `json:"service_count"`
	Revenue        float64 `json:"revenue"`
}

type ServiceRank struct {
	ServiceID   uint    `json:"service_id"`
	ServiceName string  `json:"service_name"`
	Count       int64   `json:"count"`
	Revenue     float64 `json:"revenue"`
}

type ReportData struct {
	RevenueReport         *RevenueReport          `json:"revenue_report"`
	TechnicianPerformances []TechnicianPerformance `json:"technician_performances"`
	ServiceRanks          []ServiceRank            `json:"service_ranks"`
}

func (s *ReportService) GetRevenueReport(startDate, endDate string) (*RevenueReport, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	totalRevenue, err := s.paymentRepo.GetRevenueByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	var dailyData []DailyRevenue
	currentDate := start
	for !currentDate.After(end) {
		dayEnd := currentDate
		revenue, err := s.paymentRepo.GetRevenueByDateRange(currentDate, dayEnd)
		if err != nil {
			return nil, err
		}

		dailyData = append(dailyData, DailyRevenue{
			Date:    currentDate.Format("2006-01-02"),
			Revenue: revenue,
			Count:   0,
		})

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &RevenueReport{
		TotalRevenue: totalRevenue,
		DateRange:    fmt.Sprintf("%s ~ %s", startDate, endDate),
		DailyData:    dailyData,
	}, nil
}

func (s *ReportService) GetTechnicianPerformance(startDate, endDate string) ([]TechnicianPerformance, error) {
	technicians, err := s.technicianRepo.ListAll()
	if err != nil {
		return nil, err
	}

	var performances []TechnicianPerformance
	for _, tech := range technicians {
		revenue, err := s.paymentRepo.GetRevenueByTechnician(tech.ID, parseTime(startDate), parseTime(endDate))
		if err != nil {
			return nil, err
		}

		performances = append(performances, TechnicianPerformance{
			TechnicianID:   tech.ID,
			TechnicianName: tech.Name,
			ServiceCount:   0,
			Revenue:        revenue,
		})
	}

	return performances, nil
}

func (s *ReportService) GetServiceRanking(startDate, endDate string) ([]ServiceRank, error) {
	services, err := s.serviceRepo.ListAll()
	if err != nil {
		return nil, err
	}

	var rankings []ServiceRank
	for _, svc := range services {
		revenue, err := s.paymentRepo.GetRevenueByService(svc.ID, parseTime(startDate), parseTime(endDate))
		if err != nil {
			return nil, err
		}

		rankings = append(rankings, ServiceRank{
			ServiceID:   svc.ID,
			ServiceName: svc.Name,
			Count:       0,
			Revenue:     revenue,
		})
	}

	return rankings, nil
}

func (s *ReportService) GetFullReport(startDate, endDate string) (*ReportData, error) {
	revenueReport, err := s.GetRevenueReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	techPerformance, err := s.GetTechnicianPerformance(startDate, endDate)
	if err != nil {
		return nil, err
	}

	serviceRanks, err := s.GetServiceRanking(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &ReportData{
		RevenueReport:         revenueReport,
		TechnicianPerformances: techPerformance,
		ServiceRanks:          serviceRanks,
	}, nil
}

func (s *ReportService) ExportExcel(startDate, endDate string) ([]byte, error) {
	f := excelize.NewFile()

	sheet1 := "收入报表"
	f.SetSheetName("Sheet1", sheet1)

	f.SetCellValue(sheet1, "A1", "美容美发店运营报表")
	f.SetCellValue(sheet1, "A2", fmt.Sprintf("统计时间：%s 至 %s", startDate, endDate))

	f.SetCellValue(sheet1, "A4", "日期")
	f.SetCellValue(sheet1, "B4", "收入(元)")

	revenueReport, err := s.GetRevenueReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	for i, data := range revenueReport.DailyData {
		row := i + 5
		f.SetCellValue(sheet1, fmt.Sprintf("A%d", row), data.Date)
		f.SetCellValue(sheet1, fmt.Sprintf("B%d", row), data.Revenue)
	}

	sheet2 := "技师业绩"
	f.NewSheet(sheet2)

	f.SetCellValue(sheet2, "A1", "技师名称")
	f.SetCellValue(sheet2, "B1", "业绩(元)")

	performances, err := s.GetTechnicianPerformance(startDate, endDate)
	if err != nil {
		return nil, err
	}

	for i, perf := range performances {
		row := i + 2
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), perf.TechnicianName)
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), perf.Revenue)
	}

	sheet3 := "热门服务"
	f.NewSheet(sheet3)

	f.SetCellValue(sheet3, "A1", "服务名称")
	f.SetCellValue(sheet3, "B1", "收入(元)")

	rankings, err := s.GetServiceRanking(startDate, endDate)
	if err != nil {
		return nil, err
	}

	for i, rank := range rankings {
		row := i + 2
		f.SetCellValue(sheet3, fmt.Sprintf("A%d", row), rank.ServiceName)
		f.SetCellValue(sheet3, fmt.Sprintf("B%d", row), rank.Revenue)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *ReportService) ExportPDF(startDate, endDate string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	fontPaths := []string{
		"/System/Library/Fonts/PingFang.ttc",
		"/System/Library/Fonts/STHeiti Light.ttc",
		"/System/Library/Fonts/Hiragino Sans GB.ttc",
		"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
		"/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
		"C:\\Windows\\Fonts\\msyh.ttc",
		"C:\\Windows\\Fonts\\simhei.ttf",
	}

	useChineseFont := false
	for _, fontPath := range fontPaths {
		if fileExists(fontPath) {
			pdf.AddUTF8Font("ChineseFont", "", fontPath)
			useChineseFont = true
			break
		}
	}

	pdf.AddPage()

	if useChineseFont {
		pdf.SetFont("ChineseFont", "", 16)
	} else {
		pdf.SetFont("Arial", "B", 16)
	}

	if useChineseFont {
		pdf.Cell(40, 10, "美容美发店运营报表")
	} else {
		pdf.Cell(40, 10, "Beauty Salon Report")
	}
	pdf.Ln(10)

	if useChineseFont {
		pdf.SetFont("ChineseFont", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("统计时间：%s 至 %s", startDate, endDate))
	} else {
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Date Range: %s to %s", startDate, endDate))
	}
	pdf.Ln(15)

	revenueReport, err := s.GetRevenueReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if useChineseFont {
		pdf.SetFont("ChineseFont", "B", 12)
		pdf.Cell(40, 10, fmt.Sprintf("总收入：%.2f元", revenueReport.TotalRevenue))
	} else {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Total Revenue: %.2f", revenueReport.TotalRevenue))
	}
	pdf.Ln(10)

	if useChineseFont {
		pdf.SetFont("ChineseFont", "B", 10)
		pdf.Cell(40, 8, "日期")
		pdf.Cell(40, 8, "收入(元)")
	} else {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 8, "Date")
		pdf.Cell(40, 8, "Revenue")
	}
	pdf.Ln(8)

	if useChineseFont {
		pdf.SetFont("ChineseFont", "", 10)
	} else {
		pdf.SetFont("Arial", "", 10)
	}
	for _, data := range revenueReport.DailyData {
		pdf.Cell(40, 6, data.Date)
		pdf.Cell(40, 6, fmt.Sprintf("%.2f", data.Revenue))
		pdf.Ln(6)
	}

	pdf.AddPage()

	if useChineseFont {
		pdf.SetFont("ChineseFont", "B", 12)
		pdf.Cell(40, 10, "技师业绩排行")
	} else {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Technician Performance")
	}
	pdf.Ln(10)

	performances, err := s.GetTechnicianPerformance(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if useChineseFont {
		pdf.SetFont("ChineseFont", "B", 10)
		pdf.Cell(50, 8, "技师名称")
		pdf.Cell(40, 8, "业绩(元)")
	} else {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(50, 8, "Technician")
		pdf.Cell(40, 8, "Revenue")
	}
	pdf.Ln(8)

	if useChineseFont {
		pdf.SetFont("ChineseFont", "", 10)
	} else {
		pdf.SetFont("Arial", "", 10)
	}
	for _, perf := range performances {
		pdf.Cell(50, 6, perf.TechnicianName)
		pdf.Cell(40, 6, fmt.Sprintf("%.2f", perf.Revenue))
		pdf.Ln(6)
	}

	pdf.AddPage()

	if useChineseFont {
		pdf.SetFont("ChineseFont", "B", 12)
		pdf.Cell(40, 10, "热门服务排行")
	} else {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Service Ranking")
	}
	pdf.Ln(10)

	rankings, err := s.GetServiceRanking(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if useChineseFont {
		pdf.SetFont("ChineseFont", "B", 10)
		pdf.Cell(50, 8, "服务名称")
		pdf.Cell(40, 8, "收入(元)")
	} else {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(50, 8, "Service")
		pdf.Cell(40, 8, "Revenue")
	}
	pdf.Ln(8)

	if useChineseFont {
		pdf.SetFont("ChineseFont", "", 10)
	} else {
		pdf.SetFont("Arial", "", 10)
	}
	for _, rank := range rankings {
		pdf.Cell(50, 6, rank.ServiceName)
		pdf.Cell(40, 6, fmt.Sprintf("%.2f", rank.Revenue))
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func parseTime(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
