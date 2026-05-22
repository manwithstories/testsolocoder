package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleDoctor  UserRole = "doctor"
	RolePatient UserRole = "patient"
)

type User struct {
	BaseModel
	Username     string   `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string   `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone        string   `gorm:"size:20" json:"phone"`
	PasswordHash string   `gorm:"size:255;not null" json:"-"`
	Role         UserRole `gorm:"size:20;not null;index" json:"role"`
	FullName     string   `gorm:"size:50;not null" json:"full_name"`
	Gender       string   `gorm:"size:10" json:"gender"`
	BirthDate    *time.Time `json:"birth_date"`
	AvatarURL    string   `gorm:"size:255" json:"avatar_url"`
	IsActive     bool     `gorm:"default:true" json:"is_active"`
}

type Department struct {
	BaseModel
	Name        string `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Location    string `gorm:"size:200" json:"location"`
	Doctors     []Doctor `gorm:"foreignKey:DepartmentID" json:"doctors,omitempty"`
}

type DoctorTitle string

const (
	TitleIntern       DoctorTitle = "住院医师"
	TitleResident     DoctorTitle = "主治医师"
	TitleAttending    DoctorTitle = "副主任医师"
	TitleChief        DoctorTitle = "主任医师"
	TitleProfessor    DoctorTitle = "教授"
)

type Doctor struct {
	BaseModel
	UserID         uint        `gorm:"uniqueIndex;not null" json:"user_id"`
	User           User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DepartmentID   uint        `gorm:"index;not null" json:"department_id"`
	Department     Department  `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Title          DoctorTitle `gorm:"size:50;not null" json:"title"`
	Specialty      string      `gorm:"size:200" json:"specialty"`
	Introduction   string      `gorm:"type:text" json:"introduction"`
	ConsultationFee float64    `gorm:"type:decimal(10,2);default:0" json:"consultation_fee"`
	RegistrationFee float64    `gorm:"type:decimal(10,2);default:0" json:"registration_fee"`
	AverageRating  float64     `gorm:"type:decimal(3,2);default:0" json:"average_rating"`
	ReviewCount    int         `gorm:"default:0" json:"review_count"`
	Schedules      []Schedule  `gorm:"foreignKey:DoctorID" json:"schedules,omitempty"`
	Appointments   []Appointment `gorm:"foreignKey:DoctorID" json:"appointments,omitempty"`
	Reviews        []Review    `gorm:"foreignKey:DoctorID" json:"reviews,omitempty"`
}

type DayOfWeek int

const (
	Sunday    DayOfWeek = 0
	Monday    DayOfWeek = 1
	Tuesday   DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday  DayOfWeek = 4
	Friday    DayOfWeek = 5
	Saturday  DayOfWeek = 6
)

type Schedule struct {
	BaseModel
	DoctorID      uint      `gorm:"index;not null" json:"doctor_id"`
	Doctor        Doctor    `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	DayOfWeek     DayOfWeek `gorm:"index;not null" json:"day_of_week"`
	StartTime     string    `gorm:"size:8;not null" json:"start_time"`
	EndTime       string    `gorm:"size:8;not null" json:"end_time"`
	MaxPatients   int       `gorm:"default:20;not null" json:"max_patients"`
	TimeSlot      int       `gorm:"default:15" json:"time_slot_minutes"`
	IsAvailable   bool      `gorm:"default:true" json:"is_available"`
}

type Patient struct {
	BaseModel
	UserID         uint        `gorm:"uniqueIndex;not null" json:"user_id"`
	User           User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IDCardNo       string      `gorm:"size:18;uniqueIndex" json:"id_card_no"`
	Address        string      `gorm:"size:500" json:"address"`
	EmergencyName  string      `gorm:"size:50" json:"emergency_contact_name"`
	EmergencyPhone string      `gorm:"size:20" json:"emergency_contact_phone"`
	HealthRecord   HealthRecord `gorm:"foreignKey:PatientID" json:"health_record,omitempty"`
	Appointments   []Appointment `gorm:"foreignKey:PatientID" json:"appointments,omitempty"`
}

type HealthRecord struct {
	BaseModel
	PatientID      uint           `gorm:"uniqueIndex;not null" json:"patient_id"`
	MedicalHistory JSONB          `gorm:"type:jsonb" json:"medical_history"`
	Allergies      JSONB          `gorm:"type:jsonb" json:"allergies"`
	Medications    JSONB          `gorm:"type:jsonb" json:"medications"`
	Vaccinations   JSONB          `gorm:"type:jsonb" json:"vaccinations"`
	FamilyHistory  string         `gorm:"type:text" json:"family_history"`
	LifeHabits     string         `gorm:"type:text" json:"life_habits"`
	Remarks        string         `gorm:"type:text" json:"remarks"`
}

type JSONB json.RawMessage

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSONB(result)
	return err
}

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

type AppointmentStatus string

const (
	AppointmentPending   AppointmentStatus = "pending"
	AppointmentConfirmed AppointmentStatus = "confirmed"
	AppointmentCompleted AppointmentStatus = "completed"
	AppointmentCancelled AppointmentStatus = "cancelled"
	AppointmentNoShow    AppointmentStatus = "no_show"
)

type Appointment struct {
	BaseModel
	PatientID     uint              `gorm:"index;not null" json:"patient_id"`
	Patient       Patient           `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID      uint              `gorm:"index;not null" json:"doctor_id"`
	Doctor        Doctor            `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	AppointmentDate time.Time        `gorm:"not null;index" json:"appointment_date"`
	StartTime     string            `gorm:"size:8;not null" json:"start_time"`
	EndTime       string            `gorm:"size:8;not null" json:"end_time"`
	Status        AppointmentStatus `gorm:"size:20;not null;default:pending;index" json:"status"`
	Symptoms      string            `gorm:"type:text" json:"symptoms"`
	Notes         string            `gorm:"type:text" json:"notes"`
	CancelReason  string            `gorm:"type:text" json:"cancel_reason,omitempty"`
	Consultation  *Consultation     `gorm:"foreignKey:AppointmentID" json:"consultation,omitempty"`
	Payment       *Payment          `gorm:"foreignKey:AppointmentID" json:"payment,omitempty"`
	Review        *Review           `gorm:"foreignKey:AppointmentID" json:"review,omitempty"`
}

type Consultation struct {
	BaseModel
	AppointmentID   uint               `gorm:"uniqueIndex;not null" json:"appointment_id"`
	Appointment     Appointment        `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
	Diagnosis       string             `gorm:"type:text;not null" json:"diagnosis"`
	TreatmentPlan   string             `gorm:"type:text" json:"treatment_plan"`
	DoctorNotes     string             `gorm:"type:text" json:"doctor_notes"`
	FollowUpDate    *time.Time         `json:"follow_up_date"`
	Prescription    *Prescription      `gorm:"foreignKey:ConsultationID" json:"prescription,omitempty"`
	Reports         []ExaminationReport `gorm:"foreignKey:ConsultationID" json:"reports,omitempty"`
}

type Prescription struct {
	BaseModel
	ConsultationID uint              `gorm:"uniqueIndex;not null" json:"consultation_id"`
	Consultation   Consultation      `gorm:"foreignKey:ConsultationID" json:"consultation,omitempty"`
	PrescriptionNo string            `gorm:"size:50;uniqueIndex;not null" json:"prescription_no"`
	Items          []PrescriptionItem `gorm:"foreignKey:PrescriptionID" json:"items,omitempty"`
	Notes          string            `gorm:"type:text" json:"notes"`
	IsFulfilled    bool              `gorm:"default:false" json:"is_fulfilled"`
	FulfilledAt    *time.Time        `json:"fulfilled_at"`
}

type PrescriptionItem struct {
	BaseModel
	PrescriptionID uint    `gorm:"index;not null" json:"prescription_id"`
	Prescription   Prescription `gorm:"foreignKey:PrescriptionID" json:"prescription,omitempty"`
	DrugName       string  `gorm:"size:200;not null" json:"drug_name"`
	Specification  string  `gorm:"size:200" json:"specification"`
	Dosage         string  `gorm:"size:200;not null" json:"dosage"`
	Frequency      string  `gorm:"size:100" json:"frequency"`
	Duration       string  `gorm:"size:100" json:"duration"`
	Quantity       int     `gorm:"default:1" json:"quantity"`
	UnitPrice      float64 `gorm:"type:decimal(10,2);default:0" json:"unit_price"`
	Subtotal       float64 `gorm:"type:decimal(10,2);default:0" json:"subtotal"`
	Notes          string  `gorm:"type:text" json:"notes"`
}

type ExaminationReport struct {
	BaseModel
	ConsultationID uint      `gorm:"index;not null" json:"consultation_id"`
	Consultation   Consultation `gorm:"foreignKey:ConsultationID" json:"consultation,omitempty"`
	ReportType     string    `gorm:"size:100;not null" json:"report_type"`
	ReportName     string    `gorm:"size:200;not null" json:"report_name"`
	FileURL        string    `gorm:"size:500;not null" json:"file_url"`
	FileSize       int64     `json:"file_size"`
	ContentType    string    `gorm:"size:100" json:"content_type"`
	UploadedBy     uint      `json:"uploaded_by"`
	Findings       string    `gorm:"type:text" json:"findings"`
	Conclusion     string    `gorm:"type:text" json:"conclusion"`
}

type NotificationType string

const (
	NotificationAppointmentConfirmation NotificationType = "appointment_confirmation"
	NotificationAppointmentReminder     NotificationType = "appointment_reminder"
	NotificationAppointmentCancelled    NotificationType = "appointment_cancelled"
	NotificationConsultationCompleted   NotificationType = "consultation_completed"
	NotificationPaymentSuccess          NotificationType = "payment_success"
	NotificationSystem                  NotificationType = "system"
)

type NotificationChannel string

const (
	ChannelEmail  NotificationChannel = "email"
	ChannelInApp  NotificationChannel = "in_app"
	ChannelSMS    NotificationChannel = "sms"
)

type Notification struct {
	BaseModel
	UserID         uint                `gorm:"index;not null" json:"user_id"`
	User           User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type           NotificationType    `gorm:"size:50;index" json:"type"`
	Title          string              `gorm:"size:200;not null" json:"title"`
	Content        string              `gorm:"type:text;not null" json:"content"`
	Channel        NotificationChannel `gorm:"size:20" json:"channel"`
	IsRead         bool                `gorm:"default:false;index" json:"is_read"`
	ReadAt         *time.Time          `json:"read_at"`
	RelatedID      uint                `json:"related_id"`
	RelatedType    string              `gorm:"size:50" json:"related_type"`
}

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentPaid      PaymentStatus = "paid"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentMethodWeChat  PaymentMethod = "wechat"
	PaymentMethodAlipay  PaymentMethod = "alipay"
	PaymentMethodCard    PaymentMethod = "card"
	PaymentMethodCash    PaymentMethod = "cash"
)

type Payment struct {
	BaseModel
	AppointmentID   uint          `gorm:"uniqueIndex;not null" json:"appointment_id"`
	Appointment     Appointment   `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
	TransactionNo   string        `gorm:"size:100;uniqueIndex" json:"transaction_no"`
	RegistrationFee float64       `gorm:"type:decimal(10,2);default:0" json:"registration_fee"`
	ConsultationFee float64       `gorm:"type:decimal(10,2);default:0" json:"consultation_fee"`
	DrugFee         float64       `gorm:"type:decimal(10,2);default:0" json:"drug_fee"`
	ExaminationFee  float64       `gorm:"type:decimal(10,2);default:0" json:"examination_fee"`
	OtherFee        float64       `gorm:"type:decimal(10,2);default:0" json:"other_fee"`
	TotalAmount     float64       `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status          PaymentStatus `gorm:"size:20;not null;default:pending;index" json:"status"`
	Method          PaymentMethod `gorm:"size:20" json:"method"`
	PaidAt          *time.Time    `json:"paid_at"`
	RefundedAt      *time.Time    `json:"refunded_at"`
	Notes           string        `gorm:"type:text" json:"notes"`
}

type Review struct {
	BaseModel
	AppointmentID uint      `gorm:"uniqueIndex;not null" json:"appointment_id"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
	PatientID     uint      `gorm:"index;not null" json:"patient_id"`
	Patient       Patient   `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID      uint      `gorm:"index;not null" json:"doctor_id"`
	Doctor        Doctor    `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	Rating        int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Content       string    `gorm:"type:text" json:"content"`
	IsAnonymous   bool      `gorm:"default:false" json:"is_anonymous"`
	IsVerified    bool      `gorm:"default:true" json:"is_verified"`
}
