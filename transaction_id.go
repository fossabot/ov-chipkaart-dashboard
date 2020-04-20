package main

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// TransactionID is a UUID used to trace a batch of work which is being processed.
type TransactionID uuid.UUID

// String returns the transaction id as a string
func (id TransactionID) String() (result string) {
	val := uuid.UUID(id)
	return val.String()
}

// MarshalBSONValue converts a transaction id into a string for storing and easy searching
func (id TransactionID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, id.String()), nil
}

// UnmarshalBSONValue converts a transaction from a string into a transaction ID
func (id *TransactionID) UnmarshalBSONValue(_ bsontype.Type, raw []byte) error {
	val, _, _ := bsoncore.ReadString(raw)
	uid, err := uuid.Parse(val)
	if err != nil {
		return err
	}
	*id = TransactionID(uid)
	return nil
}

// NewTransactionID generates a new UUID
func NewTransactionID() TransactionID {
	return TransactionID(uuid.New())
}
