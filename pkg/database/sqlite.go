package database

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bioform/go-web-app-template/pkg/database/find"
	"github.com/bioform/go-web-app-template/pkg/database/schema"
	"github.com/bioform/go-web-app-template/pkg/env"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initSqliteDB(dsn string) (*gorm.DB, error) {
	if strings.Contains(dsn, ":memory:") || strings.Contains(dsn, "mode=memory") {
		slog.Info("Restore schema(in-memory DB)", slog.String("dsn", dsn))

		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
		if err != nil {
			return nil, fmt.Errorf("failed to open in-memory database: %w", err)
		}
		err = schema.Restore(db)
		if err != nil {
			return nil, fmt.Errorf("failed to restore schema: %w", err)
		}
		return db, nil
	}

	dsn, err := find.File(dsn)
	if err != nil && env.IsProduction() {
		slog.Warn("cannot find db file(%s): %v", dsn, err)
	}

	slog.Info("Connecting to database", slog.String("dsn", dsn))
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
}
