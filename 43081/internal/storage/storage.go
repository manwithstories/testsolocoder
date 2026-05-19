package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"finance-tracker/internal/config"
	"finance-tracker/internal/logger"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(cfg *config.Config) error {
	dbPath := config.ExpandPath(cfg.DatabasePath)
	dbDir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	logger.Info("Database initialized successfully at %s", dbPath)
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func createTables() error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		type TEXT NOT NULL,
		balance REAL NOT NULL DEFAULT 0,
		currency TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		parent_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE,
		UNIQUE(name, type, parent_id)
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		amount REAL NOT NULL,
		category_id INTEGER NOT NULL,
		account_id INTEGER NOT NULL,
		description TEXT,
		date DATE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id),
		FOREIGN KEY (account_id) REFERENCES accounts(id)
	);

	CREATE TABLE IF NOT EXISTS transfers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_account_id INTEGER NOT NULL,
		to_account_id INTEGER NOT NULL,
		amount REAL NOT NULL,
		description TEXT,
		date DATE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (from_account_id) REFERENCES accounts(id),
		FOREIGN KEY (to_account_id) REFERENCES accounts(id)
	);

	CREATE TABLE IF NOT EXISTS budgets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		amount REAL NOT NULL,
		month TEXT NOT NULL,
		year INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id),
		UNIQUE(category_id, month, year)
	);

	CREATE TABLE IF NOT EXISTS backups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
	CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_account ON transactions(account_id);
	`

	_, err := DB.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
