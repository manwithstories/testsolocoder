package models

import "time"

type Fee struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SeasonID    uint       `json:"season_id"`
	Season      *Season    `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
	TeamID      uint       `json:"team_id"`
	Team        *Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Type        string     `gorm:"size:30" json:"type"`
	Amount      float64    `json:"amount"`
	Status      string     `gorm:"size:20;default:unpaid" json:"status"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	InvoiceNo   string     `gorm:"size:50;uniqueIndex" json:"invoice_no"`
	Note        string     `gorm:"type:text" json:"note"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title     string    `gorm:"size:200" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Type      string    `gorm:"size:30" json:"type"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type PlayerStat struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PlayerID   uint      `json:"player_id"`
	Player     *Player   `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	MatchID    uint      `json:"match_id"`
	Match      *Match    `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	SeasonID   uint      `json:"season_id"`
	TeamID     uint      `json:"team_id"`
	Goals      int       `gorm:"default:0" json:"goals"`
	Assists    int       `gorm:"default:0" json:"assists"`
	Fouls      int       `gorm:"default:0" json:"fouls"`
	YellowCard int       `gorm:"default:0" json:"yellow_card"`
	RedCard    int       `gorm:"default:0" json:"red_card"`
	Minutes    int       `gorm:"default:0" json:"minutes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Standings struct {
	SeasonID     uint    `json:"season_id"`
	TeamID       uint    `json:"team_id"`
	TeamName     string  `json:"team_name"`
	GroupName    string  `json:"group_name"`
	Played       int     `json:"played"`
	Wins         int     `json:"wins"`
	Draws        int     `json:"draws"`
	Losses       int     `json:"losses"`
	GoalsFor     int     `json:"goals_for"`
	GoalsAgainst int     `json:"goals_against"`
	GoalDiff     int     `json:"goal_diff"`
	Points       int     `json:"points"`
}
