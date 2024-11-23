package main

import (
	migration "github.com/bioform/go-web-app-template/db"
	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/pressly/goose/v3"
)

func main() {
	db := database.DefaultSqlDB
	// setup database

	goose.SetBaseFS(migration.MigrationsFS)

	if err := goose.Up(db, migration.MigrationsDir); err != nil {
		panic(err)
	}

	// run app
}
