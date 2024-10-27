package action

import (
	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/internal/user/repository"
)

type RegisterUser struct {
	repo repository.UserRepository
}

// NewRegisterUser creates a new instance of RegisterUser
func NewRegisterUser(repo repository.UserRepository) *RegisterUser {
	return &RegisterUser{repo: repo}
}

func (s *RegisterUser) Call(name, email, password string) error {
	// Business logic for registering a user
	newUser := model.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	return s.repo.Save(newUser)
}
