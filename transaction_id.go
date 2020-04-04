package main

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// TransactionID is a UUID used to trace a batch of work which is being processed.
type TransactionID uuid.UUID

// String returns the transaction id as a string
func (id TransactionID) String() (result string, err error) {
	val, err := uuid.FromBytes(id[:])
	if err != nil {
		return result, errors.Wrapf(err, "cannot convert transaction ID %s to UUID", string(id[:]))
	}
	return val.String(), nil
}

// NewTransactionID generates a new UUID
func NewTransactionID() TransactionID {
	return TransactionID(uuid.New())
}
