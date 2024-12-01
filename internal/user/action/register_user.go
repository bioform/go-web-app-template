package action

import (
	"context"
	"fmt"

	"github.com/bioform/go-web-app-template/internal/user/email"
	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/internal/user/repository"
	"github.com/bioform/go-web-app-template/pkg/action"
	"github.com/bioform/go-web-app-template/pkg/dbaction"
	validator "github.com/rezakhademix/govalidator/v2"
)

type RegisterUser struct {
	dbaction.BaseAction

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

	repo := repository.NewUserRepository()
	_, err := repo.Create(ctx, &newUser)
	if err != nil {
		return err
	}

	if err = email.SendConfirmationEmail(ctx, newUser); err != nil {
		return fmt.Errorf("send confirmation email: %w", err)
	}

	return nil
}

func (a *RegisterUser) IsValid(ctx context.Context) (bool, action.ErrorMap) {
	repo := repository.NewUserRepository()

	v := validator.New()
	v.RequiredString(a.Name, "Name", "required")
	v.RequiredString(a.Email, "Email", "required")
	v.Email(a.Email, "Email", "invalid_format")
	v.RequiredString(a.Password, "Password", "required")
	v.CustomRule(repo.IsEmailUnique(ctx, a.Email), "Email", "already_taken")
	return v.IsPassed(), v.Errors()
}
