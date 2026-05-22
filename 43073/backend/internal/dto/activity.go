package dto

type ActivityCreateRequest struct {
	Title       string  `json:"title" validate:"required,max=200"`
	Description string  `json:"description"`
	StartTime   string  `json:"startTime" validate:"required"`
	EndTime     string  `json:"endTime" validate:"required"`
	Location    string  `json:"location" validate:"required,max=500"`
	Capacity    int     `json:"capacity" validate:"required,min=1"`
	Poster      string  `json:"poster"`
}

type ActivityUpdateRequest struct {
	Title       string `json:"title" validate:"omitempty,max=200"`
	Description string `json:"description"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Location    string `json:"location" validate:"omitempty,max=500"`
	Capacity    int    `json:"capacity" validate:"omitempty,min=1"`
	Poster      string `json:"poster"`
}

type ActivityStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=draft published canceled"`
}

type ActivityListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}
