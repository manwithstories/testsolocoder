package dto

import "time"

type SaveAnswerRequest struct {
	QuestionID   uint    `json:"question_id" binding:"required"`
	OptionID     *uint   `json:"option_id"`
	TextValue    string  `json:"text_value"`
	NumericValue float64 `json:"numeric_value"`
	MatrixValues string  `json:"matrix_values"`
	RankingOrder string  `json:"ranking_order"`
}

type SubmitResponseRequest struct {
	SessionID string              `json:"session_id" binding:"required"`
	Answers   []SaveAnswerRequest `json:"answers" binding:"required"`
}

type ResponseListQuery struct {
	Page       int       `form:"page,default=1"`
	PageSize   int       `form:"page_size,default=20"`
	Status     int       `form:"status"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
	Channel    string    `form:"channel"`
	SortBy     string    `form:"sort_by,default=created_at"`
	SortOrder  string    `form:"sort_order,default=desc"`
}

type AnswerResponse struct {
	QuestionID   uint    `json:"question_id"`
	OptionID     *uint   `json:"option_id"`
	TextValue    string  `json:"text_value"`
	NumericValue float64 `json:"numeric_value"`
	MatrixValues string  `json:"matrix_values"`
	RankingOrder string  `json:"ranking_order"`
}

type ResponseDetailResponse struct {
	ID            uint             `json:"id"`
	SurveyID      uint             `json:"survey_id"`
	UserID        *uint            `json:"user_id"`
	SessionID     string           `json:"session_id"`
	Status        int              `json:"status"`
	StartTime     *time.Time       `json:"start_time"`
	CompleteTime  *time.Time       `json:"complete_time"`
	Duration      int              `json:"duration"`
	Answers       []AnswerResponse `json:"answers"`
	CreatedAt     time.Time        `json:"created_at"`
}

type AccessSurveyRequest struct {
	Password  string `json:"password"`
	SessionID string `json:"session_id"`
}
