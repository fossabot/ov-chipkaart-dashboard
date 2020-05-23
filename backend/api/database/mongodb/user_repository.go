package mongodb

import (
	"context"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository creates a new instance of the user repository
type UserRepository struct {
	repository
}

// NewUserRepository creates a new instance of the user repository
func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	return &UserRepository{repository{db, collection}}
}

// Store stores a user on the mongodb repository
func (repository *UserRepository) Store(user entities.User) error {
	_, err := repository.db.Collection(repository.collection).InsertOne(context.Background(), bson.M{})
	return err
}
