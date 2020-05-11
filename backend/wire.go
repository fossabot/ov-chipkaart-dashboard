//+build wireinject

package main

import (
	lfucache "github.com/NdoleStudio/lfu-cache"
	"github.com/google/wire"
	"github.com/labstack/gommon/log"
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

// InitializeEnrichedRecordsRepository creates an enriched record repository
func InitializeEnrichedRecordsRepository(collection string, db *mongo.Database) EnrichedRecordsRepository {
	wire.Build(
		NewMongoNSEnrichedRecordsRepository,
		wire.Bind(new(EnrichedRecordsRepository), new(*MongoNSEnrichedRecordsRepository)),
		NewBsonService,
	)
	return &MongoNSEnrichedRecordsRepository{}
}

// InitializeCache creates a new LFU cache
func InitializeCache(size int) LFUCache {
	cache, err := lfucache.New(size)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return cache
}
