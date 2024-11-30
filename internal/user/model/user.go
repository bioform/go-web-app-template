package model

import (
	"fmt"
	"log/slog"

	"github.com/bioform/go-web-app-template/pkg/util/crypt"
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
	Password     string `gorm:"-"` // Don't create a `password` column in the database
	PasswordHash string `gorm:"not null"`
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(u.Password) > 0 {
		// Password was updated, hash it
		hashedPassword, err := crypt.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.PasswordHash = hashedPassword
		u.Password = "" // Clear the plain password after hashing
	}
	return
}

func (u *User) LogValue() slog.Value {
	if u == nil {
		return slog.Value{}
	}
	return slog.GroupValue(
		slog.Uint64("id", uint64(u.ID)),
	)
}

type EmailDuplicateError struct {
	Email string
}

func (e *EmailDuplicateError) Error() string {
	return fmt.Sprintf("Email '%s' already exists", e.Email)
}
