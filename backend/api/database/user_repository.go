package database

import (
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/entities"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/id"
)

// UserRepository is an instance of the user repository
type UserRepository interface {
	Store(user entities.User) error
	FindByID(userID id.ID) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
}
