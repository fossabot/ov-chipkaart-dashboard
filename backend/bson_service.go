package main

import (
	"go.mongodb.org/mongo-driver/bson"
)

// BsonService is the data structure for the BsonM service
type BsonService struct {
}

// NewBsonService creates a new instance of the service
func NewBsonService() BsonService {
	return BsonService{}
}

// EncodeToBsonM encodes a struct to bson.M
func (service BsonService) EncodeToBsonM(value interface{}) (result bson.M, err error) {
	encoded, err := bson.Marshal(value)
	if err != nil {
		return result, err
	}

	err = bson.Unmarshal(encoded, &result)
	return result, err
}
