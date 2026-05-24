package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Password      string         `gorm:"not null" json:"-"`
	FirstName     string         `gorm:"not null" json:"firstName" binding:"required"`
	LastName      string         `gorm:"not null" json:"lastName" binding:"required"`
	Phone         string         `json:"phone"`
	AvatarURL     string         `json:"avatarUrl"`
	Role          UserRole       `gorm:"type:varchar(20);not null;default:student" json:"role"`
	Status        UserStatus     `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	Timezone      string         `gorm:"default:UTC" json:"timezone"`
	Language      string         `gorm:"default:en" json:"language"`
	EmailVerified bool           `gorm:"default:false" json:"emailVerified"`
	LastLoginAt   *time.Time     `json:"lastLoginAt"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	TeacherProfile *TeacherProfile `gorm:"foreignKey:UserID" json:"teacherProfile,omitempty"`
	StudentProfile *StudentProfile `gorm:"foreignKey:UserID" json:"studentProfile,omitempty"`
	Wallet         *Wallet         `gorm:"foreignKey:UserID" json:"wallet,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

type TeacherProfile struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Bio             string         `gorm:"type:text" json:"bio"`
	Education       string         `gorm:"type:text" json:"education"`
	Experience      string         `gorm:"type:text" json:"experience"`
	Certifications  string         `gorm:"type:text" json:"certifications"`
	HourlyRate      float64        `gorm:"not null;default:0" json:"hourlyRate"`
	Currency        string         `gorm:"default:USD" json:"currency"`
	Rating          float64        `gorm:"default:0" json:"rating"`
	ReviewCount     int            `gorm:"default:0" json:"reviewCount"`
	TotalSessions   int            `gorm:"default:0" json:"totalSessions"`
	TotalHours      float64        `gorm:"default:0" json:"totalHours"`
	IsVerified      bool           `gorm:"default:false" json:"isVerified"`
	ApprovalStatus  string         `gorm:"type:varchar(20);default:pending" json:"approvalStatus"`
	ApprovalNotes   string         `gorm:"type:text" json:"approvalNotes"`
	ApprovedAt      *time.Time     `json:"approvedAt"`
	ApprovedBy      *uuid.UUID     `gorm:"type:uuid" json:"approvedBy"`
	ResumeURL       string         `json:"resumeUrl"`
	IDCardURL       string         `json:"idCardUrl"`
	BankAccountInfo string         `json:"bankAccountInfo"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User            *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Subjects        []TeacherSubject  `gorm:"foreignKey:TeacherID" json:"subjects,omitempty"`
	Availabilities  []AvailabilitySlot `gorm:"foreignKey:TeacherID" json:"availabilities,omitempty"`
}

func (tp *TeacherProfile) BeforeCreate(tx *gorm.DB) error {
	if tp.ID == uuid.Nil {
		tp.ID = uuid.New()
	}
	return nil
}

type Subject struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Category    string         `json:"category"`
	Description string         `gorm:"type:text" json:"description"`
	IconURL     string         `json:"iconUrl"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	SortOrder   int            `gorm:"default:0" json:"sortOrder"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Subject) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type TeacherSubject struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TeacherID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacherId"`
	SubjectID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"subjectId"`
	Level        string         `gorm:"type:varchar(50);not null" json:"level"`
	CustomRate   float64        `gorm:"default:0" json:"customRate"`
	IsPrimary    bool           `gorm:"default:false" json:"isPrimary"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Teacher *TeacherProfile `gorm:"foreignKey:TeacherID" json:"-"`
	Subject *Subject        `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

func (ts *TeacherSubject) BeforeCreate(tx *gorm.DB) error {
	if ts.ID == uuid.Nil {
		ts.ID = uuid.New()
	}
	return nil
}

type AvailabilitySlot struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TeacherID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacherId"`
	DayOfWeek   int            `gorm:"not null" json:"dayOfWeek"`
	StartTime   string         `gorm:"type:varchar(10);not null" json:"startTime"`
	EndTime     string         `gorm:"type:varchar(10);not null" json:"endTime"`
	IsRecurring bool           `gorm:"default:true" json:"isRecurring"`
	StartDate   *time.Time     `json:"startDate"`
	EndDate     *time.Time     `json:"endDate"`
	IsBooked    bool           `gorm:"default:false" json:"isBooked"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Teacher *TeacherProfile `gorm:"foreignKey:TeacherID" json:"-"`
}

func (as *AvailabilitySlot) BeforeCreate(tx *gorm.DB) error {
	if as.ID == uuid.Nil {
		as.ID = uuid.New()
	}
	return nil
}
