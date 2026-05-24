package model

import (
	"encoding/json"
	"time"
)

// 订单状态常量
const (
	OrderStatusInquiry   = "inquiry"   // 询价中
	OrderStatusQuoted    = "quoted"    // 已报价
	OrderStatusConfirmed = "confirmed" // 已确认
	OrderStatusPaid      = "paid"      // 已付款
	OrderStatusProducing = "producing" // 生产中
	OrderStatusShipped   = "shipped"   // 已发货
	OrderStatusCompleted = "completed" // 已完成
	OrderStatusCancelled = "cancelled" // 已取消
)

// ValidOrderStatus 校验订单状态是否合法
func ValidOrderStatus(status string) bool {
	switch status {
	case OrderStatusInquiry, OrderStatusQuoted, OrderStatusConfirmed,
		OrderStatusPaid, OrderStatusProducing, OrderStatusShipped,
		OrderStatusCompleted, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// 订单状态流转：inquiry -> quoted -> confirmed -> paid -> producing -> shipped -> completed
// 任意非最终状态可转为 cancelled
var orderStatusFlow = map[string]map[string]bool{
	OrderStatusInquiry:   {OrderStatusQuoted: true, OrderStatusCancelled: true},
	OrderStatusQuoted:    {OrderStatusConfirmed: true, OrderStatusCancelled: true},
	OrderStatusConfirmed: {OrderStatusPaid: true, OrderStatusCancelled: true},
	OrderStatusPaid:      {OrderStatusProducing: true},
	OrderStatusProducing: {OrderStatusShipped: true},
	OrderStatusShipped:   {OrderStatusCompleted: true},
	OrderStatusCompleted: {},
	OrderStatusCancelled: {},
}

// CanTransitTo 校验订单状态是否可流转到目标状态
func CanTransitTo(from, to string) bool {
	next, ok := orderStatusFlow[from]
	if !ok {
		return false
	}
	return next[to]
}

// CustomOption 订单项定制选项（JSON 存储）
type CustomOption struct {
	OptionType  string  `json:"option_type"`
	OptionValue string  `json:"option_value"`
	PriceAdjust float64 `json:"price_adjust"`
}

// Order 订单模型
type Order struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderNo        string    `gorm:"uniqueIndex;size:64;not null" json:"order_no"`
	OwnerID        uint      `gorm:"index;not null" json:"owner_id"`
	DesignerID     uint      `gorm:"index" json:"designer_id"`
	ManufacturerID uint      `gorm:"index;not null" json:"manufacturer_id"`
	TotalAmount    float64   `gorm:"type:decimal(12,2);not null;default:0" json:"total_amount"`
	Discount       float64   `gorm:"type:decimal(12,2);not null;default:0" json:"discount"`
	FinalAmount    float64   `gorm:"type:decimal(12,2);not null;default:0" json:"final_amount"`
	Status         string    `gorm:"size:32;not null;default:inquiry;index" json:"status"`
	Address        string    `gorm:"size:255" json:"address"`
	ContactName    string    `gorm:"size:64" json:"contact_name"`
	ContactPhone   string    `gorm:"size:20" json:"contact_phone"`
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项
type OrderItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"index;not null" json:"order_id"`
	ProductID     uint      `gorm:"index;not null" json:"product_id"`
	ProductName   string    `gorm:"size:128;not null" json:"product_name"`
	BasePrice     float64   `gorm:"type:decimal(12,2);not null" json:"base_price"`
	CustomOptions string    `gorm:"type:text" json:"custom_options"` // JSON 存储
	Quantity      int       `gorm:"not null" json:"quantity"`
	Subtotal      float64   `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}

// ParseCustomOptions 解析 JSON 定制选项
func (oi *OrderItem) ParseCustomOptions() ([]CustomOption, error) {
	if oi.CustomOptions == "" {
		return nil, nil
	}
	var opts []CustomOption
	if err := json.Unmarshal([]byte(oi.CustomOptions), &opts); err != nil {
		return nil, err
	}
	return opts, nil
}

// SetCustomOptions 将定制选项序列化为 JSON
func (oi *OrderItem) SetCustomOptions(opts []CustomOption) error {
	if len(opts) == 0 {
		oi.CustomOptions = ""
		return nil
	}
	data, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	oi.CustomOptions = string(data)
	return nil
}

// OrderHistory 订单状态变更历史
type OrderHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"index;not null" json:"order_id"`
	Status       string    `gorm:"size:32;not null" json:"status"`
	OperatorID   uint      `gorm:"index;not null" json:"operator_id"`
	OperatorRole string    `gorm:"size:32;not null" json:"operator_role"`
	Remark       string    `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (OrderHistory) TableName() string {
	return "order_histories"
}
