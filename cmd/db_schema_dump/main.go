package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/database/schema"
	dbutil "github.com/bioform/go-web-app-template/pkg/database/schema"
)

func exportSchema(config dbutil.DBConfig) error {
	var cmd *exec.Cmd

	switch config.Type {
	case "postgres":
		cmd = exec.Command("pg_dump",
			"--schema-only",
			"--host", config.Host,
			"--port", config.Port,
			"--username", config.User,
			"--file", config.Output,
			config.Database)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", config.Password))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	case "mysql":
		cmd = exec.Command("mysqldump",
			"--no-data",
			"--host", config.Host,
			"--port", config.Port,
			"--user", config.User,
			fmt.Sprintf("--password=%s", config.Password),
			"--result-file", config.Output, // Added this option
			config.Database)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	case "sqlite":
		outputFile, err := os.Create(config.Output) // Create or overwrite the output file
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outputFile.Close()

		cmd = exec.Command("sqlite3", config.Database, ".schema")
		cmd.Stdout = outputFile // Redirect output to file
		cmd.Stderr = os.Stderr
	default:
		return fmt.Errorf("unsupported database type: %s", config.Type)
	}

	log.Printf("Exporting schema for %s database...", config.Type)
	return cmd.Run()
}

func addLastMigrationIdentifier(filePath, identifier string) error {
	originalContents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	schema := cleanupSchema(originalContents)

	// Add the last migration identifier to the beginning of the file.
	newContents := fmt.Sprintf("%s%s\n%s", dbutil.SchemaVersionPrefix, identifier, schema)

	// Write the updated contents back to the file.
	err = os.WriteFile(filePath, []byte(newContents), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	log.Printf("Last migration identifier added to the beginning of %s", filePath)
	return nil
}

// cleanupSchema removes lines related to `sqlite_sequence` from the schema content.
func cleanupSchema(schemaContent []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(schemaContent))
	var filteredLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Exclude lines that define or manipulate the `sqlite_sequence` table
		if strings.Contains(line, "CREATE TABLE sqlite_sequence") ||
			strings.Contains(line, "DELETE FROM sqlite_sequence") ||
			strings.Contains(line, "INSERT INTO sqlite_sequence") {
			continue
		}

		filteredLines = append(filteredLines, line)
	}

	return strings.Join(filteredLines, "\n")
}

func main() {
	db := database.GetDefault(context.Background())
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL database: %v", err)
	}

	defer database.CloseDefault()

	dbType := dbutil.DatabaseDialect(db)
	// Example configuration: Replace with your actual database details
	config, err := dbutil.ParseDSN(dbType, database.Dsn)
	if err != nil {
		log.Fatalf("Failed to parse DSN: %v", err)
	}

	if err := exportSchema(config); err != nil {
		log.Fatalf("Failed to export schema: %v", err)
	} else {
		log.Println("Schema export completed successfully!")
	}

	// Get the most recent migration version
	version, err := schema.DBVersion(sqlDb)
	if err != nil {
		log.Fatalf("Failed to get DB version: %v", err)
	}

	err = addLastMigrationIdentifier(config.Output, fmt.Sprintf("%d", version))
	if err != nil {
		log.Fatalf("Failed to add last migration identifier: %v", err)
	}
}
