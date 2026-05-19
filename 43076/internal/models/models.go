package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TicketPriorityUrgent = "urgent"
	TicketPriorityHigh   = "high"
	TicketPriorityMedium = "medium"
	TicketPriorityLow    = "low"

	TicketStatusOpen        = "open"
	TicketStatusAssigned    = "assigned"
	TicketStatusInProgress  = "in_progress"
	TicketStatusPending     = "pending"
	TicketStatusEscalated   = "escalated"
	TicketStatusResolved    = "resolved"
	TicketStatusClosed      = "closed"

	TicketTypeBug        = "bug"
	TicketTypeFeature    = "feature"
	TicketTypeSupport    = "support"
	TicketTypeQuestion   = "question"
	TicketTypeComplaint  = "complaint"

	AssignmentModeSkillGroup = "skill_group"
	AssignmentModeRoundRobin = "round_robin"
	AssignmentModeLoadBalance = "load_balance"

	CommentTypeInternal = "internal"
	CommentTypePublic   = "public"

	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleAgent    = "agent"
	RoleViewer   = "viewer"
)

var ValidPriorities = map[string]bool{
	TicketPriorityUrgent: true,
	TicketPriorityHigh:   true,
	TicketPriorityMedium: true,
	TicketPriorityLow:    true,
}

var ValidStatuses = map[string]bool{
	TicketStatusOpen:       true,
	TicketStatusAssigned:   true,
	TicketStatusInProgress: true,
	TicketStatusPending:    true,
	TicketStatusEscalated:  true,
	TicketStatusResolved:   true,
	TicketStatusClosed:     true,
}

var ValidTypes = map[string]bool{
	TicketTypeBug:       true,
	TicketTypeFeature:   true,
	TicketTypeSupport:   true,
	TicketTypeQuestion:  true,
	TicketTypeComplaint: true,
}

var ValidAssignmentModes = map[string]bool{
	AssignmentModeSkillGroup:  true,
	AssignmentModeRoundRobin:  true,
	AssignmentModeLoadBalance: true,
}

var ValidRoles = map[string]bool{
	RoleAdmin:   true,
	RoleManager: true,
	RoleAgent:   true,
	RoleViewer:  true,
}

var StatusTransitions = map[string]map[string]bool{
	TicketStatusOpen: {
		TicketStatusAssigned:   true,
		TicketStatusClosed:     true,
	},
	TicketStatusAssigned: {
		TicketStatusInProgress: true,
		TicketStatusPending:    true,
		TicketStatusEscalated:  true,
		TicketStatusClosed:     true,
	},
	TicketStatusInProgress: {
		TicketStatusPending:    true,
		TicketStatusEscalated:  true,
		TicketStatusResolved:   true,
		TicketStatusClosed:     true,
	},
	TicketStatusPending: {
		TicketStatusInProgress: true,
		TicketStatusEscalated:  true,
		TicketStatusClosed:     true,
	},
	TicketStatusEscalated: {
		TicketStatusInProgress: true,
		TicketStatusResolved:   true,
		TicketStatusClosed:     true,
	},
	TicketStatusResolved: {
		TicketStatusClosed: true,
	},
	TicketStatusClosed: {},
}

func IsValidStatusTransition(from, to string) bool {
	if transitions, ok := StatusTransitions[from]; ok {
		return transitions[to]
	}
	return false
}

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Department struct {
	BaseModel
	Name        string `gorm:"size:100;not null;unique" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	Users       []User `json:"users,omitempty"`
}

type SkillGroup struct {
	BaseModel
	Name        string `gorm:"size:100;not null;unique" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	Users       []User `gorm:"many2many:user_skill_groups;" json:"users,omitempty"`
	Tickets     []Ticket `json:"tickets,omitempty"`
}

type User struct {
	BaseModel
	Username     string      `gorm:"size:50;not null;unique" json:"username"`
	Email        string      `gorm:"size:100;not null;unique" json:"email"`
	PasswordHash string      `gorm:"size:255;not null" json:"-"`
	RealName     string      `gorm:"size:50" json:"real_name"`
	Phone        string      `gorm:"size:20" json:"phone"`
	Role         string      `gorm:"size:20;not null;default:'agent'" json:"role"`
	DepartmentID *uint       `json:"department_id"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	SkillGroups  []SkillGroup `gorm:"many2many:user_skill_groups;" json:"skill_groups,omitempty"`
	IsOnDuty     bool        `gorm:"default:false" json:"is_on_duty"`
	Tickets      []Ticket    `gorm:"foreignKey:AssigneeID" json:"tickets,omitempty"`
	CreatedTickets []Ticket  `gorm:"foreignKey:CreatorID" json:"created_tickets,omitempty"`
}

type Customer struct {
	BaseModel
	Name     string    `gorm:"size:100;not null" json:"name"`
	Email    string    `gorm:"size:100;not null;unique" json:"email"`
	Phone    string    `gorm:"size:20" json:"phone"`
	Company  string    `gorm:"size:100" json:"company"`
	Address  string    `gorm:"size:500" json:"address"`
	Tickets  []Ticket  `json:"tickets,omitempty"`
}

type TicketCounter struct {
	BaseModel
	Date       string `gorm:"size:10;not null;uniqueIndex" json:"date"`
	Count      int    `gorm:"not null;default:0" json:"count"`
	Version    int    `gorm:"not null;default:0" json:"version"`
}

type Ticket struct {
	BaseModel
	TicketNo      string     `gorm:"size:30;not null;uniqueIndex" json:"ticket_no"`
	Title         string     `gorm:"size:200;not null" json:"title"`
	Description   string     `gorm:"type:text" json:"description"`
	Priority      string     `gorm:"size:20;not null;default:'medium'" json:"priority"`
	Status        string     `gorm:"size:20;not null;default:'open'" json:"status"`
	Type          string     `gorm:"size:20;not null;default:'support'" json:"type"`
	Tags          string     `gorm:"size:500" json:"tags"`
	CustomerID    uint       `gorm:"not null" json:"customer_id"`
	Customer      *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	CreatorID     uint       `gorm:"not null" json:"creator_id"`
	Creator       *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	AssigneeID    *uint      `json:"assignee_id"`
	Assignee      *User      `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	SkillGroupID  *uint      `json:"skill_group_id"`
	SkillGroup    *SkillGroup `gorm:"foreignKey:SkillGroupID" json:"skill_group,omitempty"`
	ResponseDeadline *time.Time `json:"response_deadline"`
	ResolveDeadline  *time.Time `json:"resolve_deadline"`
	FirstResponseAt  *time.Time `json:"first_response_at"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	ClosedAt         *time.Time `json:"closed_at"`
	IsEscalated      bool       `gorm:"default:false" json:"is_escalated"`
	EscalatedAt      *time.Time `json:"escalated_at"`
	Logs             []TicketLog `json:"logs,omitempty"`
	Comments         []Comment `json:"comments,omitempty"`
	Attachments      []Attachment `json:"attachments,omitempty"`
}

type TicketLog struct {
	BaseModel
	TicketID    uint      `gorm:"not null;index" json:"ticket_id"`
	Ticket      *Ticket   `gorm:"foreignKey:TicketID" json:"-"`
	OperatorID  uint      `gorm:"not null" json:"operator_id"`
	Operator    *User     `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
	Action      string    `gorm:"size:50;not null" json:"action"`
	OldValue    string    `gorm:"type:text" json:"old_value"`
	NewValue    string    `gorm:"type:text" json:"new_value"`
	Remark      string    `gorm:"size:500" json:"remark"`
}

type Comment struct {
	BaseModel
	TicketID    uint      `gorm:"not null;index" json:"ticket_id"`
	Ticket      *Ticket   `gorm:"foreignKey:TicketID" json:"-"`
	AuthorID    uint      `gorm:"not null" json:"author_id"`
	Author      *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Type        string    `gorm:"size:20;not null;default:'internal'" json:"type"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	BaseModel
	TicketID    *uint     `gorm:"index" json:"ticket_id"`
	CommentID   *uint     `gorm:"index" json:"comment_id"`
	FileName    string    `gorm:"size:255;not null" json:"file_name"`
	FilePath    string    `gorm:"size:500;not null" json:"file_path"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `gorm:"size:100" json:"content_type"`
	UploaderID  uint      `gorm:"not null" json:"uploader_id"`
	Uploader    *User     `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
}

type SLARecord struct {
	BaseModel
	TicketID        uint      `gorm:"not null;uniqueIndex" json:"ticket_id"`
	Ticket          *Ticket   `gorm:"foreignKey:TicketID" json:"-"`
	ResponseTime    int       `json:"response_time_minutes"`
	ResolutionTime  int       `json:"resolution_time_minutes"`
	ResponseBreached bool     `gorm:"default:false" json:"response_breached"`
	ResolveBreached  bool     `gorm:"default:false" json:"resolve_breached"`
	EscalationCount int       `gorm:"default:0" json:"escalation_count"`
}

type AssignmentRule struct {
	BaseModel
	Name          string `gorm:"size:100;not null;unique" json:"name"`
	Description   string `gorm:"size:500" json:"description"`
	Mode          string `gorm:"size:20;not null;default:'skill_group'" json:"mode"`
	SkillGroupID  *uint  `json:"skill_group_id"`
	IsDefault     bool   `gorm:"default:false" json:"is_default"`
	Enabled       bool   `gorm:"default:true" json:"enabled"`
}
