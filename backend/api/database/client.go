package database

// DB is a collection of database repositories
type DB interface {
	UserRepository() UserRepository
}
