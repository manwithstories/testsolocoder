package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleCompany   UserRole = "company"
	RoleApplicant UserRole = "applicant"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email" binding:"required,email"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Role         UserRole       `gorm:"size:20;not null;default:applicant" json:"role"`
	Status       UserStatus     `gorm:"size:20;not null;default:active" json:"status"`
	CompanyID    *uint          `gorm:"index" json:"company_id,omitempty"`
	Company      *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Profile      *ApplicantProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LoginAttempts int           `gorm:"default:0" json:"-"`
	LockedUntil  *time.Time     `json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type ApplicantProfile struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	FullName    string         `gorm:"size:100;not null" json:"full_name" binding:"required"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Avatar      string         `gorm:"size:500" json:"avatar"`
	Gender      string         `gorm:"size:10" json:"gender"`
	BirthDate   *time.Time     `json:"birth_date"`
	Location    string         `gorm:"size:100" json:"location"`
	Education   string         `gorm:"size:500" json:"education"`
	Experience  string         `gorm:"size:1000" json:"experience"`
	Skills      string         `gorm:"size:500" json:"skills"`
	Summary     string         `gorm:"size:1000" json:"summary"`
	ResumeFile  string         `gorm:"size:500" json:"resume_file"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Company struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name" binding:"required"`
	Logo        string         `gorm:"size:500" json:"logo"`
	Industry    string         `gorm:"size:100" json:"industry"`
	Size        string         `gorm:"size:50" json:"size"`
	Address     string         `gorm:"size:300" json:"address"`
	Website     string         `gorm:"size:200" json:"website"`
	Description string         `gorm:"size:2000" json:"description"`
	Verified    bool           `gorm:"default:false" json:"verified"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type JobStatus string

const (
	JobStatusOpen      JobStatus = "open"
	JobStatusClosed    JobStatus = "closed"
	JobStatusPaused    JobStatus = "paused"
	JobStatusDraft     JobStatus = "draft"
)

type Job struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CompanyID     uint           `gorm:"index;not null" json:"company_id"`
	Company       *Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Title         string         `gorm:"size:200;not null" json:"title" binding:"required"`
	Description   string         `gorm:"type:text;not null" json:"description" binding:"required"`
	SalaryMin     int            `json:"salary_min"`
	SalaryMax     int            `json:"salary_max"`
	SalaryType    string         `gorm:"size:20;default:monthly" json:"salary_type"`
	Location      string         `gorm:"size:200;not null" json:"location" binding:"required"`
	JobType       string         `gorm:"size:50;not null;default:full-time" json:"job_type"`
	Experience    string         `gorm:"size:50" json:"experience"`
	Education     string         `gorm:"size:50" json:"education"`
	Skills        string         `gorm:"size:500" json:"skills"`
	Requirements  string         `gorm:"type:text" json:"requirements"`
	Benefits      string         `gorm:"size:500" json:"benefits"`
	Deadline      *time.Time     `json:"deadline"`
	Status        JobStatus      `gorm:"size:20;not null;default:draft" json:"status"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	ApplyCount    int            `gorm:"default:0" json:"apply_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Resume struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	FullName     string         `gorm:"size:100;not null" json:"full_name"`
	Email        string         `gorm:"size:255" json:"email"`
	Phone        string         `gorm:"size:20" json:"phone"`
	Location     string         `gorm:"size:200" json:"location"`
	Education    string         `gorm:"type:text" json:"education"`
	Experience   string         `gorm:"type:text" json:"experience"`
	Skills       string         `gorm:"type:text" json:"skills"`
	Summary      string         `gorm:"type:text" json:"summary"`
	Projects     string         `gorm:"type:text" json:"projects"`
	FileURL      string         `gorm:"size:500" json:"file_url"`
	FileType     string         `gorm:"size:50" json:"file_type"`
	FileSize     int64          `json:"file_size"`
	IsDefault    bool           `gorm:"default:false" json:"is_default"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type ApplicationStatus string

const (
	ApplicationStatusPending    ApplicationStatus = "pending"
	ApplicationStatusViewed     ApplicationStatus = "viewed"
	ApplicationStatusInterested ApplicationStatus = "interested"
	ApplicationStatusInterview  ApplicationStatus = "interview"
	ApplicationStatusAccepted   ApplicationStatus = "accepted"
	ApplicationStatusRejected   ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn  ApplicationStatus = "withdrawn"
)

type Application struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	JobID         uint              `gorm:"index;not null" json:"job_id"`
	Job           *Job              `gorm:"foreignKey:JobID" json:"job,omitempty"`
	ApplicantID   uint              `gorm:"index;not null" json:"applicant_id"`
	Applicant     *User             `gorm:"foreignKey:ApplicantID" json:"applicant,omitempty"`
	ResumeID      uint              `gorm:"index;not null" json:"resume_id"`
	Resume        *Resume           `gorm:"foreignKey:ResumeID" json:"resume,omitempty"`
	Status        ApplicationStatus `gorm:"size:20;not null;default:pending" json:"status"`
	CoverLetter   string            `gorm:"type:text" json:"cover_letter"`
	HRNote        string            `gorm:"type:text" json:"hr_note"`
	AppliedAt     time.Time         `json:"applied_at"`
	LastUpdateAt  time.Time         `json:"last_update_at"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
}

type ApplicationHistory struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	ApplicationID uint              `gorm:"index;not null" json:"application_id"`
	OldStatus     ApplicationStatus `gorm:"size:20" json:"old_status"`
	NewStatus     ApplicationStatus `gorm:"size:20" json:"new_status"`
	ChangedBy     uint              `json:"changed_by"`
	ChangeReason  string            `gorm:"size:500" json:"change_reason"`
	CreatedAt     time.Time         `json:"created_at"`
}

type InterviewStatus string

const (
	InterviewStatusPending  InterviewStatus = "pending"
	InterviewStatusAccepted InterviewStatus = "accepted"
	InterviewStatusRejected InterviewStatus = "rejected"
	InterviewStatusCompleted InterviewStatus = "completed"
	InterviewStatusCancelled InterviewStatus = "cancelled"
)

type Interview struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ApplicationID uint           `gorm:"index;not null" json:"application_id"`
	Application   *Application   `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	JobID         uint           `gorm:"index;not null" json:"job_id"`
	ApplicantID   uint           `gorm:"index;not null" json:"applicant_id"`
	Interviewer   string         `gorm:"size:100" json:"interviewer"`
	InterviewType string         `gorm:"size:50;default:onsite" json:"interview_type"`
	Location      string         `gorm:"size:200" json:"location"`
	MeetingLink   string         `gorm:"size:500" json:"meeting_link"`
	ScheduledAt   time.Time      `json:"scheduled_at"`
	Duration      int            `gorm:"default:60" json:"duration"`
	Status        InterviewStatus `gorm:"size:20;not null;default:pending" json:"status"`
	Notes         string         `gorm:"type:text" json:"notes"`
	Feedback      string         `gorm:"type:text" json:"feedback"`
	Rating        int            `gorm:"default:0" json:"rating"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type JobViewLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JobID     uint      `gorm:"index;not null" json:"job_id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	IPAddress string    `gorm:"size:50" json:"ip_address"`
	ViewedAt  time.Time `json:"viewed_at"`
}

type DailyStatistics struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Date            time.Time `gorm:"uniqueIndex" json:"date"`
	NewJobs         int       `gorm:"default:0" json:"new_jobs"`
	NewApplications int       `gorm:"default:0" json:"new_applications"`
	NewUsers        int       `gorm:"default:0" json:"new_users"`
	TotalViews      int       `gorm:"default:0" json:"total_views"`
	Interviews      int       `gorm:"default:0" json:"interviews"`
	Hires           int       `gorm:"default:0" json:"hires"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
