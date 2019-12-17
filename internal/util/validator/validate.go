package validator

import (
	"github.com/go-playground/validator/v10"
)

// Error stores error details
type Field struct {
	Param   string
	Message string
	Value   interface{}
	Tag     string
}

type Error struct {
	Param   string
	Message string
	Value   interface{}
}

// Val returns errors
func Val(fields ...Field) (errors []Error) {
	validate := validator.New()
	for _, field := range fields {
		err := validate.Var(field.Value, field.Tag)
		if err != nil {
			field.Tag = ""
			errors = append(errors, Error{
				field.Param,
				field.Message,
				field.Value,
			})
		}
	}
	if len(errors) > 0 {
		return
	}
	return nil
}

func Instance() *validator.Validate {
	return validator.New()
}

func InvalidValidationError(err error) bool {
	_, ok := err.(*validator.InvalidValidationError)
	return ok
}
