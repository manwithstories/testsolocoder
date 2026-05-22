package model

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleUser    UserRole = "user"
	RoleJudge   UserRole = "judge"
)

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password      string    `gorm:"size:255;not null" json:"-"`
	RealName      string    `gorm:"size:64" json:"real_name"`
	IdCard        string    `gorm:"size:32;index" json:"id_card"`
	Phone         string    `gorm:"size:32" json:"phone"`
	Email         string    `gorm:"size:128" json:"email"`
	Verified      bool      `gorm:"default:false" json:"verified"`
	Role          UserRole  `gorm:"size:16;default:user" json:"role"`
	Status        int       `gorm:"default:1" json:"status"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Event struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:128;not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	Location       string    `gorm:"size:255" json:"location"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	RegistrationDeadline time.Time `json:"registration_deadline"`
	Status         string    `gorm:"size:16;default:draft" json:"status"`
	Organizer      string    `gorm:"size:128" json:"organizer"`
	CoverImage     string    `gorm:"size:255" json:"cover_image"`
	IsPublished    bool      `gorm:"default:false" json:"is_published"`
	CreatedBy      uint      `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Items          []EventItem `gorm:"foreignKey:EventID" json:"items,omitempty"`
}

type EventItem struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	EventID       uint   `gorm:"index;not null" json:"event_id"`
	Name          string `gorm:"size:128;not null" json:"name"`
	Category      string `gorm:"size:64" json:"category"`
	Gender        string `gorm:"size:16" json:"gender"`
	MinAge        int    `gorm:"default:0" json:"min_age"`
	MaxAge        int    `gorm:"default:0" json:"max_age"`
	Quota         int    `gorm:"default:0" json:"quota"`
	WaitlistQuota int    `gorm:"default:0" json:"waitlist_quota"`
	Fee           float64 `gorm:"default:0" json:"fee"`
	Requirements  string `gorm:"type:text" json:"requirements"`
	Status        string `gorm:"size:16;default:open" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RegistrationStatus string

const (
	RegStatusPending   RegistrationStatus = "pending"
	RegStatusConfirmed RegistrationStatus = "confirmed"
	RegStatusWaitlist  RegistrationStatus = "waitlist"
	RegStatusRejected  RegistrationStatus = "rejected"
	RegStatusCancelled RegistrationStatus = "cancelled"
)

type RegistrationType string

const (
	RegTypeIndividual RegistrationType = "individual"
	RegTypeTeam       RegistrationType = "team"
)

type Registration struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	UserID       uint             `gorm:"index;not null" json:"user_id"`
	EventID      uint             `gorm:"index;not null" json:"event_id"`
	EventItemID  uint             `gorm:"index;not null" json:"event_item_id"`
	TeamName     string           `gorm:"size:128" json:"team_name"`
	TeamMembers  string           `gorm:"type:text" json:"team_members"`
	RegType      RegistrationType `gorm:"size:16;default:individual" json:"reg_type"`
	Status       RegistrationStatus `gorm:"size:16;default:pending" json:"status"`
	QueuePosition int             `gorm:"default:0" json:"queue_position"`
	PaymentStatus string          `gorm:"size:16;default:unpaid" json:"payment_status"`
	Amount       float64          `gorm:"default:0" json:"amount"`
	Remark       string           `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	User         *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Event        *Event           `gorm:"foreignKey:EventID" json:"event,omitempty"`
	EventItem    *EventItem       `gorm:"foreignKey:EventItemID" json:"event_item,omitempty"`
}

type Score struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EventID       uint      `gorm:"index;not null" json:"event_id"`
	EventItemID   uint      `gorm:"index;not null" json:"event_item_id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	RegistrationID uint     `gorm:"index;not null" json:"registration_id"`
	Score         float64   `gorm:"default:0" json:"score"`
	Rank          int       `gorm:"default:0" json:"rank"`
	Points        float64   `gorm:"default:0" json:"points"`
	TimeUsed      string    `gorm:"size:32" json:"time_used"`
	Remarks       string    `gorm:"type:text" json:"remarks"`
	IsValid       bool      `gorm:"default:true" json:"is_valid"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Certificate struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	EventID         uint      `gorm:"index;not null" json:"event_id"`
	EventItemID     uint      `gorm:"index;not null" json:"event_item_id"`
	ScoreID         uint      `gorm:"index" json:"score_id"`
	CertificateNo   string    `gorm:"size:64;uniqueIndex" json:"certificate_no"`
	CertificateName string    `gorm:"size:128" json:"certificate_name"`
	Rank            int       `json:"rank"`
	Score           float64   `json:"score"`
	FilePath        string    `gorm:"size:255" json:"file_path"`
	Status          string    `gorm:"size:16;default:generating" json:"status"`
	RetryCount      int       `gorm:"default:0" json:"retry_count"`
	GeneratedAt     *time.Time `json:"generated_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MessageType string

const (
	MsgTypeSystem        MessageType = "system"
	MsgTypeRegSuccess    MessageType = "reg_success"
	MsgTypeRegWaitlist   MessageType = "reg_waitlist"
	MsgTypeScorePublished MessageType = "score_published"
	MsgTypeCertificate   MessageType = "certificate"
)

type Message struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `gorm:"index;not null" json:"user_id"`
	Type      MessageType `gorm:"size:32" json:"type"`
	Title     string      `gorm:"size:128" json:"title"`
	Content   string      `gorm:"type:text" json:"content"`
	IsRead    bool        `gorm:"default:false" json:"is_read"`
	Extra     string      `gorm:"type:text" json:"extra"`
	CreatedAt time.Time   `json:"created_at"`
}

type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Module    string    `gorm:"size:64" json:"module"`
	Action    string    `gorm:"size:64" json:"action"`
	Target    string    `gorm:"size:255" json:"target"`
	IP        string    `gorm:"size:64" json:"ip"`
	Detail    string    `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}
