package main

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoNSStationsRepository is responsible for persisting/loading prices for NS journeys
type MongoNSStationsRepository struct {
	MongodbRepository
}

// NewMongoNSStationsRepository is used to initialize this class
func NewMongoNSStationsRepository(db *mongo.Database, collection string, bsonService BsonService) *MongoNSStationsRepository {
	return &MongoNSStationsRepository{MongodbRepository{db, collection, bsonService}}
}

// Store stores a list of  NS station to the database
func (repository *MongoNSStationsRepository) Store(stations []NSStation) (err error) {
	var documents []interface{}
	for _, stations := range stations {
		document, err := repository.bsonService.EncodeToBsonM(stations.ToLower())
		if err != nil {
			return errors.Wrap(err, "cannot convert station to bson.M")
		}
		document = repository.SetTimestampFields(document)
		documents = append(documents, document)
	}

	_, err = repository.db.Collection(repository.collection).InsertMany(context.Background(), documents)
	if err != nil {
		return errors.Wrapf(err, "cannot insert stations into the database")
	}

	return nil
}

// GetByName fetches the first NS station with a particular name
func (repository *MongoNSStationsRepository) GetByName(name string) (station NSStation, err error) {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)

	log.Println("Db fetching for name = ", name)
	err = repository.db.Collection(repository.collection).FindOne(ctx, bson.M{"name": name}).Decode(&station)
	log.Println("finished fetching")

	if err == mongo.ErrNoDocuments {
		return station, ErrNotFound
	}
	if err != nil {
		return station, errors.Wrapf(err, "cannot fetch station by name")
	}

	return station, nil
}

// GetByCode fetches the first NS station with station code
func (repository *MongoNSStationsRepository) GetByCode(code string) (station NSStation, err error) {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)

	log.Println("Db fetching for code = ", code)
	err = repository.db.Collection(repository.collection).FindOne(ctx, bson.M{"code": code}).Decode(&station)
	log.Println("finished fetching")

	if err == mongo.ErrNoDocuments {
		return station, ErrNotFound
	}
	if err != nil {
		return station, errors.Wrapf(err, "cannot fetch station by code")
	}

	return station, nil
}
