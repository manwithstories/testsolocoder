package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"museum-server/internal/dto"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid   bool              `json:"valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
	Message string            `json:"message,omitempty"`
}

var customMessages = map[string]string{
	"required":    "不能为空",
	"email":       "格式不正确",
	"min":         "长度不能小于 %s",
	"max":         "长度不能大于 %s",
	"gte":         "不能小于 %s",
	"lte":         "不能大于 %s",
	"oneof":       "必须是其中之一: %s",
	"url":         "格式不正确",
	"numeric":     "必须是数字",
	"alpha":       "只能包含字母",
	"alphanum":    "只能包含字母和数字",
	"datetime":    "日期时间格式不正确",
}

func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	if msg, ok := customMessages[tag]; ok {
		if param != "" {
			return fmt.Sprintf("%s %s", field, fmt.Sprintf(msg, param))
		}
		return fmt.Sprintf("%s %s", field, msg)
	}

	return fmt.Sprintf("%s 验证失败: %s", field, tag)
}

func ValidateStruct(s interface{}) ValidationResult {
	err := validate.Struct(s)
	if err == nil {
		return ValidationResult{Valid: true}
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: getErrorMessage(err),
		})
	}

	return ValidationResult{
		Valid:   false,
		Errors:  errors,
		Message: "数据验证失败",
	}
}

func ValidateCollectionRequest(req *dto.CollectionRequest) ValidationResult {
	return ValidateStruct(req)
}

func ValidateExhibitionRequest(req *dto.ExhibitionRequest) ValidationResult {
	result := ValidateStruct(req)
	if !result.Valid {
		return result
	}

	if req.EndDate.Before(req.StartDate) {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "end_date",
					Message: "结束日期不能早于开始日期",
				},
			},
			Message: "数据验证失败",
		}
	}

	if req.TicketPrice < 0 {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "ticket_price",
					Message: "ticket_price 不能小于 0",
				},
			},
			Message: "数据验证失败",
		}
	}

	if req.MaxVisitors < 0 {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "max_visitors",
					Message: "max_visitors 不能小于 0",
				},
			},
			Message: "数据验证失败",
		}
	}

	return result
}

func ValidateReservationRequest(req *dto.ReservationRequest) ValidationResult {
	result := ValidateStruct(req)
	if !result.Valid {
		return result
	}

	validGuideTypes := map[string]bool{
		"":         true,
		"self":     true,
		"audio":    true,
		"human":    true,
		"group":    true,
	}

	if req.GuideType != "" && !validGuideTypes[req.GuideType] {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "guide_type",
					Message: "guide_type 必须是 self、audio、human、group 之一",
				},
			},
			Message: "数据验证失败",
		}
	}

	return result
}

func ValidateTimeSlotRequest(req *dto.TimeSlotRequest) ValidationResult {
	result := ValidateStruct(req)
	if !result.Valid {
		return result
	}

	if req.StartTime >= req.EndTime {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "end_time",
					Message: "结束时间必须晚于开始时间",
				},
			},
			Message: "数据验证失败",
		}
	}

	if req.MaxCapacity <= 0 {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "max_capacity",
					Message: "max_capacity 必须大于 0",
				},
			},
			Message: "数据验证失败",
		}
	}

	return result
}

func ValidateBatchTimeSlotRequest(req *dto.BatchTimeSlotRequest) ValidationResult {
	result := ValidateStruct(req)
	if !result.Valid {
		return result
	}

	if req.EndDate.Before(req.StartDate) {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "end_date",
					Message: "结束日期不能早于开始日期",
				},
			},
			Message: "数据验证失败",
		}
	}

	if req.StartTime >= req.EndTime {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "end_time",
					Message: "结束时间必须晚于开始时间",
				},
			},
			Message: "数据验证失败",
		}
	}

	if req.Interval <= 0 {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "interval",
					Message: "interval 必须大于 0",
				},
			},
			Message: "数据验证失败",
		}
	}

	if req.MaxCapacity <= 0 {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "max_capacity",
					Message: "max_capacity 必须大于 0",
				},
			},
			Message: "数据验证失败",
		}
	}

	return result
}

func ValidateGuideScheduleRequest(req *dto.GuideScheduleRequest) ValidationResult {
	result := ValidateStruct(req)
	if !result.Valid {
		return result
	}

	if req.StartTime >= req.EndTime {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "end_time",
					Message: "结束时间必须晚于开始时间",
				},
			},
			Message: "数据验证失败",
		}
	}

	return result
}

func ValidateGuideContentRequest(req *dto.GuideContentRequest) ValidationResult {
	result := ValidateStruct(req)
	if !result.Valid {
		return result
	}

	validLanguages := map[string]bool{
		"zh": true,
		"en": true,
		"ja": true,
	}

	if !validLanguages[req.Language] {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "language",
					Message: "language 必须是 zh、en、ja 之一",
				},
			},
			Message: "数据验证失败",
		}
	}

	return result
}

func ValidateResearchApplicationRequest(req *dto.ResearchApplicationRequest) ValidationResult {
	return ValidateStruct(req)
}

func ValidateRegisterRequest(req *dto.RegisterRequest) ValidationResult {
	result := ValidateStruct(req)
	if !result.Valid {
		return result
	}

	if len(req.Password) < 6 {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "password",
					Message: "password 长度不能小于 6",
				},
			},
			Message: "数据验证失败",
		}
	}

	return result
}

func ValidateLoginRequest(req *dto.LoginRequest) ValidationResult {
	return ValidateStruct(req)
}

func ValidateUpdateUserRequest(req *dto.UpdateUserRequest) ValidationResult {
	return ValidateStruct(req)
}

func ValidateCollectionCategoryRequest(req *dto.CollectionCategoryRequest) ValidationResult {
	return ValidateStruct(req)
}

func ValidateMuseumRequest(req *dto.MuseumRequest) ValidationResult {
	return ValidateStruct(req)
}

func ValidateApplicationReviewRequest(req *dto.ApplicationReviewRequest) ValidationResult {
	return ValidateStruct(req)
}

func ValidateCollectionListQuery(query *dto.CollectionListQuery) ValidationResult {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	validSortFields := map[string]bool{
		"created_at": true,
		"name":       true,
		"code":       true,
		"era":        true,
		"view_count": true,
	}

	if query.SortBy != "" && !validSortFields[query.SortBy] {
		query.SortBy = "created_at"
	}

	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "desc"
	}

	return ValidationResult{Valid: true}
}

func ValidateExhibitionListQuery(query *dto.ExhibitionListQuery) ValidationResult {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	return ValidationResult{Valid: true}
}

func ValidateStatisticsQuery(query *dto.StatisticsQuery) ValidationResult {
	if !query.StartDate.IsZero() && !query.EndDate.IsZero() {
		if query.EndDate.Before(query.StartDate) {
			return ValidationResult{
				Valid: false,
				Errors: []ValidationError{
					{
						Field:   "end_date",
						Message: "结束日期不能早于开始日期",
					},
				},
				Message: "数据验证失败",
			}
		}
	}

	return ValidationResult{Valid: true}
}
