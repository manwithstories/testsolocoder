package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username        string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email           string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password        string    `gorm:"size:255;not null" json:"-"`
	Nickname        string    `gorm:"size:50" json:"nickname"`
	Avatar          string    `gorm:"size:255" json:"avatar"`
	Phone           string    `gorm:"size:20" json:"phone"`
	Region          string    `gorm:"size:100" json:"region"`
	ClimateZone     string    `gorm:"size:50" json:"climate_zone"`
	UserType        string    `gorm:"size:20;default:hobbyist" json:"user_type"`
	CreditScore     int       `gorm:"default:100" json:"credit_score"`
	IsVerified      bool      `gorm:"default:false" json:"is_verified"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Plots          []Plot          `gorm:"foreignKey:UserID" json:"plots,omitempty"`
	Posts          []Post          `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments       []Comment       `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	SeedExchanges  []SeedExchange  `gorm:"foreignKey:OwnerID" json:"seed_exchanges,omitempty"`
	Orders         []Order         `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	FollowerUsers  []Follow        `gorm:"foreignKey:FollowerID" json:"-"`
	FollowingUsers []Follow        `gorm:"foreignKey:FollowingID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Plot struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Description   string    `gorm:"size:500" json:"description"`
	SoilType      string    `gorm:"size:50" json:"soil_type"`
	Sunlight      string    `gorm:"size:50" json:"sunlight"`
	Area          float64   `gorm:"type:decimal(10,2)" json:"area"`
	Location      string    `gorm:"size:100" json:"location"`
	GridConfig    string    `gorm:"type:text" json:"grid_config"`
	IrrigationDevice string `gorm:"size:100" json:"irrigation_device"`
	SensorData    string    `gorm:"type:text" json:"sensor_data"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	User          User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PlantingRecords []PlantingRecord `gorm:"foreignKey:PlotID" json:"planting_records,omitempty"`
}

func (p *Plot) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Plant struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	LatinName       string    `gorm:"size:100" json:"latin_name"`
	Category        string    `gorm:"size:50;index" json:"category"`
	Description     string    `gorm:"type:text" json:"description"`
	GrowthCycle     int       `json:"growth_cycle"`
	WaterFrequency  string    `gorm:"size:50" json:"water_frequency"`
	FertilizerNeed  string    `gorm:"size:100" json:"fertilizer_need"`
	SunlightNeed    string    `gorm:"size:50" json:"sunlight_need"`
	SoilPH          string    `gorm:"size:20" json:"soil_ph"`
	PlantingSeason  string    `gorm:"size:100" json:"planting_season"`
	HarvestSeason   string    `gorm:"size:100" json:"harvest_season"`
	DiseaseInfo     string    `gorm:"type:text" json:"disease_info"`
	PestInfo        string    `gorm:"type:text" json:"pest_info"`
	ImageURL        string    `gorm:"size:255" json:"image_url"`
	SowingDepth     string    `gorm:"size:50" json:"sowing_depth"`
	Spacing         string    `gorm:"size:50" json:"spacing"`
	Difficulty      string    `gorm:"size:20" json:"difficulty"`
	ClimateZone     string    `gorm:"size:100" json:"climate_zone"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	PlantingRecords []PlantingRecord `gorm:"foreignKey:PlantID" json:"planting_records,omitempty"`
}

func (pl *Plant) BeforeCreate(tx *gorm.DB) error {
	if pl.ID == uuid.Nil {
		pl.ID = uuid.New()
	}
	return nil
}

type PlantingRecord struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PlotID        uuid.UUID `gorm:"type:uuid;not null;index" json:"plot_id"`
	PlantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"plant_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Quantity      int       `json:"quantity"`
	PlantingDate  time.Time `json:"planting_date"`
	ExpectedHarvestDate *time.Time `json:"expected_harvest_date"`
	ActualHarvestDate   *time.Time `json:"actual_harvest_date"`
	Status        string    `gorm:"size:20;default:planted" json:"status"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Plot          Plot      `gorm:"foreignKey:PlotID" json:"plot,omitempty"`
	Plant         Plant     `gorm:"foreignKey:PlantID" json:"plant,omitempty"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GrowthLogs    []GrowthLog `gorm:"foreignKey:PlantingRecordID" json:"growth_logs,omitempty"`
}

func (pr *PlantingRecord) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}

type GrowthLog struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PlantingRecordID uuid.UUID `gorm:"type:uuid;not null;index" json:"planting_record_id"`
	Title            string    `gorm:"size:100;not null" json:"title"`
	Description      string    `gorm:"type:text" json:"description"`
	ImageURL         string    `gorm:"size:255" json:"image_url"`
	LogType          string    `gorm:"size:50" json:"log_type"`
	LogDate          time.Time `json:"log_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	PlantingRecord   PlantingRecord `gorm:"foreignKey:PlantingRecordID" json:"planting_record,omitempty"`
}

func (gl *GrowthLog) BeforeCreate(tx *gorm.DB) error {
	if gl.ID == uuid.Nil {
		gl.ID = uuid.New()
	}
	return nil
}

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ImageURLs string    `gorm:"type:text" json:"image_urls"`
	Category  string    `gorm:"size:50;index" json:"category"`
	Tags      string    `gorm:"size:255" json:"tags"`
	ViewCount int       `gorm:"default:0" json:"view_count"`
	LikeCount int       `gorm:"default:0" json:"like_count"`
	IsPinned  bool      `gorm:"default:false" json:"is_pinned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments  []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Likes     []Like    `gorm:"foreignKey:PostID" json:"-"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index" json:"parent_id"`
	LikeCount int       `gorm:"default:0" json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PostID    *uuid.UUID `gorm:"type:uuid;index" json:"post_id"`
	CommentID *uuid.UUID `gorm:"type:uuid;index" json:"comment_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (l *Like) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type Follow struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`

	Follower    User      `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following   User      `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}

func (f *Follow) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type DiseaseDiagnosis struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	PlantName   string    `gorm:"size:100" json:"plant_name"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	Description string    `gorm:"type:text" json:"description"`
	Symptoms    string    `gorm:"type:text" json:"symptoms"`
	Diagnosis   string    `gorm:"type:text" json:"diagnosis"`
	Severity    string    `gorm:"size:20" json:"severity"`
	Treatment   string    `gorm:"type:text" json:"treatment"`
	Confidence  float64   `gorm:"type:decimal(5,2)" json:"confidence"`
	Status      string    `gorm:"size:20;default:pending" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (d *DiseaseDiagnosis) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type SeedExchange struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID        uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_id"`
	Title          string    `gorm:"size:200;not null" json:"title"`
	SeedName       string    `gorm:"size:100;not null" json:"seed_name"`
	Description    string    `gorm:"type:text" json:"description"`
	ImageURLs      string    `gorm:"type:text" json:"image_urls"`
	Quantity       int       `json:"quantity"`
	ExchangeType   string    `gorm:"size:20;default:exchange" json:"exchange_type"`
	WantSeeds      string    `gorm:"type:text" json:"want_seeds"`
	Location       string    `gorm:"size:100" json:"location"`
	Status         string    `gorm:"size:20;default:available" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Owner          User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	ExchangeOffers []ExchangeOffer `gorm:"foreignKey:SeedExchangeID" json:"exchange_offers,omitempty"`
}

func (se *SeedExchange) BeforeCreate(tx *gorm.DB) error {
	if se.ID == uuid.Nil {
		se.ID = uuid.New()
	}
	return nil
}

type ExchangeOffer struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SeedExchangeID uuid.UUID `gorm:"type:uuid;not null;index" json:"seed_exchange_id"`
	OffererID      uuid.UUID `gorm:"type:uuid;not null;index" json:"offerer_id"`
	OfferSeeds     string    `gorm:"type:text" json:"offer_seeds"`
	Message        string    `gorm:"type:text" json:"message"`
	Status         string    `gorm:"size:20;default:pending" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	SeedExchange   SeedExchange `gorm:"foreignKey:SeedExchangeID" json:"seed_exchange,omitempty"`
	Offerer        User      `gorm:"foreignKey:OffererID" json:"offerer,omitempty"`
}

func (eo *ExchangeOffer) BeforeCreate(tx *gorm.DB) error {
	if eo.ID == uuid.Nil {
		eo.ID = uuid.New()
	}
	return nil
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SellerID    uuid.UUID `gorm:"type:uuid;not null;index" json:"seller_id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:50;index" json:"category"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	ImageURLs   string    `gorm:"type:text" json:"image_urls"`
	Specifications string  `gorm:"type:text" json:"specifications"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Seller      User      `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity  int       `gorm:"default:1" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Order struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	OrderNo      string    `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	TotalAmount  float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status       string    `gorm:"size:20;default:pending" json:"status"`
	PaymentMethod string   `gorm:"size:50" json:"payment_method"`
	PaymentStatus string   `gorm:"size:20;default:unpaid" json:"payment_status"`
	ShippingAddress string `gorm:"type:text" json:"shipping_address"`
	ShippingPhone string   `gorm:"size:20" json:"shipping_phone"`
	ShippingName  string   `gorm:"size:50" json:"shipping_name"`
	TrackingNumber string  `gorm:"size:100" json:"tracking_number"`
	ShippingStatus string  `gorm:"size:20" json:"shipping_status"`
	Remark       string    `gorm:"size:255" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductName string   `gorm:"size:200" json:"product_name"`
	Price      float64   `gorm:"type:decimal(10,2)" json:"price"`
	Quantity   int       `json:"quantity"`
	Subtotal   float64   `gorm:"type:decimal(10,2)" json:"subtotal"`
	CreatedAt  time.Time `json:"created_at"`

	Order      Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}

type CalendarEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	EventType   string    `gorm:"size:50" json:"event_type"`
	EventDate   time.Time `json:"event_date"`
	Description string    `gorm:"type:text" json:"description"`
	PlantID     *uuid.UUID `gorm:"type:uuid;index" json:"plant_id"`
	PlotID      *uuid.UUID `gorm:"type:uuid;index" json:"plot_id"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ce *CalendarEvent) BeforeCreate(tx *gorm.DB) error {
	if ce.ID == uuid.Nil {
		ce.ID = uuid.New()
	}
	return nil
}

type OperationLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Action      string    `gorm:"size:100;not null" json:"action"`
	Resource    string    `gorm:"size:100" json:"resource"`
	ResourceID  string    `gorm:"size:100" json:"resource_id"`
	Method      string    `gorm:"size:10" json:"method"`
	Path        string    `gorm:"size:255" json:"path"`
	IPAddress   string    `gorm:"size:50" json:"ip_address"`
	UserAgent   string    `gorm:"size:255" json:"user_agent"`
	RequestBody string    `gorm:"type:text" json:"request_body"`
	ResponseStatus int    `json:"response_status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ol *OperationLog) BeforeCreate(tx *gorm.DB) error {
	if ol.ID == uuid.Nil {
		ol.ID = uuid.New()
	}
	return nil
}
