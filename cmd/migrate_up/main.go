package main

import (
	"context"
	"log"

	migration "github.com/bioform/go-web-app-template/db"
	_ "github.com/bioform/go-web-app-template/db/migrations"
	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/database/schema"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()

	gormDB := database.GetDefault(ctx)
	db, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	dialect := schema.DatabaseDialect(gormDB)
	goose.SetDialect(dialect)
	goose.SetBaseFS(migration.MigrationsFS)

	if err := goose.Up(db, migration.MigrationsDir); err != nil {
		log.Fatalf("failed to migrate up: %v", err)
	}
}