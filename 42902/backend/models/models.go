package models

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Email      string `gorm:"unique;not null" json:"email"`
	Username   string `gorm:"not null" json:"username"`
	Password   string `gorm:"not null" json:"-"`
	Verified   bool   `gorm:"default:false" json:"verified"`
	VerifyCode string `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

type Event struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Location    string    `gorm:"not null" json:"location"`
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Capacity    int       `gorm:"not null" json:"capacity"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	OrganizerID uint      `gorm:"not null" json:"organizer_id"`
	Organizer   User      `gorm:"foreignKey:OrganizerID" json:"organizer"`
	CreatedAt   time.Time `json:"created_at"`
	Registrations []Registration `gorm:"foreignKey:EventID" json:"-"`
}

type Registration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_event,unique" json:"user_id"`
	EventID   uint      `gorm:"not null;index:idx_user_event,unique" json:"event_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Event     Event     `gorm:"foreignKey:EventID" json:"-"`
	Status    string    `gorm:"default:registered" json:"status"`
	CancelCount int     `gorm:"default:0" json:"cancel_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("activity.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&User{}, &Event{}, &Registration{})
}
