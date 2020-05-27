package mongodb

import (
	"context"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/entities"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/id"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	dbRecord := map[string]interface{}{}
	err = repository.Collection().FindOne(repository.DefaultTimeoutContext(), bson.M{"id": ID.String()}).Decode(&dbRecord)

	if err == mongo.ErrNoDocuments {
		return user, database.ErrEntityNotFound
	}
	if err != nil {
		return user, errors.Wrap(err, "error fetching single user from the database by id")
	}

	return repository.hydrateUserFromDBRecord(dbRecord)
}

// FindByEmail searches a user using the email
func (repository *UserRepository) FindByEmail(email string) (user *entities.User, err error) {
	dbRecord := map[string]interface{}{}
	err = repository.Collection().FindOne(repository.DefaultTimeoutContext(), bson.M{"email": email}).Decode(&dbRecord)

	if err == mongo.ErrNoDocuments {
		return user, database.ErrEntityNotFound
	}
	if err != nil {
		return user, errors.Wrap(err, "error fetching single user from the database by email")
	}

	return repository.hydrateUserFromDBRecord(dbRecord)
}

func (repository *UserRepository) hydrateUserFromDBRecord(dbRecord map[string]interface{}) (user *entities.User, err error) {
	userID, err := id.FromString(dbRecord["id"].(string))
	if err != nil {
		return user, errors.Wrap(err, "could not decode user id form string")
	}

	return &entities.User{
		ID:        userID,
		FirstName: dbRecord["first_name"].(string),
		LastName:  dbRecord["last_name"].(string),
		Email:     dbRecord["email"].(string),
		Password:  dbRecord["password"].(string),
		CreatedAt: dbRecord["created_at"].(primitive.DateTime).Time(),
		UpdatedAt: dbRecord["updated_at"].(primitive.DateTime).Time(),
	}, err
}
