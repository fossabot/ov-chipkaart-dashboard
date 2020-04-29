//+build wireinject

package main

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitializeRawRecordsRepository creates a new RawRecords repository instance
func InitializeRawRecordsRepository(collection string, db *mongo.Database) RawRecordsRepository {
	wire.Build(
		NewMongodbRawRecordsRepository,
		wire.Bind(new(RawRecordsRepository), new(*MongodbRawRecordsRepository)),
		NewBsonService,
	)
	return &MongodbRawRecordsRepository{}
}

// InitializeNationalHolidaysRepository creates a NationalHolidaysRepository
func InitializeNationalHolidaysRepository(collection string, db *mongo.Database) NationalHolidaysRepository {
	wire.Build(
		NewMongoNationalHolidaysRepository,
		wire.Bind(new(NationalHolidaysRepository), new(*MongoNationalHolidaysRepository)),
		NewBsonService,
	)
	return &MongoNationalHolidaysRepository{}

}
