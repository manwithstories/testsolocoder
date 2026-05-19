package account

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"finance-tracker/internal/logger"
	"finance-tracker/internal/models"
	"finance-tracker/internal/storage"
)

func Create(name, accountType, currency string, initialBalance float64) (*models.Account, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}

	if initialBalance < 0 {
		return nil, fmt.Errorf("initial balance cannot be negative")
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		currency = "CNY"
	}

	accountType = strings.ToLower(strings.TrimSpace(accountType))
	validTypes := map[string]bool{"cash": true, "bank": true, "alipay": true, "wechat": true, "credit": true, "other": true}
	if !validTypes[accountType] {
		return nil, fmt.Errorf("invalid account type: %s, must be one of: cash, bank, alipay, wechat, credit, other", accountType)
	}

	now := time.Now()
	result, err := storage.DB.Exec(`
		INSERT INTO accounts (name, type, balance, currency, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, name, accountType, initialBalance, currency, now, now)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("account name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	id, _ := result.LastInsertId()
	account := &models.Account{
		ID:        id,
		Name:      name,
		Type:      accountType,
		Balance:   initialBalance,
		Currency:  currency,
		CreatedAt: now,
		UpdatedAt: now,
	}

	logger.Info("Account created: %s (ID: %d) with initial balance %.2f %s", name, id, initialBalance, currency)
	return account, nil
}

func GetByID(id int64) (*models.Account, error) {
	var account models.Account
	err := storage.DB.QueryRow(`
		SELECT id, name, type, balance, currency, created_at, updated_at
		FROM accounts WHERE id = ?
	`, id).Scan(&account.ID, &account.Name, &account.Type, &account.Balance, &account.Currency, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &account, nil
}

func GetByName(name string) (*models.Account, error) {
	var account models.Account
	err := storage.DB.QueryRow(`
		SELECT id, name, type, balance, currency, created_at, updated_at
		FROM accounts WHERE name = ?
	`, name).Scan(&account.ID, &account.Name, &account.Type, &account.Balance, &account.Currency, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &account, nil
}

func List() ([]models.Account, error) {
	rows, err := storage.DB.Query(`
		SELECT id, name, type, balance, currency, created_at, updated_at
		FROM accounts ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(&account.ID, &account.Name, &account.Type, &account.Balance, &account.Currency, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func UpdateBalance(id int64, delta float64) error {
	account, err := GetByID(id)
	if err != nil {
		return err
	}

	newBalance := account.Balance + delta
	if newBalance < 0 {
		return fmt.Errorf("insufficient balance: account '%s' has %.2f, operation requires %.2f", account.Name, account.Balance, -delta)
	}

	now := time.Now()
	_, err = storage.DB.Exec(`
		UPDATE accounts SET balance = ?, updated_at = ? WHERE id = ?
	`, newBalance, now, id)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	logger.Info("Account %s (ID: %d) balance updated: %.2f -> %.2f", account.Name, id, account.Balance, newBalance)
	return nil
}

func Transfer(fromID, toID int64, amount float64, description string, date time.Time) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	if fromID == toID {
		return fmt.Errorf("cannot transfer to the same account")
	}

	fromAccount, err := GetByID(fromID)
	if err != nil {
		return err
	}

	toAccount, err := GetByID(toID)
	if err != nil {
		return err
	}

	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient balance: account '%s' has %.2f, transfer requires %.2f", fromAccount.Name, fromAccount.Balance, amount)
	}

	tx, err := storage.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	now := time.Now()
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - ?, updated_at = ? WHERE id = ?`, amount, now, fromID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct from source account: %w", err)
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + ?, updated_at = ? WHERE id = ?`, amount, now, toID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add to target account: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO transfers (from_account_id, to_account_id, amount, description, date, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, fromID, toID, amount, description, date, now)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record transfer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transfer: %w", err)
	}

	logger.Info("Transfer completed: %.2f from %s to %s", amount, fromAccount.Name, toAccount.Name)
	return nil
}

func Delete(id int64) error {
	account, err := GetByID(id)
	if err != nil {
		return err
	}

	var count int
	err = storage.DB.QueryRow(`SELECT COUNT(*) FROM transactions WHERE account_id = ?`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check account transactions: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete account '%s' with %d transactions", account.Name, count)
	}

	_, err = storage.DB.Exec(`DELETE FROM accounts WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	logger.Info("Account deleted: %s (ID: %d)", account.Name, id)
	return nil
}
