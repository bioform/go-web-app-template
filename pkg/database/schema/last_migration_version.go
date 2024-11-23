package schema

import (
	"log"

	"github.com/pressly/goose/v3"
)

const (
	// Path to your migrations folder
	MigrationsPath = "db/migrations"
)

func LastMigrationVersion() (int64, error) {
	// Fetch all migrations
	migrations, err := goose.CollectMigrations(MigrationsPath, 0, goose.MaxVersion)
	if err != nil {
		log.Fatalf("Failed to collect migrations: %v", err)
	}

	// Check if there are any migrations
	if len(migrations) == 0 {
		log.Println("No migrations found.")
		return 0, nil
	}

	// Get the most recent migration version
	migration, err := migrations.Last()
	if err != nil {
		log.Fatalf("Failed to get last migration: %v", err)
	}

	return migration.Version, nil
}
