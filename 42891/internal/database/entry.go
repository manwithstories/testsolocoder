package database

import (
	"database/sql"
	"fmt"
	"timetrack/internal/models"
	"time"
)

func StartTimeEntry(projectID int64, startTime time.Time, tags []string) (*models.TimeEntry, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	activeEntry, err := GetActiveTimeEntry()
	if err != nil {
		return nil, err
	}
	if activeEntry != nil {
		return nil, fmt.Errorf("there is already an active time entry. Stop it first")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result, err := tx.Exec(`
		INSERT INTO time_entries (project_id, start_time, total_paused_seconds, created_at, updated_at)
		VALUES (?, ?, 0, ?, ?)
	`, projectID, startTime, now, now)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	entryID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, tagName := range tags {
		tagID, err := getOrCreateTag(tx, tagName)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		_, err = tx.Exec(`INSERT INTO time_entry_tags (time_entry_id, tag_id) VALUES (?, ?)`, entryID, tagID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.TimeEntry{
		ID:                 entryID,
		ProjectID:          projectID,
		StartTime:          startTime,
		EndTime:            nil,
		Paused:             false,
		PausedAt:           nil,
		TotalPausedSeconds:  0,
		Tags:               tags,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func GetActiveTimeEntry() (*models.TimeEntry, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var entry models.TimeEntry
	var endTime sql.NullTime
	var pausedAt sql.NullTime

	err = db.QueryRow(`
		SELECT id, project_id, start_time, end_time, paused, paused_at, total_paused_seconds, created_at, updated_at
		FROM time_entries
		WHERE end_time IS NULL
		ORDER BY start_time DESC
		LIMIT 1
	`).Scan(&entry.ID, &entry.ProjectID, &entry.StartTime, &endTime, &entry.Paused, &pausedAt, &entry.TotalPausedSeconds, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if endTime.Valid {
		entry.EndTime = &endTime.Time
	}
	if pausedAt.Valid {
		entry.PausedAt = &pausedAt.Time
	}

	entry.Tags, err = GetTagsForTimeEntry(entry.ID)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func StopTimeEntry(entryID int64, endTime time.Time) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE time_entries SET end_time = ?, paused = 0, paused_at = NULL, updated_at = ? WHERE id = ?
	`, endTime, time.Now(), entryID)
	return err
}

func PauseTimeEntry(entryID int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = db.Exec(`
		UPDATE time_entries SET paused = 1, paused_at = ?, updated_at = ? WHERE id = ?
	`, now, now, entryID)
	return err
}

func ResumeTimeEntry(entryID int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var pausedAt sql.NullTime
	var currentPausedSeconds int64
	err = db.QueryRow(`SELECT paused_at, total_paused_seconds FROM time_entries WHERE id = ?`, entryID).Scan(&pausedAt, &currentPausedSeconds)
	if err != nil {
		return err
	}

	if !pausedAt.Valid {
		return fmt.Errorf("entry is not paused")
	}

	pausedDuration := time.Since(pausedAt.Time)
	additionalPausedSeconds := int64(pausedDuration.Seconds())
	newTotalPausedSeconds := currentPausedSeconds + additionalPausedSeconds

	_, err = db.Exec(`
		UPDATE time_entries SET paused = 0, paused_at = NULL, total_paused_seconds = ?, updated_at = ? WHERE id = ?
	`, newTotalPausedSeconds, time.Now(), entryID)
	return err
}

func GetElapsedDuration(entry *models.TimeEntry) time.Duration {
	var endTime time.Time
	if entry.EndTime != nil {
		endTime = *entry.EndTime
	} else {
		if entry.Paused && entry.PausedAt != nil {
			endTime = *entry.PausedAt
		} else {
			endTime = time.Now()
		}
	}

	totalDuration := endTime.Sub(entry.StartTime)
	pausedDuration := time.Duration(entry.TotalPausedSeconds) * time.Second

	if entry.EndTime == nil && entry.Paused && entry.PausedAt != nil {
		currentPausedDuration := time.Since(*entry.PausedAt)
		pausedDuration += currentPausedDuration
	}

	elapsed := totalDuration - pausedDuration
	if elapsed < 0 {
		elapsed = 0
	}
	return elapsed
}

func GetTimeEntryByID(entryID int64) (*models.TimeEntry, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var entry models.TimeEntry
	var endTime sql.NullTime
	var pausedAt sql.NullTime

	err = db.QueryRow(`
		SELECT id, project_id, start_time, end_time, paused, paused_at, total_paused_seconds, created_at, updated_at
		FROM time_entries WHERE id = ?
	`, entryID).Scan(&entry.ID, &entry.ProjectID, &entry.StartTime, &endTime, &entry.Paused, &pausedAt, &entry.TotalPausedSeconds, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if endTime.Valid {
		entry.EndTime = &endTime.Time
	}
	if pausedAt.Valid {
		entry.PausedAt = &pausedAt.Time
	}

	entry.Tags, err = GetTagsForTimeEntry(entryID)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func ListTimeEntries(startDate, endDate time.Time, projectID *int64, tagFilter *string) ([]models.TimeEntryWithTags, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT DISTINCT te.id, te.project_id, te.start_time, te.end_time, te.paused, te.paused_at, te.total_paused_seconds, te.created_at, te.updated_at, p.name
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
	`
	args := []interface{}{}

	if tagFilter != nil {
		query += ` JOIN time_entry_tags tet ON te.id = tet.time_entry_id
		JOIN tags t ON tet.tag_id = t.id
		`
	}

	query += ` WHERE te.start_time >= ? AND (te.end_time <= ? OR te.end_time IS NULL)`
	args = append(args, startDate, endDate)

	if projectID != nil {
		query += ` AND te.project_id = ?`
		args = append(args, *projectID)
	}

	if tagFilter != nil {
		query += ` AND t.name = ?`
		args = append(args, *tagFilter)
	}

	query += ` ORDER BY te.start_time DESC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.TimeEntryWithTags
	for rows.Next() {
		var et models.TimeEntryWithTags
		var endTime sql.NullTime
		var pausedAt sql.NullTime

		err := rows.Scan(&et.ID, &et.ProjectID, &et.StartTime, &endTime, &et.Paused, &pausedAt, &et.TotalPausedSeconds, &et.CreatedAt, &et.UpdatedAt, &et.ProjectName)
		if err != nil {
			return nil, err
		}

		if endTime.Valid {
			et.EndTime = &endTime.Time
		}
		if pausedAt.Valid {
			et.PausedAt = &pausedAt.Time
		}

		et.Tags, err = GetTagsForTimeEntry(et.ID)
		if err != nil {
			return nil, err
		}

		entries = append(entries, et)
	}

	return entries, nil
}

func CheckForCrossDayConflict(entry *models.TimeEntry) (bool, time.Time, error) {
	startDay := entry.StartTime.Truncate(24 * time.Hour)
	now := time.Now().Truncate(24 * time.Hour)

	if now.After(startDay) {
		return true, startDay.Add(24 * time.Hour), nil
	}
	return false, time.Time{}, nil
}

func TruncateTimeEntry(entryID int64, truncateAt time.Time) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE time_entries SET end_time = ?, updated_at = ? WHERE id = ?
	`, truncateAt, time.Now(), entryID)
	return err
}
