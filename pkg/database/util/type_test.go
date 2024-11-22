package util

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDatabaseType(t *testing.T) {
	tests := []struct {
		name      string
		dialector gorm.Dialector
		expected  string
	}{
		{"SQLite", sqlite.Open("dsn"), "sqlite"},
		{"Unknown", nil, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var db *gorm.DB
			var err error
			if tt.dialector != nil {
				db, err = gorm.Open(tt.dialector, &gorm.Config{})
				if err != nil {
					t.Fatalf("failed to open database: %v", err)
				}
			}

			result := DatabaseType(db)
			if result != tt.expected {
				t.Errorf("expected %v; got %v", tt.expected, result)
			}
		})
	}
}
