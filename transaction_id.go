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

//
//func (id TransactionID) MarshalBSONValue() (bsontype.Type, []byte, error) {
//	log.Println("how are you doing today")
//	_, err :=  id.String()
//	if err != nil {
//		return bsontype.String, nil, err
//	}
//
//	log.Println("debugging")
//	return bsontype.String, bsoncore.AppendString(nil, `how`), nil
//}
//
//
//func (id *TransactionID) UnmarshalBSONValue(bsonType bsontype.Type, bytes []byte) error {
//	uid, err := uuid.FromBytes(bytes)
//	if err != nil {
//		return err
//	}
//
//	*id = TransactionID(uid)
//	return nil
//}

// NewTransactionID generates a new UUID
func NewTransactionID() TransactionID {
	return TransactionID(uuid.New())
}
