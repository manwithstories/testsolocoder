package transaction

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"finance-tracker/internal/account"
	"finance-tracker/internal/category"
	"finance-tracker/internal/logger"
	"finance-tracker/internal/models"
	"finance-tracker/internal/storage"
)

func Add(transType string, amount float64, categoryID, accountID int64, description string, date time.Time) (*models.Transaction, error) {
	transType = strings.ToLower(strings.TrimSpace(transType))
	if transType != "income" && transType != "expense" {
		return nil, fmt.Errorf("invalid transaction type: %s, must be 'income' or 'expense'", transType)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("transaction amount must be positive")
	}

	if !category.Exists(categoryID, transType) {
		return nil, fmt.Errorf("category with ID %d does not exist or is not of type '%s'", categoryID, transType)
	}

	if _, err := account.GetByID(accountID); err != nil {
		return nil, err
	}

	amountForAccount := amount
	if transType == "expense" {
		amountForAccount = -amount
	}

	tx, err := storage.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	now := time.Now()
	result, err := tx.Exec(`
		INSERT INTO transactions (type, amount, category_id, account_id, description, date, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, transType, amount, categoryID, accountID, description, date, now)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert transaction: %w", err)
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + ?, updated_at = ? WHERE id = ?`, amountForAccount, now, accountID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	id, _ := result.LastInsertId()
	cat, _ := category.GetByID(categoryID)
	acc, _ := account.GetByID(accountID)

	transaction := &models.Transaction{
		ID:           id,
		Type:         transType,
		Amount:       amount,
		CategoryID:   categoryID,
		AccountID:    accountID,
		Description:  description,
		Date:         date,
		CreatedAt:    now,
		CategoryName: cat.Name,
		AccountName:  acc.Name,
	}

	logger.Info("Transaction added: %s %.2f (category: %s, account: %s)", transType, amount, cat.Name, acc.Name)
	return transaction, nil
}

func List(startDate, endDate time.Time, categoryID, accountID int64, transType string) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.type, t.amount, t.category_id, t.account_id, t.description, t.date, t.created_at,
		       c.name as category_name, a.name as account_name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN accounts a ON t.account_id = a.id
		WHERE 1=1
	`
	args := []interface{}{}

	if !startDate.IsZero() {
		query += ` AND t.date >= ?`
		args = append(args, startDate.Format("2006-01-02"))
	}
	if !endDate.IsZero() {
		query += ` AND t.date <= ?`
		args = append(args, endDate.Format("2006-01-02"))
	}
	if categoryID > 0 {
		query += ` AND t.category_id = ?`
		args = append(args, categoryID)
	}
	if accountID > 0 {
		query += ` AND t.account_id = ?`
		args = append(args, accountID)
	}
	if transType != "" {
		query += ` AND t.type = ?`
		args = append(args, transType)
	}
	query += ` ORDER BY t.date DESC, t.created_at DESC`

	rows, err := storage.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var dateStr string
		err := rows.Scan(&t.ID, &t.Type, &t.Amount, &t.CategoryID, &t.AccountID, &t.Description, &dateStr, &t.CreatedAt, &t.CategoryName, &t.AccountName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		t.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			t.Date, _ = time.Parse(time.RFC3339, dateStr)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func Delete(id int64) error {
	var t models.Transaction
	var dateStr string
	err := storage.DB.QueryRow(`
		SELECT id, type, amount, category_id, account_id, description, date
		FROM transactions WHERE id = ?
	`, id).Scan(&t.ID, &t.Type, &t.Amount, &t.CategoryID, &t.AccountID, &t.Description, &dateStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("transaction with ID %d not found", id)
		}
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	amountForAccount := t.Amount
	if t.Type == "income" {
		amountForAccount = -t.Amount
	}

	tx, err := storage.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM transactions WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	now := time.Now()
	_, err = tx.Exec(`UPDATE accounts SET balance = balance + ?, updated_at = ? WHERE id = ?`, amountForAccount, now, t.AccountID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to revert account balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	logger.Info("Transaction deleted: ID %d (%.2f %s)", id, t.Amount, t.Type)
	return nil
}

func ImportCSV(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	requiredFields := []string{"date", "type", "amount", "category", "account"}
	for _, field := range requiredFields {
		if _, ok := headerMap[field]; !ok {
			return 0, fmt.Errorf("missing required column: %s", field)
		}
	}

	var imported int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return imported, fmt.Errorf("failed to read CSV record at line %d: %w", imported+2, err)
		}

		dateStr := strings.TrimSpace(record[headerMap["date"]])
		transType := strings.ToLower(strings.TrimSpace(record[headerMap["type"]]))
		amountStr := strings.TrimSpace(record[headerMap["amount"]])
		categoryName := strings.TrimSpace(record[headerMap["category"]])
		accountName := strings.TrimSpace(record[headerMap["account"]])
		description := ""
		if idx, ok := headerMap["description"]; ok {
			description = strings.TrimSpace(record[idx])
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			logger.Error("Invalid date format at line %d: %s", imported+2, dateStr)
			continue
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			logger.Error("Invalid amount at line %d: %s", imported+2, amountStr)
			continue
		}

		cat, err := category.GetByNameAndType(categoryName, transType)
		if err != nil {
			logger.Error("Category '%s' not found at line %d", categoryName, imported+2)
			continue
		}

		acc, err := account.GetByName(accountName)
		if err != nil {
			logger.Error("Account '%s' not found at line %d", accountName, imported+2)
			continue
		}

		_, err = Add(transType, amount, cat.ID, acc.ID, description, date)
		if err != nil {
			logger.Error("Failed to import line %d: %v", imported+2, err)
			continue
		}

		imported++
	}

	logger.Info("CSV import completed: %d records imported", imported)
	return imported, nil
}

func ValidateDate(dateStr string) (time.Time, error) {
	formats := []string{"2006-01-02", "2006/01/02", "01-02-2006"}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s, expected YYYY-MM-DD", dateStr)
}
