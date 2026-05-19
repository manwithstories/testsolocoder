package models

import (
	"time"

	"gorm.io/gorm"
)

type Coach struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name" binding:"required"`
	Phone         string         `gorm:"size:20;uniqueIndex;not null" json:"phone" binding:"required"`
	Specialty     string         `gorm:"size:255" json:"specialty"`
	Description   string         `gorm:"type:text" json:"description"`
	Photo         string         `gorm:"size:255" json:"photo"`
	Status        int            `gorm:"default:1" json:"status"` // 1:在职 2:离职
	Courses       []Course       `gorm:"foreignKey:CoachID" json:"courses,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type CourseType string

const (
	CourseTypeSingle   CourseType = "single"
	CourseTypeWeekly   CourseType = "weekly"
	CourseTypeMonthly  CourseType = "monthly"
)

type Course struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name" binding:"required"`
	Description   string         `gorm:"type:text" json:"description"`
	CoachID       uint           `gorm:"not null;index" json:"coach_id" binding:"required"`
	Coach         *Coach         `gorm:"foreignKey:CoachID" json:"coach,omitempty"`
	Capacity      int            `gorm:"not null" json:"capacity" binding:"required,min=1"`
	Duration      int            `gorm:"not null" json:"duration" binding:"required,min=1"` // 分钟
	Type          CourseType     `gorm:"size:20;not null;default:single" json:"type" binding:"oneof=single weekly monthly"`
	Weekdays      string         `gorm:"size:20" json:"weekdays"` // 例如: "1,3,5" 表示周一、三、五
	StartDate     time.Time      `gorm:"not null" json:"start_date"`
	EndDate       *time.Time     `json:"end_date"` // 周期课程的结束日期
	StartTime     string         `gorm:"size:10;not null" json:"start_time" binding:"required"` // HH:MM
	Location      string         `gorm:"size:100" json:"location"`
	Status        int            `gorm:"default:1" json:"status"` // 1:正常 2:取消 3:结束
	CourseSchedules []CourseSchedule `gorm:"foreignKey:CourseID" json:"schedules,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type CourseSchedule struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CourseID   uint           `gorm:"not null;index" json:"course_id"`
	Course     *Course        `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	StartTime  time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime    time.Time      `gorm:"not null" json:"end_time"`
	Capacity   int            `gorm:"not null" json:"capacity"`
	BookedCount int           `gorm:"default:0" json:"booked_count"`
	Status     int            `gorm:"default:1" json:"status"` // 1:可预约 2:已满 3:已取消 4:已结束
	Bookings   []Booking      `gorm:"foreignKey:ScheduleID" json:"bookings,omitempty"`
	Waitlists  []Waitlist     `gorm:"foreignKey:ScheduleID" json:"waitlists,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
