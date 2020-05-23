package shared

import (
	"github.com/google/uuid"
)

// TransactionID is a UUID used to trace a batch of work which is being processed.
type TransactionID uuid.UUID

// String returns the transaction id as a string
func (id TransactionID) String() (result string) {
	val := uuid.UUID(id)
	return val.String()
}

// NewTransactionID generates a new UUID
func NewTransactionID() TransactionID {
	return TransactionID(uuid.New())
}

// NewTransactionIDFromString parses a string into a transaction id
func NewTransactionIDFromString(id string) (transactionID TransactionID, err error) {
	uID, err := uuid.Parse(id)
	if err != nil {
		return transactionID, err
	}
	return TransactionID(uID), err
}
