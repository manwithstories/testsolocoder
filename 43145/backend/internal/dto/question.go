package dto

type OptionRequest struct {
	Text       string `json:"text" binding:"required"`
	OrderIndex int    `json:"order_index"`
	IsOther    bool   `json:"is_other"`
	JumpTarget string `json:"jump_target"`
	Score      int    `json:"score"`
}

type LogicJumpRequest struct {
	Condition string `json:"condition" binding:"required,oneof=equals not_equals contains greater_than less_than"`
	Value     string `json:"value" binding:"required"`
	JumpTo    string `json:"jump_to" binding:"required"`
}

type CreateQuestionRequest struct {
	Title          string              `json:"title" binding:"required"`
	Type           string              `json:"type" binding:"required,oneof=single_choice multi_choice fill_in rating ranking matrix"`
	IsRequired     bool                `json:"is_required"`
	OrderIndex     int                 `json:"order_index"`
	Description    string              `json:"description"`
	Placeholder    string              `json:"placeholder"`
	MinValue       int                 `json:"min_value"`
	MaxValue       int                 `json:"max_value"`
	DefaultValue   string              `json:"default_value"`
	ValidationRule string              `json:"validation_rule"`
	DisplayLogic   string              `json:"display_logic"`
	Options        []OptionRequest     `json:"options"`
	LogicJumps     []LogicJumpRequest  `json:"logic_jumps"`
}

type UpdateQuestionRequest struct {
	Title          string              `json:"title" binding:"omitempty"`
	Type           string              `json:"type" binding:"omitempty,oneof=single_choice multi_choice fill_in rating ranking matrix"`
	IsRequired     *bool               `json:"is_required"`
	OrderIndex     *int                `json:"order_index"`
	Description    string              `json:"description"`
	Placeholder    string              `json:"placeholder"`
	MinValue       *int                `json:"min_value"`
	MaxValue       *int                `json:"max_value"`
	DefaultValue   string              `json:"default_value"`
	ValidationRule string              `json:"validation_rule"`
	DisplayLogic   string              `json:"display_logic"`
	Options        []OptionRequest     `json:"options"`
	LogicJumps     []LogicJumpRequest  `json:"logic_jumps"`
}

type BatchUpdateQuestionsRequest struct {
	Questions []CreateQuestionRequest `json:"questions" binding:"required"`
}

type ReorderQuestionRequest struct {
	OrderIndex int `json:"order_index" binding:"required,min=0"`
}
