package database

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

var defaultDbProvider *DbProvider

func Get(ctx context.Context) *gorm.DB {
	db := defaultDbProvider.DB(ctx)
	return db.WithContext(ctx)
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
