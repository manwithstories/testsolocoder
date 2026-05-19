package pdf

import (
	"fmt"
	"os"
	"path/filepath"

	"freelancer-management/internal/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateInvoicePDF(invoice *models.Invoice) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(0, 10, "INVOICE")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 8, "Invoice #")
	pdf.SetFont("Arial", "", 14)
	pdf.Cell(0, 8, invoice.InvoiceNumber)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 7, "Issue Date:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 7, invoice.IssueDate.Format("2006-01-02"))
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 7, "Due Date:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 7, invoice.DueDate.Format("2006-01-02"))
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 7, "Status:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 7, string(invoice.Status))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "Bill To:")
	pdf.Ln(10)

	if invoice.Client != nil {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 6, invoice.Client.Name)
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(0, 5, invoice.Client.Email)
		pdf.Ln(5)
		if invoice.Client.Company != "" {
			pdf.Cell(0, 5, invoice.Client.Company)
			pdf.Ln(5)
		}
		if invoice.Client.Address != "" {
			pdf.Cell(0, 5, invoice.Client.Address)
			pdf.Ln(5)
		}
		if invoice.Client.Phone != "" {
			pdf.Cell(0, 5, invoice.Client.Phone)
			pdf.Ln(5)
		}
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(80, 8, "Description")
	pdf.Cell(30, 8, "Quantity")
	pdf.Cell(30, 8, "Unit Price")
	pdf.Cell(30, 8, "Amount")
	pdf.Ln(10)

	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(2)

	for _, item := range invoice.Items {
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(80, 7, item.Description)
		pdf.Cell(30, 7, fmt.Sprintf("%.2f", item.Quantity))
		pdf.Cell(30, 7, fmt.Sprintf("$%.2f", item.UnitPrice))
		pdf.Cell(30, 7, fmt.Sprintf("$%.2f", item.Amount))
		pdf.Ln(7)
	}

	pdf.Ln(5)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(110, 7, "")
	pdf.Cell(30, 7, "Subtotal:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(30, 7, fmt.Sprintf("$%.2f", invoice.Subtotal))
	pdf.Ln(7)

	if invoice.TaxRate > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(110, 7, "")
		pdf.Cell(30, 7, fmt.Sprintf("Tax (%.1f%%):", invoice.TaxRate))
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(30, 7, fmt.Sprintf("$%.2f", invoice.TaxAmount))
		pdf.Ln(7)
	}

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(110, 8, "")
	pdf.Cell(30, 8, "Total:")
	pdf.Cell(30, 8, fmt.Sprintf("$%.2f", invoice.Total))
	pdf.Ln(15)

	if invoice.Notes != "" {
		pdf.SetFont("Arial", "I", 11)
		pdf.MultiCell(0, 5, "Notes: "+invoice.Notes, "", "", false)
	}

	dir := "./invoices"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filename := filepath.Join(dir, fmt.Sprintf("%s.pdf", invoice.InvoiceNumber))
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}
