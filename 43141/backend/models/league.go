package models

import "time"

type League struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Sport        string         `gorm:"size:50" json:"sport"`
	Status       string         `gorm:"size:20;default:planning" json:"status"`
	CreatedBy    uint           `json:"created_by"`
	Seasons      []Season       `gorm:"foreignKey:LeagueID" json:"seasons,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Season struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	LeagueID        uint           `json:"league_id"`
	Name            string         `gorm:"size:100" json:"name"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         time.Time      `json:"end_date"`
	Format          string         `gorm:"size:30;default:round_robin" json:"format"`
	GroupCount      int            `gorm:"default:1" json:"group_count"`
	PointsForWin    int            `gorm:"default:3" json:"points_for_win"`
	PointsForDraw   int            `gorm:"default:1" json:"points_for_draw"`
	PointsForLoss   int            `gorm:"default:0" json:"points_for_loss"`
	MaxTeams        int            `gorm:"default:16" json:"max_teams"`
	RegistrationFee float64        `gorm:"default:0" json:"registration_fee"`
	VenueFee        float64        `gorm:"default:0" json:"venue_fee"`
	Status          string         `gorm:"size:20;default:planning" json:"status"`
	CustomRules     string         `gorm:"type:text" json:"custom_rules"`
	Awards          string         `gorm:"type:text" json:"awards"`
	Matches         []Match        `gorm:"foreignKey:SeasonID" json:"matches,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
