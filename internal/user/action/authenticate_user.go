package action

import (
	"context"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/internal/user/repository"
)

type AuthenticateUserAction struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func AuthenticateUser() *AuthenticateUserAction {
	repo := repository.NewUserRepository()
	return &AuthenticateUserAction{
		repo: repo,
	}
}

// AuthenticateUser authenticates a user by email and password.
func (s *AuthenticateUserAction) Call(ctx context.Context, email, password string) (*model.User, error) {

	// Find the user by email
	user, err := s.repo.FindByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, err // Other database error
	}

	return user, nil // Authentication successful
}
