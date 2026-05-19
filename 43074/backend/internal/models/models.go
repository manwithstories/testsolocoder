package models

import (
	"time"
)

type ReadingStatus string

const (
	StatusToRead    ReadingStatus = "to_read"
	StatusReading   ReadingStatus = "reading"
	StatusCompleted ReadingStatus = "completed"
	StatusAbandoned ReadingStatus = "abandoned"
)

type Book struct {
	ID              uint64        `gorm:"primaryKey" json:"id"`
	Title           string        `gorm:"size:255;not null" json:"title"`
	Author          string        `gorm:"size:255" json:"author"`
	Publisher       string        `gorm:"size:255" json:"publisher"`
	ISBN            string        `gorm:"size:20;uniqueIndex" json:"isbn"`
	CoverImage      string        `gorm:"size:500" json:"cover_image"`
	Summary         string        `gorm:"type:text" json:"summary"`
	TotalPages      int           `gorm:"default:0" json:"total_pages"`
	ReadingStatus   ReadingStatus `gorm:"size:20;default:'to_read'" json:"reading_status"`
	CurrentPage     int           `gorm:"default:0" json:"current_page"`
	ReadingProgress float64       `gorm:"default:0" json:"reading_progress"`
	StartDate       *time.Time    `json:"start_date"`
	EndDate         *time.Time    `json:"end_date"`
	TotalReadTime   int           `gorm:"default:0" json:"total_read_time"`
	Tags            []*Tag        `gorm:"many2many:book_tags;" json:"tags,omitempty"`
	Categories      []*Category   `gorm:"many2many:book_categories;" json:"categories,omitempty"`
	ReadingNotes    []ReadingNote `gorm:"foreignKey:BookID" json:"reading_notes,omitempty"`
	BorrowRecord    *BorrowRecord `gorm:"foreignKey:BookID" json:"borrow_record,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type Tag struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Color     string    `gorm:"size:20;default:'#3b82f6'" json:"color"`
	Books     []*Book   `gorm:"many2many:book_tags;" json:"books,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        uint64      `gorm:"primaryKey" json:"id"`
	Name      string      `gorm:"size:100;not null" json:"name"`
	ParentID  *uint64     `gorm:"index" json:"parent_id,omitempty"`
	Parent    *Category   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []*Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Books     []*Book     `gorm:"many2many:book_categories;" json:"books,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type ReadingNote struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	BookID    uint64    `gorm:"index;not null" json:"book_id"`
	Page      int       `json:"page"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BorrowRecord struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	BookID        uint64     `gorm:"index;unique;not null" json:"book_id"`
	BorrowerName  string     `gorm:"size:100;not null" json:"borrower_name"`
	BorrowerPhone string     `gorm:"size:20" json:"borrower_phone"`
	BorrowerEmail string     `gorm:"size:100" json:"borrower_email"`
	BorrowDate    time.Time  `gorm:"not null" json:"borrow_date"`
	ExpectedReturnDate *time.Time `json:"expected_return_date"`
	ReturnDate    *time.Time `json:"return_date,omitempty"`
	Status        string     `gorm:"size:20;default:'borrowed'" json:"status"`
	Notes         string     `gorm:"type:text" json:"notes"`
	ReminderSent  bool       `gorm:"default:false" json:"reminder_sent"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ReadingGoal struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Year        int       `gorm:"index;not null" json:"year"`
	Month       *int      `gorm:"index" json:"month,omitempty"`
	TargetBooks int       `gorm:"default:0" json:"target_books"`
	TargetPages int       `gorm:"default:0" json:"target_pages"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReadSession struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	BookID    uint64    `gorm:"index;not null" json:"book_id"`
	StartPage int       `json:"start_page"`
	EndPage   int       `json:"end_page"`
	Duration  int       `json:"duration"`
	ReadDate  time.Time `gorm:"index;not null" json:"read_date"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type BookTag struct {
	BookID uint64 `gorm:"primaryKey"`
	TagID  uint64 `gorm:"primaryKey"`
}

type BookCategory struct {
	BookID     uint64 `gorm:"primaryKey"`
	CategoryID uint64 `gorm:"primaryKey"`
}
