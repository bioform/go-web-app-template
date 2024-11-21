package action

import (
	"context"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/internal/user/repository"
)

type AuthenticateUserAction struct {
	ctx  context.Context
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func AuthenticateUser(ctx context.Context) *AuthenticateUserAction {
	repo := repository.NewUserRepository(ctx)
	return &AuthenticateUserAction{
		ctx:  ctx,
		repo: repo,
	}
}

// AuthenticateUser authenticates a user by email and password.
func (s *AuthenticateUserAction) Call(email, password string) (*model.User, error) {

	// Find the user by email
	user, err := s.repo.FindByEmailAndPassword(email, password)
	if err != nil {
		return nil, err // Other database error
	}

	return user, nil // Authentication successful
}
