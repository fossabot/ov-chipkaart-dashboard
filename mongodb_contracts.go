package main

import (
	"time"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

// DB Errors
var (
	// ErrNotFound is the error thrown when a document does not exist
	ErrNotFound = errors.New("item not found")
)

// collection names
const (
	collectionRawRecords        = "raw_records"
	collectionNSStations        = "ns_stations"
	collectionNSPrices          = "ns_journey_prices"
	collectionNSEnrichedRecords = "ns_enriched_records"
)

// db keys names
const (
	keyCreatedAt     = "created_at"
	keyUpdatedAt     = "updated_at"
	keyTransactionID = "transaction_id"
)

const dbOperationTimeout = 5 * time.Second

// MongodbRepository is the base data structure for mongodb struct
type MongodbRepository struct {
	db          *mongo.Database
	collection  string
	bsonService BsonService
}

// SetTimestampFields sets the timestamp fields for a bson.M document
func (repo MongodbRepository) SetTimestampFields(document bson.M) bson.M {
	document[keyCreatedAt] = time.Now().UTC().String()
	document[keyUpdatedAt] = time.Now().UTC().String()
	return document
}

// SetUpdatedAtField set the updated field for a bson.M document
func (repo MongodbRepository) SetUpdatedAtField(document bson.M) bson.M {
	document[keyUpdatedAt] = time.Now().UTC().String()
	return document
}
