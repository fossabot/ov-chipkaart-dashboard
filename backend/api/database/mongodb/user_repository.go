package mongodb

import (
	"context"
	"time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/entities"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/pkg/errors"
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
	_, err := repository.Collection().InsertOne(context.Background(), bson.M{
		"id":         user.ID.String(),
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
	return err
}

// FindByID finds a user in the database using it's ID
func (repository *UserRepository) FindByID(ID id.ID) (user *entities.User, err error) {
	dbUser := map[string]interface{}{}
	err = repository.Collection().FindOne(repository.DefaultTimeoutContext(), bson.M{"id": ID.String()}).Decode(&dbUser)
	if err != nil {
		return user, errors.Wrap(err, "error fetching single user from the database")
	}

	userID, err := id.FromString(dbUser["id"].(string))
	if err != nil {
		return user, errors.Wrap(err, "could not decode user id form string")
	}

	return &entities.User{
		ID:        userID,
		FirstName: dbUser["first_name"].(string),
		LastName:  dbUser["last_name"].(string),
		Email:     dbUser["email"].(string),
		Password:  dbUser["password"].(string),
		CreatedAt: dbUser["created_at"].(time.Time),
		UpdatedAt: dbUser["updated_at"].(time.Time),
	}, nil
}
