package validator

import "github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/model"

// ValidationResult stores the result of a validation
type ValidationResult struct {
	HasError bool
	Error    error
}

// Validator represents a validator
type Validator interface {
	ValidateCreateUserInput(input model.CreateUserInput) ValidationResult
}
