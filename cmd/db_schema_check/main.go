package main

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/bioform/go-web-app-template/db"
	"github.com/bioform/go-web-app-template/pkg/database/schema"
)

func checkSchemaConsistency(fsys fs.FS) (bool, error) {
	// Get the most recent migration file identifier.
	latestMigration, err := schema.LastMigrationVersion()
	if err != nil {
		return false, fmt.Errorf("failed to fetch latest migration: %w", err)
	}

	// Get the schema version from the schema file.
	schemaVersion, err := schema.SchemaVersion(fsys)
	if err != nil {
		return false, fmt.Errorf("failed to fetch schema version: %w", err)
	}

	// Expect the schema version to match the latest migration.
	return schemaVersion == latestMigration, nil
}

func main() {
	ok, err := checkSchemaConsistency(db.MigrationsFS)
	if !ok || err != nil {
		log.Fatalf("Schema is inconsistent: %v", err)
	}
}
