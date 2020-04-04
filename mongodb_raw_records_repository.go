package main

import (
	"context"
	"time"

	"github.com/AchoArnold/homework/services/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "raw_records"

// MongoRawRecordsRepository is responsible for persisting/loading transactions which have not been processed
type MongoRawRecordsRepository struct {
	db *mongo.Database
}

const dbTimeout = 5 * time.Second

// NewRawRecordsRepository is used to initialize this class
func NewRawRecordsRepository(db *mongo.Database) *MongoRawRecordsRepository {
	return &MongoRawRecordsRepository{db: db}
}

// Store is responsible for storing the raw records in the database.
func (repository *MongoRawRecordsRepository) Store(records []Record, id TransactionID) (err error) {
	idString, err := id.String()
	if err != nil {
		return errors.Wrapf(err, "cannot convert transaction id to string")
	}

	var documents []interface{}
	for _, record := range records {
		valAsMap, err := json.ToInterfaceMap(record)
		if err != nil {
			return errors.Wrapf(err, "cannot convert record to map")
		}

		valAsMap["transactionId"] = idString
		documents = append(documents, bson.M(valAsMap))
	}
	_, err = repository.db.Collection(collection).InsertMany(context.Background(), documents)
	return errors.Wrapf(err, "cannot insert documents into DB")
}
