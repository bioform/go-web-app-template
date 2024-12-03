package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/bioform/go-web-app-template/pkg/api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var defaultDB *gorm.DB

func Default() *gorm.DB {
	return defaultDB
}

func Use(ctx context.Context, tx *sql.Tx) (db *gorm.DB, err error) {
	api, err := api.From(ctx)
	if err != nil {
		return nil, err
	}

	dialect := api.DB().Dialector.Name()

	switch dialect {
	case "postgres":
		//db, err = gorm.Open(postgres.New(postgres.Config{Conn: tx}), &gorm.Config{})
		return nil, fmt.Errorf("not implemented")

	case "mysql":
		return nil, fmt.Errorf("not implemented")

	case "sqlite":
		// SQLite DSN is just the file path
		db, err = gorm.Open(sqlite.New(sqlite.Config{Conn: tx}), &gorm.Config{})

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dialect)
	}

	return db.WithContext(ctx), err
}

func CloseDefault() {
	db, err := defaultDB.DB()
	if err != nil {
		slog.Error("get sql.DB connection to close", slog.String("db", Dsn), slog.Any("error", err))
		return
	}

	if err := db.Close(); err != nil {
		slog.Error("close DB connection", slog.String("db", Dsn), slog.Any("error", err))
	} else {
		slog.Info("DB connection gracefully closed", slog.String("db", Dsn))
	}
}
