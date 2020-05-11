package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongodbRawRecordsRepository is responsible for persisting/loading transactions which have not been processed
type MongodbRawRecordsRepository struct {
	MongodbRepository
}

// NewMongodbRawRecordsRepository is used to initialize this class
func NewMongodbRawRecordsRepository(db *mongo.Database, collection string, bsonService BsonService) *MongodbRawRecordsRepository {
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

// First returns the first raw record in the repository
func (repository *MongodbRawRecordsRepository) First() (rawRecord RawRecord, err error) {
	ctx := context.Background()
	err = repository.db.Collection(repository.collection).FindOne(ctx, bson.M{}).Decode(&rawRecord)
	return rawRecord, err
}

// GetByTransactionID returns the price of an NS journey repository based on the journey hash
func (repository *MongodbRawRecordsRepository) GetByTransactionID(getOptions GetRawRecordsOptions) (rawRecords []RawRecord, err error) {
	order := 1
	if getOptions.SortDirection == "DESC" {
		order = -1
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{getOptions.SortBy, order}})

	ctx := context.Background()
	cursor, err := repository.db.Collection(repository.collection).Find(ctx, bson.M{"transaction_id": getOptions.TransactionID.String()}, findOptions)
	if err != nil {
		return rawRecords, err
	}
	defer func() { _ = cursor.Close(context.Background()) }()

	for cursor.Next(ctx) {
		var record RawRecord
		err := cursor.Decode(&record)
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
