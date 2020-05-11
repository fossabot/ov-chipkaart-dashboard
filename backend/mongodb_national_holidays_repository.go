package main

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoNationalHolidaysRepository is responsible for persisting/loading prices for national holidays
type MongoNationalHolidaysRepository struct {
	MongodbRepository
}

// NewMongoNationalHolidaysRepository is used to initialize this class
func NewMongoNationalHolidaysRepository(db *mongo.Database, collection string, bsonService BsonService) *MongoNationalHolidaysRepository {
	return &MongoNationalHolidaysRepository{MongodbRepository{db, collection, bsonService}}
}

// Store stores a slice of national holidays into the mongodb repository
func (repository *MongoNationalHolidaysRepository) Store(records []Holiday) (err error) {
	var documents []interface{}
	for _, record := range records {
		document, err := repository.bsonService.EncodeToBsonM(record)
		if err != nil {
			return errors.Wrapf(err, "cannot convert record to map")
		}
		document = repository.SetTimestampFields(document)
		documents = append(documents, repository.SetTimestampFields(document))
	}
	_, err = repository.db.Collection(repository.collection).InsertMany(context.Background(), documents)
	if err != nil {
		return errors.Wrapf(err, "cannot insert documents into db")
	}

	return nil
}

// HasHoliday checks if there is a national holiday for a given timestamp
func (repository *MongoNationalHolidaysRepository) HasHoliday(timestamp time.Time) (result bool, err error) {
	_, err = repository.GetByTimestamp(timestamp)
	if err == ErrNotFound {
		return false, nil
	}
	if err != nil {
		return result, errors.Wrapf(err, "cannot fetch holiday for timestamp")
	}

	return true, nil
}

// GetByTimestamp fetches the national holiday for a specific timestamp.
func (repository *MongoNationalHolidaysRepository) GetByTimestamp(timestamp time.Time) (holiday Holiday, err error) {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)

	err = repository.db.Collection(repository.collection).FindOne(ctx, bson.M{"date": timestamp.Format(dateFormat)}).Decode(&holiday)
	if err == mongo.ErrNoDocuments {
		return holiday, ErrNotFound
	}

	if err != nil {
		return holiday, errors.Wrapf(err, "cannot fetch holiday for timestamp by code")
	}

	return holiday, nil
}
