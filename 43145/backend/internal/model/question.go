package model

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	SurveyID       uint           `json:"survey_id" gorm:"index;not null"`
	Survey         *Survey        `json:"survey,omitempty" gorm:"foreignKey:SurveyID"`
	Title          string         `json:"title" gorm:"type:text;not null"`
	Type           string         `json:"type" gorm:"size:20;not null;comment:single_choice,multi_choice,fill_in,rating,ranking,matrix"`
	IsRequired     bool           `json:"is_required" gorm:"default:true"`
	OrderIndex     int            `json:"order_index" gorm:"default:0"`
	Description    string         `json:"description" gorm:"type:text"`
	Placeholder    string         `json:"placeholder" gorm:"size:200"`
	MinValue       int            `json:"min_value" gorm:"default:0"`
	MaxValue       int            `json:"max_value" gorm:"default:0"`
	DefaultValue   string         `json:"default_value" gorm:"size:100"`
	ValidationRule string         `json:"validation_rule" gorm:"size:100;comment:regex or validation type"`
	DisplayLogic   string         `json:"display_logic" gorm:"type:text;comment:JSON condition config"`
	Options        []Option       `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
	LogicJumps     []LogicJump    `json:"logic_jumps,omitempty" gorm:"foreignKey:QuestionID"`
	Status         int            `json:"status" gorm:"default:1;comment:1=active,2=deleted"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type Option struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	QuestionID   uint           `json:"question_id" gorm:"index;not null"`
	Question     *Question      `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Text         string         `json:"text" gorm:"type:text;not null"`
	OrderIndex   int            `json:"order_index" gorm:"default:0"`
	IsOther      bool           `json:"is_other" gorm:"default:false"`
	JumpTarget   string         `json:"jump_target" gorm:"size:50;comment:question_id or end"`
	Score        int            `json:"score" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type LogicJump struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	QuestionID   uint           `json:"question_id" gorm:"index;not null"`
	Question     *Question      `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Condition    string         `json:"condition" gorm:"size:50;comment:equals,not_equals,contains,greater_than,less_than"`
	Value        string         `json:"value" gorm:"size:200"`
	JumpTo       string         `json:"jump_to" gorm:"size:50;comment:question_id or end"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
