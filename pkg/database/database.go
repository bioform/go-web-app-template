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

	db       *gorm.DB
	SQL_DB   *sql.DB
	MIGRATOR gorm.Migrator
)

func Get(ctx context.Context) *gorm.DB {
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

	SQL_DB, _ = db.DB()
	MIGRATOR = db.Migrator()
}

func Close() {
	if err := SQL_DB.Close(); err != nil {
		slog.Error("Error closing DB connection", slog.String("db", Dsn), slog.Any("error", err))
	} else {
		slog.Info("DB connection gracefully closed", slog.String("db", Dsn))
	}
}
