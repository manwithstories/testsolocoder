package models

import "time"

type Team struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Logo         string         `gorm:"size:255" json:"logo"`
	CaptainID    uint           `json:"captain_id"`
	Captain      *User          `gorm:"foreignKey:CaptainID" json:"captain,omitempty"`
	Description  string         `gorm:"type:text" json:"description"`
	ContactEmail string         `gorm:"size:100" json:"contact_email"`
	ContactPhone string         `gorm:"size:20" json:"contact_phone"`
	Players      []Player       `gorm:"foreignKey:TeamID" json:"players,omitempty"`
	Registrations []Registration `gorm:"foreignKey:TeamID" json:"registrations,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Player struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TeamID     uint      `json:"team_id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Number     int       `json:"number"`
	Position   string    `gorm:"size:30" json:"position"`
	BirthDate  time.Time `json:"birth_date"`
	UserID     *uint     `json:"user_id,omitempty"`
	User       *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Registration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SeasonID  uint      `json:"season_id"`
	Season    *Season   `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
	TeamID    uint      `json:"team_id"`
	Team      *Team     `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	GroupName string    `gorm:"size:20" json:"group_name"`
	Status    string    `gorm:"size:20;default:pending" json:"status"`
	Paid      bool      `gorm:"default:false" json:"paid"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	Note      string    `gorm:"type:text" json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
