package govalidator

import (
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/validator"
)

// GoValidator is a validator using the govalidator package
type GoValidator struct {
	db database.DB
}

// New creates a new go validator
func New(db database.DB) validator.Validator {
	return &GoValidator{db}
}

// ValidateCreateUserInput validates the create user input request
func (validator GoValidator) ValidateCreateUserInput(input model.CreateUserInput) (result validator.ValidationResult) {
	return result
}
