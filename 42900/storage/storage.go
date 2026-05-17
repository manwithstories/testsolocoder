package storage

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Vault struct {
	ID         int64
	Name       string
	Salt       []byte
	Verifier   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Entry struct {
	ID         int64
	VaultID    int64
	Name       string
	Username   string
	Password   string
	URL        string
	Tags       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS vaults (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		salt BLOB NOT NULL,
		verifier TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vault_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		url TEXT,
		tags TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_entries_vault_id ON entries(vault_id);
	CREATE INDEX IF NOT EXISTS idx_entries_name ON entries(name);
	`

	_, err := db.Exec(schema)
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) CreateVault(name string, salt []byte, verifier string) (*Vault, error) {
	now := time.Now()
	result, err := s.db.Exec(
		"INSERT INTO vaults (name, salt, verifier, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		name, salt, verifier, now, now,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Vault{
		ID:        id,
		Name:      name,
		Salt:      salt,
		Verifier:  verifier,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *Store) GetVaultByName(name string) (*Vault, error) {
	var vault Vault
	err := s.db.QueryRow(
		"SELECT id, name, salt, verifier, created_at, updated_at FROM vaults WHERE name = ?",
		name,
	).Scan(&vault.ID, &vault.Name, &vault.Salt, &vault.Verifier, &vault.CreatedAt, &vault.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}

func (s *Store) GetVaultByID(id int64) (*Vault, error) {
	var vault Vault
	err := s.db.QueryRow(
		"SELECT id, name, salt, verifier, created_at, updated_at FROM vaults WHERE id = ?",
		id,
	).Scan(&vault.ID, &vault.Name, &vault.Salt, &vault.Verifier, &vault.CreatedAt, &vault.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}

func (s *Store) ListVaults() ([]*Vault, error) {
	rows, err := s.db.Query("SELECT id, name, salt, verifier, created_at, updated_at FROM vaults ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vaults []*Vault
	for rows.Next() {
		var vault Vault
		err := rows.Scan(&vault.ID, &vault.Name, &vault.Salt, &vault.Verifier, &vault.CreatedAt, &vault.UpdatedAt)
		if err != nil {
			return nil, err
		}
		vaults = append(vaults, &vault)
	}
	return vaults, nil
}

func (s *Store) DeleteVault(id int64) error {
	_, err := s.db.Exec("DELETE FROM vaults WHERE id = ?", id)
	return err
}

func (s *Store) UpdateVaultSaltAndVerifier(id int64, salt []byte, verifier string) error {
	_, err := s.db.Exec(
		"UPDATE vaults SET salt = ?, verifier = ?, updated_at = ? WHERE id = ?",
		salt, verifier, time.Now(), id,
	)
	return err
}

func (s *Store) CreateEntry(entry *Entry) (*Entry, error) {
	now := time.Now()
	result, err := s.db.Exec(
		"INSERT INTO entries (vault_id, name, username, password, url, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		entry.VaultID, entry.Name, entry.Username, entry.Password, entry.URL, entry.Tags, now, now,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	entry.ID = id
	entry.CreatedAt = now
	entry.UpdatedAt = now
	return entry, nil
}

func (s *Store) GetEntry(id int64) (*Entry, error) {
	var entry Entry
	err := s.db.QueryRow(
		"SELECT id, vault_id, name, username, password, url, tags, created_at, updated_at FROM entries WHERE id = ?",
		id,
	).Scan(&entry.ID, &entry.VaultID, &entry.Name, &entry.Username, &entry.Password, &entry.URL, &entry.Tags, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Store) ListEntries(vaultID int64) ([]*Entry, error) {
	rows, err := s.db.Query(
		"SELECT id, vault_id, name, username, password, url, tags, created_at, updated_at FROM entries WHERE vault_id = ? ORDER BY name",
		vaultID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*Entry
	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.VaultID, &entry.Name, &entry.Username, &entry.Password, &entry.URL, &entry.Tags, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}

func (s *Store) SearchEntries(vaultID int64, nameQuery string, tagQuery string) ([]*Entry, error) {
	var rows *sql.Rows
	var err error

	query := "SELECT id, vault_id, name, username, password, url, tags, created_at, updated_at FROM entries WHERE vault_id = ?"
	args := []interface{}{vaultID}

	if nameQuery != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+nameQuery+"%")
	}

	if tagQuery != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+tagQuery+"%")
	}

	query += " ORDER BY name"

	rows, err = s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*Entry
	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.VaultID, &entry.Name, &entry.Username, &entry.Password, &entry.URL, &entry.Tags, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}

func (s *Store) UpdateEntry(entry *Entry) error {
	_, err := s.db.Exec(
		"UPDATE entries SET name = ?, username = ?, password = ?, url = ?, tags = ?, updated_at = ? WHERE id = ?",
		entry.Name, entry.Username, entry.Password, entry.URL, entry.Tags, time.Now(), entry.ID,
	)
	return err
}

func (s *Store) DeleteEntry(id int64) error {
	_, err := s.db.Exec("DELETE FROM entries WHERE id = ?", id)
	return err
}

func (s *Store) SetConfig(key, value string) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)",
		key, value,
	)
	return err
}

func (s *Store) GetConfig(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *Store) DeleteConfig(key string) error {
	_, err := s.db.Exec("DELETE FROM config WHERE key = ?", key)
	return err
}

func (s *Store) ReencryptEntries(vaultID int64, reencryptFunc func(string) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id, password FROM entries WHERE vault_id = ?", vaultID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type entryPass struct {
		ID       int64
		Password string
	}

	var entries []entryPass
	for rows.Next() {
		var e entryPass
		if err := rows.Scan(&e.ID, &e.Password); err != nil {
			return err
		}
		entries = append(entries, e)
	}

	for _, e := range entries {
		_, err := tx.Exec("UPDATE entries SET password = ?, updated_at = ? WHERE id = ?",
			e.Password, time.Now(), e.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
