package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"luxury-trading-platform/internal/model"

	"github.com/signintech/gopdf"
)

type PDFService struct {
	uploadPath string
}

func NewPDFService(uploadPath string) *PDFService {
	return &PDFService{uploadPath: uploadPath}
}

type AuthenticationReportData struct {
	ReportID          string
	ProductTitle      string
	ProductCategory   string
	BrandName         string
	AuthenticatorName string
	AuthenticatorID   string
	Result            string
	ResultCN          string
	ReportContent     string
	Notes             string
	Date              string
}

func (s *PDFService) GenerateAuthenticationReport(data *AuthenticationReportData) (string, error) {
	pdfDir := filepath.Join(s.uploadPath, "reports")
	if err := os.MkdirAll(pdfDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create reports directory: %w", err)
	}

	filename := fmt.Sprintf("auth_report_%s_%s.pdf", data.ReportID, time.Now().Format("20060102150405"))
	filePath := filepath.Join(pdfDir, filename)

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 24)
	pdf.SetY(40)
	pdf.Cell(nil, "奢侈品鉴定报告")

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetY(70)
	pdf.Cell(nil, fmt.Sprintf("报告编号: %s", data.ReportID))
	pdf.SetY(85)
	pdf.Cell(nil, fmt.Sprintf("生成日期: %s", data.Date))

	pdf.SetY(105)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(nil, "商品信息")

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetY(125)
	pdf.Cell(nil, fmt.Sprintf("商品名称: %s", data.ProductTitle))
	pdf.SetY(140)
	pdf.Cell(nil, fmt.Sprintf("商品类别: %s", data.ProductCategory))
	pdf.SetY(155)
	pdf.Cell(nil, fmt.Sprintf("品牌名称: %s", data.BrandName))

	pdf.SetY(180)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(nil, "鉴定师信息")

	pdf.SetFont("Helvetica", "", 12)
	pdf.SetY(200)
	pdf.Cell(nil, fmt.Sprintf("鉴定师: %s", data.AuthenticatorName))
	pdf.SetY(215)
	pdf.Cell(nil, fmt.Sprintf("资质编号: %s", data.AuthenticatorID))

	pdf.SetY(240)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(nil, "鉴定结果")

	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetY(260)
	pdf.Cell(nil, fmt.Sprintf("结论: %s", data.ResultCN))

	pdf.SetY(285)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(nil, "鉴定详情:")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetY(300)
	pdf.Cell(nil, truncateString(data.ReportContent, 500))

	if data.Notes != "" {
		pdf.SetY(400)
		pdf.SetFont("Helvetica", "B", 12)
		pdf.Cell(nil, "备注:")
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetY(415)
		pdf.Cell(nil, truncateString(data.Notes, 300))
	}

	pdf.SetY(500)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(nil, "本报告由奢侈品交易平台生成，具有法律效力")
	pdf.SetY(515)
	pdf.Cell(nil, fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))

	err := pdf.WritePdf(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to write PDF: %w", err)
	}

	return fmt.Sprintf("/uploads/reports/%s", filename), nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func GetResultCN(result model.AuthenticationResult) string {
	switch result {
	case model.AuthenticationResultGenuine:
		return "正品"
	case model.AuthenticationResultCounterfeit:
		return "赝品"
	case model.AuthenticationResultInconclusive:
		return "无法鉴定"
	default:
		return "未知"
	}
}
