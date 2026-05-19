package stats

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"finance-tracker/internal/logger"
	"finance-tracker/internal/models"
	"finance-tracker/internal/storage"
	"finance-tracker/internal/transaction"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetProfit    float64 `json:"net_profit"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

type CategoryStats struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type AccountStats struct {
	AccountID   int64   `json:"account_id"`
	AccountName string  `json:"account_name"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
	Net         float64 `json:"net"`
}

func GetSummary(startDate, endDate time.Time, accountID int64) (*Summary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense
		FROM transactions
		WHERE date >= ? AND date <= ?
	`
	args := []interface{}{startDate.Format("2006-01-02"), endDate.Format("2006-01-02")}

	if accountID > 0 {
		query += ` AND account_id = ?`
		args = append(args, accountID)
	}

	var income, expense float64
	err := storage.DB.QueryRow(query, args...).Scan(&income, &expense)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}

	return &Summary{
		TotalIncome:  income,
		TotalExpense: expense,
		NetProfit:    income - expense,
		StartDate:    startDate,
		EndDate:      endDate,
	}, nil
}

func GetByCategory(startDate, endDate time.Time, transType string, accountID int64) ([]CategoryStats, error) {
	query := `
		SELECT c.id, c.name, COALESCE(SUM(t.amount), 0) as total
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.date >= ? AND t.date <= ? AND t.type = ?
	`
	args := []interface{}{startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), transType}

	if accountID > 0 {
		query += ` AND t.account_id = ?`
		args = append(args, accountID)
	}
	query += ` GROUP BY c.id, c.name ORDER BY total DESC`

	rows, err := storage.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", err)
	}
	defer rows.Close()

	var total float64
	var stats []CategoryStats
	for rows.Next() {
		var s CategoryStats
		err := rows.Scan(&s.CategoryID, &s.CategoryName, &s.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category stats: %w", err)
		}
		total += s.Amount
		stats = append(stats, s)
	}

	for i := range stats {
		if total > 0 {
			stats[i].Percentage = (stats[i].Amount / total) * 100
		}
	}

	return stats, nil
}

func GetByAccount(startDate, endDate time.Time) ([]AccountStats, error) {
	query := `
		SELECT a.id, a.name,
			COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0) as expense
		FROM accounts a
		LEFT JOIN transactions t ON a.id = t.account_id
			AND t.date >= ? AND t.date <= ?
		GROUP BY a.id, a.name
		ORDER BY a.name
	`

	rows, err := storage.DB.Query(query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to get account stats: %w", err)
	}
	defer rows.Close()

	var stats []AccountStats
	for rows.Next() {
		var s AccountStats
		err := rows.Scan(&s.AccountID, &s.AccountName, &s.Income, &s.Expense)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account stats: %w", err)
		}
		s.Net = s.Income - s.Expense
		stats = append(stats, s)
	}
	return stats, nil
}

func GenerateReport(startDate, endDate time.Time, accountID int64, outputFile string) error {
	summary, err := GetSummary(startDate, endDate, accountID)
	if err != nil {
		return err
	}

	expenseByCategory, err := GetByCategory(startDate, endDate, "expense", accountID)
	if err != nil {
		return err
	}

	incomeByCategory, err := GetByCategory(startDate, endDate, "income", accountID)
	if err != nil {
		return err
	}

	accountStats, err := GetByAccount(startDate, endDate)
	if err != nil {
		return err
	}

	transactions, err := transaction.List(startDate, endDate, 0, accountID, "")
	if err != nil {
		return err
	}

	ext := filepath.Ext(outputFile)
	switch ext {
	case ".csv":
		return exportCSV(outputFile, summary, expenseByCategory, incomeByCategory, accountStats, transactions)
	case ".pdf":
		return exportPDF(outputFile, summary, expenseByCategory, incomeByCategory, accountStats, transactions)
	case ".xlsx", ".xls":
		return exportExcel(outputFile, summary, expenseByCategory, incomeByCategory, accountStats, transactions)
	default:
		return exportText(outputFile, summary, expenseByCategory, incomeByCategory, accountStats, transactions)
	}
}

func exportText(filename string, summary *Summary, expenses, incomes []CategoryStats, accounts []AccountStats, transactions []models.Transaction) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "========================================\n")
	fmt.Fprintf(f, "FINANCIAL REPORT\n")
	fmt.Fprintf(f, "========================================\n")
	fmt.Fprintf(f, "Period: %s to %s\n\n", summary.StartDate.Format("2006-01-02"), summary.EndDate.Format("2006-01-02"))

	fmt.Fprintf(f, "SUMMARY\n")
	fmt.Fprintf(f, "----------------------------------------\n")
	fmt.Fprintf(f, "Total Income:  %.2f\n", summary.TotalIncome)
	fmt.Fprintf(f, "Total Expense: %.2f\n", summary.TotalExpense)
	fmt.Fprintf(f, "Net Profit:    %.2f\n\n", summary.NetProfit)

	fmt.Fprintf(f, "EXPENSES BY CATEGORY\n")
	fmt.Fprintf(f, "----------------------------------------\n")
	for _, s := range expenses {
		fmt.Fprintf(f, "%-20s %10.2f (%.1f%%)\n", s.CategoryName, s.Amount, s.Percentage)
	}
	fmt.Fprintln(f)

	fmt.Fprintf(f, "INCOME BY CATEGORY\n")
	fmt.Fprintf(f, "----------------------------------------\n")
	for _, s := range incomes {
		fmt.Fprintf(f, "%-20s %10.2f (%.1f%%)\n", s.CategoryName, s.Amount, s.Percentage)
	}
	fmt.Fprintln(f)

	fmt.Fprintf(f, "ACCOUNT SUMMARY\n")
	fmt.Fprintf(f, "----------------------------------------\n")
	for _, a := range accounts {
		fmt.Fprintf(f, "%-20s Income: %10.2f  Expense: %10.2f  Net: %10.2f\n", a.AccountName, a.Income, a.Expense, a.Net)
	}
	fmt.Fprintln(f)

	fmt.Fprintf(f, "TRANSACTIONS\n")
	fmt.Fprintf(f, "----------------------------------------\n")
	for _, t := range transactions {
		fmt.Fprintf(f, "%s  %-7s %10.2f  %-15s %-15s %s\n",
			t.Date.Format("2006-01-02"), t.Type, t.Amount, t.CategoryName, t.AccountName, t.Description)
	}

	logger.Info("Report generated: %s", filename)
	return nil
}

func exportCSV(filename string, summary *Summary, expenses, incomes []CategoryStats, accounts []AccountStats, transactions []models.Transaction) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV report: %w", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "Financial Report")
	fmt.Fprintf(f, "Period,%s - %s\n", summary.StartDate.Format("2006-01-02"), summary.EndDate.Format("2006-01-02"))
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Summary")
	fmt.Fprintf(f, "Total Income,%.2f\n", summary.TotalIncome)
	fmt.Fprintf(f, "Total Expense,%.2f\n", summary.TotalExpense)
	fmt.Fprintf(f, "Net Profit,%.2f\n", summary.NetProfit)
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Expenses by Category")
	fmt.Fprintln(f, "Category,Amount,Percentage")
	for _, s := range expenses {
		fmt.Fprintf(f, "%s,%.2f,%.1f%%\n", s.CategoryName, s.Amount, s.Percentage)
	}
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Income by Category")
	fmt.Fprintln(f, "Category,Amount,Percentage")
	for _, s := range incomes {
		fmt.Fprintf(f, "%s,%.2f,%.1f%%\n", s.CategoryName, s.Amount, s.Percentage)
	}
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Transactions")
	fmt.Fprintln(f, "Date,Type,Amount,Category,Account,Description")
	for _, t := range transactions {
		fmt.Fprintf(f, "%s,%s,%.2f,%s,%s,%s\n",
			t.Date.Format("2006-01-02"), t.Type, t.Amount, t.CategoryName, t.AccountName, t.Description)
	}

	logger.Info("CSV Report generated: %s", filename)
	return nil
}

func exportPDF(filename string, summary *Summary, expenses, incomes []CategoryStats, accounts []AccountStats, transactions []models.Transaction) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Financial Report", true)
	pdf.SetAuthor("Finance Tracker", true)

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "FINANCIAL REPORT")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, fmt.Sprintf("Period: %s to %s", summary.StartDate.Format("2006-01-02"), summary.EndDate.Format("2006-01-02")))
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Summary")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, fmt.Sprintf("Total Income:  %.2f", summary.TotalIncome))
	pdf.Ln(8)
	pdf.Cell(40, 8, fmt.Sprintf("Total Expense: %.2f", summary.TotalExpense))
	pdf.Ln(8)
	pdf.Cell(40, 8, fmt.Sprintf("Net Profit:    %.2f", summary.NetProfit))
	pdf.Ln(12)

	if len(expenses) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Expenses by Category")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 12)
		for _, s := range expenses {
			pdf.Cell(80, 8, s.CategoryName)
			pdf.Cell(30, 8, fmt.Sprintf("%.2f", s.Amount))
			pdf.Cell(30, 8, fmt.Sprintf("%.1f%%", s.Percentage))
			pdf.Ln(8)
		}
		pdf.Ln(5)
	}

	if len(incomes) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Income by Category")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 12)
		for _, s := range incomes {
			pdf.Cell(80, 8, s.CategoryName)
			pdf.Cell(30, 8, fmt.Sprintf("%.2f", s.Amount))
			pdf.Cell(30, 8, fmt.Sprintf("%.1f%%", s.Percentage))
			pdf.Ln(8)
		}
		pdf.Ln(5)
	}

	if len(accounts) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Account Summary")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 12)
		for _, a := range accounts {
			pdf.Cell(50, 8, a.AccountName)
			pdf.Cell(30, 8, fmt.Sprintf("Income: %.2f", a.Income))
			pdf.Cell(30, 8, fmt.Sprintf("Expense: %.2f", a.Expense))
			pdf.Cell(30, 8, fmt.Sprintf("Net: %.2f", a.Net))
			pdf.Ln(8)
		}
		pdf.Ln(5)
	}

	if len(transactions) > 0 {
		if pdf.GetY() > 200 {
			pdf.AddPage()
		}
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Transactions")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, t := range transactions {
			if pdf.GetY() > 270 {
				pdf.AddPage()
			}
			pdf.Cell(25, 6, t.Date.Format("2006-01-02"))
			pdf.Cell(20, 6, t.Type)
			pdf.Cell(25, 6, fmt.Sprintf("%.2f", t.Amount))
			pdf.Cell(35, 6, t.CategoryName)
			pdf.Cell(30, 6, t.AccountName)
			pdf.Cell(50, 6, t.Description)
			pdf.Ln(6)
		}
	}

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	logger.Info("PDF Report generated: %s", filename)
	return nil
}

func exportExcel(filename string, summary *Summary, expenses, incomes []CategoryStats, accounts []AccountStats, transactions []models.Transaction) error {
	f := excelize.NewFile()

	sheetName := "Report"
	f.SetSheetName("Sheet1", sheetName)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16},
	})

	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)
	f.SetCellValue(sheetName, "A1", "FINANCIAL REPORT")

	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Period: %s to %s", summary.StartDate.Format("2006-01-02"), summary.EndDate.Format("2006-01-02")))

	row := 4
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), headerStyle)
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Summary")
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Income")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), summary.TotalIncome)
	row++
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Expense")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), summary.TotalExpense)
	row++
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Net Profit")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), summary.NetProfit)
	row += 2

	if len(expenses) > 0 {
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), headerStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Expenses by Category")
		row++

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Category")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Amount")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Percentage")
		row++

		for _, s := range expenses {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), s.CategoryName)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), s.Amount)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("%.1f%%", s.Percentage))
			row++
		}
		row++
	}

	if len(incomes) > 0 {
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), headerStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Income by Category")
		row++

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Category")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Amount")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Percentage")
		row++

		for _, s := range incomes {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), s.CategoryName)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), s.Amount)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("%.1f%%", s.Percentage))
			row++
		}
		row++
	}

	if len(accounts) > 0 {
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Account Summary")
		row++

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Account")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Income")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Expense")
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "Net")
		row++

		for _, a := range accounts {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), a.AccountName)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), a.Income)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), a.Expense)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), a.Net)
			row++
		}
		row++
	}

	if len(transactions) > 0 {
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row), headerStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Transactions")
		row++

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Date")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Type")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Amount")
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "Category")
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "Account")
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "Description")
		row++

		for _, t := range transactions {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), t.Date.Format("2006-01-02"))
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), t.Type)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), t.Amount)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), t.CategoryName)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), t.AccountName)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), t.Description)
			row++
		}
	}

	f.SetColWidth(sheetName, "A", "F", 15)

	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	logger.Info("Excel Report generated: %s", filename)
	return nil
}
