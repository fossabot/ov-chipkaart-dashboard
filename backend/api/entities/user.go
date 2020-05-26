package entities

import (
	"time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/id"
)

// User is the user entity
type User struct {
	ID        id.ID
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
