package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus int

const (
	BookingStatusBooked   BookingStatus = 1 // 已预约
	BookingStatusCancelled BookingStatus = 2 // 已取消
	BookingStatusCheckedIn BookingStatus = 3 // 已签到
	BookingStatusMissed   BookingStatus = 4 // 未到场
)

type Booking struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	MemberID     uint           `gorm:"not null;index:idx_member_schedule,unique" json:"member_id" binding:"required"`
	Member       *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	ScheduleID   uint           `gorm:"not null;index:idx_member_schedule,unique" json:"schedule_id" binding:"required"`
	Schedule     *CourseSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Status       BookingStatus  `gorm:"default:1;not null" json:"status"`
	BookingTime  time.Time      `json:"booking_time"`
	CancelTime   *time.Time     `json:"cancel_time,omitempty"`
	CheckInID    *uint          `json:"check_in_id,omitempty"`
	CheckIn      *CheckIn       `gorm:"foreignKey:CheckInID" json:"check_in,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Waitlist struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	MemberID   uint           `gorm:"not null;index:idx_waitlist_member_schedule,unique" json:"member_id"`
	Member     *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	ScheduleID uint           `gorm:"not null;index:idx_waitlist_member_schedule,unique" json:"schedule_id"`
	Schedule   *CourseSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Position   int            `gorm:"not null" json:"position"`
	Notified   bool           `gorm:"default:false" json:"notified"`
	Status     int            `gorm:"default:1" json:"status"` // 1:等待中 2:已转预约 3:已取消
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CheckIn struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	MemberID   uint           `gorm:"not null;index" json:"member_id"`
	Member     *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	ScheduleID *uint          `gorm:"index" json:"schedule_id,omitempty"`
	Schedule   *CourseSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	CheckInTime time.Time     `gorm:"not null" json:"check_in_time"`
	CheckType  int            `gorm:"default:1" json:"check_type"` // 1:正常签到 2:补签
	Remark     string         `gorm:"size:255" json:"remark"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type Reminder struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MemberID       uint           `gorm:"not null;index" json:"member_id"`
	Member         *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	ScheduleID     *uint          `gorm:"index" json:"schedule_id,omitempty"`
	Type           int            `gorm:"not null" json:"type"` // 1:会员卡到期提醒 2:课程开始前1天提醒 3:课程开始前2小时提醒
	Title          string         `gorm:"size:255;not null" json:"title"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	ScheduleTime   time.Time      `gorm:"not null;index" json:"schedule_time"`
	SentTime       *time.Time     `json:"sent_time,omitempty"`
	Status         int            `gorm:"default:1" json:"status"` // 1:待发送 2:已发送 3:发送失败
	RetryCount     int            `gorm:"default:0" json:"retry_count"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type OperationLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	MemberID     *uint          `gorm:"index" json:"member_id,omitempty"`
	Member       *Member        `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	OperatorID   uint           `json:"operator_id"`
	OperatorType string         `gorm:"size:20" json:"operator_type"` // admin, member, system
	Action       string         `gorm:"size:50;not null" json:"action"`
	ResourceType string         `gorm:"size:50;not null" json:"resource_type"`
	ResourceID   uint           `json:"resource_id"`
	Detail       string         `gorm:"type:text" json:"detail"`
	IPAddress    string         `gorm:"size:50" json:"ip_address"`
	UserAgent    string         `gorm:"size:255" json:"user_agent"`
	CreatedAt    time.Time      `json:"created_at"`
}
