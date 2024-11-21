package handler

import (
	"context"

	"github.com/bioform/go-web-app-template/internal/jwt"
	"github.com/bioform/go-web-app-template/internal/user/action"
	"github.com/bioform/go-web-app-template/pkg/server/session"
	"github.com/danielgtaylor/huma/v2"
)

type LoginInput struct {
	Body struct {
		Email    string `json:"email" maxLength:"255" example:"username@example.com"`
		Password string `json:"password" maxLength:"255" example:"complex_password"`
	}
}

type LoginOutput struct {
	Body struct {
		Token string
	}
}

func LoginHandler(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := action.AuthenticateUser(ctx).Call(input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, huma.Error404NotFound("User not found")
	}

	token, err := jwt.UserToken(user)
	if err != nil {
		return nil, err
	}

	session.Manager.Put(ctx, session.UserIdKey, int64(user.ID))

	result := &LoginOutput{}
	result.Body.Token = token

	return result, nil
}
