package budget

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"finance-tracker/internal/category"
	"finance-tracker/internal/logger"
	"finance-tracker/internal/models"
	"finance-tracker/internal/storage"
)

func Set(categoryID int64, amount float64, month string, year int) (*models.Budget, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("budget amount must be positive")
	}

	if _, err := category.GetByID(categoryID); err != nil {
		return nil, err
	}

	now := time.Now()
	result, err := storage.DB.Exec(`
		INSERT INTO budgets (category_id, amount, month, year, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(category_id, month, year) DO UPDATE SET
			amount = excluded.amount,
			updated_at = excluded.updated_at
	`, categoryID, amount, month, year, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to set budget: %w", err)
	}

	id, _ := result.LastInsertId()
	budget := &models.Budget{
		ID:         id,
		CategoryID: categoryID,
		Amount:     amount,
		Month:      month,
		Year:       year,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	logger.Info("Budget set: category %d, %s %d, amount %.2f", categoryID, month, year, amount)
	return budget, nil
}

func Get(categoryID int64, month string, year int) (*models.Budget, error) {
	var budget models.Budget
	err := storage.DB.QueryRow(`
		SELECT b.id, b.category_id, b.amount, b.month, b.year, b.created_at, b.updated_at, c.name
		FROM budgets b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE b.category_id = ? AND b.month = ? AND b.year = ?
	`, categoryID, month, year).Scan(&budget.ID, &budget.CategoryID, &budget.Amount, &budget.Month, &budget.Year, &budget.CreatedAt, &budget.UpdatedAt, &budget.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	spent, _ := getCategorySpent(categoryID, month, year)
	budget.Spent = spent
	budget.Remaining = budget.Amount - spent

	return &budget, nil
}

func List(month string, year int) ([]models.Budget, error) {
	query := `
		SELECT b.id, b.category_id, b.amount, b.month, b.year, b.created_at, b.updated_at, c.name
		FROM budgets b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE 1=1
	`
	args := []interface{}{}

	if month != "" {
		query += ` AND b.month = ?`
		args = append(args, month)
	}
	if year > 0 {
		query += ` AND b.year = ?`
		args = append(args, year)
	}
	query += ` ORDER BY c.name`

	rows, err := storage.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list budgets: %w", err)
	}
	defer rows.Close()

	var budgets []models.Budget
	for rows.Next() {
		var b models.Budget
		err := rows.Scan(&b.ID, &b.CategoryID, &b.Amount, &b.Month, &b.Year, &b.CreatedAt, &b.UpdatedAt, &b.CategoryName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		spent, _ := getCategorySpent(b.CategoryID, b.Month, b.Year)
		b.Spent = spent
		b.Remaining = b.Amount - spent
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func Delete(categoryID int64, month string, year int) error {
	_, err := storage.DB.Exec(`
		DELETE FROM budgets WHERE category_id = ? AND month = ? AND year = ?
	`, categoryID, month, year)
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}
	logger.Info("Budget deleted: category %d, %s %d", categoryID, month, year)
	return nil
}

func getCategorySpent(categoryID int64, month string, year int) (float64, error) {
	childIDs, err := getAllChildCategoryIDs(categoryID)
	if err != nil {
		return 0, err
	}

	allIDs := append(childIDs, categoryID)

	placeholders := make([]string, len(allIDs))
	args := make([]interface{}, 0, len(allIDs)+2)
	for i, id := range allIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}
	args = append(args, fmt.Sprintf("%d", year), month)

	var spent float64
	query := fmt.Sprintf(`
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE category_id IN (%s) AND type = 'expense'
		  AND substr(date, 1, 4) = ?
		  AND substr(date, 6, 2) = ?
	`, strings.Join(placeholders, ","))

	err = storage.DB.QueryRow(query, args...).Scan(&spent)
	if err != nil {
		return 0, err
	}
	return spent, nil
}

func getAllChildCategoryIDs(parentID int64) ([]int64, error) {
	var childIDs []int64

	rows, err := storage.DB.Query(`
		SELECT id FROM categories WHERE parent_id = ?
	`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		childIDs = append(childIDs, id)

		subChildren, err := getAllChildCategoryIDs(id)
		if err != nil {
			return nil, err
		}
		childIDs = append(childIDs, subChildren...)
	}

	return childIDs, nil
}

func CheckOverspent(month string, year int) ([]models.Budget, error) {
	budgets, err := List(month, year)
	if err != nil {
		return nil, err
	}

	var overspent []models.Budget
	for _, b := range budgets {
		if b.Remaining < 0 {
			overspent = append(overspent, b)
		}
	}
	return overspent, nil
}

func GetWarnings(month string, year int) []string {
	overspent, err := CheckOverspent(month, year)
	if err != nil {
		return nil
	}

	var warnings []string
	for _, b := range overspent {
		percent := (b.Spent / b.Amount) * 100
		warnings = append(warnings, fmt.Sprintf("WARNING: Category '%s' budget exceeded by %.2f (%.1f%% over budget)",
			b.CategoryName, -b.Remaining, percent-100))
	}
	return warnings
}
