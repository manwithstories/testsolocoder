package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderNo       string         `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	BuyerID       uint           `gorm:"index;not null" json:"buyer_id"`
	SellerID      uint           `gorm:"index;not null" json:"seller_id"`
	ProductID     uint           `gorm:"index;not null" json:"product_id"`
	ProductTitle  string         `gorm:"size:200;not null" json:"product_title"`
	ProductImage  string         `gorm:"size:255" json:"product_image"`
	OriginalPrice float64        `gorm:"type:decimal(10,2);not null" json:"original_price"`
	FinalPrice    float64        `gorm:"type:decimal(10,2);not null" json:"final_price"`
	Negotiated    bool           `gorm:"default:false" json:"negotiated"`
	Status        int            `gorm:"default:1;index" json:"status"`
	PaymentMethod string         `gorm:"size:20" json:"payment_method"`
	PaidAt        *time.Time     `json:"paid_at"`
	ShippedAt     *time.Time     `json:"shipped_at"`
	DeliveredAt   *time.Time     `json:"delivered_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	CancelledAt   *time.Time     `json:"cancelled_at"`
	RefundedAt    *time.Time     `json:"refunded_at"`
	TrackingNo    string         `gorm:"size:100" json:"tracking_no"`
	TrackingCompany string       `gorm:"size:50" json:"tracking_company"`
	ReceiverName  string         `gorm:"size:50" json:"receiver_name"`
	ReceiverPhone string         `gorm:"size:20" json:"receiver_phone"`
	ReceiverAddress string       `gorm:"size:255" json:"receiver_address"`
	WarrantyDays  int            `gorm:"default:30" json:"warranty_days"`
	WarrantyUntil *time.Time     `json:"warranty_until"`
	Commission    float64        `gorm:"type:decimal(10,2);default:0" json:"commission"`
	Remark        string         `gorm:"size:255" json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Buyer    User     `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller   User     `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Product  Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Reviews  []Review `gorm:"foreignKey:OrderID" json:"reviews,omitempty"`
	Negotiations []Negotiation `gorm:"foreignKey:OrderID" json:"negotiations,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

const (
	OrderStatusPending    = 1
	OrderStatusPaid       = 2
	OrderStatusShipped    = 3
	OrderStatusDelivered  = 4
	OrderStatusCompleted  = 5
	OrderStatusCancelled  = 6
	OrderStatusRefunding  = 7
	OrderStatusRefunded   = 8
	OrderStatusNegotiating = 9
)

var OrderStatusText = map[int]string{
	1: "待支付",
	2: "已支付",
	3: "已发货",
	4: "已送达",
	5: "已完成",
	6: "已取消",
	7: "退款中",
	8: "已退款",
	9: "议价中",
}

type Negotiation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderID       uint           `gorm:"index;not null" json:"order_id"`
	BuyerID       uint           `gorm:"index;not null" json:"buyer_id"`
	SellerID      uint           `gorm:"index;not null" json:"seller_id"`
	OfferedPrice  float64        `gorm:"type:decimal(10,2);not null" json:"offered_price"`
	CounterPrice  float64        `gorm:"type:decimal(10,2)" json:"counter_price"`
	Status        int            `gorm:"default:1" json:"status"`
	BuyerMessage  string         `gorm:"size:255" json:"buyer_message"`
	SellerMessage string         `gorm:"size:255" json:"seller_message"`
	ExpireAt      *time.Time     `json:"expire_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Negotiation) TableName() string {
	return "negotiations"
}

const (
	NegotiationStatusPending  = 1
	NegotiationStatusAccepted = 2
	NegotiationStatusRejected = 3
	NegotiationStatusExpired  = 4
)

type WalletLog struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	Type        string         `gorm:"size:20;not null" json:"type"`
	Amount      float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Balance     float64        `gorm:"type:decimal(10,2);not null" json:"balance"`
	OrderNo     string         `gorm:"size:50" json:"order_no"`
	Description string         `gorm:"size:255" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WalletLog) TableName() string {
	return "wallet_logs"
}

const (
	WalletTypeRecharge  = "recharge"
	WalletTypeWithdraw  = "withdraw"
	WalletTypePayment   = "payment"
	WalletTypeIncome    = "income"
	WalletTypeRefund    = "refund"
	WalletTypeCommission = "commission"
)
