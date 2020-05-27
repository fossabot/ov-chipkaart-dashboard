package validator

import (
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/pkg/errors"
)

var (
	//ErrInvalidEmailOrPassword is thrown when the user's email/password is wrong.
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

// ValidationResult stores the result of a validation
type ValidationResult struct {
	HasError bool
	Error    error
}

// Validator represents a validator
type Validator interface {
	ValidateCreateUserInput(input model.CreateUserInput) ValidationResult
	ValidateLoginInput(input model.LoginInput) ValidationResult
}
