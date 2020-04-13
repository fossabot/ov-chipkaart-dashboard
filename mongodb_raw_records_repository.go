package main

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongodbRawRecordsRepository is responsible for persisting/loading transactions which have not been processed
type MongodbRawRecordsRepository struct {
	MongodbRepository
}

// NewRawRecordsRepository is used to initialize this class
func NewRawRecordsRepository(db *mongo.Database, collection string, bsonService BsonService) *MongodbRawRecordsRepository {
	return &MongodbRawRecordsRepository{MongodbRepository{db, collection, bsonService}}
}

// Store is responsible for storing the raw records in the database.
func (repository *MongodbRawRecordsRepository) Store(records []RawRecord, id TransactionID) (err error) {
	idString, err := id.String()
	if err != nil {
		return errors.Wrapf(err, "cannot convert transaction id to string")
	}

	var documents []interface{}
	for _, record := range records {
		document, err := repository.bsonService.EncodeToBsonM(record)
		if err != nil {
			return errors.Wrapf(err, "cannot convert record to map")
		}

		document = repository.SetTimestampFields(document)
		document = repository.SetTransactionIDField(document, idString)

		documents = append(documents, document)
	}
	_, err = repository.db.Collection(repository.collection).InsertMany(context.Background(), documents)
	if err != nil {
		return errors.Wrapf(err, "cannot insert documents into db")
	}

	return nil
}
