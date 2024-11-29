package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	// "gorm.io/driver/mysql"
	// "gorm.io/driver/postgres"
	"github.com/bioform/go-web-app-template/config"
	"github.com/bioform/go-web-app-template/pkg/database/schema"
	"github.com/bioform/go-web-app-template/pkg/env"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Identify the database dialect by inspecting the DSN
func detectDialect(dsn string) (string, error) {
	if strings.HasPrefix(dsn, "mysql://") || strings.Contains(dsn, "@tcp(") {
		return "mysql", nil
	} else if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		return "postgres", nil
	} else if strings.HasSuffix(dsn, ".db") || strings.HasPrefix(dsn, "file:") {
		return "sqlite", nil
	}
	return "", errors.New("unable to detect database dialect")
}

// Create a new database and connect to it
func setupDatabase(dsn string) (*gorm.DB, error) {
	dialect, err := detectDialect(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to detect database dialect: %w", err)
	}

	var masterDSN, dbName string

	switch dialect {
	case "mysql":
		// Extract database name from DSN
		parts := strings.Split(dsn, "/")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid MySQL DSN format")
		}
		dbName = parts[len(parts)-1]
		if idx := strings.Index(dbName, "?"); idx != -1 {
			dbName = dbName[:idx]
		}
		// Use the "mysql" database as the master
		masterDSN = strings.Join(parts[:len(parts)-1], "/") + "/mysql"

	case "postgres":
		// Parse the DSN using url.Parse
		u, err := url.Parse(dsn)
		if err != nil {
			return nil, fmt.Errorf("invalid PostgreSQL DSN format: %w", err)
		}

		dbName = strings.TrimPrefix(u.Path, "/")
		if dbName == "" {
			return nil, errors.New("no database name found in PostgreSQL DSN")
		}
		// Replace the database name with "postgres"
		u.Path = "/postgres"
		masterDSN = u.String()

	case "sqlite":
		// Delete the database and recreate it
		if strings.Contains(dsn, ":memory:") {
			return nil, errors.New("cannot create schema for in-memory SQLite database")
		}
		err := os.Remove(dsn)
		if err != nil && os.IsExist(err) {
			return nil, fmt.Errorf("failed to delete SQLite database: %w", err)
		}
		// SQLite doesn't have a "master" DB; the database file will be created automatically
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
		}
		return db, nil

	default:
		return nil, errors.New("unsupported dialect")
	}

	// Connect to the master database
	sqlDB, err := sql.Open(dialect, masterDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master database: %w", err)
	}
	defer sqlDB.Close()

	// Drop the database if it already exists
	_, err = sqlDB.Exec("DROP DATABASE IF EXISTS \"?\"", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to drop database: %w", err)
	}

	// Create the database
	_, err = sqlDB.Exec("CREATE DATABASE ?", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	log.Printf("Database '%s' created successfully.\n", dbName)

	// Connect to the target database
	var db *gorm.DB
	switch dialect {
	case "mysql":
		//db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		db, err = nil, fmt.Errorf("MySQL is not supported yet")
	case "postgres":
		//db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		db, err = nil, fmt.Errorf("PostgreSQL is not supported yet")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to target database: %w", err)
	}

	return db, nil
}

func main() {
	if env.App() != "test" {
		log.Fatal("This command is only intended for use in the test environment.")
	}

	dsn := config.App.Database.Dsn

	db, err := setupDatabase(dsn)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	log.Println("Successfully connected to the database.")
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Import the schema.sql file.
	if err := schema.Restore(db); err != nil {
		log.Fatalf("failed to import schema: %v", err)
	}

	log.Printf("Schema imported into database %s successfully.", db.Name())
}
