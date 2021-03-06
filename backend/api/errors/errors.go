package errors

import "github.com/pkg/errors"

var (
	// ErrInternalServerError is thrown when there's a server error
	ErrInternalServerError = errors.New("internal server error")

	// ErrValidationError the input is invalid
	ErrValidationError = errors.New("input validation errors")
)
