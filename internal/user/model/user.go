package model

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

const (
	UniqueConstraintUsername = "users_username_key"
	UniqueConstraintEmail    = "users_email_key"
)

type User struct {
	gorm.Model

	Name         string
	Email        string `gorm:"unique;not null"`
	Password     string `gorm:"-:migration"` // Don't create a `password` column in the database
	PasswordHash string `gorm:"not null"`
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("id", uint64(u.ID)),
		slog.String("email", u.Email),
	)
}

type EmailDuplicateError struct {
	Email string
}

func (e *EmailDuplicateError) Error() string {
	return fmt.Sprintf("Email '%s' already exists", e.Email)
}
