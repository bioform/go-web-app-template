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

	err := database.MIGRATOR.CreateTable(&User{})
	if err != nil {
		return err
	}

	db := database.GetDefault(ctx)

	hashedPassword, err := crypt.HashPassword("password")
	if err != nil {
		return err
	}
	_ = db.Create(&User{Name: "admin", Email: "admin@example.com", PasswordHash: hashedPassword})

	return err
}

func downAddSomeColumn(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE users;")
	return err
}
