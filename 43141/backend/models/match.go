package models

import "time"

type Match struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SeasonID     uint       `json:"season_id"`
	Season       *Season    `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
	Round        int        `json:"round"`
	GroupName    string     `gorm:"size:20" json:"group_name"`
	HomeTeamID   uint       `json:"home_team_id"`
	HomeTeam     *Team      `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeamID   uint       `json:"away_team_id"`
	AwayTeam     *Team      `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	VenueID      *uint      `json:"venue_id,omitempty"`
	Venue        *Venue     `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	RefereeID    *uint      `json:"referee_id,omitempty"`
	Referee      *User      `gorm:"foreignKey:RefereeID" json:"referee,omitempty"`
	MatchTime    *time.Time `json:"match_time,omitempty"`
	HomeScore    *int       `json:"home_score,omitempty"`
	AwayScore    *int       `json:"away_score,omitempty"`
	HasOT        bool       `gorm:"default:false" json:"has_ot"`
	OTHomeScore  *int       `json:"ot_home_score,omitempty"`
	OTAwayScore  *int       `json:"ot_away_score,omitempty"`
	HasPenalty   bool       `gorm:"default:false" json:"has_penalty"`
	PenHomeScore *int       `json:"pen_home_score,omitempty"`
	PenAwayScore *int       `json:"pen_away_score,omitempty"`
	Status       string     `gorm:"size:20;default:scheduled" json:"status"`
	KnockoutStage string   `gorm:"size:30" json:"knockout_stage"`
	WinnerTeamID *uint      `json:"winner_team_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Venue struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Address   string    `gorm:"size:255" json:"address"`
	Capacity  int       `json:"capacity"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RefereeAssignment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MatchID    uint      `json:"match_id"`
	Match      *Match    `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	RefereeID  uint      `json:"referee_id"`
	Referee    *User     `gorm:"foreignKey:RefereeID" json:"referee,omitempty"`
	Status     string    `gorm:"size:20;default:assigned" json:"status"`
	AssignedAt time.Time `json:"assigned_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
