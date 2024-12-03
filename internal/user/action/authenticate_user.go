package action

import (
	"context"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/internal/user/repository"
	"github.com/bioform/go-web-app-template/pkg/api"
)

type AuthenticateUserAction struct{}

// NewUserService creates a new instance of UserService
func AuthenticateUser() *AuthenticateUserAction {
	return &AuthenticateUserAction{}
}

// AuthenticateUser authenticates a user by email and password.
func (s *AuthenticateUserAction) Call(ctx context.Context, email, password string) (*model.User, error) {
	api, err := api.From(ctx)
	if err != nil {
		return nil, err
	}

	repo := repository.NewUserRepository(api.DB())
	// Find the user by email
	user, err := repo.FindByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, err // Other database error
	}

	return user, nil // Authentication successful
}
