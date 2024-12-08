package action

import (
	"context"
	"fmt"

	"github.com/bioform/go-web-app-template/internal/user/email"
	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/internal/user/repository"
	"github.com/bioform/go-web-app-template/pkg/action"
	"github.com/bioform/go-web-app-template/pkg/api"
	validator "github.com/rezakhademix/govalidator/v2"
)

type RegisterUser struct {
	api.BaseAction

	Name     string
	Email    string
	Password string
	User     *model.User
	repo     repository.UserRepository
}

// The Init method initializes the action with the provided parameters.
func (a *RegisterUser) Init() {
	a.User = &model.User{
		Name:     a.Name,
		Email:    a.Email,
		Password: a.Password,
	}
}

// Update user repository each time the context is changed.
// This happens when the action is created or when the action is performed.
func (a *RegisterUser) SetContext(ctx context.Context) {
	a.BaseAction.SetContext(ctx) // Always call the parent method to set the context

	a.repo = repository.NewUserRepository(a.API().DB())
}

func (a *RegisterUser) Perform(ctx context.Context) error {
	_, err := a.repo.Create(ctx, a.User)
	if err != nil {
		return err
	}

	if err = email.SendConfirmationEmail(ctx, a.User); err != nil {
		return fmt.Errorf("send confirmation email: %w", err)
	}

	return nil
}

func (a *RegisterUser) IsValid(ctx context.Context) (bool, error) {

	v := validator.New()
	v.RequiredString(a.Name, "Name", "required")
	v.RequiredString(a.Email, "Email", "required")
	v.Email(a.Email, "Email", "invalid_format")
	v.RequiredString(a.Password, "Password", "required")
	v.CustomRule(a.repo.IsEmailUnique(ctx, a.Email), "Email", "already_taken")
	return v.IsPassed(), action.ErrorMap(v.Errors())
}
