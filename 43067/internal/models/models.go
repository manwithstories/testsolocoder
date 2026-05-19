package models

import "time"

type Snippet struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Language    string    `json:"language"`
	Tags        []string  `json:"tags"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Encrypted   bool      `json:"encrypted"`
}

type Vault struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Snippets  []Snippet `json:"snippets"`
}

type Config struct {
	DefaultVault   string `json:"default_vault"`
	DefaultEditor  string `json:"default_editor"`
	HighlightTheme string `json:"highlight_theme"`
	EncryptionKey  string `json:"encryption_key,omitempty"`
}

type ExportData struct {
	Version string  `json:"version"`
	Vaults  []Vault `json:"vaults"`
}

type SearchQuery struct {
	Keyword  string
	Tags     []string
	Language string
	Fields   []string
}
