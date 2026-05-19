package models

import (
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type ChannelType string

const (
	ChannelTypeEmail   ChannelType = "email"
	ChannelTypeSMS     ChannelType = "sms"
	ChannelTypeWeChat  ChannelType = "wechat"
	ChannelTypeDingTalk ChannelType = "dingtalk"
	ChannelTypeWebhook ChannelType = "webhook"
)

type Channel struct {
	BaseModel
	Name        string                 `gorm:"type:varchar(100);not null;uniqueIndex" json:"name" validate:"required,max=100"`
	Type        ChannelType            `gorm:"type:varchar(50);not null;index" json:"type" validate:"required,oneof=email sms wechat dingtalk webhook"`
	Description string                 `gorm:"type:varchar(500)" json:"description" validate:"max=500"`
	Config      string                 `gorm:"type:text;not null" json:"config"`
	Enabled     bool                   `gorm:"default:true;index" json:"enabled"`
	Priority    int                    `gorm:"default:0;index" json:"priority"`
	RateLimit   *ChannelRateLimit      `gorm:"embedded;embeddedPrefix:rate_limit_" json:"rate_limit,omitempty"`
	Templates   []Template             `gorm:"foreignKey:ChannelID" json:"templates,omitempty"`
}

type ChannelRateLimit struct {
	Enabled           bool    `json:"enabled"`
	RequestsPerSecond float64 `json:"requests_per_second"`
	Burst             int     `json:"burst"`
}

type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	FromAddress  string `json:"from_address"`
	FromName     string `json:"from_name"`
	UseTLS       bool   `json:"use_tls"`
	UseSSL       bool   `json:"use_ssl"`
}

type SMSConfig struct {
	Provider    string `json:"provider"`
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	SignName    string `json:"sign_name"`
	TemplateCode string `json:"template_code"`
	Endpoint    string `json:"endpoint"`
}

type WeChatConfig struct {
	AppID      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
	Token      string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
}

type DingTalkConfig struct {
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret"`
	AccessToken string `json:"access_token"`
}

type WebhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Timeout int               `json:"timeout"`
}

type Template struct {
	BaseModel
	ChannelID   uint                `gorm:"not null;index" json:"channel_id" validate:"required"`
	Channel     Channel             `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	Name        string              `gorm:"type:varchar(100);not null" json:"name" validate:"required,max=100"`
	Description string              `gorm:"type:varchar(500)" json:"description" validate:"max=500"`
	Content     string              `gorm:"type:text;not null" json:"content" validate:"required"`
	Subject     string              `gorm:"type:varchar(500)" json:"subject" validate:"max=500"`
	Language    string              `gorm:"type:varchar(10);default:'zh-CN';index" json:"language" validate:"required"`
	Variables   []TemplateVariable  `gorm:"serializer:json" json:"variables,omitempty"`
	IsDefault   bool                `gorm:"default:false;index:idx_channel_lang_default" json:"is_default"`
}

type TemplateVariable struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=string int float bool date"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Pattern     string `json:"pattern"`
}

type Recipient struct {
	BaseModel
	Email       string            `gorm:"type:varchar(200);uniqueIndex" json:"email" validate:"omitempty,email"`
	Phone       string            `gorm:"type:varchar(50);uniqueIndex" json:"phone" validate:"omitempty,max=50"`
	WeChatOpenID string           `gorm:"type:varchar(100);uniqueIndex" json:"wechat_open_id"`
	DingTalkUserID string         `gorm:"type:varchar(100);uniqueIndex" json:"dingtalk_user_id"`
	Name        string            `gorm:"type:varchar(100)" json:"name" validate:"max=100"`
	Description string            `gorm:"type:varchar(500)" json:"description" validate:"max=500"`
	MetaData    map[string]string `gorm:"serializer:json" json:"meta_data,omitempty"`
	Enabled     bool              `gorm:"default:true;index" json:"enabled"`
	Tags        []Tag             `gorm:"many2many:recipient_tags;" json:"tags,omitempty"`
	Groups      []RecipientGroup  `gorm:"many2many:recipient_group_members;" json:"groups,omitempty"`
}

type Tag struct {
	BaseModel
	Name        string      `gorm:"type:varchar(50);not null;uniqueIndex" json:"name" validate:"required,max=50"`
	Description string      `gorm:"type:varchar(200)" json:"description" validate:"max=200"`
	Color       string      `gorm:"type:varchar(20)" json:"color" validate:"max=20"`
	Recipients  []Recipient `gorm:"many2many:recipient_tags;" json:"recipients,omitempty"`
}

type RecipientGroup struct {
	BaseModel
	Name        string      `gorm:"type:varchar(100);not null" json:"name" validate:"required,max=100"`
	Description string      `gorm:"type:varchar(500)" json:"description" validate:"max=500"`
	IsSystem    bool        `gorm:"default:false" json:"is_system"`
	Recipients  []Recipient `gorm:"many2many:recipient_group_members;" json:"recipients,omitempty"`
	Count       int         `gorm:"-" json:"count"`
}

type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusQueued    MessageStatus = "queued"
	MessageStatusSending   MessageStatus = "sending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusCancelled MessageStatus = "cancelled"
	MessageStatusRetrying  MessageStatus = "retrying"
)

type MessagePriority int

const (
	PriorityLow    MessagePriority = 1
	PriorityNormal MessagePriority = 5
	PriorityHigh   MessagePriority = 9
)

type Message struct {
	BaseModel
	MessageID    string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"message_id"`
	ChannelID    uint           `gorm:"not null;index" json:"channel_id"`
	Channel      Channel        `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	TemplateID   *uint          `gorm:"index" json:"template_id,omitempty"`
	Template     *Template      `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Recipient    string         `gorm:"type:varchar(500);not null" json:"recipient"`
	RecipientID  *uint          `gorm:"index" json:"recipient_id,omitempty"`
	Subject      string         `gorm:"type:varchar(500)" json:"subject"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	Variables    map[string]any `gorm:"serializer:json" json:"variables,omitempty"`
	Status       MessageStatus  `gorm:"type:varchar(50);not null;index" json:"status"`
	Priority     MessagePriority `gorm:"default:5;index" json:"priority"`
	ScheduledAt  *time.Time     `gorm:"index" json:"scheduled_at,omitempty"`
	SentAt       *time.Time     `json:"sent_at,omitempty"`
	DeliveredAt  *time.Time     `json:"delivered_at,omitempty"`
	DurationMs   int64          `json:"duration_ms"`
	RetryCount   int            `gorm:"default:0" json:"retry_count"`
	MaxRetries   int            `gorm:"default:5" json:"max_retries"`
	NextRetryAt  *time.Time     `gorm:"index" json:"next_retry_at,omitempty"`
	LastError    string         `gorm:"type:text" json:"last_error,omitempty"`
	ErrorStack   string         `gorm:"type:text" json:"error_stack,omitempty"`
	WebhookURL   string         `gorm:"type:varchar(500)" json:"webhook_url,omitempty"`
	RequestID    string         `gorm:"type:varchar(100);index" json:"request_id,omitempty"`
}

type MessageQueue struct {
	BaseModel
	MessageID   string        `gorm:"type:varchar(36);uniqueIndex;not null" json:"message_id"`
	Message     Message       `gorm:"foreignKey:MessageID;references:MessageID" json:"message,omitempty"`
	Priority    MessagePriority `gorm:"default:5;index" json:"priority"`
	ScheduledAt time.Time     `gorm:"not null;index" json:"scheduled_at"`
	PickedAt    *time.Time    `gorm:"index" json:"picked_at,omitempty"`
	WorkerID    string        `gorm:"type:varchar(100);index" json:"worker_id,omitempty"`
	RetryCount  int           `gorm:"default:0" json:"retry_count"`
}

type Webhook struct {
	BaseModel
	Name        string            `gorm:"type:varchar(100);not null" json:"name" validate:"required,max=100"`
	URL         string            `gorm:"type:varchar(500);not null" json:"url" validate:"required,url,max=500"`
	Method      string            `gorm:"type:varchar(10);default:'POST'" json:"method" validate:"oneof=GET POST PUT"`
	Headers     map[string]string `gorm:"serializer:json" json:"headers,omitempty"`
	Secret      string            `gorm:"type:varchar(200)" json:"secret" validate:"max=200"`
	Timeout     int               `gorm:"default:5000" json:"timeout"`
	MaxRetries  int               `gorm:"default:3" json:"max_retries"`
	Enabled     bool              `gorm:"default:true;index" json:"enabled"`
	Description string            `gorm:"type:varchar(500)" json:"description" validate:"max=500"`
	Events      []string          `gorm:"serializer:json" json:"events,omitempty"`
	ChannelIDs  []uint            `gorm:"serializer:json" json:"channel_ids,omitempty"`
}

type WebhookLog struct {
	BaseModel
	WebhookID   uint   `gorm:"not null;index" json:"webhook_id"`
	Webhook     Webhook `gorm:"foreignKey:WebhookID" json:"webhook,omitempty"`
	MessageID   string `gorm:"type:varchar(36);index" json:"message_id"`
	Event       string `gorm:"type:varchar(100);index" json:"event"`
	Status      string `gorm:"type:varchar(50);index" json:"status"`
	StatusCode  int    `json:"status_code"`
	Request     string `gorm:"type:text" json:"request,omitempty"`
	Response    string `gorm:"type:text" json:"response,omitempty"`
	Error       string `gorm:"type:text" json:"error,omitempty"`
	DurationMs  int64  `json:"duration_ms"`
	RetryCount  int    `gorm:"default:0" json:"retry_count"`
}

func GenerateMessageID() string {
	return uuid.New().String()
}

func (c *Channel) IsValid() bool {
	return c.Enabled
}

func (m *Message) CanRetry() bool {
	return m.Status == MessageStatusFailed && m.RetryCount < m.MaxRetries
}

func (m *Message) IsTerminal() bool {
	return m.Status == MessageStatusSent || m.Status == MessageStatusCancelled ||
		(m.Status == MessageStatusFailed && !m.CanRetry())
}
