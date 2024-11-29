package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bioform/go-web-app-template/pkg/logging"
)

type DBHealthStats struct {
	Status            string
	Message           string
	OpenConnections   int
	InUse             int
	Idle              int
	WaitCount         int64
	WaitDuration      string
	MaxIdleClosed     int64
	MaxLifetimeClosed int64
	Error             string
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func Health(requestContext context.Context) DBHealthStats {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	log := logging.Logger(requestContext)

	stats := DBHealthStats{}

	// Ping the database
	db, err := GetDefault(requestContext).DB()
	if err != nil {
		setDownStatus(&stats, err)
		log.Error("db down", slog.Any("error", err))

		return stats
	}
	err = db.PingContext(ctx)
	if err != nil {
		setDownStatus(&stats, err)
		log.Error("db down", slog.Any("error", err))

		return stats
	}

	// Database is up, add more statistics
	stats.Status = "up"
	stats.Message = "It's healthy"

	// stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	// stats["in_use"] = strconv.Itoa(dbStats.InUse)
	// stats["idle"] = strconv.Itoa(dbStats.Idle)
	// stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	// stats["wait_duration"] = dbStats.WaitDuration.String()
	// stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	// stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := db.Stats()
	stats.OpenConnections = dbStats.OpenConnections
	stats.InUse = dbStats.InUse
	stats.Idle = dbStats.Idle
	stats.WaitCount = dbStats.WaitCount
	stats.WaitDuration = dbStats.WaitDuration.String()
	stats.MaxIdleClosed = dbStats.MaxIdleClosed
	stats.MaxLifetimeClosed = dbStats.MaxLifetimeClosed

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats.Message = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats.Message = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func setDownStatus(stats *DBHealthStats, err error) {
	stats.Status = "down"
	stats.Error = fmt.Sprintf("db down: %v", err)
}
