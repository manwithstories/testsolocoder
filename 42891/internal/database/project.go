package database

import (
	"database/sql"
	"fmt"
	"timetrack/internal/models"
	"time"
)

func CreateProject(name string, hourlyRate float64, currency string) (*models.Project, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result, err := db.Exec(`
		INSERT INTO projects (name, hourly_rate, currency, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, name, hourlyRate, currency, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Project{
		ID:         id,
		Name:       name,
		HourlyRate: hourlyRate,
		Currency:   currency,
		Archived:   false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func GetProjectByName(name string) (*models.Project, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var p models.Project
	err = db.QueryRow(`
		SELECT id, name, hourly_rate, currency, archived, created_at, updated_at
		FROM projects WHERE name = ?
	`, name).Scan(&p.ID, &p.Name, &p.HourlyRate, &p.Currency, &p.Archived, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func GetProjectByID(id int64) (*models.Project, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var p models.Project
	err = db.QueryRow(`
		SELECT id, name, hourly_rate, currency, archived, created_at, updated_at
		FROM projects WHERE id = ?
	`, id).Scan(&p.ID, &p.Name, &p.HourlyRate, &p.Currency, &p.Archived, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func ListProjects(includeArchived bool) ([]models.Project, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, hourly_rate, currency, archived, created_at, updated_at
		FROM projects
	`
	if !includeArchived {
		query += ` WHERE archived = 0`
	}
	query += ` ORDER BY name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(&p.ID, &p.Name, &p.HourlyRate, &p.Currency, &p.Archived, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func ArchiveProject(id int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE projects SET archived = 1, updated_at = ? WHERE id = ?
	`, time.Now(), id)
	return err
}

func UnarchiveProject(id int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE projects SET archived = 0, updated_at = ? WHERE id = ?
	`, time.Now(), id)
	return err
}

func DeleteProject(id int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM time_entry_tags WHERE time_entry_id IN (SELECT id FROM time_entries WHERE project_id = ?)`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM time_entries WHERE project_id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func SetProjectRate(id int64, hourlyRate float64, currency string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE projects SET hourly_rate = ?, currency = ?, updated_at = ? WHERE id = ?
	`, hourlyRate, currency, time.Now(), id)
	return err
}
