package database

import (
	"context"
	"database/sql"
	"log"
	"log/slog"

	"github.com/bioform/go-web-app-template/config"
	"gorm.io/gorm"
)

var (
	Dsn string

	db           *gorm.DB // should be retrieved via GetDefault() only
	DefaultSqlDB *sql.DB
	MIGRATOR     gorm.Migrator
)

func GetDefault(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func init() {
	Dsn = config.App.Database.Dsn

	var err error
	db, err = initSqliteDB(Dsn)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Panicf("cannot open db(%s): %v", Dsn, err)
	}

	DefaultSqlDB, _ = db.DB()
	MIGRATOR = db.Migrator()
}

func CloseDefault() {
	if err := DefaultSqlDB.Close(); err != nil {
		slog.Error("Error closing DB connection", slog.String("db", Dsn), slog.Any("error", err))
	} else {
		slog.Info("DB connection gracefully closed", slog.String("db", Dsn))
	}
}

func Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		slog.Error("Error closing DB connection", slog.Any("error", err))
	} else {
		slog.Info("DB connection gracefully closed")
	}
}
