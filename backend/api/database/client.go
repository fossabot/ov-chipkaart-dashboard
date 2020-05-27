package database

import "github.com/pkg/errors"

var (
	// ErrEntityNotFound is returned when an entity does not exist in the database
	ErrEntityNotFound = errors.New("entity not found")
)

// DB is a collection of database repositories
type DB interface {
	UserRepository() UserRepository
}
