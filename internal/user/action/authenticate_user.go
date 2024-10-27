package action

import (
	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/internal/user/repository"
)

type AuthenticateUser struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewAuthenticateUser(repo repository.UserRepository) *AuthenticateUser {
	return &AuthenticateUser{repo: repo}
}

// AuthenticateUser authenticates a user by email and password.
func (s *AuthenticateUser) Call(email, password string) (*model.User, error) {

	// Find the user by email
	user, err := s.repo.FindByEmailAndPassword(email, password)
	if err != nil {
		return nil, err // Other database error
	}

	return user, nil // Authentication successful
}
