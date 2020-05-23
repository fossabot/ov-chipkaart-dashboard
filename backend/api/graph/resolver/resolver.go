package resolver

import "github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver resolves
type Resolver struct {
	db database.DB
}

// NewResolver creates a new instance of the resolver
func NewResolver(db database.DB) *Resolver {
	return &Resolver{db: db}
}
