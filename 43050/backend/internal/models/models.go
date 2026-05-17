package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Groups    []Group        `gorm:"many2many:group_members;" json:"groups,omitempty"`
	Expenses  []Expense      `gorm:"foreignKey:PaidBy" json:"-"`
}

type Group struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	CreatorID   uint           `gorm:"not null" json:"creatorId"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Members     []User         `gorm:"many2many:group_members;" json:"members,omitempty"`
	Expenses    []Expense      `json:"expenses,omitempty"`
	InviteCode  string         `gorm:"uniqueIndex;size:20;not null" json:"inviteCode"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type GroupMember struct {
	GroupID   uint      `gorm:"primaryKey" json:"groupId"`
	UserID    uint      `gorm:"primaryKey" json:"userId"`
	JoinedAt  time.Time `json:"joinedAt"`
	IsActive  bool      `gorm:"default:true" json:"isActive"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type SplitType string

const (
	SplitEqual   SplitType = "equal"
	SplitRatio   SplitType = "ratio"
	SplitCustom  SplitType = "custom"
)

type Expense struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	GroupID     uint           `gorm:"not null;index" json:"groupId"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Amount      float64        `gorm:"not null" json:"amount"`
	PaidBy      uint           `gorm:"not null" json:"paidBy"`
	Payer       User           `gorm:"foreignKey:PaidBy" json:"payer,omitempty"`
	SplitType   SplitType      `gorm:"size:20;not null" json:"splitType"`
	CreatedBy   uint           `gorm:"not null" json:"createdBy"`
	Version     int            `gorm:"default:1" json:"version"`
	ExpenseDate time.Time      `json:"expenseDate"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Participants []ExpenseParticipant `json:"participants,omitempty"`
}

type ExpenseParticipant struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	ExpenseID  uint    `gorm:"not null;index" json:"expenseId"`
	UserID     uint    `gorm:"not null" json:"userId"`
	User       User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Amount     float64 `gorm:"not null" json:"amount"`
	Ratio      float64 `json:"ratio,omitempty"`
	IsSettled  bool    `gorm:"default:false" json:"isSettled"`
}

type Settlement struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	GroupID     uint           `gorm:"not null;index" json:"groupId"`
	FromUserID  uint           `gorm:"not null" json:"fromUserId"`
	FromUser    User           `gorm:"foreignKey:FromUserID" json:"fromUser,omitempty"`
	ToUserID    uint           `gorm:"not null" json:"toUserId"`
	ToUser      User           `gorm:"foreignKey:ToUserID" json:"toUser,omitempty"`
	Amount      float64        `gorm:"not null" json:"amount"`
	IsPaid      bool           `gorm:"default:false" json:"isPaid"`
	PaidAt      *time.Time     `json:"paidAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
