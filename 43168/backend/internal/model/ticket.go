package model

import (
	"encoding/json"
	"time"
)

// 售后工单类型常量
const (
	TicketTypeProductQuality = "product_quality" // 产品质量
	TicketTypeDelivery       = "delivery"        // 配送问题
	TicketTypeService        = "service"         // 服务问题
	TicketTypeOther          = "other"           // 其他
)

// ValidTicketType 校验工单类型是否合法
func ValidTicketType(t string) bool {
	switch t {
	case TicketTypeProductQuality, TicketTypeDelivery,
		TicketTypeService, TicketTypeOther:
		return true
	default:
		return false
	}
}

// 售后工单状态常量
const (
	TicketStatusOpen       = "open"       // 待处理
	TicketStatusProcessing = "processing" // 处理中
	TicketStatusResolved   = "resolved"   // 已解决
	TicketStatusClosed     = "closed"     // 已关闭
)

// ValidTicketStatus 校验工单状态是否合法
func ValidTicketStatus(s string) bool {
	switch s {
	case TicketStatusOpen, TicketStatusProcessing,
		TicketStatusResolved, TicketStatusClosed:
		return true
	default:
		return false
	}
}

// Ticket 售后工单模型
type Ticket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	OwnerID   uint      `gorm:"index;not null" json:"owner_id"`
	Type      string    `gorm:"size:32;not null;index" json:"type"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Status    string    `gorm:"size:32;not null;default:open;index" json:"status"`
	Images    string    `gorm:"type:text" json:"images"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Ticket) TableName() string {
	return "tickets"
}

// ParseImages 解析图片 JSON 数组
func (t *Ticket) ParseImages() ([]string, error) {
	if t.Images == "" {
		return nil, nil
	}
	var list []string
	if err := json.Unmarshal([]byte(t.Images), &list); err != nil {
		return nil, err
	}
	return list, nil
}

// SetImages 将图片列表序列化为 JSON
func (t *Ticket) SetImages(images []string) error {
	if len(images) == 0 {
		t.Images = ""
		return nil
	}
	data, err := json.Marshal(images)
	if err != nil {
		return err
	}
	t.Images = string(data)
	return nil
}
