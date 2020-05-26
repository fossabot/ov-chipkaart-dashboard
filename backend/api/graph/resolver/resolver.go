package resolver

import (
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/validator"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/services/jwt"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/services/password"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/errorhandler"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver resolves
type Resolver struct {
	db              database.DB
	validator       validator.Validator
	passwordService password.Service
	errorHandler    errorhandler.ErrorHandler
	logger          logger.Logger
	jwtService      jwt.Service
}

// NewResolver creates a new instance of the resolver
func NewResolver(
	db database.DB,
	validator validator.Validator,
	passwordService password.Service,
	errorHandler errorhandler.ErrorHandler,
	logger logger.Logger,
	jwtService jwt.Service,
) *Resolver {

	return &Resolver{
		db:              db,
		validator:       validator,
		passwordService: passwordService,
		errorHandler:    errorHandler,
		logger:          logger,
		jwtService:      jwtService,
	}
}
