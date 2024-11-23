package db

import "embed"

//go:embed migrations
var MigrationsFS embed.FS
var MigrationsDir = "migrations"
