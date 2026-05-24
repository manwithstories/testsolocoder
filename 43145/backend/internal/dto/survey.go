package dto

import "time"

type CreateSurveyRequest struct {
	Title         string     `json:"title" binding:"required,min=1,max=200"`
	Description   string     `json:"description"`
	CoverImage    string     `json:"cover_image"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	Anonymous     bool       `json:"anonymous"`
	Password      string     `json:"password"`
	MaxResponses  int        `json:"max_responses"`
	MaxPerUser    int        `json:"max_per_user"`
	RequiresLogin bool       `json:"requires_login"`
	AllowResume   bool       `json:"allow_resume"`
	Category      string     `json:"category"`
	Tags          string     `json:"tags"`
}

type UpdateSurveyRequest struct {
	Title         string     `json:"title" binding:"omitempty,min=1,max=200"`
	Description   string     `json:"description"`
	CoverImage    string     `json:"cover_image"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	Anonymous     *bool      `json:"anonymous"`
	Password      string     `json:"password"`
	MaxResponses  int        `json:"max_responses"`
	MaxPerUser    int        `json:"max_per_user"`
	RequiresLogin *bool      `json:"requires_login"`
	AllowResume   *bool      `json:"allow_resume"`
	Category      string     `json:"category"`
	Tags          string     `json:"tags"`
}

type SurveyListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   int    `form:"status"`
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
	SortBy   string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}

type SurveyResponse struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CoverImage    string     `json:"cover_image"`
	UserID        uint       `json:"user_id"`
	Status        int        `json:"status"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	Anonymous     bool       `json:"anonymous"`
	HasPassword   bool       `json:"has_password"`
	MaxResponses  int        `json:"max_responses"`
	MaxPerUser    int        `json:"max_per_user"`
	RequiresLogin bool       `json:"requires_login"`
	AllowResume   bool       `json:"allow_resume"`
	ResponseCount int        `json:"response_count"`
	Category      string     `json:"category"`
	Tags          string     `json:"tags"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type PublishSurveyRequest struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

type CopySurveyRequest struct {
	Title string `json:"title"`
}
