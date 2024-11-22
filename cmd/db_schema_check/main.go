package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pressly/goose/v3"
)

func checkSchemaConsistency(schemaFilePath string) (bool, error) {
	// Get the most recent migration identifier.
	latestMigration, err := getLastMigrationIdentifier()
	if err != nil {
		return false, fmt.Errorf("failed to fetch latest migration: %w", err)
	}

	// Open the schema file for reading.
	file, err := os.Open(schemaFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to open schema file: %w", err)
	}
	defer file.Close()

	// Read only the first few lines of the file.
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "-- Latest Migration: ") {
			migrationID, err := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "-- Latest Migration: ")), 10, 64)
			if err != nil {
				return false, fmt.Errorf("invalid migration ID in schema file: %w", err)
			}
			return migrationID == latestMigration, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("failed to read schema file: %w", err)
	}

	return false, fmt.Errorf("latest migration identifier not found in schema file")
}

func getLastMigrationIdentifier() (int64, error) {
	// Path to your migrations folder
	migrationsPath := "./db/migrations"

	// Fetch all migrations
	migrations, err := goose.CollectMigrations(migrationsPath, 0, goose.MaxVersion)
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

func main() {
	ok, err := checkSchemaConsistency("db/schema.sql")
	if !ok || err != nil {
		log.Fatalf("Schema is inconsistent: %v", err)
	}
}
