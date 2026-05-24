package dto

import "time"

type StatisticsQuery struct {
	SurveyID  uint       `form:"survey_id" binding:"required"`
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
	Channel   string     `form:"channel"`
	QuestionID uint      `form:"question_id"`
	Dimension string     `form:"dimension" binding:"omitempty,oneof=time channel user_type"`
}

type StatisticsResponse struct {
	SurveyID        uint                `json:"survey_id"`
	TotalResponses  int                 `json:"total_responses"`
	CompletedCount  int                 `json:"completed_count"`
	InProgressCount int                 `json:"in_progress_count"`
	AbandonedCount  int                 `json:"abandoned_count"`
	CompletionRate  float64             `json:"completion_rate"`
	AvgDuration     float64             `json:"avg_duration"`
	QuestionStats   []QuestionStat      `json:"question_stats"`
	TimeDistribution []TimeDataPoint    `json:"time_distribution"`
	ChannelStats    []ChannelStat       `json:"channel_stats"`
}

type QuestionStat struct {
	QuestionID    uint           `json:"question_id"`
	QuestionTitle string         `json:"question_title"`
	QuestionType  string         `json:"question_type"`
	ResponseCount int            `json:"response_count"`
	OptionStats   []OptionStat   `json:"option_stats,omitempty"`
	TextStats     *TextStat      `json:"text_stats,omitempty"`
	RatingStats   *RatingStat    `json:"rating_stats,omitempty"`
	WordCloud     []WordCloudItem `json:"word_cloud,omitempty"`
}

type OptionStat struct {
	OptionID uint    `json:"option_id"`
	Text     string  `json:"text"`
	Count    int     `json:"count"`
	Percent  float64 `json:"percent"`
}

type TextStat struct {
	AvgLength float64 `json:"avg_length"`
	MaxLength int     `json:"max_length"`
	MinLength int     `json:"min_length"`
}

type RatingStat struct {
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
	StdDev  float64 `json:"std_dev"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Distribution map[int]int `json:"distribution"`
}

type WordCloudItem struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type TimeDataPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ChannelStat struct {
	Channel string  `json:"channel"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type CrossAnalysisQuery struct {
	SurveyID       uint `form:"survey_id" binding:"required"`
	RowQuestionID  uint `form:"row_question_id" binding:"required"`
	ColQuestionID  uint `form:"col_question_id" binding:"required"`
}

type CrossAnalysisResponse struct {
	RowQuestionID    uint              `json:"row_question_id"`
	ColQuestionID    uint              `json:"col_question_id"`
	RowOptions       []string          `json:"row_options"`
	ColOptions       []string          `json:"col_options"`
	Matrix           [][]int           `json:"matrix"`
	PercentMatrix    [][]float64       `json:"percent_matrix"`
	ChiSquare        float64           `json:"chi_square"`
	PValue           float64           `json:"p_value"`
	Significant      bool              `json:"significant"`
}
