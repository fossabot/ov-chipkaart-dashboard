package database

import "github.com/pkg/errors"

var (
	// ErrEntityNotFound is thrown when the entity does not exist
	ErrEntityNotFound = errors.New("entity not found")
)

// DB is a collection of database repositories
type DB interface {
	UserRepository() UserRepository
}
