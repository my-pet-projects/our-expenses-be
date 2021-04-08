package validator

import (
	"github.com/go-playground/validator/v10"
)

// ValidatorInterface defines validation methods.
type ValidatorInterface interface {
	ValidateStruct(object interface{}) []ValidationError
}

// Validator is a wrapper around "Package validator".
type Validator struct {
	*validator.Validate
}

// ValidationError is a struct with validation error details.
type ValidationError struct {
	Field   string `json:"field"`
	Details string `json:"details"`
}

var validationErrors = map[string]string{
	"required": "required, but was not received",
	"min":      "value or length is less than allowed",
	"max":      "value or length is bigger than allowed",
}

// ProvideValidator returns a Validator.
func ProvideValidator() *Validator {
	return &Validator{
		validator.New(),
	}
}

// ValidateStruct validates a struct.
func (v *Validator) ValidateStruct(object interface{}) []ValidationError {
	validationError := v.Struct(object)
	var errMsg []ValidationError
	if validationError != nil {
		for _, validationError := range validationError.(validator.ValidationErrors) {
			errMsg = append(errMsg, ValidationError{
				Field:   validationError.Field(),
				Details: getValidationMessage(validationError.ActualTag()),
			})
		}
	}
	return errMsg
}

func getValidationMessage(s string) string {
	if v, ok := validationErrors[s]; ok {
		return v
	}
	return s
}
