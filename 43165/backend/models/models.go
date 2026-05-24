package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleEmployer  UserRole = "employer"
	RoleAgent     UserRole = "agent"
	RoleTemporary UserRole = "temporary"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Username      string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email         string     `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone         string     `gorm:"size:20" json:"phone"`
	Password      string     `gorm:"size:255;not null" json:"-"`
	RealName      string     `gorm:"size:50" json:"real_name"`
	Role          UserRole   `gorm:"size:20;not null" json:"role"`
	Avatar        string     `gorm:"size:255" json:"avatar"`
	CreditScore   int        `gorm:"default:100" json:"credit_score"`
	Status        UserStatus `gorm:"size:20;default:active" json:"status"`
	Company       string     `gorm:"size:100" json:"company"`
	Address       string     `gorm:"size:255" json:"address"`
	IDCard        string     `gorm:"size:20" json:"id_card"`
	FaceData      string     `gorm:"type:text" json:"-"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	LastLoginIP   string     `gorm:"size:50" json:"last_login_ip"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type JobPosting struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	EmployerID     uuid.UUID    `gorm:"type:uuid;index;not null" json:"employer_id"`
	Employer       *User        `gorm:"foreignKey:EmployerID" json:"employer,omitempty"`
	Title          string       `gorm:"size:100;not null" json:"title"`
	Description    string       `gorm:"type:text;not null" json:"description"`
	ActivityType   string       `gorm:"size:50;index" json:"activity_type"`
	Position       string       `gorm:"size:50;not null" json:"position"`
	Location       string       `gorm:"size:255;not null" json:"location"`
	Latitude       float64      `gorm:"default:0" json:"latitude"`
	Longitude      float64      `gorm:"default:0" json:"longitude"`
	StartDate      time.Time    `gorm:"not null" json:"start_date"`
	EndDate        time.Time    `gorm:"not null" json:"end_date"`
	SalaryPerHour  float64      `gorm:"not null" json:"salary_per_hour"`
	SalaryType     string       `gorm:"size:20;default:hourly" json:"salary_type"`
	WorkHours      string       `gorm:"size:100" json:"work_hours"`
	Headcount      int          `gorm:"not null;default:1" json:"headcount"`
	Applicants     int          `gorm:"default:0" json:"applicants"`
	HiredCount     int          `gorm:"default:0" json:"hired_count"`
	Requirements   string       `gorm:"type:text" json:"requirements"`
	Benefits       string       `gorm:"type:text" json:"benefits"`
	ContactPerson  string       `gorm:"size:50" json:"contact_person"`
	ContactPhone   string       `gorm:"size:20" json:"contact_phone"`
	Status         string       `gorm:"size:20;default:recruiting" json:"status"`
	Tags           string       `gorm:"size:255" json:"tags"`
	IsUrgent       bool         `gorm:"default:false" json:"is_urgent"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func (j *JobPosting) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

type JobApplication struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	JobID         uuid.UUID   `gorm:"type:uuid;index;not null" json:"job_id"`
	JobPosting    *JobPosting `gorm:"foreignKey:JobID" json:"job_posting,omitempty"`
	TemporaryID   uuid.UUID   `gorm:"type:uuid;index;not null" json:"temporary_id"`
	Temporary     *User       `gorm:"foreignKey:TemporaryID" json:"temporary,omitempty"`
	AgentID       *uuid.UUID  `gorm:"type:uuid;index" json:"agent_id"`
	Agent         *User       `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Message       string      `gorm:"type:text" json:"message"`
	Status        string      `gorm:"size:20;default:pending" json:"status"`
	AppliedAt     time.Time   `json:"applied_at"`
	ReviewedAt    *time.Time  `json:"reviewed_at"`
	ReviewNote    string      `gorm:"type:text" json:"review_note"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (a *JobApplication) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type Schedule struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	JobID        uuid.UUID   `gorm:"type:uuid;index;not null" json:"job_id"`
	JobPosting   *JobPosting `gorm:"foreignKey:JobID" json:"job_posting,omitempty"`
	TemporaryID  uuid.UUID   `gorm:"type:uuid;index;not null" json:"temporary_id"`
	Temporary    *User       `gorm:"foreignKey:TemporaryID" json:"temporary,omitempty"`
	ShiftDate    time.Time   `gorm:"index;not null" json:"shift_date"`
	StartTime    string      `gorm:"size:10;not null" json:"start_time"`
	EndTime      string      `gorm:"size:10;not null" json:"end_time"`
	Location     string      `gorm:"size:255" json:"location"`
	Notes        string      `gorm:"type:text" json:"notes"`
	Status       string      `gorm:"size:20;default:scheduled" json:"status"`
	CreatedBy    uuid.UUID   `gorm:"type:uuid" json:"created_by"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type CheckIn struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ScheduleID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"schedule_id"`
	Schedule     *Schedule  `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	TemporaryID  uuid.UUID  `gorm:"type:uuid;index;not null" json:"temporary_id"`
	Temporary    *User      `gorm:"foreignKey:TemporaryID" json:"temporary,omitempty"`
	CheckInType  string     `gorm:"size:20;not null" json:"check_in_type"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	Location     string     `gorm:"size:255" json:"location"`
	FaceVerified bool       `gorm:"default:false" json:"face_verified"`
	QRCode       string     `gorm:"size:100" json:"qr_code"`
	Status       string     `gorm:"size:20;default:checked_in" json:"status"`
	Remarks      string     `gorm:"type:text" json:"remarks"`
	WorkHours    float64    `gorm:"default:0" json:"work_hours"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (c *CheckIn) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type SalaryRecord struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	TemporaryID   uuid.UUID   `gorm:"type:uuid;index;not null" json:"temporary_id"`
	Temporary     *User       `gorm:"foreignKey:TemporaryID" json:"temporary,omitempty"`
	EmployerID    uuid.UUID   `gorm:"type:uuid;index;not null" json:"employer_id"`
	Employer      *User       `gorm:"foreignKey:EmployerID" json:"employer,omitempty"`
	JobID         uuid.UUID   `gorm:"type:uuid;index;not null" json:"job_id"`
	JobPosting    *JobPosting `gorm:"foreignKey:JobID" json:"job_posting,omitempty"`
	PeriodStart   time.Time   `gorm:"not null" json:"period_start"`
	PeriodEnd     time.Time   `gorm:"not null" json:"period_end"`
	TotalHours    float64     `gorm:"default:0" json:"total_hours"`
	BaseSalary    float64     `gorm:"default:0" json:"base_salary"`
	OvertimeHours float64     `gorm:"default:0" json:"overtime_hours"`
	OvertimePay   float64     `gorm:"default:0" json:"overtime_pay"`
	Deductions    float64     `gorm:"default:0" json:"deductions"`
	TotalSalary   float64     `gorm:"default:0" json:"total_salary"`
	Status        string      `gorm:"size:20;default:pending" json:"status"`
	PaymentMethod string      `gorm:"size:20" json:"payment_method"`
	PaymentAt     *time.Time  `json:"payment_at"`
	TransactionID string      `gorm:"size:100" json:"transaction_id"`
	Remark        string      `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (s *SalaryRecord) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type SalaryDetail struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	SalaryID      uuid.UUID     `gorm:"type:uuid;index;not null" json:"salary_id"`
	SalaryRecord  *SalaryRecord `gorm:"foreignKey:SalaryID" json:"salary_record,omitempty"`
	CheckInID     uuid.UUID     `gorm:"type:uuid;index" json:"check_in_id"`
	CheckIn       *CheckIn      `gorm:"foreignKey:CheckInID" json:"check_in,omitempty"`
	Date          time.Time     `json:"date"`
	WorkHours     float64       `gorm:"default:0" json:"work_hours"`
	HourlyRate    float64       `gorm:"default:0" json:"hourly_rate"`
	Amount        float64       `gorm:"default:0" json:"amount"`
	Type          string        `gorm:"size:20" json:"type"`
	Description   string        `gorm:"size:255" json:"description"`
	CreatedAt     time.Time     `json:"created_at"`
}

func (d *SalaryDetail) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type Evaluation struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	JobID         uuid.UUID  `gorm:"type:uuid;index;not null" json:"job_id"`
	JobPosting    *JobPosting `gorm:"foreignKey:JobID" json:"job_posting,omitempty"`
	FromUserID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"from_user_id"`
	FromUser      *User      `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUserID      uuid.UUID  `gorm:"type:uuid;index;not null" json:"to_user_id"`
	ToUser        *User      `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
	Rating        int        `gorm:"not null;default:5" json:"rating"`
	Content       string     `gorm:"type:text" json:"content"`
	Tags          string     `gorm:"size:255" json:"tags"`
	Type          string     `gorm:"size:20;not null" json:"type"`
	IsAnonymous   bool       `gorm:"default:false" json:"is_anonymous"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (e *Evaluation) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type JobTemplate struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	EmployerID    uuid.UUID `gorm:"type:uuid;index;not null" json:"employer_id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	ActivityType  string    `gorm:"size:50" json:"activity_type"`
	Position      string    `gorm:"size:50;not null" json:"position"`
	Description   string    `gorm:"type:text" json:"description"`
	SalaryPerHour float64   `gorm:"default:0" json:"salary_per_hour"`
	WorkHours     string    `gorm:"size:100" json:"work_hours"`
	Requirements  string    `gorm:"type:text" json:"requirements"`
	Benefits      string    `gorm:"type:text" json:"benefits"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (t *JobTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Title     string    `gorm:"size:100;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Type      string    `gorm:"size:20;default:system" json:"type"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	RelatedID *uuid.UUID `gorm:"type:uuid" json:"related_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
