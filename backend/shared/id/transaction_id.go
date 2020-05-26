package id

import (
	"github.com/google/uuid"
)

// ID is a UUID used to trace a batch of work which is being processed.
type ID uuid.UUID

// String returns the transaction id as a string
func (id ID) String() (result string) {
	val := uuid.UUID(id)
	return val.String()
}

// New generates a new UUID
func New() ID {
	return ID(uuid.New())
}

// FromString parses a string into a transaction id
func FromString(idString string) (id ID, err error) {
	uID, err := uuid.Parse(idString)
	if err != nil {
		return id, err
	}
	return ID(uID), err
}
