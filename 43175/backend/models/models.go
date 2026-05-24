package models

import (
	"log"
	"smart-energy-platform/config"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Username     string         `json:"username" gorm:"size:50;not null"`
	Password     string         `json:"-" gorm:"not null"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	Phone        string         `json:"phone" gorm:"size:20"`
	Status       int            `json:"status" gorm:"default:1"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Family struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"size:500"`
	OwnerID     uint           `json:"ownerId" gorm:"not null"`
	Owner       User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Members     []FamilyMember `json:"members" gorm:"foreignKey:FamilyID"`
}

type FamilyMember struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FamilyID  uint           `json:"familyId" gorm:"not null;index"`
	UserID    uint           `json:"userId" gorm:"not null;index"`
	Role      string         `json:"role" gorm:"size:20;not null;default:member"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Family    Family         `json:"-" gorm:"foreignKey:FamilyID"`
}

type Device struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	FamilyID       uint           `json:"familyId" gorm:"not null;index"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	DeviceType     string         `json:"deviceType" gorm:"size:50;not null"`
	Vendor         string         `json:"vendor" gorm:"size:100"`
	Location       string         `json:"location" gorm:"size:100"`
	Power          float64        `json:"power" gorm:"default:0"`
	Protocol       string         `json:"protocol" gorm:"size:50"`
	Status         string         `json:"status" gorm:"size:20;default:offline"`
	LastOnlineTime *time.Time     `json:"lastOnlineTime"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	EnergyData     []EnergyData   `json:"-" gorm:"foreignKey:DeviceID"`
}

type EnergyData struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DeviceID    uint      `json:"deviceId" gorm:"not null;index"`
	FamilyID    uint      `json:"familyId" gorm:"not null;index"`
	Power       float64   `json:"power" gorm:"default:0"`
	Voltage     float64   `json:"voltage" gorm:"default:0"`
	Current     float64   `json:"current" gorm:"default:0"`
	EnergyUsed  float64   `json:"energyUsed" gorm:"default:0"`
	Timestamp   time.Time `json:"timestamp" gorm:"index"`
	Date        string    `json:"date" gorm:"size:10;index"`
	Hour        int       `json:"hour" gorm:"index"`
}

type EnergyAlert struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FamilyID    uint      `json:"familyId" gorm:"not null;index"`
	DeviceID    uint      `json:"deviceId" gorm:"index"`
	AlertType   string    `json:"alertType" gorm:"size:50;not null"`
	Level       string    `json:"level" gorm:"size:20;not null;default:warning"`
	Message     string    `json:"message" gorm:"size:500"`
	Value       float64   `json:"value" gorm:"default:0"`
	Threshold   float64   `json:"threshold" gorm:"default:0"`
	Resolved    bool      `json:"resolved" gorm:"default:false"`
	CreatedAt   time.Time `json:"createdAt"`
	ResolvedAt  *time.Time `json:"resolvedAt"`
}

type DeviceGroup struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FamilyID    uint           `json:"familyId" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"size:500"`
	Type        string         `json:"type" gorm:"size:20;default:room"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Devices     []Device       `json:"devices" gorm:"many2many:group_devices;"`
}

type Scene struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FamilyID    uint           `json:"familyId" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"size:500"`
	Icon        string         `json:"icon" gorm:"size:50"`
	IsActive    bool           `json:"isActive" gorm:"default:true"`
	Conditions  []SceneCondition `json:"conditions" gorm:"foreignKey:SceneID"`
	Actions     []SceneAction  `json:"actions" gorm:"foreignKey:SceneID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type SceneCondition struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SceneID     uint      `json:"sceneId" gorm:"not null;index"`
	Type        string    `json:"type" gorm:"size:20;not null"`
	DeviceID    *uint     `json:"deviceId" gorm:"index"`
	Operator    string    `json:"operator" gorm:"size:10"`
	Value       string    `json:"value" gorm:"size:100"`
	TimeExpr    string    `json:"timeExpr" gorm:"size:50"`
}

type SceneAction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SceneID     uint      `json:"sceneId" gorm:"not null;index"`
	DeviceID    uint      `json:"deviceId" gorm:"not null;index"`
	Action      string    `json:"action" gorm:"size:50;not null"`
	Value       string    `json:"value" gorm:"size:100"`
}

type Schedule struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FamilyID    uint           `json:"familyId" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"size:500"`
	DeviceID    uint           `json:"deviceId" gorm:"not null;index"`
	Action      string         `json:"action" gorm:"size:50;not null"`
	Value       string         `json:"value" gorm:"size:100"`
	CronExpr    string         `json:"cronExpr" gorm:"size:100;not null"`
	IsEnabled   bool           `json:"isEnabled" gorm:"default:true"`
	LastRun     *time.Time     `json:"lastRun"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ScheduleLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ScheduleID  uint      `json:"scheduleId" gorm:"not null;index"`
	DeviceID    uint      `json:"deviceId" gorm:"not null;index"`
	Action      string    `json:"action" gorm:"size:50"`
	Value       string    `json:"value" gorm:"size:100"`
	Success     bool      `json:"success"`
	Message     string    `json:"message" gorm:"size:500"`
	EnergyDelta float64   `json:"energyDelta" gorm:"default:0"`
	ExecutedAt  time.Time `json:"executedAt"`
}

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	Type      string    `json:"type" gorm:"size:20;not null"`
	Title     string    `json:"title" gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"size:1000"`
	IsRead    bool      `json:"isRead" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
}

type Invitation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FamilyID  uint      `json:"familyId" gorm:"not null;index"`
	Family    Family    `json:"family" gorm:"foreignKey:FamilyID"`
	InviterID uint      `json:"inviterId" gorm:"not null"`
	Inviter   User      `json:"inviter" gorm:"foreignKey:InviterID"`
	Email     string    `json:"email" gorm:"size:100;not null"`
	Role      string    `json:"role" gorm:"size:20;default:member"`
	Status    string    `json:"status" gorm:"size:20;default:pending"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func InitDB(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Database connected successfully")

	err = DB.AutoMigrate(
		&User{},
		&Family{},
		&FamilyMember{},
		&Device{},
		&EnergyData{},
		&EnergyAlert{},
		&DeviceGroup{},
		&Scene{},
		&SceneCondition{},
		&SceneAction{},
		&Schedule{},
		&ScheduleLog{},
		&Notification{},
		&Invitation{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
}
