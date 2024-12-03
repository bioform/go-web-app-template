package database

import (
	"log"

	"github.com/bioform/go-web-app-template/config"
)

var (
	Dsn string // default DB DSN
)

func init() {
	Dsn = config.App.Database.Dsn

	db, err := initSqliteDB(Dsn)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Panicf("cannot open db(%s): %v", Dsn, err)
	}
	defaultDB = db
}
