package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".timetrack")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "timetrack.db")
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		hourly_rate REAL NOT NULL DEFAULT 0,
		currency TEXT NOT NULL DEFAULT 'CNY',
		archived INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS time_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		paused INTEGER NOT NULL DEFAULT 0,
		paused_at DATETIME,
		total_paused_seconds INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS time_entry_tags (
		time_entry_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (time_entry_id, tag_id),
		FOREIGN KEY (time_entry_id) REFERENCES time_entries(id),
		FOREIGN KEY (tag_id) REFERENCES tags(id)
	);

	CREATE INDEX IF NOT EXISTS idx_time_entries_project_id ON time_entries(project_id);
	CREATE INDEX IF NOT EXISTS idx_time_entries_start_time ON time_entries(start_time);
	CREATE INDEX IF NOT EXISTS idx_time_entries_end_time ON time_entries(end_time);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	if err := migrateDB(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func migrateDB(db *sql.DB) error {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('time_entries') WHERE name = 'total_paused_seconds'`).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec(`ALTER TABLE time_entries ADD COLUMN total_paused_seconds INTEGER NOT NULL DEFAULT 0`)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".timetrack", "timetrack.db"), nil
}
