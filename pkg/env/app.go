package env

import (
	"os"
	"testing"
)

var appEnv string

func init() {
	if testing.Testing() {
		appEnv = "test"
		return
	}

	appEnv = Get("APP_ENV", "development")
}

func App() string {
	return appEnv
}

func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func IsProduction() bool {
	return App() == "production"
}

func IsTest() bool {
	return App() == "test"
}
