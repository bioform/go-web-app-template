package migrations

import (
	"context"
	"database/sql"
	"time"

	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	type User struct {
		*gorm.Model
		ConfirmedAt time.Time
	}
	confirmedAt := "ConfirmedAt"

	upAddConfirmedAtToUser := func(ctx context.Context, tx *sql.Tx) error {
		db, err := database.Use(ctx, tx)
		if err != nil {
			return err
		}

		err = db.Migrator().AddColumn(&User{}, confirmedAt)
		if err != nil {
			return err
		}

		err = db.Model(&User{}).Where("confirmed_at IS NULL").Update("confirmed_at", time.Now()).Error

		return err
	}

	downAddConfirmedAtToUser := func(ctx context.Context, tx *sql.Tx) error {
		db, err := database.Use(ctx, tx)
		if err != nil {
			return err
		}

		err = db.Migrator().DropColumn(&User{}, confirmedAt)
		return err
	}

	goose.AddMigrationContext(upAddConfirmedAtToUser, downAddConfirmedAtToUser)
}
