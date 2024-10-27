package logging

import (
	"log/slog"

	"github.com/bioform/go-web-app-template/pkg/env"
)

func logLevel() (slog.Level, error) {
	var logLevel string
	if env.IsProduction() {
		logLevel = env.Get("LOG_LEVEL", "info")
	} else {
		logLevel = env.Get("LOG_LEVEL", "debug")
	}

	// parse string, this is built-in feature of logrus
	level, err := parseLevel(logLevel)
	if err != nil {
		return slog.LevelDebug, err
	}

	return level, nil
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
