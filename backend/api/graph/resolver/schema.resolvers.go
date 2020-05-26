package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/id"
	internalTime "github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/entities"
	internalErrors "github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/errors"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/generated"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/model"
	pkgErrors "github.com/pkg/errors"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.CreateUserOutput, error) {
	validationResult := r.validator.ValidateCreateUserInput(input)
	if validationResult.HasError {
		return nil, validationResult.Error
	}

	hashedPassword, err := r.passwordService.HashPassword(input.Password)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "could not hash password"))
		return nil, internalErrors.ErrInternalServerError
	}

	user := entities.User{
		ID:        id.New(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = r.db.UserRepository().Store(user)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "cannot save user in the database"))
		return nil, internalErrors.ErrInternalServerError
	}

	token, err := r.jwtService.GenerateTokenForUserID(user.ID)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrapf(err, "cannot generate jwt token for user with ID: %s", user.ID.String()))
		return nil, internalErrors.ErrInternalServerError
	}

	return &model.CreateUserOutput{
		User: &model.User{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(internalTime.DefaultFormat),
			UpdatedAt: user.UpdatedAt.Format(internalTime.DefaultFormat),
		},
		Token: &model.Token{
			Value: token,
		},
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CancelToken(ctx context.Context, input model.CancelTokenInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	return &model.User{}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
