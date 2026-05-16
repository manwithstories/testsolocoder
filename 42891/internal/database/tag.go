package database

import (
	"database/sql"
	"timetrack/internal/models"
)

func getOrCreateTag(tx *sql.Tx, name string) (int64, error) {
	var id int64
	err := tx.QueryRow(`SELECT id FROM tags WHERE name = ?`, name).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	result, err := tx.Exec(`INSERT INTO tags (name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetTagsForTimeEntry(entryID int64) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT t.name FROM tags t
		JOIN time_entry_tags tet ON t.id = tet.tag_id
		WHERE tet.time_entry_id = ?
	`, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func ListTags() ([]models.Tag, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT id, name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}
