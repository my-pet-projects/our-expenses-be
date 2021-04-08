package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideValidator_ReturnsValidator(t *testing.T) {
	results := ProvideValidator()

	assert.NotNil(t, results, "Validator should not be nil.")
}

func TestValidateStruct_ReturnsValidationErrors(t *testing.T) {
	category := &testRequest{
		Name:  "test",
		Email: "test",
	}

	validator := ProvideValidator()

	results := validator.ValidateStruct(category)

	assert.NotNil(t, results, "Should return validation errors.")
	assert.Len(t, results, 3, "Should contain 3 validation errors.")
	assert.Equal(t, results[0].Field, "ID", "Should fail 'ID' field validaton.")
	assert.Equal(t, results[1].Field, "Path", "Should fail 'Path' field validation.")
	assert.Equal(t, results[2].Field, "Email", "email")
}

func TestValidateStruct_ReturnsEmptyErrors(t *testing.T) {
	category := &testRequest{
		ID:    "12345",
		Path:  "path",
		Name:  "name",
		Email: "test@test.com",
	}

	validator := ProvideValidator()

	results := validator.ValidateStruct(category)

	assert.Nil(t, results, "Should return validation errors.")
	assert.Len(t, results, 0, "Should contain 0 validation errors.")
}

type testRequest struct {
	ID    string `json:"id" validate:"required,len=5"`
	Name  string `json:"name" validate:"required"`
	Path  string `json:"path" validate:"required"`
	Email string `json:"email" validate:"email"`
}
