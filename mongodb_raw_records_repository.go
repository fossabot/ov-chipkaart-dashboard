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
func (repository *MongodbRawRecordsRepository) Store(records []RawRecord) (err error) {
	var documents []interface{}
	for _, record := range records {
		document, err := repository.bsonService.EncodeToBsonM(record)
		if err != nil {
			return errors.Wrapf(err, "cannot convert record to map")
		}

		document = repository.SetTimestampFields(document)

		documents = append(documents, document)
	}
	_, err = repository.db.Collection(repository.collection).InsertMany(context.Background(), documents)
	if err != nil {
		return errors.Wrapf(err, "cannot insert documents into db")
	}

	return nil
}

// GetByTransactionID returns the price of an NS journey repository based on the journey hash
func (repository *MongodbRawRecordsRepository) GetByTransactionID(options GetRawRecordsOptions) (rawRecords []RawRecord, err error) {
	filter, err := repository.bsonService.EncodeToBsonM(options)
	if err != nil {
		return rawRecords, err
	}

	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	cursor, err := repository.db.Collection(repository.collection).Find(ctx, filter)
	if err != nil {
		return rawRecords, err
	}
	defer func() { _ = cursor.Close(context.Background()) }()

	for cursor.Next(ctx) {
		var record RawRecord
		err := cursor.Decode(&RawRecord{})
		if err != nil {
			return rawRecords, errors.Wrap(err, "cannot decode to bson.M to enriched record")
		}
		rawRecords = append(rawRecords, record)
	}

	err = cursor.Err()
	if err != nil {
		return rawRecords, errors.Wrap(err, "DB error")
	}

	return rawRecords, nil
}
