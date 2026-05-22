package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password     string    `json:"-" gorm:"size:255;not null"`
	Nickname     string    `json:"nickname" gorm:"size:50"`
	Avatar       string    `json:"avatar" gorm:"size:255"`
	Phone        string    `json:"phone" gorm:"size:20"`
	Level        int       `json:"level" gorm:"default:1"`
	Points       int       `json:"points" gorm:"default:0"`
	IsExpert     bool      `json:"isExpert" gorm:"default:false"`
	ExpertStatus string    `json:"expertStatus" gorm:"size:20;default:'none'"`
	Role         string    `json:"role" gorm:"size:20;default:'user'"`
	Status       string    `json:"status" gorm:"size:20;default:'active'"`
	Bio          string    `json:"bio" gorm:"size:500"`
	LastLoginAt  time.Time `json:"lastLoginAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Question struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"userId" gorm:"index;not null"`
	User            User      `json:"user" gorm:"foreignKey:UserID"`
	Title           string    `json:"title" gorm:"size:200;not null"`
	Content         string    `json:"content" gorm:"type:text;not null"`
	CategoryID      uint      `json:"categoryId" gorm:"index"`
	Category        Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Tags            []Tag     `json:"tags" gorm:"many2many:question_tags;"`
	Views           int       `json:"views" gorm:"default:0"`
	AnswerCount     int       `json:"answerCount" gorm:"default:0"`
	LikeCount       int       `json:"likeCount" gorm:"default:0"`
	CollectCount    int       `json:"collectCount" gorm:"default:0"`
	HotScore        float64   `json:"hotScore" gorm:"default:0"`
	RewardPoints    int       `json:"rewardPoints" gorm:"default:0"`
	IsSolved        bool      `json:"isSolved" gorm:"default:false"`
	AcceptedAnswerID *uint    `json:"acceptedAnswerId" gorm:"index"`
	Status          string    `json:"status" gorm:"size:20;default:'published'"`
	AuditStatus     string    `json:"auditStatus" gorm:"size:20;default:'pending'"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Answer struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	QuestionID   uint      `json:"questionId" gorm:"index;not null"`
	Question     Question  `json:"question" gorm:"foreignKey:QuestionID"`
	UserID       uint      `json:"userId" gorm:"index;not null"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	LikeCount    int       `json:"likeCount" gorm:"default:0"`
	DislikeCount int       `json:"dislikeCount" gorm:"default:0"`
	CollectCount int       `json:"collectCount" gorm:"default:0"`
	IsAccepted   bool      `json:"isAccepted" gorm:"default:false"`
	Status       string    `json:"status" gorm:"size:20;default:'published'"`
	AuditStatus  string    `json:"auditStatus" gorm:"size:20;default:'pending'"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Comment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"userId" gorm:"index;not null"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	AnswerID    uint      `json:"answerId" gorm:"index;not null"`
	Answer      Answer    `json:"answer" gorm:"foreignKey:AnswerID"`
	Content     string    `json:"content" gorm:"size:500;not null"`
	LikeCount   int       `json:"likeCount" gorm:"default:0"`
	Status      string    `json:"status" gorm:"size:20;default:'published'"`
	AuditStatus string    `json:"auditStatus" gorm:"size:20;default:'pending'"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	Description string    `json:"description" gorm:"size:200"`
	Icon        string    `json:"icon" gorm:"size:255"`
	SortOrder   int       `json:"sortOrder" gorm:"default:0"`
	Status      string    `json:"status" gorm:"size:20;default:'active'"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Tag struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Description string    `json:"description" gorm:"size:200"`
	UsageCount  int       `json:"usageCount" gorm:"default:0"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type QuestionTag struct {
	QuestionID uint `json:"questionId" gorm:"primaryKey"`
	TagID      uint `json:"tagId" gorm:"primaryKey"`
}

type AuditRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AdminID      uint      `json:"adminId" gorm:"index"`
	Admin        User      `json:"admin" gorm:"foreignKey:AdminID"`
	TargetType   string    `json:"targetType" gorm:"size:20;not null"`
	TargetID     uint      `json:"targetId" gorm:"not null"`
	Action       string    `json:"action" gorm:"size:20;not null"`
	Reason       string    `json:"reason" gorm:"size:500"`
	Status       string    `json:"status" gorm:"size:20;not null"`
	OperatorIP   string    `json:"operatorIp" gorm:"size:50"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Report struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ReporterID  uint      `json:"reporterId" gorm:"index;not null"`
	Reporter    User      `json:"reporter" gorm:"foreignKey:ReporterID"`
	TargetType  string    `json:"targetType" gorm:"size:20;not null"`
	TargetID    uint      `json:"targetId" gorm:"not null"`
	Reason      string    `json:"reason" gorm:"size:20;not null"`
	Description string    `json:"description" gorm:"size:500"`
	Status      string    `json:"status" gorm:"size:20;default:'pending'"`
	HandlerID   *uint     `json:"handlerId" gorm:"index"`
	Handler     User      `json:"handler" gorm:"foreignKey:HandlerID"`
	HandleResult string   `json:"handleResult" gorm:"size:500"`
	HandledAt   *time.Time `json:"handledAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SensitiveWord struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Word     string    `json:"word" gorm:"uniqueIndex;size:100;not null"`
	Category string    `json:"category" gorm:"size:50"`
	Level    int       `json:"level" gorm:"default:1"`
	ReplaceTo string   `json:"replaceTo" gorm:"size:50"`
	CreatedAt time.Time `json:"createdAt"`
}

type PointLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"userId" gorm:"index;not null"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Type        string    `json:"type" gorm:"size:50;not null"`
	Points      int       `json:"points" gorm:"not null"`
	Balance     int       `json:"balance" gorm:"not null"`
	Description string    `json:"description" gorm:"size:200"`
	RefType     string    `json:"refType" gorm:"size:20"`
	RefID       uint      `json:"refId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Reward struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:500"`
	Image       string    `json:"image" gorm:"size:255"`
	PointsCost  int       `json:"pointsCost" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"default:-1"`
	Status      string    `json:"status" gorm:"size:20;default:'active'"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RewardExchange struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"userId" gorm:"index;not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	RewardID   uint      `json:"rewardId" gorm:"index;not null"`
	Reward     Reward    `json:"reward" gorm:"foreignKey:RewardID"`
	PointsCost int       `json:"pointsCost" gorm:"not null"`
	Status     string    `json:"status" gorm:"size:20;default:'pending'"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Favorite struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"userId" gorm:"index;not null"`
	TargetType string    `json:"targetType" gorm:"size:20;not null"`
	TargetID   uint      `json:"targetId" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Follow struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FollowerID   uint      `json:"followerId" gorm:"index;not null"`
	Follower     User      `json:"follower" gorm:"foreignKey:FollowerID"`
	FollowingType string   `json:"followingType" gorm:"size:20;not null"`
	FollowingID  uint      `json:"followingId" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"index;not null"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	Title     string    `json:"title" gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"size:500"`
	RefType   string    `json:"refType" gorm:"size:20"`
	RefID     uint      `json:"refId"`
	IsRead    bool      `json:"isRead" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
}

type ExpertApplication struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"userId" gorm:"index;not null"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Field        string    `json:"field" gorm:"size:100;not null"`
	Description  string    `json:"description" gorm:"size:500"`
	Status       string    `json:"status" gorm:"size:20;default:'pending'"`
	ReviewerID   *uint     `json:"reviewerId" gorm:"index"`
	Reviewer     User      `json:"reviewer" gorm:"foreignKey:ReviewerID"`
	ReviewRemark string    `json:"reviewRemark" gorm:"size:500"`
	ReviewedAt   *time.Time `json:"reviewedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}
