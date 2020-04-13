package main

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoNSEnrichedRecordsRepository is responsible for persisting/loading prices for NS journeys
type MongoNSEnrichedRecordsRepository struct {
	MongodbRepository
}

// NewMongoNSEnrichedRecordsRepository is used to initialize this class
func NewMongoNSEnrichedRecordsRepository(db *mongo.Database, collection string, bsonService BsonService) *MongoNSEnrichedRecordsRepository {
	return &MongoNSEnrichedRecordsRepository{MongodbRepository{db, collection, bsonService}}
}

// Store stores an NSJourneyPrice object into the mongodb repository
func (repository *MongoNSEnrichedRecordsRepository) Store(price NSJourneyPrice) (err error) {
	document, err := repository.bsonService.EncodeToBsonM(price)
	if err != nil {
		return errors.Wrap(err, "cannot convert struct to bson.M")
	}

	document = repository.SetTimestampFields(document)

	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)

	_, err = repository.db.Collection(repository.collection).InsertOne(ctx, document)
	if err != nil {
		return errors.Wrapf(err, "cannot insert document into db")
	}

	return nil
}

// GetByTransactionID returns the price of an NS journey repository based on the journey hash
func (repository *MongoNSEnrichedRecordsRepository) GetByTransactionID(id TransactionID) (enrichedRecords []EnrichedRecord, err error) {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	cursor, err := repository.db.Collection(repository.collection).Find(ctx, bson.M{"transaction_id": id.String()})
	if err != nil {
		return enrichedRecords, err
	}
	defer func() { _ = cursor.Close(context.Background()) }()

	for cursor.Next(ctx) {
		var record EnrichedRecord
		err := cursor.Decode(&EnrichedRecord{})
		if err != nil {
			return enrichedRecords, errors.Wrap(err, "cannot decode to bson.M to enriched record")
		}
		enrichedRecords = append(enrichedRecords, record)
	}

	err = cursor.Err()
	if err != nil {
		return enrichedRecords, errors.Wrap(err, "DB error")
	}

	return enrichedRecords, nil
}
