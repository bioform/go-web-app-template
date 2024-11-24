package database

import (
	"log"

	"github.com/bioform/go-web-app-template/config"
	"github.com/bioform/go-web-app-template/pkg/action"
	"gorm.io/gorm"
)

var (
	Dsn string // default DB DSN

	MIGRATOR gorm.Migrator
)

func init() {
	Dsn = config.App.Database.Dsn

	db, err := initSqliteDB(Dsn)
	db.Transaction(func(tx *gorm.DB) error {
		return nil
	})
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Panicf("cannot open db(%s): %v", Dsn, err)
	}
	defaultDbProvider = &DbProvider{db: db}

	// inject default db provider to action package
	action.SetTransactionProvider(defaultDbProvider)

	MIGRATOR = db.Migrator()
}
