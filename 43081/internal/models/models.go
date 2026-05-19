package models

import "time"

type Account struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Balance  float64   `json:"balance"`
	Currency string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	ParentID *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Children []Category `json:"children,omitempty"`
}

type Transaction struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	CategoryID  int64     `json:"category_id"`
	AccountID   int64     `json:"account_id"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	CategoryName string   `json:"category_name,omitempty"`
	AccountName  string   `json:"account_name,omitempty"`
}

type Transfer struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
}

type Budget struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	Amount     float64   `json:"amount"`
	Month      string    `json:"month"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CategoryName string  `json:"category_name,omitempty"`
	Spent      float64   `json:"spent,omitempty"`
	Remaining  float64   `json:"remaining,omitempty"`
}

type Backup struct {
	ID        int64     `json:"id"`
	FilePath  string    `json:"file_path"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionImportRow struct {
	Date        string  `csv:"date"`
	Type        string  `csv:"type"`
	Amount      float64 `csv:"amount"`
	Category    string  `csv:"category"`
	Account     string  `csv:"account"`
	Description string  `csv:"description"`
}
