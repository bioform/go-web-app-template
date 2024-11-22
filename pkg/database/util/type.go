package util

import (
	"gorm.io/gorm"
)

func DatabaseType(db *gorm.DB) string {
	if db == nil {
		return "Unknown"
	}

	// Possible values: "sqlite", "mysql", "postgres", "Unknown"
	return db.Dialector.Name()
}
