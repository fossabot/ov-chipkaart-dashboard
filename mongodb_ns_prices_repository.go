package main

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoNSPricesRepository is responsible for persisting/loading prices for NS journeys
type MongoNSPricesRepository struct {
	MongodbRepository
}

// NewMongoNSPricesRepository is used to initialize this class
func NewMongoNSPricesRepository(db *mongo.Database, collection string, bsonService BsonService) *MongoNSPricesRepository {
	return &MongoNSPricesRepository{MongodbRepository{db, collection, bsonService}}
}

// Store stores an NSJourneyPrice object into the mongodb repository
func (repository *MongoNSPricesRepository) Store(price NSJourneyPrice) (err error) {
	document, err := repository.bsonService.EncodeToBsonM(price)
	if err != nil {
		return errors.Wrap(err, "cannot convert struct to map")
	}

	document = repository.SetTimestampFields(document)

	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)

	_, err = repository.db.Collection(repository.collection).InsertOne(ctx, document)
	if err != nil {
		return errors.Wrapf(err, "cannot insert document into db")
	}

	return nil
}

// GetByHash returns the price of an NS journey repository based on the journey hash
func (repository *MongoNSPricesRepository) GetByHash(hash string) (price NSJourneyPrice, err error) {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)

	err = repository.db.Collection(repository.collection).FindOne(ctx, bson.M{"hash": hash}).Decode(&price)

	if err == mongo.ErrNoDocuments {
		return price, ErrNotFound
	}
	if err != nil {
		return price, errors.Wrapf(err, "cannot fetch document by hash")
	}

	return price, nil
}
