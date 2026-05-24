package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin       UserRole = "admin"
	RoleEntrepreneur UserRole = "entrepreneur"
	RoleAgent       UserRole = "agent"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:100;not null"`
	Password     string         `json:"-" gorm:"size:255;not null"`
	RealName     string         `json:"realName" gorm:"size:100"`
	Email        string         `json:"email" gorm:"size:100"`
	Phone        string         `json:"phone" gorm:"size:20"`
	Role         UserRole       `json:"role" gorm:"size:20;not null;index"`
	Status       UserStatus     `json:"status" gorm:"size:20;default:active"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	LastLoginAt  *time.Time     `json:"lastLoginAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	AgentProfile *AgentProfile  `json:"agentProfile,omitempty" gorm:"foreignKey:UserID"`
}

type AgentProfile struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"userId" gorm:"uniqueIndex;not null"`
	EmployeeNo      string         `json:"employeeNo" gorm:"uniqueIndex;size:50;not null"`
	SpecialtyTags   string         `json:"specialtyTags" gorm:"size:500"`
	MaxApplications int            `json:"maxApplications" gorm:"default:5"`
	CurrentApps     int            `json:"currentApps" gorm:"default:0"`
	WorkStartTime   string         `json:"workStartTime" gorm:"size:10"`
	WorkEndTime     string         `json:"workEndTime" gorm:"size:10"`
	Status          string         `json:"status" gorm:"size:20;default:available"`
	PerformanceScore float64       `json:"performanceScore" gorm:"default:0"`
	TotalHandled    int            `json:"totalHandled" gorm:"default:0"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type ApplicationStatus string

const (
	AppStatusDraft        ApplicationStatus = "draft"
	AppStatusPendingReview ApplicationStatus = "pending_review"
	AppStatusReviewing    ApplicationStatus = "reviewing"
	AppStatusProcessing   ApplicationStatus = "processing"
	AppStatusCompleted    ApplicationStatus = "completed"
	AppStatusRejected     ApplicationStatus = "rejected"
	AppStatusCancelled    ApplicationStatus = "cancelled"
	AppStatusPaymentPending ApplicationStatus = "payment_pending"
)

type CompanyType string

const (
	CompanyTypeLLC    CompanyType = "llc"
	CompanyTypeJointStock CompanyType = "joint_stock"
	CompanyTypeSole   CompanyType = "sole"
	CompanyTypePartnership CompanyType = "partnership"
)

type Application struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	ApplicationNo    string            `json:"applicationNo" gorm:"uniqueIndex;size:50;not null"`
	EntrepreneurID   uint              `json:"entrepreneurId" gorm:"index;not null"`
	AgentID          *uint             `json:"agentId" gorm:"index"`
	CompanyName      string            `json:"companyName" gorm:"size:200;not null"`
	CompanyType      CompanyType       `json:"companyType" gorm:"size:50;not null"`
	RegisteredCapital float64           `json:"registeredCapital" gorm:"default:0"`
	BusinessScope    string            `json:"businessScope" gorm:"type:text"`
	RegisteredAddress string           `json:"registeredAddress" gorm:"size:500"`
	ShareholderInfo  string            `json:"shareholderInfo" gorm:"type:text"`
	IDCardFront      string            `json:"idCardFront" gorm:"size:255"`
	IDCardBack       string            `json:"idCardBack" gorm:"size:255"`
	LicensePreview   string            `json:"licensePreview" gorm:"size:255"`
	OtherMaterials   string            `json:"otherMaterials" gorm:"size:1000"`
	Status           ApplicationStatus `json:"status" gorm:"size:50;index;not null"`
	ReviewComments   string            `json:"reviewComments" gorm:"type:text"`
	CurrentStep      string            `json:"currentStep" gorm:"size:50"`
	ProgressPercent  int               `json:"progressPercent" gorm:"default:0"`
	SubmittedAt      *time.Time        `json:"submittedAt"`
	CompletedAt      *time.Time        `json:"completedAt"`
	RejectedAt       *time.Time        `json:"rejectedAt"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt    `json:"-" gorm:"index"`
	Entrepreneur     *User             `json:"entrepreneur,omitempty" gorm:"foreignKey:EntrepreneurID"`
	Agent            *User             `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	ProcessSteps     []ProcessStep     `json:"processSteps,omitempty" gorm:"foreignKey:ApplicationID"`
	Fee              *ApplicationFee   `json:"fee,omitempty" gorm:"foreignKey:ApplicationID"`
	Notifications    []Notification    `json:"notifications,omitempty" gorm:"foreignKey:ApplicationID"`
}

type ProcessStepType string

const (
	StepTypeNaming      ProcessStepType = "naming"
	StepTypeRegistration ProcessStepType = "registration"
	StepTypeTax         ProcessStepType = "tax"
	StepTypeBank        ProcessStepType = "bank"
	StepTypeSeal        ProcessStepType = "seal"
	StepTypeCompletion  ProcessStepType = "completion"
)

type ProcessStepStatus string

const (
	StepStatusPending   ProcessStepStatus = "pending"
	StepStatusInProgress ProcessStepStatus = "in_progress"
	StepStatusCompleted ProcessStepStatus = "completed"
	StepStatusFailed    ProcessStepStatus = "failed"
	StepStatusSkipped   ProcessStepStatus = "skipped"
)

type ProcessStep struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	ApplicationID uint              `json:"applicationId" gorm:"index;not null"`
	StepType      ProcessStepType   `json:"stepType" gorm:"size:50;not null"`
	StepName      string            `json:"stepName" gorm:"size:100;not null"`
	StepOrder     int               `json:"stepOrder" gorm:"default:0"`
	Status        ProcessStepStatus `json:"status" gorm:"size:50;index;not null"`
	Description   string            `json:"description" gorm:"type:text"`
	Remark        string            `json:"remark" gorm:"type:text"`
	CertificateFile string          `json:"certificateFile" gorm:"size:255"`
	HandlerID     *uint             `json:"handlerId" gorm:"index"`
	StartedAt     *time.Time        `json:"startedAt"`
	CompletedAt   *time.Time        `json:"completedAt"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
	Application   *Application      `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
	Handler       *User             `json:"handler,omitempty" gorm:"foreignKey:HandlerID"`
}

type FeeStatus string

const (
	FeeStatusPending  FeeStatus = "pending"
	FeeStatusPaid     FeeStatus = "paid"
	FeeStatusRefunded FeeStatus = "refunded"
)

type FeeItem struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ApplicationFeeID uint         `json:"applicationFeeId" gorm:"index;not null"`
	ItemName       string         `json:"itemName" gorm:"size:100;not null"`
	Amount         float64        `json:"amount" gorm:"default:0"`
	Description    string         `json:"description" gorm:"size:500"`
	CreatedAt      time.Time      `json:"createdAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type ApplicationFee struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ApplicationID  uint           `json:"applicationId" gorm:"uniqueIndex;not null"`
	TotalAmount    float64        `json:"totalAmount" gorm:"default:0"`
	DiscountAmount float64        `json:"discountAmount" gorm:"default:0"`
	PaidAmount     float64        `json:"paidAmount" gorm:"default:0"`
	Status         FeeStatus      `json:"status" gorm:"size:50;index;not null"`
	PaymentMethod  string         `json:"paymentMethod" gorm:"size:50"`
	PaymentTime    *time.Time     `json:"paymentTime"`
	TransactionNo  string         `json:"transactionNo" gorm:"size:100"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	Application    *Application   `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
	FeeItems       []FeeItem      `json:"feeItems,omitempty" gorm:"foreignKey:ApplicationFeeID"`
}

type FeeStandard struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	CompanyType    CompanyType    `json:"companyType" gorm:"uniqueIndex;size:50;not null"`
	NamingFee      float64        `json:"namingFee" gorm:"default:0"`
	RegistrationFee float64       `json:"registrationFee" gorm:"default:0"`
	TaxFee         float64        `json:"taxFee" gorm:"default:0"`
	BankFee        float64        `json:"bankFee" gorm:"default:0"`
	SealFee        float64        `json:"sealFee" gorm:"default:0"`
	ServiceFee     float64        `json:"serviceFee" gorm:"default:0"`
	CapitalRate    float64        `json:"capitalRate" gorm:"default:0"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type DiscountPolicy struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Code           string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Type           string         `json:"type" gorm:"size:20"`
	Value          float64        `json:"value" gorm:"default:0"`
	MinAmount      float64        `json:"minAmount" gorm:"default:0"`
	MaxDiscount    float64        `json:"maxDiscount" gorm:"default:0"`
	StartDate      *time.Time     `json:"startDate"`
	EndDate        *time.Time     `json:"endDate"`
	IsActive       bool           `json:"isActive" gorm:"default:true"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type NotificationType string

const (
	NotificationTypeSystem  NotificationType = "system"
	NotificationTypeEmail   NotificationType = "email"
	NotificationTypeSMS     NotificationType = "sms"
)

type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusRead      NotificationStatus = "read"
	NotificationStatusFailed    NotificationStatus = "failed"
)

type Notification struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	UserID        uint              `json:"userId" gorm:"index;not null"`
	ApplicationID *uint             `json:"applicationId" gorm:"index"`
	Type          NotificationType  `json:"type" gorm:"size:20;not null"`
	Title         string            `json:"title" gorm:"size:200;not null"`
	Content       string            `json:"content" gorm:"type:text"`
	Status        NotificationStatus `json:"status" gorm:"size:20;index;not null"`
	IsRead        bool              `json:"isRead" gorm:"default:false"`
	SentAt        *time.Time        `json:"sentAt"`
	ReadAt        *time.Time        `json:"readAt"`
	CreatedAt     time.Time         `json:"createdAt"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
	User          *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Application   *Application      `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
}

type NotificationTemplate struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Type      string         `json:"type" gorm:"size:20"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Variables string         `json:"variables" gorm:"size:500"`
	IsActive  bool           `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ExportTask struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"userId" gorm:"index;not null"`
	Type        string         `json:"type" gorm:"size:50;not null"`
	FileName    string         `json:"fileName" gorm:"size:200"`
	FilePath    string         `json:"filePath" gorm:"size:500"`
	Status      string         `json:"status" gorm:"size:20;index;not null"`
	ExpiresAt   *time.Time     `json:"expiresAt"`
	Downloaded  bool           `json:"downloaded" gorm:"default:false"`
	DownloadedAt *time.Time    `json:"downloadedAt"`
	Params      string         `json:"params" gorm:"type:text"`
	ErrorMsg    string         `json:"errorMsg" gorm:"size:500"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type DownloadLog struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ExportID   uint           `json:"exportId" gorm:"index;not null"`
	UserID     uint           `json:"userId" gorm:"index;not null"`
	IP         string         `json:"ip" gorm:"size:50"`
	UserAgent  string         `json:"userAgent" gorm:"size:500"`
	DownloadedAt time.Time    `json:"downloadedAt"`
	CreatedAt  time.Time      `json:"createdAt"`
}

type OperationLog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"userId" gorm:"index"`
	Module      string         `json:"module" gorm:"size:50;index"`
	Action      string         `json:"action" gorm:"size:100"`
	TargetType  string         `json:"targetType" gorm:"size:50"`
	TargetID    *uint          `json:"targetId" gorm:"index"`
	Content     string         `json:"content" gorm:"type:text"`
	IP          string         `json:"ip" gorm:"size:50"`
	UserAgent   string         `json:"userAgent" gorm:"size:500"`
	Result      string         `json:"result" gorm:"size:20"`
	CreatedAt   time.Time      `json:"createdAt"`
}
