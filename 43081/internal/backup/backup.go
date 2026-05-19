package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"finance-tracker/internal/config"
	"finance-tracker/internal/logger"
	"finance-tracker/internal/models"
	"finance-tracker/internal/storage"
)

func CreateLocal(cfg *config.Config) (string, error) {
	dbPath := config.ExpandPath(cfg.DatabasePath)
	backupDir := config.ExpandPath(cfg.Backup.LocalPath)

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("finance_%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	src, err := os.Open(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source database: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(backupPath)
		return "", fmt.Errorf("failed to copy database: %w", err)
	}

	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat backup file: %w", err)
	}

	_, err = storage.DB.Exec(`
		INSERT INTO backups (file_path, file_size, created_at)
		VALUES (?, ?, ?)
	`, backupPath, fileInfo.Size(), time.Now())
	if err != nil {
		logger.Error("Failed to record backup in database: %v", err)
	}

	if err := cleanupOldBackups(cfg); err != nil {
		logger.Error("Failed to cleanup old backups: %v", err)
	}

	logger.Info("Backup created: %s (size: %d bytes)", backupPath, fileInfo.Size())
	return backupPath, nil
}

func List() ([]models.Backup, error) {
	rows, err := storage.DB.Query(`
		SELECT id, file_path, file_size, created_at
		FROM backups ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}
	defer rows.Close()

	var backups []models.Backup
	for rows.Next() {
		var b models.Backup
		err := rows.Scan(&b.ID, &b.FilePath, &b.FileSize, &b.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan backup: %w", err)
		}
		backups = append(backups, b)
	}
	return backups, nil
}

func Restore(backupPath string, cfg *config.Config) error {
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	dbPath := config.ExpandPath(cfg.DatabasePath)

	storage.Close()

	src, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to write to database: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to restore database: %w", err)
	}

	if err := storage.Init(cfg); err != nil {
		return fmt.Errorf("failed to reinitialize database: %w", err)
	}

	logger.Info("Database restored from: %s", backupPath)
	return nil
}

func cleanupOldBackups(cfg *config.Config) error {
	backupDir := config.ExpandPath(cfg.Backup.LocalPath)
	maxBackups := cfg.Backup.MaxBackups
	retentionDays := cfg.Backup.RetentionDays

	files, err := filepath.Glob(filepath.Join(backupDir, "finance_*.db"))
	if err != nil {
		return err
	}

	sort.Strings(files)

	if len(files) > maxBackups {
		for _, f := range files[:len(files)-maxBackups] {
			if err := os.Remove(f); err != nil {
				logger.Error("Failed to remove old backup: %s", f)
			} else {
				logger.Info("Removed old backup: %s", f)
			}
		}
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(f); err != nil {
				logger.Error("Failed to remove expired backup: %s", f)
			} else {
				logger.Info("Removed expired backup: %s", f)
			}
		}
	}

	return nil
}

func CheckAndAutoBackup(cfg *config.Config) {
	if !cfg.Backup.Enabled {
		return
	}

	lastBackupTime, err := getLastBackupTime()
	if err != nil {
		logger.Error("Failed to get last backup time: %v", err)
		return
	}

	now := time.Now()
	var backupInterval time.Duration
	switch cfg.Backup.Schedule {
	case "hourly":
		backupInterval = 1 * time.Hour
	case "daily":
		backupInterval = 24 * time.Hour
	case "weekly":
		backupInterval = 7 * 24 * time.Hour
	default:
		backupInterval = 24 * time.Hour
	}

	if now.Sub(lastBackupTime) >= backupInterval {
		logger.Info("Auto backup triggered (last backup: %s)", lastBackupTime.Format("2006-01-02 15:04:05"))
		if _, err := CreateLocal(cfg); err != nil {
			logger.Error("Auto backup failed: %v", err)
		}
	}
}

func getLastBackupTime() (time.Time, error) {
	var lastBackupStr string
	err := storage.DB.QueryRow(`
		SELECT COALESCE(MAX(created_at), '1970-01-01 00:00:00')
		FROM backups
	`).Scan(&lastBackupStr)
	if err != nil {
		return time.Time{}, err
	}

	if idx := strings.Index(lastBackupStr, " m="); idx != -1 {
		lastBackupStr = lastBackupStr[:idx]
	}

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.000000 -0700 MST",
		"2006-01-02 15:04:05.000000 -0700",
		"2006-01-02 15:04:05.000000",
		time.RFC3339,
	}

	var lastBackup time.Time
	var parseErr error
	for _, format := range formats {
		lastBackup, parseErr = time.Parse(format, lastBackupStr)
		if parseErr == nil {
			return lastBackup, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse backup time '%s': %w", lastBackupStr, parseErr)
}
