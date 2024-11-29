# Define a reusable command for loading .env variables
build:
	go build -o dist/server cmd/api/main.go
	go build -o dist/migrate-up cmd/migrate_up/main.go
run:
	go run cmd/api/main.go
test:
	APP_ENV=test \
	ginkgo ./...
db-schema-dump:
	go run cmd/db_schema_dump/main.go
db-schema-check:
	go run cmd/db_schema_check/main.go
migrate:
	go run cmd/migrate/main.go $(filter-out $@,$(MAKECMDGOALS))
db-test-prepare:
	APP_ENV=test \
	go run cmd/db_test_prepare/main.go
migrate-up:
	dist/migrate-up
vet:
	go vet ./...


# Prevents Make from treating the arguments as targets
%:
	@:
