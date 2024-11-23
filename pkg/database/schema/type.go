package schema

import (
	"gorm.io/gorm"
)

func DatabaseDialect(db *gorm.DB) string {
	if db == nil {
		return "Unknown"
	}

	// Possible values: "sqlite", "mysql", "postgres", "Unknown"
	return db.Dialector.Name()
}
