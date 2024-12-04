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
}

func (a *RegisterUser) Perform(ctx context.Context) error {
	newUser := model.User{
		Name:     a.Name,
		Email:    a.Email,
		Password: a.Password,
	}

	api, err := api.From(ctx)
	if err != nil {
		return err
	}

	repo := repository.NewUserRepository(api.DB())
	_, err = repo.Create(ctx, &newUser)
	if err != nil {
		return err
	}

	if err = email.SendConfirmationEmail(ctx, newUser); err != nil {
		return fmt.Errorf("send confirmation email: %w", err)
	}

	return nil
}

func (a *RegisterUser) IsValid(ctx context.Context) (bool, error) {
	api, err := api.From(ctx)
	if err != nil {
		return false, err
	}

	repo := repository.NewUserRepository(api.DB())

	v := validator.New()
	v.RequiredString(a.Name, "Name", "required")
	v.RequiredString(a.Email, "Email", "required")
	v.Email(a.Email, "Email", "invalid_format")
	v.RequiredString(a.Password, "Password", "required")
	v.CustomRule(repo.IsEmailUnique(ctx, a.Email), "Email", "already_taken")
	return v.IsPassed(), action.ErrorMap(v.Errors())
}
