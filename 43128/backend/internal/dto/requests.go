package dto

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	RealName string `json:"real_name" binding:"required"`
	IdCard   string `json:"id_card" binding:"required,len=18"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	RealName string `json:"real_name"`
	Verified bool   `json:"verified"`
}

type VerifyReq struct {
	RealName string `json:"real_name" binding:"required"`
	IdCard   string `json:"id_card" binding:"required,len=18"`
}

type EventCreateReq struct {
	Name                  string         `json:"name" binding:"required"`
	Description           string         `json:"description"`
	Location              string         `json:"location"`
	StartDate             string         `json:"start_date" binding:"required"`
	EndDate               string         `json:"end_date" binding:"required"`
	RegistrationDeadline  string         `json:"registration_deadline" binding:"required"`
	Organizer             string         `json:"organizer"`
	CoverImage            string         `json:"cover_image"`
	Items                 []EventItemReq `json:"items" binding:"required,min=1"`
}

type EventItemReq struct {
	Name          string  `json:"name" binding:"required"`
	Category      string  `json:"category"`
	Gender        string  `json:"gender"`
	MinAge        int     `json:"min_age"`
	MaxAge        int     `json:"max_age"`
	Quota         int     `json:"quota" binding:"required,gt=0"`
	WaitlistQuota int     `json:"waitlist_quota"`
	Fee           float64 `json:"fee"`
	Requirements  string  `json:"requirements"`
}

type EventUpdateReq struct {
	EventCreateReq
	ID   uint   `json:"id" binding:"required"`
}

type RegistrationReq struct {
	EventItemID uint           `json:"event_item_id" binding:"required"`
	RegType     string         `json:"reg_type" binding:"required,oneof=individual team"`
	TeamName    string         `json:"team_name"`
	TeamMembers []TeamMember   `json:"team_members"`
}

type TeamMember struct {
	Name     string `json:"name" binding:"required"`
	IdCard   string `json:"id_card" binding:"required"`
	Phone    string `json:"phone"`
}

type ScoreEntryReq struct {
	EventID      uint          `json:"event_id" binding:"required"`
	EventItemID  uint          `json:"event_item_id" binding:"required"`
	Scores       []ScoreItem   `json:"scores" binding:"required"`
}

type ScoreItem struct {
	UserID   uint    `json:"user_id" binding:"required"`
	Score    float64 `json:"score" binding:"required"`
	TimeUsed string  `json:"time_used"`
	Remarks  string  `json:"remarks"`
}

type MessagePushReq struct {
	UserIDs []uint `json:"user_ids"`
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type Pagination struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=200"`
}

type PagedData struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     interface{} `json:"list"`
}
