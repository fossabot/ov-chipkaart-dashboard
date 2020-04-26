// +build ignore

package main

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func initializeRawRepositoryRecords(collection string, db *mongo.Database) RawRecordsRepository {
	wire.Build(
		NewMongodbRawRecordsRepository,
		wire.Bind(new(RawRecordsRepository), new(*MongodbRawRecordsRepository)),
		NewBsonService,
	)
	return &MongodbRawRecordsRepository{}
}
