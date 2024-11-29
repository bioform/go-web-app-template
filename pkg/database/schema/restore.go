package schema

import (
	"fmt"

	"github.com/bioform/go-web-app-template/pkg/logging"
	"gorm.io/gorm"
)

func Restore(db *gorm.DB) error {
	log := logging.Logger(db.Statement.Context)
	log.Info("restoring schemafrom file", "file", dbSchemaFile)

	// Read schema.sql.
	schemaSQL, err := ReadSchemaFile()

	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Import the schema.sql file.
	if err := db.Exec(string(schemaSQL)).Error; err != nil {
		return fmt.Errorf("failed to import schema file: %w", err)
	}

	return nil
}
