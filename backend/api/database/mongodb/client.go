package mongodb

import (
	"context"
	"time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbOperationTimeout = 5 * time.Second
)

type repository struct {
	db         *mongo.Database
	collection string
}

func (repository repository) Collection() *mongo.Collection {
	return repository.db.Collection(repository.collection)
}

func (repository repository) DefaultTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	return ctx
}

// MongoDB is the struct for mongodb
type MongoDB struct {
	client *mongo.Database
}

// NewMongoDB creates a new instance of the mongodb client
func NewMongoDB(client *mongo.Database) database.DB {
	return &MongoDB{
		client: client,
	}
}

// UserRepository returns the user repository
func (db *MongoDB) UserRepository() database.UserRepository {
	return NewUserRepository(db.client, "users")
}
