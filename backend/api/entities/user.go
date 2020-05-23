package entities

import "github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared"

// User is the user entity
type User struct {
	ID        shared.TransactionID
	FirstName string
	LastName  string
	Email     string
	Password  string
}
