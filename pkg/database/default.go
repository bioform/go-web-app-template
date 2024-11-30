package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/bioform/go-web-app-template/pkg/database/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var defaultDbProvider *DbProvider

func Get(ctx context.Context) *gorm.DB {
	db := defaultDbProvider.DB(ctx)
	return db.WithContext(ctx)
}

func Use(ctx context.Context, tx *sql.Tx) (db *gorm.DB, err error) {
	dialect := schema.GormDialect(Get(ctx))

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

func With(ctx context.Context, db *gorm.DB) context.Context {
	return defaultDbProvider.SetDB(ctx, db)
}

func CloseDefault() {
	db, err := Get(context.Background()).DB()
	if err != nil {
		slog.Error("Error getting DB connection", slog.String("db", Dsn), slog.Any("error", err))
		return
	}

	if err := db.Close(); err != nil {
		slog.Error("Error closing DB connection", slog.String("db", Dsn), slog.Any("error", err))
	} else {
		slog.Info("DB connection gracefully closed", slog.String("db", Dsn))
	}
}
