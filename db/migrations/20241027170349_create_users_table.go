package migrations

import (
	"context"
	"database/sql"

	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/util/crypt"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationContext(upAddSomeColumn, downAddSomeColumn)
}

func upAddSomeColumn(ctx context.Context, tx *sql.Tx) error {
	type User struct {
		*gorm.Model

		Name         string
		Email        string `gorm:"unique;not null"`
		PasswordHash string `gorm:"not null"`
	}
	db, err := database.Use(ctx, tx)
	if err != nil {
		return err
	}

	err = db.Migrator().CreateTable(&User{})
	if err != nil {
		return err
	}

	hashedPassword, err := crypt.HashPassword("password")
	if err != nil {
		return err
	}
	_ = db.Create(&User{Name: "admin", Email: "username@example.com", PasswordHash: hashedPassword})

	return err
}

func downAddSomeColumn(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE users;")
	return err
}
