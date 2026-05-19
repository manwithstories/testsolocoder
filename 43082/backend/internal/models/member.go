package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Member struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:100;not null" json:"name" binding:"required"`
	Phone           string         `gorm:"size:20;uniqueIndex;not null" json:"phone" binding:"required,len=11"`
	Email           string         `gorm:"size:100;uniqueIndex" json:"email"`
	Gender          string         `gorm:"size:10" json:"gender"`
	Birthday        *time.Time     `json:"birthday"`
	Address         string         `gorm:"size:255" json:"address"`
	ProfilePhoto    string         `gorm:"size:255" json:"profile_photo"`
	Password        string         `gorm:"size:255;not null" json:"-"`
	Status          int            `gorm:"default:1" json:"status"` // 1:正常 2:冻结
	MembershipID    *uint          `json:"membership_id"`
	Membership      *Membership    `gorm:"foreignKey:MembershipID" json:"membership,omitempty"`
	Bookings        []Booking      `gorm:"foreignKey:MemberID" json:"bookings,omitempty"`
	CheckIns        []CheckIn      `gorm:"foreignKey:MemberID" json:"check_ins,omitempty"`
	OperationLogs   []OperationLog `gorm:"foreignKey:MemberID" json:"operation_logs,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Member) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hashed)
	return nil
}

func (m *Member) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	return err == nil
}

type MembershipType string

const (
	MembershipTypeMonthly MembershipType = "monthly"
	MembershipTypeQuarter MembershipType = "quarter"
	MembershipTypeYearly  MembershipType = "yearly"
)

var ValidMembershipTypes = map[MembershipType]int{
	MembershipTypeMonthly: 30,
	MembershipTypeQuarter: 90,
	MembershipTypeYearly:  365,
}

type Membership struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	MemberID   uint           `gorm:"not null;index" json:"member_id"`
	Member     *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Type       MembershipType `gorm:"size:20;not null" json:"type" binding:"required,oneof=monthly quarter yearly"`
	StartDate  time.Time      `gorm:"not null" json:"start_date"`
	EndDate    time.Time      `gorm:"not null" json:"end_date"`
	Price      float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Status     int            `gorm:"default:1" json:"status"` // 1:有效 2:已过期 3:已冻结
	AutoRenew  bool           `gorm:"default:false" json:"auto_renew"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Membership) IsValid() bool {
	now := time.Now()
	return m.Status == 1 && now.After(m.StartDate) && now.Before(m.EndDate)
}

func (m *Membership) DaysRemaining() int {
	now := time.Now()
	if now.After(m.EndDate) {
		return 0
	}
	return int(m.EndDate.Sub(now).Hours() / 24)
}
