package database

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bioform/go-web-app-template/pkg/env"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initSqliteDB(dsn string) (*gorm.DB, error) {
	sqliteDbPath, err := findFile(dsn)
	if err != nil && env.IsProduction() {
		log.Panicf("cannot find db file(%s): %v", dsn, err)
	}

	slog.Info("Connecting to database", slog.String("dsn", dsn))
	return gorm.Open(sqlite.Open(sqliteDbPath), &gorm.Config{TranslateError: true})
}

func findFile(path string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		absolutePath := filepath.Join(dir, path)
		if _, err := os.Stat(absolutePath); err == nil {
			return absolutePath, nil // found the root
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist // reached the filesystem root
		}
		dir = parent
	}
}
