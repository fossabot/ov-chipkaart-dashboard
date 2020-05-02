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
func (repository *MongoNSEnrichedRecordsRepository) Store(records []EnrichedRecord) (err error) {
	var documents []interface{}
	for _, record := range records {
		document, err := repository.bsonService.EncodeToBsonM(record)
		if err != nil {
			return errors.Wrapf(err, "cannot convert record to map")
		}
		documents = append(documents, repository.SetTimestampFields(document))
	}
	_, err = repository.db.Collection(repository.collection).InsertMany(context.Background(), documents)
	if err != nil {
		return errors.Wrapf(err, "cannot insert documents into db")
	}

	return nil
}

// FetchAllForTransactionID returns []EnrichedRecord based on on the transaction id
func (repository *MongoNSEnrichedRecordsRepository) FetchAllForTransactionID(id TransactionID) (enrichedRecords []EnrichedRecord, err error) {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	cursor, err := repository.db.Collection(repository.collection).Find(ctx, bson.M{"transaction_id": id.String()})
	if err != nil {
		return enrichedRecords, err
	}
	defer func() { _ = cursor.Close(context.Background()) }()

	for cursor.Next(ctx) {
		var record EnrichedRecord
		err := cursor.Decode(&record)
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
