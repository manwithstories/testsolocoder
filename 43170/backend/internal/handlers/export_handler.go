package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"photo-rental/internal/models"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type ExportHandler struct{}

func NewExportHandler() *ExportHandler {
	return &ExportHandler{}
}

type RentalRecord struct {
	OrderID     uint
	EquipmentName string
	RenterName  string
	StartDate   string
	EndDate     string
	TotalRent   float64
	Deposit     float64
	Status      string
}

func (h *ExportHandler) ExportRentalRecords(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")
	format := c.DefaultQuery("format", "csv")

	if userRole != "owner" && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only owners can export rental records")
		return
	}

	var orders []models.Order
	db := database.DB.Preload("Equipment").Preload("Renter")

	if userRole == "owner" {
		db = db.Where("owner_id = ?", userID)
	}

	if err := db.Order("created_at DESC").Find(&orders).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch rental records")
		return
	}

	records := make([]RentalRecord, len(orders))
	for i, order := range orders {
		records[i] = RentalRecord{
			OrderID:       order.ID,
			EquipmentName: order.Equipment.Name,
			RenterName:    order.Renter.RealName,
			StartDate:     order.StartDate,
			EndDate:       order.EndDate,
			TotalRent:     order.TotalRent,
			Deposit:       order.Deposit,
			Status:        order.Status,
		}
	}

	switch format {
	case "csv":
		exportCSV(c, records)
	case "pdf":
		exportPDF(c, records)
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid format. Use csv or pdf")
	}
}

func (h *ExportHandler) ExportRevenueReport(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")
	format := c.DefaultQuery("format", "csv")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if userRole != "owner" && userRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only owners can export revenue reports")
		return
	}

	var orders []models.Order
	db := database.DB.Preload("Equipment").Preload("Renter").
		Where("status = ?", "completed")

	if userRole == "owner" {
		db = db.Where("owner_id = ?", userID)
	}

	if startDate != "" {
		db = db.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("created_at <= ?", endDate)
	}

	if err := db.Order("created_at DESC").Find(&orders).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch revenue data")
		return
	}

	records := make([]RentalRecord, len(orders))
	for i, order := range orders {
		records[i] = RentalRecord{
			OrderID:       order.ID,
			EquipmentName: order.Equipment.Name,
			RenterName:    order.Renter.RealName,
			StartDate:     order.StartDate,
			EndDate:       order.EndDate,
			TotalRent:     order.TotalRent,
			Deposit:       order.Deposit,
			Status:        order.Status,
		}
	}

	switch format {
	case "csv":
		exportCSV(c, records)
	case "pdf":
		exportPDF(c, records)
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid format. Use csv or pdf")
	}
}

func exportCSV(c *gin.Context, records []RentalRecord) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"Order ID", "Equipment Name", "Renter Name", "Start Date", "End Date", "Total Rent", "Deposit", "Status"}
	writer.Write(headers)

	for _, record := range records {
		row := []string{
			fmt.Sprintf("%d", record.OrderID),
			record.EquipmentName,
			record.RenterName,
			record.StartDate,
			record.EndDate,
			fmt.Sprintf("%.2f", record.TotalRent),
			fmt.Sprintf("%.2f", record.Deposit),
			record.Status,
		}
		writer.Write(row)
	}

	writer.Flush()

	filename := fmt.Sprintf("rental_records_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv", buf.Bytes())
}

func exportPDF(c *gin.Context, records []RentalRecord) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Rental Records Report")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"Order ID", "Equipment", "Renter", "Start", "End", "Rent", "Deposit", "Status"}
	colWidths := []float64{20, 40, 35, 25, 25, 20, 20, 25}

	for i, header := range headers {
		pdf.Cell(colWidths[i], 7, header)
	}
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 8)
	for _, record := range records {
		pdf.Cell(colWidths[0], 6, fmt.Sprintf("%d", record.OrderID))
		pdf.Cell(colWidths[1], 6, truncateString(record.EquipmentName, 15))
		pdf.Cell(colWidths[2], 6, truncateString(record.RenterName, 12))
		pdf.Cell(colWidths[3], 6, record.StartDate)
		pdf.Cell(colWidths[4], 6, record.EndDate)
		pdf.Cell(colWidths[5], 6, fmt.Sprintf("%.2f", record.TotalRent))
		pdf.Cell(colWidths[6], 6, fmt.Sprintf("%.2f", record.Deposit))
		pdf.Cell(colWidths[7], 6, record.Status)
		pdf.Ln(6)
	}

	filename := fmt.Sprintf("rental_records_%s.pdf", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	var pdfBuffer bytes.Buffer
	if err := pdf.Output(&pdfBuffer); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}
	c.Data(http.StatusOK, "application/pdf", pdfBuffer.Bytes())
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
