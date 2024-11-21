package jwt

import (
	"time"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/pkg/env"
	"github.com/golang-jwt/jwt/v4"
)

func UserToken(user *model.User) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"name":   user.Name,
		"tenant": "default",
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := env.Get("JWT_SECRET", "")
	signedString, err := token.SignedString([]byte(jwtSecret)) // Generate encoded token and send it as response.
	if err != nil {
		return "", err
	}

	return signedString, nil
}
