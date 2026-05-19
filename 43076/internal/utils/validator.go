package utils

import (
	"ticket-system/internal/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()

	validate.RegisterValidation("priority", func(fl validator.FieldLevel) bool {
		return models.ValidPriorities[fl.Field().String()]
	})

	validate.RegisterValidation("ticket_status", func(fl validator.FieldLevel) bool {
		return models.ValidStatuses[fl.Field().String()]
	})

	validate.RegisterValidation("ticket_type", func(fl validator.FieldLevel) bool {
		return models.ValidTypes[fl.Field().String()]
	})

	validate.RegisterValidation("assignment_mode", func(fl validator.FieldLevel) bool {
		return models.ValidAssignmentModes[fl.Field().String()]
	})

	validate.RegisterValidation("user_role", func(fl validator.FieldLevel) bool {
		return models.ValidRoles[fl.Field().String()]
	})

	validate.RegisterValidation("comment_type", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == models.CommentTypeInternal || fl.Field().String() == models.CommentTypePublic
	})
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
