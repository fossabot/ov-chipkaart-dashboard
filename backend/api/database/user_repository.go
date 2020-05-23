package database

import "github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/entities"

// UserRepository is an instance of the user repository
type UserRepository interface {
	Store(user entities.User) error
}
