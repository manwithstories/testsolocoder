package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Role     string `json:"role" binding:"required,oneof=owner store keeper"`
	RealName string `json:"real_name"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string    `json:"token"`
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserProfile struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Role       string    `json:"role"`
	AvatarURL  string    `json:"avatar_url"`
	RealName   string    `json:"real_name"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	StoreInfo  *StoreInfoDTO  `json:"store_info,omitempty"`
	KeeperInfo *KeeperInfoDTO `json:"keeper_info,omitempty"`
}

type StoreInfoDTO struct {
	ID            uuid.UUID `json:"id"`
	StoreName     string    `json:"store_name"`
	Address       string    `json:"address"`
	LicenseNo     string    `json:"license_no"`
	BusinessHours string    `json:"business_hours"`
	Description   string    `json:"description"`
	Rating        float64   `json:"rating"`
	ReviewCount   int       `json:"review_count"`
}

type KeeperInfoDTO struct {
	ID            uuid.UUID `json:"id"`
	NickName      string    `json:"nick_name"`
	Experience    int       `json:"experience"`
	Specialty     string    `json:"specialty"`
	Rating        float64   `json:"rating"`
	ReviewCount   int       `json:"review_count"`
	Certification string    `json:"certification"`
	StoreID       uuid.UUID `json:"store_id,omitempty"`
}

type UpdateUserRequest struct {
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
	RealName  string `json:"real_name"`
}

type UpdateStoreRequest struct {
	StoreName     string `json:"store_name"`
	Address       string `json:"address"`
	LicenseNo     string `json:"license_no"`
	BusinessHours string `json:"business_hours"`
	Description   string `json:"description"`
}

type UpdateKeeperRequest struct {
	NickName      string `json:"nick_name"`
	Experience    int    `json:"experience"`
	Specialty     string `json:"specialty"`
	Certification string `json:"certification"`
	StoreID       string `json:"store_id"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

type PetRequest struct {
	Name        string    `json:"name" binding:"required,max=50"`
	Species     string    `json:"species" binding:"required,oneof=dog cat rabbit bird other"`
	Breed       string    `json:"breed"`
	Gender      string    `json:"gender" binding:"required,oneof=male female"`
	BirthDate   *time.Time `json:"birth_date"`
	Weight      float64   `json:"weight"`
	Color       string    `json:"color"`
	AvatarURL   string    `json:"avatar_url"`
	Allergies   string    `json:"allergies"`
	DietHabit   string    `json:"diet_habit"`
	Temperament string    `json:"temperament"`
}

type VaccineRecordRequest struct {
	PetID        string    `json:"pet_id" binding:"required"`
	VaccineName  string    `json:"vaccine_name" binding:"required"`
	VaccinatedAt time.Time `json:"vaccinated_at" binding:"required"`
	ExpireAt     time.Time `json:"expire_at" binding:"required"`
	Hospital     string    `json:"hospital"`
	ProofURL     string    `json:"proof_url"`
}

type DewormRecordRequest struct {
	PetID       string    `json:"pet_id" binding:"required"`
	DewormType  string    `json:"deworm_type" binding:"required"`
	DewormedAt  time.Time `json:"dewormed_at" binding:"required"`
	ExpireAt    time.Time `json:"expire_at" binding:"required"`
	Medicine    string    `json:"medicine"`
}

type BoardingPackageRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Type        string  `json:"type" binding:"required,oneof=daycare boarding"`
	Description string  `json:"description"`
	PricePerDay float64 `json:"price_per_day" binding:"required,gt=0"`
	Capacity    int     `json:"capacity" binding:"required,gt=0"`
	Features    string  `json:"features"`
	SortOrder   int     `json:"sort_order"`
}

type ReservationRequest struct {
	PetID        string    `json:"pet_id" binding:"required"`
	StoreID      string    `json:"store_id" binding:"required"`
	PackageID    string    `json:"package_id" binding:"required"`
	CheckInDate  time.Time `json:"check_in_date" binding:"required"`
	CheckOutDate time.Time `json:"check_out_date" binding:"required"`
	KeeperID     string    `json:"keeper_id"`
	Remark       string    `json:"remark"`
}

type ReservationConfirmRequest struct {
	Status string `json:"status" binding:"required,oneof=confirmed cancelled"`
	Reason string `json:"reason"`
}

type CheckInOutRequest struct {
	Remark string `json:"remark"`
}

type DailyRecordRequest struct {
	ReservationID string   `json:"reservation_id" binding:"required"`
	PetID         string   `json:"pet_id" binding:"required"`
	RecordDate    time.Time `json:"record_date" binding:"required"`
	FeedStatus    string   `json:"feed_status"`
	Activity      string   `json:"activity"`
	HealthStatus  string   `json:"health_status"`
	Mood          string   `json:"mood"`
	Photos        string   `json:"photos"`
	Remark        string   `json:"remark"`
}

type ReviewRequest struct {
	ReservationID string `json:"reservation_id" binding:"required"`
	StoreRating   int    `json:"store_rating" binding:"required,gte=1,lte=5"`
	KeeperRating  int    `json:"keeper_rating" binding:"gte=1,lte=5"`
	Content       string `json:"content" binding:"max=1000"`
	Images        string `json:"images"`
}

type ReviewReplyRequest struct {
	Reply string `json:"reply" binding:"required,max=1000"`
}

type OrderRequest struct {
	ReservationID string  `json:"reservation_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	AmountHash    string  `json:"amount_hash" binding:"required"`
	PayMethod     string  `json:"pay_method" binding:"required,oneof=wechat alipay bank_transfer cash"`
	Remark        string  `json:"remark"`
}

type OrderSettlementRequest struct {
	ReservationID string  `json:"reservation_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	AmountHash    string  `json:"amount_hash" binding:"required"`
	PayMethod     string  `json:"pay_method" binding:"required,oneof=wechat alipay bank_transfer cash"`
	Remark        string  `json:"remark"`
}

type StatisticsQuery struct {
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
	StoreID   string     `form:"store_id"`
	GroupBy   string     `form:"group_by"`
}

type Pagination struct {
	Page     int `form:"page" binding:"gte=1"`
	PageSize int `form:"page_size" binding:"gte=1,lte=100"`
}

type PagedResult struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Items    interface{} `json:"items"`
}

func NewPagination() Pagination {
	return Pagination{
		Page:     1,
		PageSize: 10,
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func NewPagedResult(total int64, page, pageSize int, items interface{}) *PagedResult {
	return &PagedResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}
}
